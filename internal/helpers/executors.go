package helpers

import (
	"os/exec"
)

func GitExec(path string, args ...string) ([]byte, error) {
	cmd := exec.Command("git", args...)
	cmd.Dir = path

	return cmd.Output()
}

func ScriptExec(path string, args []string) ([]byte, error) {
	args = append([]string{"-e", path}, args...)
	cmd := exec.Command("sh", args...)

	return cmd.CombinedOutput()
}

func GetExitCode(err error) int {
	if exitError, ok := err.(*exec.ExitError); ok {
		return exitError.ExitCode()
	}

	return 0
}
