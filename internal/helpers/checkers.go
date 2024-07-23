package helpers

import (
	"os"
	"os/exec"
)

func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)

	return err == nil
}

func DirExists(path string) bool {
	_, err := os.Stat(path)

	return !os.IsNotExist(err)
}
