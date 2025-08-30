package cmd

import (
	"strings"
	"testing"

	"gitstuff/internal/scm"
)

func TestFindRepositoryByPath_ExactMatch(t *testing.T) {
	repos := []*scm.Repository{
		{
			ID:       "1",
			Name:     "exact-repo",
			FullPath: "group/exact-repo",
			Provider: "gitlab",
		},
		{
			ID:       "2",
			Name:     "other-repo",
			FullPath: "group/other-repo",
			Provider: "gitlab",
		},
	}

	mockClient := &mockSCMClient{
		providerType: "gitlab",
		repos:        repos,
	}

	repo, err := findRepositoryByPath(mockClient, "group/exact-repo")
	if err != nil {
		t.Errorf("findRepositoryByPath failed: %v", err)
	}
	if repo == nil {
		t.Fatal("Expected repository to be found")
	}
	if repo.FullPath != "group/exact-repo" {
		t.Errorf("Expected repo with path 'group/exact-repo', got: %s", repo.FullPath)
	}
}

func TestFindRepositoryByPath_PartialMatch(t *testing.T) {
	repos := []*scm.Repository{
		{
			ID:       "1",
			Name:     "partial-repo",
			FullPath: "group/subgroup/partial-repo",
			Provider: "gitlab",
		},
	}

	mockClient := &mockSCMClient{
		providerType: "gitlab",
		repos:        repos,
	}

	repo, err := findRepositoryByPath(mockClient, "partial-repo")
	if err != nil {
		t.Errorf("findRepositoryByPath failed: %v", err)
	}
	if repo == nil {
		t.Fatal("Expected repository to be found")
	}
	if repo.FullPath != "group/subgroup/partial-repo" {
		t.Errorf("Expected repo with path 'group/subgroup/partial-repo', got: %s", repo.FullPath)
	}
}

func TestFindRepositoryByPath_NotFound(t *testing.T) {
	repos := []*scm.Repository{
		{
			ID:       "1",
			Name:     "existing-repo",
			FullPath: "group/existing-repo",
			Provider: "gitlab",
		},
	}

	mockClient := &mockSCMClient{
		providerType: "gitlab",
		repos:        repos,
	}

	repo, err := findRepositoryByPath(mockClient, "nonexistent-repo")
	if err == nil {
		t.Error("Expected error for nonexistent repository")
	}
	if repo != nil {
		t.Error("Expected no repository to be found")
	}
	if !strings.Contains(err.Error(), "repository not found") {
		t.Errorf("Expected 'repository not found' error, got: %v", err)
	}
}

func TestGroupRepositoryFiltering(t *testing.T) {
	groupRepos := []*scm.Repository{
		{
			ID:       "1",
			Name:     "group-repo1",
			FullPath: "testgroup/group-repo1",
			Provider: "gitlab",
		},
		{
			ID:       "2",
			Name:     "group-repo2",
			FullPath: "testgroup/group-repo2",
			Provider: "gitlab",
		},
	}

	mockClient := &mockSCMClient{
		providerType: "gitlab",
		repos:        []*scm.Repository{}, // Empty overall repos
		groupRepos: map[string][]*scm.Repository{
			"testgroup": groupRepos,
		},
	}

	// Test that the client can return group-specific repositories
	repos, err := mockClient.ListRepositoriesInGroup("testgroup")
	if err != nil {
		t.Errorf("ListRepositoriesInGroup failed: %v", err)
	}

	if len(repos) != 2 {
		t.Errorf("Expected 2 repositories in group, got %d", len(repos))
	}

	if repos[0].FullPath != "testgroup/group-repo1" {
		t.Errorf("Expected first repo to be 'testgroup/group-repo1', got: %s", repos[0].FullPath)
	}

	if repos[1].FullPath != "testgroup/group-repo2" {
		t.Errorf("Expected second repo to be 'testgroup/group-repo2', got: %s", repos[1].FullPath)
	}
}

