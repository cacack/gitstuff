package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"gitstuff/internal/config"
	"gitstuff/internal/gitlab"
)

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestDisplayRepositoryList_WithoutVerbose(t *testing.T) {
	// Mock client and config
	cfg := &config.Config{
		Local: config.LocalConfig{
			BaseDir: "/tmp/test",
		},
	}
	
	// Create mock client with test data
	repos := []*gitlab.Repository{
		{
			ID:       1,
			Name:     "test-repo",
			FullPath: "group/test-repo",
			WebURL:   "https://gitlab.com/group/test-repo",
		},
		{
			ID:       2,
			Name:     "another-repo",
			FullPath: "group/another-repo",
			WebURL:   "https://gitlab.com/group/another-repo",
		},
	}
	
	mockClient := &mockGitLabClient{repos: repos}
	
	output := captureOutput(func() {
		displayRepositoryList(mockClient, cfg, false, false) // showStatus=false, showVerbose=false
	})
	
	// Should contain repository names
	if !strings.Contains(output, "group/test-repo") {
		t.Errorf("Expected output to contain repository path 'group/test-repo'")
	}
	
	if !strings.Contains(output, "group/another-repo") {
		t.Errorf("Expected output to contain repository path 'group/another-repo'")
	}
	
	// Should NOT contain URLs when verbose=false
	if strings.Contains(output, "https://gitlab.com/group/test-repo") {
		t.Errorf("Expected output to NOT contain URL when verbose=false")
	}
	
	if strings.Contains(output, "https://gitlab.com/group/another-repo") {
		t.Errorf("Expected output to NOT contain URL when verbose=false")
	}
}

func TestDisplayRepositoryList_WithVerbose(t *testing.T) {
	cfg := &config.Config{
		Local: config.LocalConfig{
			BaseDir: "/tmp/test",
		},
	}
	
	repos := []*gitlab.Repository{
		{
			ID:       1,
			Name:     "test-repo",
			FullPath: "group/test-repo",
			WebURL:   "https://gitlab.com/group/test-repo",
		},
		{
			ID:       2,
			Name:     "another-repo",
			FullPath: "group/another-repo",
			WebURL:   "https://gitlab.com/group/another-repo",
		},
	}
	
	mockClient := &mockGitLabClient{repos: repos}
	
	output := captureOutput(func() {
		displayRepositoryList(mockClient, cfg, false, true) // showStatus=false, showVerbose=true
	})
	
	// Should contain repository names
	if !strings.Contains(output, "group/test-repo") {
		t.Errorf("Expected output to contain repository path 'group/test-repo'")
	}
	
	if !strings.Contains(output, "group/another-repo") {
		t.Errorf("Expected output to contain repository path 'group/another-repo'")
	}
	
	// Should contain URLs when verbose=true
	if !strings.Contains(output, "https://gitlab.com/group/test-repo") {
		t.Errorf("Expected output to contain URL when verbose=true")
	}
	
	if !strings.Contains(output, "https://gitlab.com/group/another-repo") {
		t.Errorf("Expected output to contain URL when verbose=true")
	}
}

func TestDisplayRepositoryTree_WithoutVerbose(t *testing.T) {
	cfg := &config.Config{
		Local: config.LocalConfig{
			BaseDir: "/tmp/test",
		},
	}
	
	// Create a mock tree structure
	tree := &gitlab.RepositoryTree{
		Groups: map[string]*gitlab.GroupNode{
			"group1": {
				Group: &gitlab.Group{
					Name:     "group1",
					FullPath: "group1",
				},
				SubGroups: make(map[string]*gitlab.GroupNode),
				Repositories: []*gitlab.Repository{
					{
						ID:       1,
						Name:     "repo1",
						FullPath: "group1/repo1",
						WebURL:   "https://gitlab.com/group1/repo1",
					},
				},
			},
		},
		Repositories: []*gitlab.Repository{},
	}
	
	mockClient := &mockGitLabClient{tree: tree}
	
	output := captureOutput(func() {
		displayRepositoryTree(mockClient, cfg, false, false) // showStatus=false, showVerbose=false
	})
	
	// Should contain group and repo names
	if !strings.Contains(output, "group1") {
		t.Errorf("Expected output to contain group name 'group1'")
	}
	
	if !strings.Contains(output, "repo1") {
		t.Errorf("Expected output to contain repository name 'repo1'")
	}
	
	// Should NOT contain URLs when verbose=false
	if strings.Contains(output, "https://gitlab.com/group1/repo1") {
		t.Errorf("Expected output to NOT contain URL when verbose=false")
	}
}

