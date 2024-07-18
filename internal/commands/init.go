package commands

import (
	"aagh/internal/helpers"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "init",
		Short: "Initialize the repository in the current directory",
		Run: func(_ *cobra.Command, _ []string) {
			hooks := helpers.Hooks(helpers.ProjectRoot())

			hooks.Config().Set()
			helpers.CreateDir(hooks.Directory().FullPath())
		},
	})
}