func TestEmptyGroupFiltering(t *testing.T) {
	mockClient := &mockSCMClient{
		providerType: "gitlab",
		repos:        []*scm.Repository{},            // Empty overall repos
		groupRepos:   map[string][]*scm.Repository{}, // No groups
	}

	// Test that empty group returns empty list
	repos, err := mockClient.ListRepositoriesInGroup("nonexistent")
	if err != nil {
		t.Errorf("ListRepositoriesInGroup failed: %v", err)
	}

	if len(repos) != 0 {
		t.Errorf("Expected 0 repositories for nonexistent group, got %d", len(repos))
	}
}

func TestSubgroupFiltering(t *testing.T) {
	subgroupRepos := []*scm.Repository{
		{
			ID:       "1",
			Name:     "subgroup-repo",
			FullPath: "group/subgroup/subgroup-repo",
			Provider: "gitlab",
		},
	}

	mockClient := &mockSCMClient{
		providerType: "gitlab",
		groupRepos: map[string][]*scm.Repository{
			"group/subgroup": subgroupRepos,
		},
	}

	// Test that subgroup filtering works
	repos, err := mockClient.ListRepositoriesInGroup("group/subgroup")
	if err != nil {
		t.Errorf("ListRepositoriesInGroup failed: %v", err)
	}

	if len(repos) != 1 {
		t.Errorf("Expected 1 repository in subgroup, got %d", len(repos))
	}

	if repos[0].FullPath != "group/subgroup/subgroup-repo" {
		t.Errorf("Expected repo to be 'group/subgroup/subgroup-repo', got: %s", repos[0].FullPath)
	}
}

func TestMultipleProviderSupport(t *testing.T) {
	gitlabRepos := []*scm.Repository{
		{
			ID:       "1",
			Name:     "gitlab-repo",
			FullPath: "gitlab-group/gitlab-repo",
			Provider: "gitlab",
		},
	}

	githubRepos := []*scm.Repository{
		{
			ID:       "2",
			Name:     "github-repo",
			FullPath: "github-org/github-repo",
			Provider: "github",
		},
	}

	gitlabClient := &mockSCMClient{
		providerType: "gitlab",
		repos:        gitlabRepos,
		groupRepos: map[string][]*scm.Repository{
			"gitlab-group": gitlabRepos,
		},
	}

	githubClient := &mockSCMClient{
		providerType: "github",
		repos:        githubRepos,
		groupRepos: map[string][]*scm.Repository{
			"github-org": githubRepos,
		},
	}

	clients := []scm.Client{gitlabClient, githubClient}

	// Test finding repository across providers
	var foundRepo *scm.Repository
	for _, client := range clients {
		repo, err := findRepositoryByPath(client, "gitlab-repo")
		if err == nil && repo != nil {
			foundRepo = repo
			break
		}
	}

	if foundRepo == nil {
		t.Fatal("Expected to find GitLab repository")
	}
	if foundRepo.Provider != "gitlab" {
		t.Errorf("Expected GitLab provider, got: %s", foundRepo.Provider)
	}

	// Test group filtering across providers
	var allGroupRepos []*scm.Repository
	for _, client := range clients {
		repos, err := client.ListRepositoriesInGroup("gitlab-group")
		if err != nil {
			t.Errorf("ListRepositoriesInGroup failed for %s: %v", client.GetProviderType(), err)
			continue
		}
		allGroupRepos = append(allGroupRepos, repos...)
	}

	// Should only find the GitLab repo in gitlab-group
	if len(allGroupRepos) != 1 {
		t.Errorf("Expected 1 repository in gitlab-group across all providers, got %d", len(allGroupRepos))
	}
	if allGroupRepos[0].Provider != "gitlab" {
		t.Errorf("Expected GitLab provider in gitlab-group, got: %s", allGroupRepos[0].Provider)
	}
}
