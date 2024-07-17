package main

import (
	"aagh/internal/helpers"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func TestCommands(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "")

	if err != nil {
		t.Fatalf("Could not create temporary directory: %s", err)
	}

	t.Logf("Temporary directory: %s", tempDir)

	defer func() {
		_ = os.Remove(tempDir)
	}()

	_, file, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(file)

	for _, command := range []string{"check", "init", "setup"} {
		for _, workDir := range []struct {
			isGit bool
			path  string
		}{
			{isGit: false, path: tempDir},
			{isGit: true, path: currentDir},
		} {
			args := []string{"run", filepath.Join(currentDir, "main.go"), command}

			if command == "setup" {
				args = append(args, "test")
			}

			cmd := exec.Command("go", args...)
			cmd.Dir = workDir.path
			_, err := cmd.Output()

			if workDir.isGit && err != nil {
				t.Logf("Command '%s' should not have failed in '%s'", command, workDir.path)
			}

			if !workDir.isGit && err == nil {
				t.Logf("Command '%s' should have failed in '%s'", command, workDir.path)
			}

			if command == "setup" {
				os.Remove(filepath.Join(currentDir, "../..", helpers.DIR, "test"))
			}
		}
	}
}
