package commands

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"aagh/internal/helpers"

	"github.com/kermage/GO-Mods/pathinfo"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "run [hook]",
		Short: "Run a hook in the repository of current directory",
		Args:  cobra.ExactArgs(1),
		RunE: func(_ *cobra.Command, args []string) error {
			hooks := helpers.Hooks(helpers.ProjectRoot())

			if !hooks.ValidRoot() {
				cobra.CheckErr(fmt.Sprintf("'%s' is not initialized.\n\nRun '%s init' first before setting up hooks.\n", hooks.Project().FullPath(), helpers.NAME))
			}

			hookName := args[0]
			hookPath := pathinfo.Get(filepath.Join(hooks.Runner().FullPath(), hookName))

			if !hookPath.Exists() {
				cobra.CheckErr(fmt.Sprintf("'%[1]s' does not exists.\n\nRun '%[2]s setup %[1]s' first before running this command.\n", hookName, helpers.NAME))
			}

			cmd := exec.Command(hookPath.FullPath())
			cmd.Dir = helpers.ProjectRoot()
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err := cmd.Run()

			if err != nil {
				os.Exit(helpers.GetExitCode(err))
			}

			return nil
		},
	})
}
