package integration

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/chhlga/sh_adow/internal/config"
	"github.com/chhlga/sh_adow/internal/repo"
	"github.com/chhlga/sh_adow/internal/shadow"
)

func setupTestEnv(t *testing.T) (string, string) {
	t.Helper()

	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(testFile, []byte("initial content"), 0644)

	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	t.Cleanup(func() {
		os.Setenv("HOME", originalHome)
	})

	return tmpDir, testFile
}

func TestWorkflow_SaveListRestore(t *testing.T) {
	_, testFile := setupTestEnv(t)

	cfg, _ := config.Load()
	shadowPath, _ := repo.ResolveShadowPath(testFile, cfg)
	repo.EnsureShadowDir(shadowPath)

	os.WriteFile(testFile, []byte("version 1"), 0644)

	list := &shadow.List{Files: []shadow.FileEntry{}}
	content, _ := os.ReadFile(testFile)
	versionID := shadow.GenerateVersionID(content)

	snapshotPath := filepath.Join(shadowPath, "snapshots", versionID)
	shadow.CopyFile(testFile, snapshotPath)

	version := shadow.Version{
		ID:    versionID,
		Tags:  []string{"v1"},
		Notes: "first version",
		Size:  9,
	}

	absPath, _ := filepath.Abs(testFile)
	list.AddVersion(absPath, version)
	list.Save(shadowPath)

	loadedList, err := shadow.LoadList(shadowPath)
	if err != nil {
		t.Fatalf("failed to load list: %v", err)
	}

	entry := loadedList.FindFile(absPath)
	if entry == nil {
		t.Fatal("file not found in list")
	}

	if len(entry.Versions) != 1 {
		t.Fatalf("expected 1 version, got %d", len(entry.Versions))
	}

	if entry.Versions[0].ID != versionID {
		t.Errorf("version ID mismatch: expected %s, got %s", versionID, entry.Versions[0].ID)
	}

	os.WriteFile(testFile, []byte("version 2 - modified"), 0644)

	if err := shadow.CopyFile(snapshotPath, testFile); err != nil {
		t.Fatalf("failed to restore: %v", err)
	}

	restoredContent, _ := os.ReadFile(testFile)
	if string(restoredContent) != "version 1" {
		t.Errorf("expected 'version 1', got '%s'", string(restoredContent))
	}
}

func TestWorkflow_MultipleVersions(t *testing.T) {
	_, testFile := setupTestEnv(t)

	cfg, _ := config.Load()
	shadowPath, _ := repo.ResolveShadowPath(testFile, cfg)
	repo.EnsureShadowDir(shadowPath)

	list := &shadow.List{Files: []shadow.FileEntry{}}
	absPath, _ := filepath.Abs(testFile)

	versions := []string{"content v1", "content v2", "content v3"}

	for i, content := range versions {
		os.WriteFile(testFile, []byte(content), 0644)

		fileContent, _ := os.ReadFile(testFile)
		versionID := shadow.GenerateVersionID(fileContent)

		snapshotPath := filepath.Join(shadowPath, "snapshots", versionID)
		shadow.CopyFile(testFile, snapshotPath)

		version := shadow.Version{
			ID:    versionID,
			Tags:  []string{string(rune('A' + i))},
			Notes: content,
			Size:  int64(len(content)),
		}

		list.AddVersion(absPath, version)
	}

	list.Save(shadowPath)

	loadedList, _ := shadow.LoadList(shadowPath)
	entry := loadedList.FindFile(absPath)

	if len(entry.Versions) != 3 {
		t.Fatalf("expected 3 versions, got %d", len(entry.Versions))
	}

	if entry.Versions[0].Tags[0] != "C" {
		t.Error("newest version should be first")
	}
}

