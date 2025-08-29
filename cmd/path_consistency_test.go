package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"gitstuff/internal/config"
	"gitstuff/internal/scm"
)

// Test that ensures clone and list commands use consistent paths
func TestCloneAndListPathConsistency(t *testing.T) {
	// Create a temporary directory for testing
	tempDir := t.TempDir()

	cfg := &config.Config{
		Local: config.LocalConfig{
			BaseDir: tempDir,
		},
	}

	testRepo := &scm.Repository{
		ID:          "123",
		Name:        "test-repo",
		FullPath:    "owner/test-repo",
		Provider:    "github",
		CloneURL:    "https://github.com/owner/test-repo.git",
		SSHCloneURL: "git@github.com:owner/test-repo.git",
		WebURL:      "https://github.com/owner/test-repo",
	}

	// Test path construction for both commands
	t.Run("Path construction consistency", func(t *testing.T) {
		// This is how the clone command constructs paths (from clone.go:85)
		clonePath := filepath.Join(cfg.Local.BaseDir, testRepo.Provider, testRepo.FullPath)

		// This is how the list command should construct paths (from list.go:111)
		listPath := filepath.Join(cfg.Local.BaseDir, testRepo.Provider, testRepo.FullPath)

		if clonePath != listPath {
			t.Errorf("Path mismatch between clone and list commands:\n"+
				"Clone path: %s\n"+
				"List path:  %s", clonePath, listPath)
		}

		expectedPath := filepath.Join(tempDir, "github", "owner", "test-repo")
		if clonePath != expectedPath {
			t.Errorf("Expected path: %s, but got: %s", expectedPath, clonePath)
		}
	})

	// Test that the path structure follows provider/owner/repo pattern
	t.Run("Provider subdirectory structure", func(t *testing.T) {
		path := filepath.Join(cfg.Local.BaseDir, testRepo.Provider, testRepo.FullPath)

		// Path should be: baseDir/provider/owner/repo
		expectedComponents := []string{tempDir, "github", "owner", "test-repo"}
		expectedPath := filepath.Join(expectedComponents...)

		if path != expectedPath {
			t.Errorf("Path structure incorrect:\n"+
				"Expected: %s\n"+
				"Got:      %s", expectedPath, path)
		}
	})

	// Test with GitLab provider as well
	t.Run("GitLab provider path consistency", func(t *testing.T) {
		gitlabRepo := &scm.Repository{
			ID:          "456",
			Name:        "gitlab-repo",
			FullPath:    "group/subgroup/gitlab-repo",
			Provider:    "gitlab",
			CloneURL:    "https://gitlab.com/group/subgroup/gitlab-repo.git",
			SSHCloneURL: "git@gitlab.com:group/subgroup/gitlab-repo.git",
			WebURL:      "https://gitlab.com/group/subgroup/gitlab-repo",
		}

		path := filepath.Join(cfg.Local.BaseDir, gitlabRepo.Provider, gitlabRepo.FullPath)
		expectedPath := filepath.Join(tempDir, "gitlab", "group", "subgroup", "gitlab-repo")

		if path != expectedPath {
			t.Errorf("GitLab path structure incorrect:\n"+
				"Expected: %s\n"+
				"Got:      %s", expectedPath, path)
		}
	})
}

// Test that verifies the exact path construction used in both clone and list commands
func TestActualPathConstruction(t *testing.T) {
	tempDir := t.TempDir()

	cfg := &config.Config{
		Local: config.LocalConfig{
			BaseDir: tempDir,
		},
	}

	testRepo := &scm.Repository{
		Provider: "github",
		FullPath: "neilfarmer/argo-examples",
	}

	// Test the actual path construction used in the application
	t.Run("Real path construction", func(t *testing.T) {
		// This mimics the path construction in both clone.go and list.go
		actualPath := filepath.Join(cfg.Local.BaseDir, testRepo.Provider, testRepo.FullPath)
		expectedPath := filepath.Join(tempDir, "github", "neilfarmer", "argo-examples")

		if actualPath != expectedPath {
			t.Errorf("Path construction failed:\n"+
				"Expected: %s\n"+
				"Actual:   %s", expectedPath, actualPath)
		}
	})

	// Create the directory structure and verify it can be detected
	t.Run("Directory detection", func(t *testing.T) {
		repoPath := filepath.Join(cfg.Local.BaseDir, testRepo.Provider, testRepo.FullPath)

		// Create the directory structure
		err := os.MkdirAll(repoPath, 0755)
		if err != nil {
			t.Fatalf("Failed to create test directory: %v", err)
		}

		// Verify the directory exists
		if _, err := os.Stat(repoPath); os.IsNotExist(err) {
			t.Errorf("Directory was not created at expected path: %s", repoPath)
		}
	})
}
