package helpers

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

func GitDir(path string) string {
	cmd := exec.Command("git", "rev-parse", "--absolute-git-dir")
	cmd.Dir = path
	res, err := cmd.Output()

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
