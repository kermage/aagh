package commands

import (
	"fmt"
	"os"
	"path/filepath"

	runner "aagh/cmd"
	"aagh/internal/helpers"

	"github.com/kermage/GO-Mods/pathinfo"
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

			return setupHooks(cmd, args)
		},
	})
}

func setupHooks(cmd *cobra.Command, args []string) error {
	hooks := helpers.Hooks(helpers.ProjectRoot())

	for _, hook := range args {
		cmd.Printf("Setting up '%s' hook\n", hook)

		cmd.Printf("Writing script to '%s'\n", filepath.Join(hooks.Directory().FullPath(), hook))
		err := os.WriteFile(filepath.Join(hooks.Runner().FullPath(), hook), runner.Executable, 0755)

		if err != nil {
			return err
		}

		hookPath := pathinfo.Get(filepath.Join(hooks.Directory().FullPath(), hook))

		if hookPath.Exists() {
			continue
		}

		err = os.WriteFile(hookPath.FullPath(), []byte(helpers.SCRIPT), 0644)

		if err != nil {
			return err
		}
	}

	return nil
}
