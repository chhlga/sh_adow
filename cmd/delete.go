package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/chhlga/sh_adow/internal/config"
	"github.com/chhlga/sh_adow/internal/repo"
	"github.com/chhlga/sh_adow/internal/shadow"
	"github.com/spf13/cobra"
)

var (
	deleteForce bool
)

var deleteCmd = &cobra.Command{
	Use:   "delete <file> <version-id>",
	Short: "Delete a specific version of a file",
	Args:  cobra.ExactArgs(2),
	RunE:  runDelete,
}

func init() {
	deleteCmd.Flags().BoolVarP(&deleteForce, "force", "f", false, "Skip confirmation prompt")
}

func runDelete(cmd *cobra.Command, args []string) error {
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

	fmt.Printf("Version %s of %s\n", version.ID, filePath)
	fmt.Printf("  Created: %s\n", version.CreatedAt.Format("2006-01-02 15:04:05"))
	if len(version.Tags) > 0 {
		fmt.Printf("  Tags: %s\n", joinStrings(version.Tags, ", "))
	}
	if version.Notes != "" {
		fmt.Printf("  Notes: %s\n", version.Notes)
	}
	fmt.Printf("  Size: %s\n", formatSize(version.Size))
	fmt.Println()

	var confirm bool
	if !deleteForce {
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewConfirm().
					Title("Delete this version?").
					Value(&confirm),
			),
		)

		if err := form.Run(); err != nil {
			return err
		}
	} else {
		confirm = true
	}

	if !confirm {
		fmt.Println("Cancelled")
		return nil
	}

	if !list.RemoveVersion(absPath, versionID) {
		return fmt.Errorf("failed to remove version from list")
	}

	if err := list.Save(shadowPath); err != nil {
		return fmt.Errorf("failed to save list: %w", err)
	}

	snapshotPath := filepath.Join(shadowPath, "snapshots", versionID)
	if err := os.Remove(snapshotPath); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to delete snapshot: %w", err)
	}

	fmt.Printf("✓ Deleted version %s\n", versionID)
	return nil
}
