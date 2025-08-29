package cmd

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"

	"gitstuff/internal/config"
	"gitstuff/internal/scm"
)

func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = io.Copy(&buf, r)
	return buf.String()
}

// Mock SCM client for testing
type mockSCMClient struct {
	providerType string
	repos        []*scm.Repository
	groupRepos   map[string][]*scm.Repository
	tree         *scm.RepositoryTree
}

func (m *mockSCMClient) ListAllRepositories() ([]*scm.Repository, error) {
	return m.repos, nil
}

func (m *mockSCMClient) ListRepositoriesInGroup(groupPath string) ([]*scm.Repository, error) {
	if repos, exists := m.groupRepos[groupPath]; exists {
		return repos, nil
	}
	return []*scm.Repository{}, nil
}

func (m *mockSCMClient) BuildRepositoryTree() (*scm.RepositoryTree, error) {
	return m.tree, nil
}

func (m *mockSCMClient) GetProviderType() string {
	return m.providerType
}

func TestDisplayRepositoryList_WithoutVerbose(t *testing.T) {
	// Mock config
	cfg := &config.Config{
		Local: config.LocalConfig{
			BaseDir: "/tmp/test",
		},
	}

	// Create mock client with test data
	repos := []*scm.Repository{
		{
			ID:       "1",
			Name:     "test-repo",
			FullPath: "group/test-repo",
			WebURL:   "https://gitlab.com/group/test-repo",
			Provider: "gitlab",
		},
		{
			ID:       "2",
			Name:     "another-repo",
			FullPath: "group/another-repo",
			WebURL:   "https://gitlab.com/group/another-repo",
			Provider: "gitlab",
		},
	}

	mockClient := &mockSCMClient{
		providerType: "gitlab",
		repos:        repos,
	}

	clients := []scm.Client{mockClient}

	output := captureOutput(func() {
		_ = displayRepositoryList(clients, cfg, false, false, "")
	})

	// Check output contains repository names
	if !strings.Contains(output, "test-repo") {
		t.Errorf("Expected output to contain 'test-repo', got: %s", output)
	}
	if !strings.Contains(output, "another-repo") {
		t.Errorf("Expected output to contain 'another-repo', got: %s", output)
	}
	if !strings.Contains(output, "[gitlab]") {
		t.Errorf("Expected output to contain '[gitlab]', got: %s", output)
	}
}

func TestDisplayRepositoryList_WithVerbose(t *testing.T) {
	// Mock config
	cfg := &config.Config{
		Local: config.LocalConfig{
			BaseDir: "/tmp/test",
		},
	}

	// Create mock client with test data
	repos := []*scm.Repository{
		{
			ID:       "1",
			Name:     "test-repo",
			FullPath: "group/test-repo",
			WebURL:   "https://gitlab.com/group/test-repo",
			Provider: "gitlab",
		},
	}

	mockClient := &mockSCMClient{
		providerType: "gitlab",
		repos:        repos,
	}

	clients := []scm.Client{mockClient}

	output := captureOutput(func() {
		_ = displayRepositoryList(clients, cfg, false, true, "")
	})

	// Check output contains URLs when verbose
	if !strings.Contains(output, "https://gitlab.com/group/test-repo") {
		t.Errorf("Expected verbose output to contain URL, got: %s", output)
	}
}

func TestDisplayRepositoryTree_MultipleProviders(t *testing.T) {
	// Mock config
	cfg := &config.Config{
		Local: config.LocalConfig{
			BaseDir: "/tmp/test",
		},
	}

	// Create mock GitLab client
	gitlabRepos := []*scm.Repository{
		{
			ID:       "1",
			Name:     "gitlab-repo",
			FullPath: "gitlab-group/gitlab-repo",
			Provider: "gitlab",
		},
	}

	gitlabTree := &scm.RepositoryTree{
		Groups: map[string]*scm.GroupNode{
			"gitlab-group": {
				Group: &scm.Group{
					Name:     "gitlab-group",
					FullPath: "gitlab-group",
					Provider: "gitlab",
				},
				SubGroups:    make(map[string]*scm.GroupNode),
				Repositories: gitlabRepos,
			},
		},
		Repositories: []*scm.Repository{},
	}

	gitlabClient := &mockSCMClient{
		providerType: "gitlab",
		repos:        gitlabRepos,
		tree:         gitlabTree,
	}

	// Create mock GitHub client
	githubRepos := []*scm.Repository{
		{
			ID:       "2",
			Name:     "github-repo",
			FullPath: "github-org/github-repo",
			Provider: "github",
		},
	}

	githubTree := &scm.RepositoryTree{
		Groups: map[string]*scm.GroupNode{
			"github-org": {
				Group: &scm.Group{
					Name:     "github-org",
					FullPath: "github-org",
					Provider: "github",
				},
				SubGroups:    make(map[string]*scm.GroupNode),
				Repositories: githubRepos,
			},
		},
		Repositories: []*scm.Repository{},
	}

	githubClient := &mockSCMClient{
		providerType: "github",
		repos:        githubRepos,
		tree:         githubTree,
	}

	clients := []scm.Client{gitlabClient, githubClient}

	output := captureOutput(func() {
		_ = displayRepositoryTree(clients, cfg, false, false, "")
	})

	// Check output contains both providers
	if !strings.Contains(output, "GITLAB Provider") {
		t.Errorf("Expected output to contain 'GITLAB Provider', got: %s", output)
	}
	if !strings.Contains(output, "GITHUB Provider") {
		t.Errorf("Expected output to contain 'GITHUB Provider', got: %s", output)
	}
	if !strings.Contains(output, "gitlab-repo") {
		t.Errorf("Expected output to contain 'gitlab-repo', got: %s", output)
	}
	if !strings.Contains(output, "github-repo") {
		t.Errorf("Expected output to contain 'github-repo', got: %s", output)
	}
}

func TestCreateClient_GitLab(t *testing.T) {
	providerConfig := config.ProviderConfig{
		Name:     "test-gitlab",
		Type:     "gitlab",
		URL:      "https://gitlab.com",
		Token:    "test-token",
		Insecure: false,
		Group:    "",
	}

	client, err := createClient(providerConfig)
	if err != nil {
		t.Fatalf("createClient failed: %v", err)
	}

	if client.GetProviderType() != "gitlab" {
		t.Errorf("Expected provider type 'gitlab', got '%s'", client.GetProviderType())
	}
}

func TestCreateClient_GitHub(t *testing.T) {
	providerConfig := config.ProviderConfig{
		Name:     "test-github",
		Type:     "github",
		URL:      "https://github.com",
		Token:    "test-token",
		Insecure: false,
		Group:    "",
	}

	client, err := createClient(providerConfig)
	if err != nil {
		t.Fatalf("createClient failed: %v", err)
	}

	if client.GetProviderType() != "github" {
		t.Errorf("Expected provider type 'github', got '%s'", client.GetProviderType())
	}
}

func TestCreateClient_UnsupportedProvider(t *testing.T) {
	providerConfig := config.ProviderConfig{
		Name:     "test-bitbucket",
		Type:     "bitbucket",
		URL:      "https://bitbucket.org",
		Token:    "test-token",
		Insecure: false,
		Group:    "",
	}

	_, err := createClient(providerConfig)
	if err == nil {
		t.Fatal("Expected error for unsupported provider type")
	}

	expectedErr := "unsupported provider type: bitbucket"
	if !strings.Contains(err.Error(), expectedErr) {
		t.Errorf("Expected error to contain '%s', got: %s", expectedErr, err.Error())
	}
}
