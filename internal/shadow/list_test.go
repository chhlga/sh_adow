package shadow

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestLoadList_EmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	shadowPath := filepath.Join(tmpDir, ".shadow")
	os.MkdirAll(shadowPath, 0755)

	list, err := LoadList(shadowPath)
	if err != nil {
		t.Fatalf("LoadList failed: %v", err)
	}

	if len(list.Files) != 0 {
		t.Errorf("expected empty list, got %d files", len(list.Files))
	}
}

func TestLoadList_WithData(t *testing.T) {
	tmpDir := t.TempDir()
	shadowPath := filepath.Join(tmpDir, ".shadow")
	os.MkdirAll(shadowPath, 0755)

	testData := List{
		Files: []FileEntry{
			{
				Path: "/tmp/test.txt",
				Versions: []Version{
					{
						ID:        "abc123",
						CreatedAt: time.Now(),
						Tags:      []string{"test"},
						Notes:     "test version",
						Size:      100,
						Hash:      "hash123",
					},
				},
			},
		},
	}

	data, _ := json.Marshal(testData)
	os.WriteFile(filepath.Join(shadowPath, "list.json"), data, 0644)

	list, err := LoadList(shadowPath)
	if err != nil {
		t.Fatalf("LoadList failed: %v", err)
	}

	if len(list.Files) != 1 {
		t.Errorf("expected 1 file, got %d", len(list.Files))
	}

	if list.Files[0].Path != "/tmp/test.txt" {
		t.Errorf("expected path /tmp/test.txt, got %s", list.Files[0].Path)
	}

	if len(list.Files[0].Versions) != 1 {
		t.Errorf("expected 1 version, got %d", len(list.Files[0].Versions))
	}

	v := list.Files[0].Versions[0]
	if v.ID != "abc123" {
		t.Errorf("expected version ID abc123, got %s", v.ID)
	}
	if v.Tags[0] != "test" {
		t.Errorf("expected tag 'test', got %s", v.Tags[0])
	}
}

func TestSaveList(t *testing.T) {
	tmpDir := t.TempDir()
	shadowPath := filepath.Join(tmpDir, ".shadow")
	os.MkdirAll(shadowPath, 0755)

	list := &List{
		Files: []FileEntry{
			{
				Path: "/tmp/test.txt",
				Versions: []Version{
					{
						ID:        "abc123",
						CreatedAt: time.Now(),
						Tags:      []string{"v1"},
						Notes:     "first version",
						Size:      50,
						Hash:      "hash1",
					},
				},
			},
		},
	}

	err := list.Save(shadowPath)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	listPath := filepath.Join(shadowPath, "list.json")
	if _, err := os.Stat(listPath); os.IsNotExist(err) {
		t.Fatal("list.json not created")
	}

	loaded, err := LoadList(shadowPath)
	if err != nil {
		t.Fatalf("LoadList after save failed: %v", err)
	}

	if len(loaded.Files) != 1 {
		t.Errorf("expected 1 file after save, got %d", len(loaded.Files))
	}
}

func TestFindFile(t *testing.T) {
	list := &List{
		Files: []FileEntry{
			{Path: "/tmp/file1.txt"},
			{Path: "/tmp/file2.txt"},
		},
	}

	entry := list.FindFile("/tmp/file1.txt")
	if entry == nil {
		t.Fatal("expected to find file1.txt")
	}
	if entry.Path != "/tmp/file1.txt" {
		t.Errorf("expected /tmp/file1.txt, got %s", entry.Path)
	}

	notFound := list.FindFile("/tmp/nonexistent.txt")
	if notFound != nil {
		t.Error("expected nil for nonexistent file")
	}
}

func TestAddVersion_NewFile(t *testing.T) {
	list := &List{Files: []FileEntry{}}

	version := Version{
		ID:        "v1",
		CreatedAt: time.Now(),
		Tags:      []string{"new"},
		Size:      100,
	}

	list.AddVersion("/tmp/newfile.txt", version)

	if len(list.Files) != 1 {
		t.Fatalf("expected 1 file, got %d", len(list.Files))
	}

	if list.Files[0].Path != "/tmp/newfile.txt" {
		t.Errorf("expected path /tmp/newfile.txt, got %s", list.Files[0].Path)
	}

	if len(list.Files[0].Versions) != 1 {
		t.Errorf("expected 1 version, got %d", len(list.Files[0].Versions))
	}

	if list.Files[0].Versions[0].ID != "v1" {
		t.Errorf("expected version ID v1, got %s", list.Files[0].Versions[0].ID)
	}
}

