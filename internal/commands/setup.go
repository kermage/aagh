package commands

import (
	"fmt"
	"os"
	"path/filepath"

	"aagh/internal/helpers"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "setup [[hook]...]",
		Short: "Setup a hook in the repository of current directory",
		Args:  cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			path := helpers.ProjectRoot()

			validatePath(path)

			for _, hook := range args {
				cmd.Printf("Setting up '%s' hook\n", hook)

				err := os.WriteFile(filepath.Join(path, helpers.DIR, hook), []byte(helpers.SCRIPT), helpers.PERM)

				if err != nil {
					return err
				}
			}

			return nil
		},
	})
}

func validatePath(path string) {
	if helpers.HooksPathCorrect(path) && helpers.HooksDirExists(path) {
		return
	}

	cobra.CheckErr(fmt.Sprintf("'%s' is not initialized.\n\nRun '%s init' first before setting up hooks.\n", path, helpers.NAME))
}
