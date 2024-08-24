package commands

import (
	"fmt"

	"aagh/internal/helpers"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "check",
		Short: "Check the repository status in the current directory",
		RunE: func(_ *cobra.Command, _ []string) error {
			hooks := helpers.Hooks(helpers.ProjectRoot())

			fmt.Printf("Project Root: %s\n", hooks.Project().FullPath())
			fmt.Printf("Config Path Set: %v\n", hooks.Config().Correct())
			fmt.Printf("Hooks Directory: %v\n", hooks.IsReady())

			return nil
		},
	})
}