func TestAddVersion_ExistingFile(t *testing.T) {
	list := &List{
		Files: []FileEntry{
			{
				Path: "/tmp/existing.txt",
				Versions: []Version{
					{ID: "old", CreatedAt: time.Now().Add(-time.Hour)},
				},
			},
		},
	}

	newVersion := Version{
		ID:        "new",
		CreatedAt: time.Now(),
		Size:      200,
	}

	list.AddVersion("/tmp/existing.txt", newVersion)

	if len(list.Files) != 1 {
		t.Errorf("expected 1 file, got %d", len(list.Files))
	}

	if len(list.Files[0].Versions) != 2 {
		t.Fatalf("expected 2 versions, got %d", len(list.Files[0].Versions))
	}

	if list.Files[0].Versions[0].ID != "new" {
		t.Errorf("expected newest version first, got %s", list.Files[0].Versions[0].ID)
	}

	if list.Files[0].Versions[1].ID != "old" {
		t.Errorf("expected old version second, got %s", list.Files[0].Versions[1].ID)
	}
}

func TestRemoveVersion(t *testing.T) {
	list := &List{
		Files: []FileEntry{
			{
				Path: "/tmp/test.txt",
				Versions: []Version{
					{ID: "v1"},
					{ID: "v2"},
					{ID: "v3"},
				},
			},
		},
	}

	removed := list.RemoveVersion("/tmp/test.txt", "v2")
	if !removed {
		t.Fatal("expected RemoveVersion to return true")
	}

	if len(list.Files[0].Versions) != 2 {
		t.Fatalf("expected 2 versions remaining, got %d", len(list.Files[0].Versions))
	}

	ids := []string{list.Files[0].Versions[0].ID, list.Files[0].Versions[1].ID}
	if ids[0] == "v2" || ids[1] == "v2" {
		t.Error("v2 should have been removed")
	}
}

func TestRemoveVersion_LastVersion(t *testing.T) {
	list := &List{
		Files: []FileEntry{
			{
				Path: "/tmp/test.txt",
				Versions: []Version{
					{ID: "only"},
				},
			},
		},
	}

	removed := list.RemoveVersion("/tmp/test.txt", "only")
	if !removed {
		t.Fatal("expected RemoveVersion to return true")
	}

	if len(list.Files) != 0 {
		t.Errorf("expected file entry removed when last version deleted, got %d files", len(list.Files))
	}
}

func TestRemoveVersion_NotFound(t *testing.T) {
	list := &List{
		Files: []FileEntry{
			{
				Path:     "/tmp/test.txt",
				Versions: []Version{{ID: "v1"}},
			},
		},
	}

	removed := list.RemoveVersion("/tmp/test.txt", "nonexistent")
	if removed {
		t.Error("expected RemoveVersion to return false for nonexistent version")
	}

	removed = list.RemoveVersion("/tmp/nonexistent.txt", "v1")
	if removed {
		t.Error("expected RemoveVersion to return false for nonexistent file")
	}
}

func TestGenerateVersionID(t *testing.T) {
	content1 := []byte("test content")
	content2 := []byte("different content")

	id1 := GenerateVersionID(content1)
	id2 := GenerateVersionID(content2)

	if id1 == id2 {
		t.Error("different content should generate different IDs")
	}

	if len(id1) != 8 {
		t.Errorf("expected version ID length 8, got %d", len(id1))
	}

	id1Again := GenerateVersionID(content1)
	if id1 != id1Again {
		t.Error("same content should generate same ID")
	}
}

func TestCopyFile(t *testing.T) {
	tmpDir := t.TempDir()

	srcPath := filepath.Join(tmpDir, "source.txt")
	dstPath := filepath.Join(tmpDir, "subdir", "dest.txt")

	content := "test file content"
	os.WriteFile(srcPath, []byte(content), 0644)

	err := CopyFile(srcPath, dstPath)
	if err != nil {
		t.Fatalf("CopyFile failed: %v", err)
	}

	readContent, err := os.ReadFile(dstPath)
	if err != nil {
		t.Fatalf("failed to read destination file: %v", err)
	}

	if string(readContent) != content {
		t.Errorf("expected content %q, got %q", content, string(readContent))
	}
}

func TestCopyFile_SourceNotFound(t *testing.T) {
	tmpDir := t.TempDir()
	err := CopyFile("/nonexistent/source.txt", filepath.Join(tmpDir, "dest.txt"))
	if err == nil {
		t.Error("expected error when source file doesn't exist")
	}
}

func TestHashFile(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")

	content := "test content for hashing"
	os.WriteFile(filePath, []byte(content), 0644)

	hash, err := HashFile(filePath)
	if err != nil {
		t.Fatalf("HashFile failed: %v", err)
	}

	if len(hash) != 64 {
		t.Errorf("expected SHA256 hex length 64, got %d", len(hash))
	}

	hashAgain, _ := HashFile(filePath)
	if hash != hashAgain {
		t.Error("same file should produce same hash")
	}

	os.WriteFile(filePath, []byte("modified content"), 0644)
	hashModified, _ := HashFile(filePath)
	if hash == hashModified {
		t.Error("modified file should produce different hash")
	}
}

func TestHashFile_NotFound(t *testing.T) {
	_, err := HashFile("/nonexistent/file.txt")
	if err == nil {
		t.Error("expected error when file doesn't exist")
	}
}
