package repo

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/chhlga/sh_adow/internal/config"
)

func ResolveShadowPath(filePath string, cfg config.Config) (string, error) {
	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return "", err
	}

	fileDir := filepath.Dir(absFilePath)
	stat, err := os.Stat(absFilePath)
	if err == nil && stat.IsDir() {
		fileDir = absFilePath
	}

	repoPath := cfg.RepoPath
	if strings.HasPrefix(repoPath, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		repoPath = filepath.Join(home, repoPath[2:])
	}

	var shadowBase string
	if filepath.IsAbs(repoPath) {
		shadowBase = repoPath
	} else if repoPath == "./" {
		shadowBase = fileDir
	} else {
		shadowBase = filepath.Join(fileDir, repoPath)
	}

	return filepath.Join(shadowBase, ".shadow"), nil
}

func EnsureShadowDir(shadowPath string) error {
	snapshotsDir := filepath.Join(shadowPath, "snapshots")
	return os.MkdirAll(snapshotsDir, 0755)
}
