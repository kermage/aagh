package helpers

import (
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func setupGitRepo(t *testing.T, dir string) {
	t.Helper()

	cmd := exec.Command("git", "init")
	cmd.Dir = dir
	if err := cmd.Run(); err != nil {
		t.Fatalf("failed to init git repo: %v", err)
	}

	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = dir
	_ = cmd.Run()

	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = dir
	_ = cmd.Run()
}

func TestHooksConfig(t *testing.T) {
	tmpDir := t.TempDir()
	setupGitRepo(t, tmpDir)

	h := Hooks(tmpDir)

	t.Run("returns config instance", func(t *testing.T) {
		cfg := h.Config()
		if cfg == nil {
			t.Fatal("Config() returned nil")
		}
		if cfg.hooks != h {
			t.Error("Config().hooks does not point to parent hooks instance")
		}
	})

	t.Run("runner path is set correctly", func(t *testing.T) {
		cfg := h.Config()
		expectedRunner := filepath.Join(DIR, RUNNER)
		if cfg.runner != expectedRunner {
			t.Errorf("Config().runner = %s; want %s", cfg.runner, expectedRunner)
		}
	})
}

func TestConfigGet(t *testing.T) {
	tmpDir := t.TempDir()
	setupGitRepo(t, tmpDir)

	h := Hooks(tmpDir)
	cfg := h.Config()

	t.Run("get unset config", func(t *testing.T) {
		value := cfg.Get()
		if value != "" {
			t.Errorf("Get() = %q; want empty string for unset config", value)
		}
	})

	t.Run("get set config", func(t *testing.T) {
		// Set a value directly using git
		cmd := exec.Command("git", "config", KEY, "test-value")
		cmd.Dir = tmpDir
		if err := cmd.Run(); err != nil {
			t.Fatalf("failed to set git config: %v", err)
		}

		value := cfg.Get()
		if value != "test-value" {
			t.Errorf("Get() = %q; want %q", value, "test-value")
		}
	})
}

func TestConfigSet(t *testing.T) {
	tmpDir := t.TempDir()
	setupGitRepo(t, tmpDir)

	h := Hooks(tmpDir)
	cfg := h.Config()

	t.Run("set config value", func(t *testing.T) {
		err := cfg.Set()
		if err != nil {
			t.Fatalf("Set() error = %v; want nil", err)
		}

		// Verify it was set correctly
		value := cfg.Get()
		expectedValue := filepath.Join(DIR, RUNNER)
		if value != expectedValue {
			t.Errorf("After Set(), Get() = %q; want %q", value, expectedValue)
		}
	})

	t.Run("set config multiple times", func(t *testing.T) {
		// Should be idempotent
		err := cfg.Set()
		if err != nil {
			t.Fatalf("First Set() error = %v", err)
		}

		err = cfg.Set()
		if err != nil {
			t.Fatalf("Second Set() error = %v", err)
		}

		value := cfg.Get()
		expectedValue := filepath.Join(DIR, RUNNER)
		if value != expectedValue {
			t.Errorf("After multiple Set(), Get() = %q; want %q", value, expectedValue)
		}
	})
}

func TestConfigCorrect(t *testing.T) {
	tmpDir := t.TempDir()
	setupGitRepo(t, tmpDir)

	h := Hooks(tmpDir)
	cfg := h.Config()

	tests := []struct {
		name     string
		setup    func()
		expected bool
	}{
		{
			name: "correct when unset",
			setup: func() {
				// Do nothing - config is unset
			},
			expected: false,
		},
		{
			name: "correct when set to correct value",
			setup: func() {
				_ = cfg.Set()
			},
			expected: true,
		},
		{
			name: "incorrect when set to wrong value",
			setup: func() {
				cmd := exec.Command("git", "config", KEY, "wrong-value")
				cmd.Dir = tmpDir
				_ = cmd.Run()
			},
			expected: false,
		},
		{
			name: "correct after fixing wrong value",
			setup: func() {
				// First set wrong value
				cmd := exec.Command("git", "config", KEY, "wrong-value")
				cmd.Dir = tmpDir
				_ = cmd.Run()

				// Then fix it
				_ = cfg.Set()
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset git config before each test
			cmd := exec.Command("git", "config", "--unset", KEY)
			cmd.Dir = tmpDir
			_ = cmd.Run()

			tt.setup()

			got := cfg.Correct()
			if got != tt.expected {
				currentValue := cfg.Get()
				t.Errorf("Correct() = %v; want %v (current value: %q)", got, tt.expected, currentValue)
			}
		})
	}
}

func TestConfigIntegration(t *testing.T) {
	tmpDir := t.TempDir()
	setupGitRepo(t, tmpDir)

	h := Hooks(tmpDir)

	t.Run("full workflow", func(t *testing.T) {
		cfg := h.Config()

		// Initially should not be correct
		if cfg.Correct() {
			t.Error("Config should not be correct initially")
		}

		// Set the config
		if err := cfg.Set(); err != nil {
			t.Fatalf("Set() failed: %v", err)
		}

		// Now should be correct
		if !cfg.Correct() {
			t.Error("Config should be correct after Set()")
		}

		// Get should return the correct value
		value := cfg.Get()
		expectedValue := filepath.Join(DIR, RUNNER)
		if value != expectedValue {
			t.Errorf("Get() = %q; want %q", value, expectedValue)
		}

		// Verify using git command directly
		cmd := exec.Command("git", "config", "--get", KEY)
		cmd.Dir = tmpDir
		output, err := cmd.Output()
		if err != nil {
			t.Fatalf("git config --get failed: %v", err)
		}

		gitValue := strings.TrimSpace(string(output))
		if gitValue != expectedValue {
			t.Errorf("git config value = %q; want %q", gitValue, expectedValue)
		}
	})
}
