package repo

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/chhlga/sh_adow/internal/config"
)

func TestResolveShadowPath_LocalRepo(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(filePath, []byte("test"), 0644)

	cfg := config.Config{RepoPath: "./"}

	shadowPath, err := ResolveShadowPath(filePath, cfg)
	if err != nil {
		t.Fatalf("ResolveShadowPath failed: %v", err)
	}

	expected := filepath.Join(tmpDir, ".shadow")
	if shadowPath != expected {
		t.Errorf("expected %s, got %s", expected, shadowPath)
	}
}

func TestResolveShadowPath_AbsolutePath(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "subdir", "test.txt")
	os.MkdirAll(filepath.Dir(filePath), 0755)
	os.WriteFile(filePath, []byte("test"), 0644)

	cacheDir := filepath.Join(tmpDir, "cache")
	cfg := config.Config{RepoPath: cacheDir}

	shadowPath, err := ResolveShadowPath(filePath, cfg)
	if err != nil {
		t.Fatalf("ResolveShadowPath failed: %v", err)
	}

	expected := filepath.Join(cacheDir, ".shadow")
	if shadowPath != expected {
		t.Errorf("expected %s, got %s", expected, shadowPath)
	}
}

func TestResolveShadowPath_RelativePath(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "project", "subdir", "test.txt")
	os.MkdirAll(filepath.Dir(filePath), 0755)
	os.WriteFile(filePath, []byte("test"), 0644)

	cfg := config.Config{RepoPath: "../cache"}

	shadowPath, err := ResolveShadowPath(filePath, cfg)
	if err != nil {
		t.Fatalf("ResolveShadowPath failed: %v", err)
	}

	expected := filepath.Join(tmpDir, "project", "cache", ".shadow")
	if shadowPath != expected {
		t.Errorf("expected %s, got %s", expected, shadowPath)
	}
}

func TestResolveShadowPath_TildeExpansion(t *testing.T) {
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	t.Cleanup(func() {
		os.Setenv("HOME", originalHome)
	})

	filePath := filepath.Join(tmpDir, "test.txt")
	os.WriteFile(filePath, []byte("test"), 0644)

	cfg := config.Config{RepoPath: "~/.shadow_backups"}

	shadowPath, err := ResolveShadowPath(filePath, cfg)
	if err != nil {
		t.Fatalf("ResolveShadowPath failed: %v", err)
	}

	expected := filepath.Join(tmpDir, ".shadow_backups", ".shadow")
	if shadowPath != expected {
		t.Errorf("expected %s, got %s", expected, shadowPath)
	}
}

func TestResolveShadowPath_Directory(t *testing.T) {
	tmpDir := t.TempDir()
	dirPath := filepath.Join(tmpDir, "subdir")
	os.MkdirAll(dirPath, 0755)

	cfg := config.Config{RepoPath: "./"}

	shadowPath, err := ResolveShadowPath(dirPath, cfg)
	if err != nil {
		t.Fatalf("ResolveShadowPath failed: %v", err)
	}

	expected := filepath.Join(dirPath, ".shadow")
	if shadowPath != expected {
		t.Errorf("expected %s, got %s", expected, shadowPath)
	}
}

func TestResolveShadowPath_NonexistentFile(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "nonexistent.txt")

	cfg := config.Config{RepoPath: "./"}

	shadowPath, err := ResolveShadowPath(filePath, cfg)
	if err != nil {
		t.Fatalf("ResolveShadowPath failed: %v", err)
	}

	expected := filepath.Join(tmpDir, ".shadow")
	if shadowPath != expected {
		t.Errorf("expected %s, got %s", expected, shadowPath)
	}
}

func TestEnsureShadowDir(t *testing.T) {
	tmpDir := t.TempDir()
	shadowPath := filepath.Join(tmpDir, ".shadow")

	err := EnsureShadowDir(shadowPath)
	if err != nil {
		t.Fatalf("EnsureShadowDir failed: %v", err)
	}

	stat, err := os.Stat(shadowPath)
	if err != nil {
		t.Fatalf("shadow directory not created: %v", err)
	}
	if !stat.IsDir() {
		t.Error("expected directory, got file")
	}

	snapshotsDir := filepath.Join(shadowPath, "snapshots")
	stat, err = os.Stat(snapshotsDir)
	if err != nil {
		t.Fatalf("snapshots directory not created: %v", err)
	}
	if !stat.IsDir() {
		t.Error("expected snapshots to be directory")
	}
}

func TestEnsureShadowDir_AlreadyExists(t *testing.T) {
	tmpDir := t.TempDir()
	shadowPath := filepath.Join(tmpDir, ".shadow")

	os.MkdirAll(filepath.Join(shadowPath, "snapshots"), 0755)

	err := EnsureShadowDir(shadowPath)
	if err != nil {
		t.Fatalf("EnsureShadowDir failed on existing directory: %v", err)
	}

	stat, err := os.Stat(shadowPath)
	if err != nil || !stat.IsDir() {
		t.Error("shadow directory should still exist")
	}
}

func TestEnsureShadowDir_NestedPath(t *testing.T) {
	tmpDir := t.TempDir()
	shadowPath := filepath.Join(tmpDir, "deep", "nested", "path", ".shadow")

	err := EnsureShadowDir(shadowPath)
	if err != nil {
		t.Fatalf("EnsureShadowDir failed with nested path: %v", err)
	}

	stat, err := os.Stat(filepath.Join(shadowPath, "snapshots"))
	if err != nil || !stat.IsDir() {
		t.Error("nested snapshots directory not created")
	}
}
