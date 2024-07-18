package commands

import (
	"os"

	"aagh/internal/helpers"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: helpers.NAME,
	// SilenceErrors: true,
	// SilenceUsage: true,
}

func init() {
	if !helpers.CommandExists("git") {
		cobra.CheckErr("'git' command does not exists.\n\nMake sure it is globally available")
	}
}

func Commander() {
	cobra.EnableCommandSorting = false

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
