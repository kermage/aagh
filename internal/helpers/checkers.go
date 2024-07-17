package helpers

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func HooksPathCorrect(path string) bool {
	cmd := exec.Command("git", "config", "--get", "core.hooksPath", DIR)
	cmd.Dir = path
	val, err := cmd.Output()

	if err != nil {
		return false
	}

	return strings.TrimSpace(string(val)) == DIR
}

func HooksDirExists(path string) bool {
	_, err := os.Stat(filepath.Join(path, DIR))

	return !os.IsNotExist(err)
}
