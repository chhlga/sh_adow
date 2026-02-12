package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "shadow",
	Short: "Simple file versioning tool",
	Long:  `Shadow is a lightweight file versioning tool that creates snapshots of files without git complexity.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(saveCmd)
	rootCmd.AddCommand(listCmd)
	rootCmd.AddCommand(restoreCmd)
	rootCmd.AddCommand(deleteCmd)
}
