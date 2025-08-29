package github

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"gitstuff/internal/scm"
)

func TestNewClient(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		token    string
		insecure bool
		wantErr  bool
	}{
		{
			name:    "valid github.com client",
			url:     "https://github.com",
			token:   "test-token",
			wantErr: false,
		},
		{
			name:    "valid enterprise client",
			url:     "https://github.enterprise.com",
			token:   "test-token",
			wantErr: false,
		},
		{
			name:    "empty token",
			url:     "https://github.com",
			token:   "",
			wantErr: true,
		},
		{
			name:    "empty url",
			url:     "",
			token:   "test-token",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.url, tt.token, tt.insecure)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("NewClient() returned nil client without error")
			}
		})
	}
}

func TestClient_GetProviderType(t *testing.T) {
	client, err := NewClient("https://github.com", "test-token", false)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	if got := client.GetProviderType(); got != "github" {
		t.Errorf("GetProviderType() = %v, want %v", got, "github")
	}
}

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "github.com without protocol",
			input: "github.com",
			want:  "https://github.com",
		},
		{
			name:  "github.com with https",
			input: "https://github.com",
			want:  "https://github.com",
		},
		{
			name:  "enterprise without protocol",
			input: "github.enterprise.com",
			want:  "https://github.enterprise.com/api/v3",
		},
		{
			name:  "enterprise with https",
			input: "https://github.enterprise.com",
			want:  "https://github.enterprise.com/api/v3",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizeURL(tt.input)
			if err != nil {
				t.Errorf("normalizeURL() error = %v", err)
				return
			}
			if got != tt.want {
				t.Errorf("normalizeURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_ListAllRepositories_MockResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v3/user/repos" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`[
				{
					"id": 123,
					"name": "test-repo",
					"full_name": "testuser/test-repo",
					"clone_url": "https://github.com/testuser/test-repo.git",
					"ssh_url": "git@github.com:testuser/test-repo.git",
					"html_url": "https://github.com/testuser/test-repo",
					"default_branch": "main",
					"owner": {
						"login": "testuser",
						"type": "User"
					},
					"private": false,
					"permissions": {
						"pull": true
					}
				}
			]`))
		}
	}))
	defer server.Close()

	client, err := NewClient(server.URL+"/api/v3", "test-token", false)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	repos, err := client.ListAllRepositories()
	if err != nil {
		t.Fatalf("ListAllRepositories() error = %v", err)
	}

	if len(repos) != 1 {
		t.Errorf("Expected 1 repository, got %d", len(repos))
	}

	repo := repos[0]
	if repo.ID != "123" {
		t.Errorf("Expected ID '123', got '%s'", repo.ID)
	}
	if repo.Name != "test-repo" {
		t.Errorf("Expected name 'test-repo', got '%s'", repo.Name)
	}
	if repo.FullPath != "testuser/test-repo" {
		t.Errorf("Expected full path 'testuser/test-repo', got '%s'", repo.FullPath)
	}
	if repo.Provider != "github" {
		t.Errorf("Expected provider 'github', got '%s'", repo.Provider)
	}
}

