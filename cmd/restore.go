package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/huh"
	"github.com/chhlga/sh_adow/internal/config"
	"github.com/chhlga/sh_adow/internal/repo"
	"github.com/chhlga/sh_adow/internal/shadow"
	"github.com/spf13/cobra"
)

var (
	restoreNoSave bool
)

var restoreCmd = &cobra.Command{
	Use:   "restore <file> <version-id>",
	Short: "Restore a file to a specific version",
	Args:  cobra.ExactArgs(2),
	RunE:  runRestore,
}

func init() {
	restoreCmd.Flags().BoolVar(&restoreNoSave, "no-save", false, "Don't save current state before restoring")
}

func runRestore(cmd *cobra.Command, args []string) error {
	filePath := args[0]
	versionID := args[1]

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	shadowPath, err := repo.ResolveShadowPath(filePath, cfg)
	if err != nil {
		return fmt.Errorf("failed to resolve shadow path: %w", err)
	}

	list, err := shadow.LoadList(shadowPath)
	if err != nil {
		return fmt.Errorf("failed to load list: %w", err)
	}

	absPath, _ := filepath.Abs(filePath)
	entry := list.FindFile(absPath)
	if entry == nil {
		return fmt.Errorf("file not tracked: %s", filePath)
	}

	var version *shadow.Version
	for i := range entry.Versions {
		if entry.Versions[i].ID == versionID {
			version = &entry.Versions[i]
			break
		}
	}
	if version == nil {
		return fmt.Errorf("version not found: %s", versionID)
	}

	var saveFirst bool
	if !restoreNoSave {
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Save current state before restoring?").
					Value(&saveFirst),
			),
		)

		if err := form.Run(); err != nil {
			return err
		}
	}

	if saveFirst {
		if _, err := os.Stat(filePath); err == nil {
			content, _ := os.ReadFile(filePath)
			newVersionID := shadow.GenerateVersionID(content)

			snapshotPath := filepath.Join(shadowPath, "snapshots", newVersionID)
			if err := shadow.CopyFile(filePath, snapshotPath); err != nil {
				return fmt.Errorf("failed to save current state: %w", err)
			}

			fileHash, _ := shadow.HashFile(filePath)
			stat, _ := os.Stat(filePath)

			newVersion := shadow.Version{
				ID:        newVersionID,
				CreatedAt: time.Now(),
				Tags:      []string{"auto-save"},
				Notes:     "Saved before restore",
				Size:      stat.Size(),
				Hash:      fileHash,
			}

			list.AddVersion(absPath, newVersion)
			if err := list.Save(shadowPath); err != nil {
				return fmt.Errorf("failed to save list: %w", err)
			}

			fmt.Printf("✓ Saved current state as %s\n", newVersionID)
		}
	}

	snapshotPath := filepath.Join(shadowPath, "snapshots", versionID)
	if err := shadow.CopyFile(snapshotPath, filePath); err != nil {
		return fmt.Errorf("failed to restore file: %w", err)
	}

	fmt.Printf("✓ Restored %s to version %s\n", filePath, versionID)
	return nil
}
