package helpers

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestGitExec(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name    string
		path    string
		args    []string
		setup   func() string
		wantErr bool
		check   func(t *testing.T, output []byte)
	}{
		{
			name: "git version command",
			setup: func() string {
				cwd, _ := os.Getwd()
				return cwd
			},
			args:    []string{"--version"},
			wantErr: false,
			check: func(t *testing.T, output []byte) {
				if !strings.Contains(string(output), "git version") {
					t.Errorf("expected 'git version' in output, got: %s", output)
				}
			},
		},
		{
			name: "git status in valid repo",
			setup: func() string {
				cwd, _ := os.Getwd()
				return cwd
			},
			args:    []string{"status", "--short"},
			wantErr: false,
		},
		{
			name: "git command in non-git directory",
			setup: func() string {
				return tmpDir
			},
			args:    []string{"status"},
			wantErr: true,
		},
		{
			name: "invalid git command",
			setup: func() string {
				cwd, _ := os.Getwd()
				return cwd
			},
			args:    []string{"invalidcommand"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup()
			output, err := GitExec(path, tt.args...)

			if (err != nil) != tt.wantErr {
				t.Errorf("GitExec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.check != nil && !tt.wantErr {
				tt.check(t, output)
			}
		})
	}
}

func TestScriptExec(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name    string
		script  string
		args    []string
		setup   func() string
		wantErr bool
		check   func(t *testing.T, output []byte)
	}{
		{
			name:   "simple echo script",
			script: "#!/bin/sh\necho 'Hello World'",
			args:   []string{},
			setup: func() string {
				scriptPath := filepath.Join(tmpDir, "echo.sh")
				_ = os.WriteFile(scriptPath, []byte("#!/bin/sh\necho 'Hello World'"), PermExecutable)
				return scriptPath
			},
			wantErr: false,
			check: func(t *testing.T, output []byte) {
				if !strings.Contains(string(output), "Hello World") {
					t.Errorf("expected 'Hello World' in output, got: %s", output)
				}
			},
		},
		{
			name:   "script with arguments",
			script: "#!/bin/sh\necho $1 $2",
			args:   []string{"arg1", "arg2"},
			setup: func() string {
				scriptPath := filepath.Join(tmpDir, "args.sh")
				_ = os.WriteFile(scriptPath, []byte("#!/bin/sh\necho $1 $2"), PermExecutable)
				return scriptPath
			},
			wantErr: false,
			check: func(t *testing.T, output []byte) {
				output_str := strings.TrimSpace(string(output))
				if output_str != "arg1 arg2" {
					t.Errorf("expected 'arg1 arg2' in output, got: %s", output_str)
				}
			},
		},
		{
			name:   "script with exit code",
			script: "#!/bin/sh\nexit 42",
			args:   []string{},
			setup: func() string {
				scriptPath := filepath.Join(tmpDir, "exit.sh")
				_ = os.WriteFile(scriptPath, []byte("#!/bin/sh\nexit 42"), PermExecutable)
				return scriptPath
			},
			wantErr: true,
		},
		{
			name: "non-existent script",
			setup: func() string {
				return filepath.Join(tmpDir, "nonexistent.sh")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			scriptPath := tt.setup()
			output, err := ScriptExec(scriptPath, tt.args)

			if (err != nil) != tt.wantErr {
				t.Errorf("ScriptExec() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.check != nil && !tt.wantErr {
				tt.check(t, output)
			}
		})
	}
}

func TestGetExitCode(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected int
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: 0,
		},
		{
			name: "exit code 1",
			err: &exec.ExitError{
				ProcessState: &os.ProcessState{},
			},
			expected: 0, // Will be set by the actual error
		},
		{
			name:     "non-exit error",
			err:      exec.ErrNotFound,
			expected: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// For actual exit codes, we need to run a real command
			if tt.name == "exit code 1" {
				cmd := exec.Command("sh", "-c", "exit 42")
				err := cmd.Run()
				got := GetExitCode(err)
				if got != 42 {
					t.Errorf("GetExitCode() = %d; want 42", got)
				}
				return
			}

			got := GetExitCode(tt.err)
			if got != tt.expected {
				t.Errorf("GetExitCode() = %d; want %d", got, tt.expected)
			}
		})
	}
}
