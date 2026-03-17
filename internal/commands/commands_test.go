package commands

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"aagh/internal/helpers"
)

func setupTestGitRepo(t *testing.T, dir string) {
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

func setupTestRunner(t *testing.T, dir string, full bool) {
	t.Helper()

	_ = os.MkdirAll(filepath.Join(dir, helpers.DIR, helpers.RUNNER), helpers.PermExecutable)

	if !full {
		return
	}

	cmd := exec.Command("git", "config", "core.hooksPath", filepath.Join(helpers.DIR, helpers.RUNNER))
	cmd.Dir = dir
	_ = cmd.Run()
}

func TestCheckCommandLogic(t *testing.T) {
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	tests := []struct {
		name          string
		setup         func(dir string)
		expectReady   bool
		expectCorrect bool
	}{
		{
			name: "uninitialized repository",
			setup: func(dir string) {
				setupTestGitRepo(t, dir)
			},
			expectReady:   false,
			expectCorrect: false,
		},
		{
			name: "initialized repository",
			setup: func(dir string) {
				setupTestGitRepo(t, dir)
				setupTestRunner(t, dir, true)
			},
			expectReady:   true,
			expectCorrect: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			tt.setup(tmpDir)

			if err := os.Chdir(tmpDir); err != nil {
				t.Fatalf("failed to change directory: %v", err)
			}

			hooks := helpers.ProjectHooks()

			if hooks.IsReady() != tt.expectReady {
				t.Errorf("IsReady() = %v; want %v", hooks.IsReady(), tt.expectReady)
			}

			if hooks.Config().Correct() != tt.expectCorrect {
				t.Errorf("Config().Correct() = %v; want %v", hooks.Config().Correct(), tt.expectCorrect)
			}
		})
	}
}

func TestInitCommandLogic(t *testing.T) {
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	tmpDir := t.TempDir()
	setupTestGitRepo(t, tmpDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	hooks := helpers.ProjectHooks()

	hooks.Config().Set()
	helpers.CreateDir(hooks.Directory().FullPath())
	helpers.CreateDir(hooks.Runner().FullPath())

	err := os.WriteFile(filepath.Join(hooks.Runner().FullPath(), ".gitignore"), []byte("*"), helpers.PermReadWrite)
	if err != nil {
		t.Fatalf("failed to create .gitignore: %v", err)
	}

	if !hooks.Directory().Exists() {
		t.Error(".aagh directory was not created")
	}

	if !hooks.Runner().Exists() {
		t.Error(".aagh/_ directory was not created")
	}

	if !hooks.Config().Correct() {
		t.Error("git config was not set correctly")
	}

	gitignorePath := filepath.Join(hooks.Runner().FullPath(), ".gitignore")
	if !helpers.FileExists(gitignorePath) {
		t.Error(".gitignore was not created in runner directory")
	}
}

func TestSetupHooksLogic(t *testing.T) {
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	tmpDir := t.TempDir()
	setupTestGitRepo(t, tmpDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	hooks := helpers.ProjectHooks()

	hooks.Config().Set()
	helpers.CreateDir(hooks.Directory().FullPath())
	helpers.CreateDir(hooks.Runner().FullPath())

	// Simulate setup command logic for pre-commit hook
	hookName := "pre-commit"

	// Create runner hook (would normally be the runner binary)
	runnerHookPath := filepath.Join(hooks.Runner().FullPath(), hookName)
	err := os.WriteFile(runnerHookPath, []byte("#!/bin/sh\necho test"), helpers.PermExecutable)
	if err != nil {
		t.Fatalf("failed to create runner hook: %v", err)
	}

	// Create hook script
	hookPath := filepath.Join(hooks.Directory().FullPath(), hookName)
	err = os.WriteFile(hookPath, []byte(helpers.SCRIPT), helpers.PermReadWrite)
	if err != nil {
		t.Fatalf("failed to create hook script: %v", err)
	}

	// Verify setup
	if !helpers.FileExists(runnerHookPath) {
		t.Error("runner hook was not created")
	}

	if !helpers.FileExists(hookPath) {
		t.Error("hook script was not created")
	}
}

func TestRunCommandLogic(t *testing.T) {
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	tmpDir := t.TempDir()
	setupTestGitRepo(t, tmpDir)

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("failed to change directory: %v", err)
	}

	hooks := helpers.ProjectHooks()

	// Initialize and setup
	hooks.Config().Set()
	helpers.CreateDir(hooks.Directory().FullPath())
	helpers.CreateDir(hooks.Runner().FullPath())

	runnerHookPath := filepath.Join(hooks.Runner().FullPath(), "pre-commit")

	// Create a simple executable hook
	err := os.WriteFile(runnerHookPath, []byte("#!/bin/sh\necho 'Hook executed'\nexit 0"), helpers.PermExecutable)
	if err != nil {
		t.Fatalf("failed to create runner hook: %v", err)
	}

	if !helpers.FileExists(runnerHookPath) {
		t.Fatal("runner hook does not exist")
	}

	cmd := exec.Command(runnerHookPath)
	cmd.Dir = hooks.Project().FullPath()
	output, err := cmd.CombinedOutput()

	if err != nil {
		t.Errorf("hook execution failed: %v\nOutput: %s", err, output)
	}
}

func TestValidRootCheck(t *testing.T) {
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	tests := []struct {
		name      string
		setup     func(dir string)
		wantValid bool
	}{
		{
			name: "not initialized",
			setup: func(dir string) {
				setupTestGitRepo(t, dir)
			},
			wantValid: false,
		},
		{
			name: "partially initialized - no config",
			setup: func(dir string) {
				setupTestGitRepo(t, dir)
				setupTestRunner(t, dir, false)
				_ = os.MkdirAll(filepath.Join(dir, helpers.DIR, helpers.RUNNER), helpers.PermExecutable)
			},
			wantValid: false,
		},
		{
			name: "fully initialized",
			setup: func(dir string) {
				setupTestGitRepo(t, dir)
				setupTestRunner(t, dir, true)
			},
			wantValid: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			tt.setup(tmpDir)

			if err := os.Chdir(tmpDir); err != nil {
				t.Fatalf("failed to change directory: %v", err)
			}

			hooks := helpers.ProjectHooks()
			got := hooks.ValidRoot()

			if got != tt.wantValid {
				t.Errorf("ValidRoot() = %v; want %v", got, tt.wantValid)
			}
		})
	}
}