func TestWorkflow_DeleteVersion(t *testing.T) {
	_, testFile := setupTestEnv(t)

	cfg, _ := config.Load()
	shadowPath, _ := repo.ResolveShadowPath(testFile, cfg)
	repo.EnsureShadowDir(shadowPath)

	list := &shadow.List{Files: []shadow.FileEntry{}}
	absPath, _ := filepath.Abs(testFile)

	content := "test content"
	os.WriteFile(testFile, []byte(content), 0644)

	fileContent, _ := os.ReadFile(testFile)
	versionID := shadow.GenerateVersionID(fileContent)

	snapshotPath := filepath.Join(shadowPath, "snapshots", versionID)
	shadow.CopyFile(testFile, snapshotPath)

	version := shadow.Version{
		ID:   versionID,
		Tags: []string{"deleteme"},
		Size: int64(len(content)),
	}

	list.AddVersion(absPath, version)
	list.Save(shadowPath)

	removed := list.RemoveVersion(absPath, versionID)
	if !removed {
		t.Fatal("failed to remove version")
	}

	list.Save(shadowPath)

	os.Remove(snapshotPath)

	if _, err := os.Stat(snapshotPath); !os.IsNotExist(err) {
		t.Error("snapshot file should be deleted")
	}

	loadedList, _ := shadow.LoadList(shadowPath)
	if len(loadedList.Files) != 0 {
		t.Error("file entry should be removed when last version deleted")
	}
}

func TestWorkflow_ConfigRepoPath(t *testing.T) {
	tmpDir := t.TempDir()
	testFile := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(testFile, []byte("test"), 0644)

	cacheDir := filepath.Join(tmpDir, "cache")
	cfg := config.Config{RepoPath: cacheDir}

	shadowPath, err := repo.ResolveShadowPath(testFile, cfg)
	if err != nil {
		t.Fatalf("ResolveShadowPath failed: %v", err)
	}

	expected := filepath.Join(cacheDir, ".shadow")
	if shadowPath != expected {
		t.Errorf("expected %s, got %s", expected, shadowPath)
	}

	repo.EnsureShadowDir(shadowPath)

	if !fileExists(filepath.Join(shadowPath, "snapshots")) {
		t.Error("snapshots directory not created in custom cache location")
	}
}

func TestWorkflow_HashConsistency(t *testing.T) {
	_, testFile := setupTestEnv(t)

	content := "consistent content"
	os.WriteFile(testFile, []byte(content), 0644)

	hash1, _ := shadow.HashFile(testFile)
	hash2, _ := shadow.HashFile(testFile)

	if hash1 != hash2 {
		t.Error("same file should produce same hash")
	}

	versionID1 := shadow.GenerateVersionID([]byte(content))
	versionID2 := shadow.GenerateVersionID([]byte(content))

	if versionID1 != versionID2 {
		t.Error("same content should produce same version ID")
	}

	os.WriteFile(testFile, []byte("different content"), 0644)
	hash3, _ := shadow.HashFile(testFile)

	if hash1 == hash3 {
		t.Error("different content should produce different hash")
	}
}

func TestWorkflow_ConcurrentFiles(t *testing.T) {
	tmpDir := t.TempDir()

	file1 := filepath.Join(tmpDir, "file1.txt")
	file2 := filepath.Join(tmpDir, "file2.txt")
	os.WriteFile(file1, []byte("content1"), 0644)
	os.WriteFile(file2, []byte("content2"), 0644)

	cfg := config.Config{RepoPath: "./"}
	shadowPath, _ := repo.ResolveShadowPath(tmpDir, cfg)
	repo.EnsureShadowDir(shadowPath)

	list := &shadow.List{Files: []shadow.FileEntry{}}

	for _, file := range []string{file1, file2} {
		content, _ := os.ReadFile(file)
		versionID := shadow.GenerateVersionID(content)

		snapshotPath := filepath.Join(shadowPath, "snapshots", versionID)
		shadow.CopyFile(file, snapshotPath)

		absPath, _ := filepath.Abs(file)
		version := shadow.Version{
			ID:   versionID,
			Size: int64(len(content)),
		}
		list.AddVersion(absPath, version)
	}

	list.Save(shadowPath)

	loadedList, _ := shadow.LoadList(shadowPath)
	if len(loadedList.Files) != 2 {
		t.Fatalf("expected 2 files tracked, got %d", len(loadedList.Files))
	}
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}
