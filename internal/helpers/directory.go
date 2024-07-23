package helpers

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func CreateDir(path string) error {
	_, err := os.Stat(path)

	if os.IsNotExist(err) {
		return os.Mkdir(path, 0755)
	}

	return err
}

func GitDir(path string) string {
	res, err := GitExec(path, "rev-parse", "--absolute-git-dir")

	if err != nil {
		cobra.CheckErr("not a git repository (or any of the parent directories)\n\nRun 'git init' first before running this command")
	}

	return strings.TrimSpace(string(res))
}

func ProjectRoot() string {
	path, err := os.Getwd()

	if err != nil {
		cobra.CheckErr("cannot determine current directory (permission denied)\n\nRun 'chmod +x ./' first before running this command")
	}

	return filepath.Dir(GitDir(path))
}
