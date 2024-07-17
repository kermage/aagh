package commands

import (
	"os"
	"os/exec"
	"path/filepath"

	"aagh/internal/helpers"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "init",
		Short: "Initialize the repository in the current directory",
		Run: func(_ *cobra.Command, _ []string) {
			path := helpers.ProjectRoot()

			configHooksPath(path)
			createHooksDir(path)
		},
	})
}

func configHooksPath(path string) {
	cmd := exec.Command("git", "config", "core.hooksPath", helpers.DIR)
	cmd.Dir = path
	_, err := cmd.Output()

	if err != nil {
		cobra.CheckErr("running 'git config' failed.\n\nMake sure 'git' command is globally available")
	}
}

func createHooksDir(path string) error {
	hooksPath := filepath.Join(path, helpers.DIR)
	_, err := os.Stat(hooksPath)

	if os.IsNotExist(err) {
		return os.Mkdir(hooksPath, helpers.PERM)
	}

	return err
}
