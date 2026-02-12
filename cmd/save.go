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
	saveTags  []string
	saveNotes string
)

var saveCmd = &cobra.Command{
	Use:   "save <file>",
	Short: "Save a version of a file",
	Args:  cobra.ExactArgs(1),
	RunE:  runSave,
}

func init() {
	saveCmd.Flags().StringSliceVarP(&saveTags, "tag", "t", []string{}, "Tags for this version")
	saveCmd.Flags().StringVarP(&saveNotes, "note", "n", "", "Notes for this version")
}

func runSave(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	if _, err := os.Stat(filePath); err != nil {
		return fmt.Errorf("file not found: %s", filePath)
	}

	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	shadowPath, err := repo.ResolveShadowPath(filePath, cfg)
	if err != nil {
		return fmt.Errorf("failed to resolve shadow path: %w", err)
	}

	if err := repo.EnsureShadowDir(shadowPath); err != nil {
		return fmt.Errorf("failed to create shadow directory: %w", err)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}

	versionID := shadow.GenerateVersionID(content)

	snapshotPath := filepath.Join(shadowPath, "snapshots", versionID)
	if err := shadow.CopyFile(filePath, snapshotPath); err != nil {
		return fmt.Errorf("failed to copy file: %w", err)
	}

	fileHash, err := shadow.HashFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to hash file: %w", err)
	}

	stat, _ := os.Stat(filePath)

	if len(saveTags) == 0 && saveNotes == "" {
		var tagsInput string
		form := huh.NewForm(
			huh.NewGroup(
				huh.NewInput().
					Title("Tags (comma-separated)").
					Value(&tagsInput),
				huh.NewText().
					Title("Notes").
					Value(&saveNotes),
			),
		)

		if err := form.Run(); err == nil {
			if tagsInput != "" {
				for _, tag := range splitTags(tagsInput) {
					saveTags = append(saveTags, tag)
				}
			}
		}
	}

	version := shadow.Version{
		ID:        versionID,
		CreatedAt: time.Now(),
		Tags:      saveTags,
		Notes:     saveNotes,
		Size:      stat.Size(),
		Hash:      fileHash,
	}

	list, err := shadow.LoadList(shadowPath)
	if err != nil {
		return fmt.Errorf("failed to load list: %w", err)
	}

	absPath, _ := filepath.Abs(filePath)
	list.AddVersion(absPath, version)

	if err := list.Save(shadowPath); err != nil {
		return fmt.Errorf("failed to save list: %w", err)
	}

	fmt.Printf("✓ Saved version %s of %s\n", versionID, filePath)
	return nil
}

func splitTags(input string) []string {
	var tags []string
	for _, tag := range splitByComma(input) {
		if tag != "" {
			tags = append(tags, tag)
		}
	}
	return tags
}

func splitByComma(s string) []string {
	var result []string
	var current string
	for _, c := range s {
		if c == ',' {
			result = append(result, trimSpace(current))
			current = ""
		} else {
			current += string(c)
		}
	}
	if current != "" {
		result = append(result, trimSpace(current))
	}
	return result
}

func trimSpace(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t') {
		end--
	}
	return s[start:end]
}
