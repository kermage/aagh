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
			hooks := helpers.Hooks(helpers.ProjectRoot())

			if !hooks.ValidRoot() {
				cobra.CheckErr(fmt.Sprintf("'%s' is not initialized.\n\nRun '%s init' first before setting up hooks.\n", hooks.Project().FullPath(), helpers.NAME))
			}

			for _, hook := range args {
				cmd.Printf("Setting up '%s' hook\n", hook)

				err := os.WriteFile(filepath.Join(hooks.Directory().FullPath(), hook), []byte(helpers.SCRIPT), helpers.PERM)

				if err != nil {
					return err
				}
			}

			return nil
		},
	})
}