func TestDisplayRepositoryTree_WithVerbose(t *testing.T) {
	cfg := &config.Config{
		Local: config.LocalConfig{
			BaseDir: "/tmp/test",
		},
	}
	
	tree := &gitlab.RepositoryTree{
		Groups: map[string]*gitlab.GroupNode{
			"group1": {
				Group: &gitlab.Group{
					Name:     "group1",
					FullPath: "group1",
				},
				SubGroups: make(map[string]*gitlab.GroupNode),
				Repositories: []*gitlab.Repository{
					{
						ID:       1,
						Name:     "repo1",
						FullPath: "group1/repo1",
						WebURL:   "https://gitlab.com/group1/repo1",
					},
				},
			},
		},
		Repositories: []*gitlab.Repository{},
	}
	
	mockClient := &mockGitLabClient{tree: tree}
	
	output := captureOutput(func() {
		displayRepositoryTree(mockClient, cfg, false, true) // showStatus=false, showVerbose=true
	})
	
	// Should contain group and repo names
	if !strings.Contains(output, "group1") {
		t.Errorf("Expected output to contain group name 'group1'")
	}
	
	if !strings.Contains(output, "repo1") {
		t.Errorf("Expected output to contain repository name 'repo1'")
	}
	
	// Should contain URLs when verbose=true
	if !strings.Contains(output, "https://gitlab.com/group1/repo1") {
		t.Errorf("Expected output to contain URL when verbose=true")
	}
}

// Mock GitLab client for testing
type mockGitLabClient struct {
	repos []*gitlab.Repository
	tree  *gitlab.RepositoryTree
}

func (m *mockGitLabClient) ListAllRepositories() ([]*gitlab.Repository, error) {
	return m.repos, nil
}

func (m *mockGitLabClient) BuildRepositoryTree() (*gitlab.RepositoryTree, error) {
	if m.tree != nil {
		return m.tree, nil
	}
	return &gitlab.RepositoryTree{
		Groups:       make(map[string]*gitlab.GroupNode),
		Repositories: m.repos,
	}, nil
}

func TestListCommandFlags(t *testing.T) {
	// Test that the verbose flag is properly registered
	cmd := listCmd
	
	// Check if verbose flag exists
	verboseFlag := cmd.Flags().Lookup("verbose")
	if verboseFlag == nil {
		t.Error("Expected --verbose flag to be registered")
	}
	
	// Check if the short flag maps to the same flag (Cobra handles this automatically)
	if verboseFlag.Shorthand != "v" {
		t.Errorf("Expected shorthand for verbose flag to be 'v', got '%s'", verboseFlag.Shorthand)
	}
	
	// Test default value
	defaultValue := verboseFlag.DefValue
	if defaultValue != "false" {
		t.Errorf("Expected default value for verbose flag to be 'false', got '%s'", defaultValue)
	}
	
	// Test flag usage
	usage := verboseFlag.Usage
	expectedUsage := "Show additional details like URLs"
	if usage != expectedUsage {
		t.Errorf("Expected flag usage '%s', got '%s'", expectedUsage, usage)
	}
}

func TestListCommandFlagCombinations(t *testing.T) {
	// Test that flags can be combined
	cmd := listCmd
	
	// Reset flags for clean test
	cmd.Flags().Set("verbose", "false")
	cmd.Flags().Set("tree", "false")
	cmd.Flags().Set("status", "true")
	
	// Test setting verbose flag
	err := cmd.Flags().Set("verbose", "true")
	if err != nil {
		t.Errorf("Failed to set verbose flag: %v", err)
	}
	
	verboseValue, err := cmd.Flags().GetBool("verbose")
	if err != nil {
		t.Errorf("Failed to get verbose flag value: %v", err)
	}
	
	if !verboseValue {
		t.Error("Expected verbose flag to be true after setting")
	}
	
	// Test combining with tree flag
	err = cmd.Flags().Set("tree", "true")
	if err != nil {
		t.Errorf("Failed to set tree flag: %v", err)
	}
	
	treeValue, err := cmd.Flags().GetBool("tree")
	if err != nil {
		t.Errorf("Failed to get tree flag value: %v", err)
	}
	
	if !treeValue {
		t.Error("Expected tree flag to be true after setting")
	}
	
	// Both flags should still be set correctly
	verboseValue, _ = cmd.Flags().GetBool("verbose")
	if !verboseValue {
		t.Error("Expected verbose flag to remain true when combining with other flags")
	}
}