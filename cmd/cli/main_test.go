package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"aagh/internal/helpers"
)

func TestCommands(t *testing.T) {
	tempDir := t.TempDir()
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
				if workDir.path == currentDir {
					defer func() {
						hooks := helpers.ProjectHooks()

						os.Remove(filepath.Join(hooks.Directory().FullPath(), "test"))
						os.Remove(filepath.Join(hooks.Runner().FullPath(), "test"))
					}()
				}

				args = append(args, "test")
			}

			cmd := exec.Command("go", args...)
			cmd.Dir = workDir.path
			_, err := cmd.Output()

			if workDir.isGit && err != nil {
				t.Errorf("Command '%s' should not have failed in '%s': %v", command, workDir.path, err)
			}

			if !workDir.isGit && err == nil {
				t.Errorf("Command '%s' should have failed in '%s'", command, workDir.path)
			}
		}
	}
}
