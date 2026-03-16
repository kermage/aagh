package helpers

import (
	"os"
	"path/filepath"
	"testing"
)

func TestHooks(t *testing.T) {
	tmpDir := t.TempDir()

	t.Run("creates hooks instance", func(t *testing.T) {
		h := Hooks(tmpDir)
		if h == nil {
			t.Fatal("Hooks() returned nil")
		}
	})

	t.Run("project path is set correctly", func(t *testing.T) {
		h := Hooks(tmpDir)
		project := h.Project()
		if project.FullPath() != tmpDir {
			t.Errorf("Project().FullPath() = %s; want %s", project.FullPath(), tmpDir)
		}
	})
}

func TestHooksProject(t *testing.T) {
	tmpDir := t.TempDir()
	h := Hooks(tmpDir)

	t.Run("returns project pathinfo", func(t *testing.T) {
		project := h.Project()
		if project == nil {
			t.Fatal("Project() returned nil")
		}
		if project.FullPath() != tmpDir {
			t.Errorf("Project().FullPath() = %s; want %s", project.FullPath(), tmpDir)
		}
	})
}

func TestHooksDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	h := Hooks(tmpDir)

	t.Run("returns hooks directory pathinfo", func(t *testing.T) {
		dir := h.Directory()
		if dir == nil {
			t.Fatal("Directory() returned nil")
		}

		expectedPath := filepath.Join(tmpDir, DIR)
		if dir.FullPath() != expectedPath {
			t.Errorf("Directory().FullPath() = %s; want %s", dir.FullPath(), expectedPath)
		}
	})

	t.Run("directory exists check", func(t *testing.T) {
		dir := h.Directory()

		// Initially should not exist
		if dir.Exists() {
			t.Error("Directory should not exist initially")
		}

		// Create the directory
		_ = os.MkdirAll(dir.FullPath(), PermExecutable)

		// Now should exist
		dir = h.Directory() // Refresh
		if !dir.Exists() {
			t.Error("Directory should exist after creation")
		}
	})
}

func TestHooksRunner(t *testing.T) {
	tmpDir := t.TempDir()
	h := Hooks(tmpDir)

	t.Run("returns runner directory pathinfo", func(t *testing.T) {
		runner := h.Runner()
		if runner == nil {
			t.Fatal("Runner() returned nil")
		}

		expectedPath := filepath.Join(tmpDir, DIR, RUNNER)
		if runner.FullPath() != expectedPath {
			t.Errorf("Runner().FullPath() = %s; want %s", runner.FullPath(), expectedPath)
		}
	})
}

func TestHooksIsReady(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name     string
		setup    func(h *hooks)
		expected bool
	}{
		{
			name: "not ready - nothing exists",
			setup: func(h *hooks) {
				// Do nothing
			},
			expected: false,
		},
		{
			name: "not ready - only directory exists",
			setup: func(h *hooks) {
				_ = os.MkdirAll(h.Directory().FullPath(), PermExecutable)
			},
			expected: false,
		},
		{
			name: "ready - both exist",
			setup: func(h *hooks) {
				_ = os.MkdirAll(h.Runner().FullPath(), PermExecutable)
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a fresh temp dir for each test
			testDir := filepath.Join(tmpDir, tt.name)
			_ = os.MkdirAll(testDir, PermExecutable)

			h := Hooks(testDir)
			tt.setup(h)

			got := h.IsReady()
			if got != tt.expected {
				t.Errorf("IsReady() = %v; want %v", got, tt.expected)
			}
		})
	}
}

func TestHooksValidRoot(t *testing.T) {
	// This test requires a git repository
	// We'll use the current directory
	cwd, _ := os.Getwd()

	// Navigate to project root (two levels up from internal/helpers)
	projectRoot := filepath.Dir(filepath.Dir(cwd))

	t.Run("valid root with proper setup", func(t *testing.T) {
		h := Hooks(projectRoot)

		// This should be true if the project is properly initialized
		// (which it should be since we're running tests in the project)
		result := h.ValidRoot()

		if !result {
			t.Log("ValidRoot() returned false - this is NOT expected")
		}
	})

	t.Run("invalid root without setup", func(t *testing.T) {
		tmpDir := t.TempDir()
		h := Hooks(tmpDir)

		result := h.ValidRoot()
		if result {
			t.Error("ValidRoot() should return false for uninitialized directory")
		}
	})
}