func TestClient_ListRepositoriesInGroup_MockResponse(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v3/orgs/testorg/repos" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`[
				{
					"id": 456,
					"name": "org-repo",
					"full_name": "testorg/org-repo",
					"clone_url": "https://github.com/testorg/org-repo.git",
					"ssh_url": "git@github.com:testorg/org-repo.git",
					"html_url": "https://github.com/testorg/org-repo",
					"default_branch": "main",
					"owner": {
						"login": "testorg",
						"type": "Organization"
					}
				}
			]`))
		}
	}))
	defer server.Close()

	client, err := NewClient(server.URL+"/api/v3", "test-token", false)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	repos, err := client.ListRepositoriesInGroup("testorg")
	if err != nil {
		t.Fatalf("ListRepositoriesInGroup() error = %v", err)
	}

	if len(repos) != 1 {
		t.Errorf("Expected 1 repository, got %d", len(repos))
	}

	repo := repos[0]
	if repo.FullPath != "testorg/org-repo" {
		t.Errorf("Expected full path 'testorg/org-repo', got '%s'", repo.FullPath)
	}
}

func TestClient_BuildRepositoryTree(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v3/user/repos" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`[
				{
					"id": 123,
					"name": "personal-repo",
					"full_name": "testuser/personal-repo",
					"clone_url": "https://github.com/testuser/personal-repo.git",
					"ssh_url": "git@github.com:testuser/personal-repo.git",
					"html_url": "https://github.com/testuser/personal-repo",
					"default_branch": "main",
					"owner": {
						"login": "testuser",
						"type": "User"
					},
					"private": false,
					"permissions": {
						"pull": true
					}
				},
				{
					"id": 456,
					"name": "org-repo",
					"full_name": "testorg/org-repo",
					"clone_url": "https://github.com/testorg/org-repo.git",
					"ssh_url": "git@github.com:testorg/org-repo.git",
					"html_url": "https://github.com/testorg/org-repo",
					"default_branch": "main",
					"owner": {
						"login": "testorg",
						"type": "Organization"
					},
					"private": false,
					"permissions": {
						"pull": true
					}
				}
			]`))
		}
	}))
	defer server.Close()

	client, err := NewClient(server.URL+"/api/v3", "test-token", false)
	if err != nil {
		t.Fatalf("Failed to create client: %v", err)
	}

	tree, err := client.BuildRepositoryTree()
	if err != nil {
		t.Fatalf("BuildRepositoryTree() error = %v", err)
	}

	if tree == nil {
		t.Fatal("BuildRepositoryTree() returned nil tree")
	}

	if len(tree.Groups) != 2 {
		t.Errorf("Expected 2 groups, got %d", len(tree.Groups))
	}

	if len(tree.Repositories) != 0 {
		t.Errorf("Expected 0 root repositories, got %d", len(tree.Repositories))
	}

	// Check organization group
	orgGroup, exists := tree.Groups["testorg"]
	if !exists {
		t.Error("Expected 'testorg' group to exist")
	}
	if orgGroup != nil && len(orgGroup.Repositories) != 1 {
		t.Errorf("Expected 1 repository in testorg group, got %d", len(orgGroup.Repositories))
	}

	// Check user group
	userGroup, exists := tree.Groups["testuser"]
	if !exists {
		t.Error("Expected 'testuser' group to exist")
	}
	if userGroup != nil && len(userGroup.Repositories) != 1 {
		t.Errorf("Expected 1 repository in testuser group, got %d", len(userGroup.Repositories))
	}
}

func TestBuildTreeFromRepos(t *testing.T) {
	repos := []*scm.Repository{
		{
			ID:            "123",
			Name:          "personal-repo",
			FullPath:      "testuser/personal-repo",
			CloneURL:      "https://github.com/testuser/personal-repo.git",
			SSHCloneURL:   "git@github.com:testuser/personal-repo.git",
			DefaultBranch: "main",
			WebURL:        "https://github.com/testuser/personal-repo",
			Provider:      "github",
		},
		{
			ID:            "456",
			Name:          "org-repo",
			FullPath:      "testorg/org-repo",
			CloneURL:      "https://github.com/testorg/org-repo.git",
			SSHCloneURL:   "git@github.com:testorg/org-repo.git",
			DefaultBranch: "main",
			WebURL:        "https://github.com/testorg/org-repo",
			Provider:      "github",
		},
	}

	tree := buildTreeFromRepos(repos)

	if tree == nil {
		t.Fatal("buildTreeFromRepos() returned nil")
	}

	if len(tree.Groups) != 2 {
		t.Errorf("Expected 2 groups, got %d", len(tree.Groups))
	}

	// Check testuser group
	userGroup, exists := tree.Groups["testuser"]
	if !exists {
		t.Error("Expected 'testuser' group to exist")
	}
	if userGroup != nil && len(userGroup.Repositories) != 1 {
		t.Errorf("Expected 1 repository in testuser group, got %d", len(userGroup.Repositories))
	}

	// Check testorg group
	orgGroup, exists := tree.Groups["testorg"]
	if !exists {
		t.Error("Expected 'testorg' group to exist")
	}
	if orgGroup != nil && len(orgGroup.Repositories) != 1 {
		t.Errorf("Expected 1 repository in testorg group, got %d", len(orgGroup.Repositories))
	}

	if len(tree.Repositories) != 0 {
		t.Errorf("Expected 0 root repositories, got %d", len(tree.Repositories))
	}
}
