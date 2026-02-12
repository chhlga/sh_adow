package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	RepoPath string `yaml:"repo_path"`
}

// DefaultConfig returns default configuration
func DefaultConfig() Config {
	return Config{
		RepoPath: "./",
	}
}

// Load reads config from ~/.config/sh_adow/config.yml
func Load() (Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return DefaultConfig(), nil
	}

	configPath := filepath.Join(home, ".config", "sh_adow", "config.yml")

	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return Config{}, err
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	if cfg.RepoPath == "" {
		cfg.RepoPath = "./"
	}

	return cfg, nil
}
