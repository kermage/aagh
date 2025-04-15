package commands

import (
	"os"
	"path/filepath"

	"aagh/internal/helpers"

	"github.com/spf13/cobra"
)

func init() {
	var apply bool

	command := &cobra.Command{
		Use:   "init",
		Short: "Initialize the repository in the current directory",
		RunE: func(cmd *cobra.Command, _ []string) error {
			hooks := helpers.Hooks(helpers.ProjectRoot())

			hooks.Config().Set()
			helpers.CreateDir(hooks.Directory().FullPath())
			helpers.CreateDir(hooks.Runner().FullPath())

			err := os.WriteFile(filepath.Join(hooks.Runner().FullPath(), ".gitignore"), []byte("*"), 0644)

			if err != nil {
				return err
			}

			if apply {
				files, _ := os.ReadDir(hooks.Directory().FullPath())
				list := make([]string, 0)

				for _, file := range files {
					if !file.IsDir() {
						list = append(list, file.Name())
					}
				}

				return setupHooks(cmd, list)
			}

			return nil
		},
	}

	command.PersistentFlags().BoolVar(&apply, "apply", false, "Setup existing hooks after")

	rootCmd.AddCommand(command)
}
