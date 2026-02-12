package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/chhlga/sh_adow/internal/config"
	"github.com/chhlga/sh_adow/internal/repo"
	"github.com/chhlga/sh_adow/internal/shadow"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list [file]",
	Short: "List tracked files or versions of a specific file",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runList,
}

func runList(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	var shadowPath string
	if len(args) == 0 {
		wd, _ := os.Getwd()
		shadowPath, err = repo.ResolveShadowPath(wd, cfg)
	} else {
		shadowPath, err = repo.ResolveShadowPath(args[0], cfg)
	}
	if err != nil {
		return fmt.Errorf("failed to resolve shadow path: %w", err)
	}

	list, err := shadow.LoadList(shadowPath)
	if err != nil {
		return fmt.Errorf("failed to load list: %w", err)
	}

	if len(args) == 0 {
		return listAllFiles(list, shadowPath)
	}

	absPath, _ := filepath.Abs(args[0])
	return listFileVersions(list, absPath)
}

func listAllFiles(list *shadow.List, shadowPath string) error {
	if len(list.Files) == 0 {
		fmt.Println("No files tracked yet")
		return nil
	}

	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6"))
	fmt.Println(headerStyle.Render(fmt.Sprintf("Files tracked in shadow (%s):", shadowPath)))

	for _, file := range list.Files {
		var totalSize int64
		for _, v := range file.Versions {
			totalSize += v.Size
		}
		fmt.Printf("  • %s (%d versions, %s)\n", file.Path, len(file.Versions), formatSize(totalSize))
	}

	return nil
}

func listFileVersions(list *shadow.List, filePath string) error {
	entry := list.FindFile(filePath)
	if entry == nil {
		return fmt.Errorf("file not tracked: %s", filePath)
	}

	headerStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("6"))
	virtualStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	versionStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("3"))

	fmt.Println(headerStyle.Render(entry.Path))

	stat, err := os.Stat(filePath)
	if err == nil {
		fmt.Println(virtualStyle.Render(fmt.Sprintf("  → VIRTUAL HEAD (current: %s)", formatSize(stat.Size()))))
	} else {
		fmt.Println(virtualStyle.Render("  → VIRTUAL HEAD (file not found)"))
	}

	for _, v := range entry.Versions {
		age := time.Since(v.CreatedAt)
		tags := ""
		if len(v.Tags) > 0 {
			tags = fmt.Sprintf(" - \"%s\"", joinStrings(v.Tags, ", "))
		}
		fmt.Println(versionStyle.Render(fmt.Sprintf("  • %s - %s%s (%s)",
			v.ID, formatDuration(age), tags, formatSize(v.Size))))
		if v.Notes != "" {
			fmt.Printf("    %s\n", v.Notes)
		}
	}

	return nil
}

func formatSize(bytes int64) string {
	const unit = 1024
	if bytes < unit {
		return fmt.Sprintf("%d B", bytes)
	}
	div, exp := int64(unit), 0
	for n := bytes / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB", float64(bytes)/float64(div), "KMGTPE"[exp])
}

func formatDuration(d time.Duration) string {
	if d < time.Minute {
		return "just now"
	}
	if d < time.Hour {
		m := int(d.Minutes())
		return fmt.Sprintf("%dm ago", m)
	}
	if d < 24*time.Hour {
		h := int(d.Hours())
		return fmt.Sprintf("%dh ago", h)
	}
	days := int(d.Hours() / 24)
	return fmt.Sprintf("%dd ago", days)
}

func joinStrings(strs []string, sep string) string {
	if len(strs) == 0 {
		return ""
	}
	result := strs[0]
	for i := 1; i < len(strs); i++ {
		result += sep + strs[i]
	}
	return result
}
