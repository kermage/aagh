package helpers

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateDir(t *testing.T) {
	tmpDir := t.TempDir()

	tests := []struct {
		name      string
		path      string
		setup     func() string
		wantErr   bool
		checkFunc func(t *testing.T, path string)
	}{
		{
			name: "create new directory",
			setup: func() string {
				return filepath.Join(tmpDir, "newdir")
			},
			wantErr: false,
			checkFunc: func(t *testing.T, path string) {
				if !DirExists(path) {
					t.Errorf("directory was not created: %s", path)
				}
			},
		},
		{
			name: "directory already exists",
			setup: func() string {
				dir := filepath.Join(tmpDir, "existing")
				_ = os.MkdirAll(dir, PermExecutable)
				return dir
			},
			wantErr: false,
			checkFunc: func(t *testing.T, path string) {
				if !DirExists(path) {
					t.Errorf("directory should still exist: %s", path)
				}
			},
		},
		{
			name: "file exists at path",
			setup: func() string {
				file := filepath.Join(tmpDir, "file.txt")
				_ = os.WriteFile(file, []byte("test"), PermReadWrite)
				return file
			},
			wantErr: true, // MkdirAll fails if a file exists at the path
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup()
			err := CreateDir(path)

			if (err != nil) != tt.wantErr {
				t.Errorf("CreateDir() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if tt.checkFunc != nil && !tt.wantErr {
				tt.checkFunc(t, path)
			}
		})
	}
}

func TestGitDir(t *testing.T) {
	tests := []struct {
		name    string
		setup   func() string
		wantErr bool
	}{
		{
			name: "valid git repository",
			setup: func() string {
				// Use current directory which should be a git repo
				cwd, _ := os.Getwd()
				return cwd
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setup()
			result := GitDir(path)

			if result == "" {
				t.Error("GitDir() returned empty string for valid git repo")
			}

			// Verify the result is a valid path
			if !filepath.IsAbs(result) {
				t.Errorf("GitDir() returned non-absolute path: %s", result)
			}
		})
	}
}

func TestProjectRoot(t *testing.T) {
	t.Run("returns project root", func(t *testing.T) {
		root := ProjectRoot()

		if root == "" {
			t.Error("ProjectRoot() returned empty string")
		}

		// Verify the returned path exists
		if !DirExists(root) {
			t.Errorf("ProjectRoot() returned non-existent path: %s", root)
		}

		// Verify it's an absolute path
		if !filepath.IsAbs(root) {
			t.Errorf("ProjectRoot() returned non-absolute path: %s", root)
		}
	})
}
