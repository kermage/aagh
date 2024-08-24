package commands

import (
	"os"
	"path/filepath"

	"aagh/internal/helpers"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "init",
		Short: "Initialize the repository in the current directory",
		RunE: func(_ *cobra.Command, _ []string) error {
			hooks := helpers.Hooks(helpers.ProjectRoot())

			hooks.Config().Set()
			helpers.CreateDir(hooks.Directory().FullPath())
			helpers.CreateDir(hooks.Runner().FullPath())

			err := os.WriteFile(filepath.Join(hooks.Runner().FullPath(), ".gitignore"), []byte("*"), 0644)

			if err != nil {
				return err
			}

			return nil
		},
	})
}
