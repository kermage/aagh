package cmd

import (
	"fmt"

	"aagh/helpers"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "check",
		Short: "Check the repository status in the current directory",
		RunE: func(_ *cobra.Command, _ []string) error {
			path := helpers.ProjectRoot()

			fmt.Printf("Project Root: %s\n", path)
			fmt.Printf("Config Path Set: %v\n", helpers.HooksPathCorrect(path))
			fmt.Printf("Hooks Directory: %v\n", helpers.HooksDirExists(path))

			return nil
		},
	})
}
