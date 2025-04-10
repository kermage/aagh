package commands

import (
	"os"

	"aagh/internal/helpers"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     helpers.NAME,
	Version: helpers.VERSION,
	Long:    "A cross-platform executable for handling Git hooks",
	// SilenceErrors: true,
	// SilenceUsage: true,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if _, found := map[string]struct{}{
			"check": {},
			"init":  {},
			"setup": {},
		}[cmd.Name()]; !found {
			return
		}

		if !helpers.CommandExists("git") {
			cobra.CheckErr("'git' command does not exists.\n\nMake sure it is globally available")
		}
	},
}

func Commander() {
	cobra.EnableCommandSorting = false

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
