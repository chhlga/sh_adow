package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	cfg := DefaultConfig()

	if cfg.RepoPath != "./" {
		t.Errorf("expected default RepoPath './'. got %s", cfg.RepoPath)
	}
}

func TestLoad_NoConfigFile(t *testing.T) {
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	t.Cleanup(func() {
		os.Setenv("HOME", originalHome)
	})

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.RepoPath != "./" {
		t.Errorf("expected default RepoPath when config missing, got %s", cfg.RepoPath)
	}
}

func TestLoad_WithConfigFile(t *testing.T) {
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	t.Cleanup(func() {
		os.Setenv("HOME", originalHome)
	})

	configDir := filepath.Join(tmpDir, ".config", "sh_adow")
	os.MkdirAll(configDir, 0755)

	configPath := filepath.Join(configDir, "config.yml")
	configContent := "repo_path: \"~/.shadow_backups/\"\n"
	os.WriteFile(configPath, []byte(configContent), 0644)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.RepoPath != "~/.shadow_backups/" {
		t.Errorf("expected RepoPath '~/.shadow_backups/', got %s", cfg.RepoPath)
	}
}

func TestLoad_EmptyRepoPath(t *testing.T) {
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	t.Cleanup(func() {
		os.Setenv("HOME", originalHome)
	})

	configDir := filepath.Join(tmpDir, ".config", "sh_adow")
	os.MkdirAll(configDir, 0755)

	configPath := filepath.Join(configDir, "config.yml")
	configContent := "repo_path: \"\"\n"
	os.WriteFile(configPath, []byte(configContent), 0644)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.RepoPath != "./" {
		t.Errorf("expected default RepoPath when empty, got %s", cfg.RepoPath)
	}
}

func TestLoad_RelativePath(t *testing.T) {
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	t.Cleanup(func() {
		os.Setenv("HOME", originalHome)
	})

	configDir := filepath.Join(tmpDir, ".config", "sh_adow")
	os.MkdirAll(configDir, 0755)

	configPath := filepath.Join(configDir, "config.yml")
	configContent := "repo_path: \"../cache/\"\n"
	os.WriteFile(configPath, []byte(configContent), 0644)

	cfg, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if cfg.RepoPath != "../cache/" {
		t.Errorf("expected RepoPath '../cache/', got %s", cfg.RepoPath)
	}
}

func TestLoad_InvalidYAML(t *testing.T) {
	tmpDir := t.TempDir()
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tmpDir)
	t.Cleanup(func() {
		os.Setenv("HOME", originalHome)
	})

	configDir := filepath.Join(tmpDir, ".config", "sh_adow")
	os.MkdirAll(configDir, 0755)

	configPath := filepath.Join(configDir, "config.yml")
	invalidContent := "repo_path: [invalid: yaml\n"
	os.WriteFile(configPath, []byte(invalidContent), 0644)

	_, err := Load()
	if err == nil {
		t.Error("expected error when YAML is invalid")
	}
}
