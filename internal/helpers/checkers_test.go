package helpers

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCommandExists(t *testing.T) {
	tests := []struct {
		name     string
		cmd      string
		expected bool
	}{
		{
			name:     "existing command - sh",
			cmd:      "sh",
			expected: true,
		},
		{
			name:     "existing command - git",
			cmd:      "git",
			expected: true,
		},
		{
			name:     "non-existing command",
			cmd:      "nonexistentcommand12345",
			expected: false,
		},
		{
			name:     "empty command",
			cmd:      "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CommandExists(tt.cmd)
			if got != tt.expected {
				t.Errorf("CommandExists(%q) = %v; want %v", tt.cmd, got, tt.expected)
			}
		})
	}
}

func TestDirExists(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name     string
		path     string
		setup    func() string
		expected bool
	}{
		{
			name: "existing directory",
			setup: func() string {
				dir := filepath.Join(tmpDir, "existing")
				_ = os.MkdirAll(dir, PermExecutable)
				return dir
			},
			expected: true,
		},
		{
			name: "non-existing directory",
			setup: func() string {
				return filepath.Join(tmpDir, "nonexistent")
			},
			expected: false,
		},
		{
			name: "empty path",
			setup: func() string {
				return ""
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup()
			got := DirExists(path)
			if got != tt.expected {
				t.Errorf("DirExists(%q) = %v; want %v", path, got, tt.expected)
			}
		})
	}
}
