package scm

import (
	"testing"
)

func TestRepositoryStructure(t *testing.T) {
	repo := &Repository{
		ID:            "123",
		Name:          "test-repo",
		FullPath:      "org/test-repo",
		CloneURL:      "https://example.com/org/test-repo.git",
		SSHCloneURL:   "git@example.com:org/test-repo.git",
		DefaultBranch: "main",
		WebURL:        "https://example.com/org/test-repo",
		Provider:      "github",
	}

	if repo.Provider != "github" {
		t.Errorf("Expected provider to be 'github', got %s", repo.Provider)
	}

	if repo.FullPath != "org/test-repo" {
		t.Errorf("Expected full path to be 'org/test-repo', got %s", repo.FullPath)
	}
}

func TestGroupStructure(t *testing.T) {
	group := &Group{
		ID:       "456",
		Name:     "test-org",
		FullPath: "test-org",
		Provider: "gitlab",
	}

	if group.Provider != "gitlab" {
		t.Errorf("Expected provider to be 'gitlab', got %s", group.Provider)
	}

	if group.Name != "test-org" {
		t.Errorf("Expected name to be 'test-org', got %s", group.Name)
	}
}

func TestRepositoryTree(t *testing.T) {
	tree := &RepositoryTree{
		Groups:       make(map[string]*GroupNode),
		Repositories: []*Repository{},
	}

	if tree.Groups == nil {
		t.Error("Expected Groups map to be initialized")
	}

	if tree.Repositories == nil {
		t.Error("Expected Repositories slice to be initialized")
	}

	if len(tree.Groups) != 0 {
		t.Errorf("Expected empty Groups map, got %d entries", len(tree.Groups))
	}

	if len(tree.Repositories) != 0 {
		t.Errorf("Expected empty Repositories slice, got %d entries", len(tree.Repositories))
	}
}
