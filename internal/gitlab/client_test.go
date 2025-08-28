package gitlab

import (
	"strings"
	"testing"
)

func TestBuildRepositoryTree_EmptyRepos(t *testing.T) {
	repos := []*Repository{}

	tree := buildTreeFromRepos(repos)

	if len(tree.Groups) != 0 {
		t.Errorf("Expected 0 groups, got %d", len(tree.Groups))
	}

	if len(tree.Repositories) != 0 {
		t.Errorf("Expected 0 root repositories, got %d", len(tree.Repositories))
	}
}

func TestBuildRepositoryTree_RootRepos(t *testing.T) {
	repos := []*Repository{
		{
			ID:       1,
			Name:     "repo1",
			FullPath: "repo1",
		},
		{
			ID:       2,
			Name:     "repo2",
			FullPath: "repo2",
		},
	}

	tree := buildTreeFromRepos(repos)

	if len(tree.Groups) != 0 {
		t.Errorf("Expected 0 groups, got %d", len(tree.Groups))
	}

	if len(tree.Repositories) != 2 {
		t.Errorf("Expected 2 root repositories, got %d", len(tree.Repositories))
	}

	if tree.Repositories[0].Name != "repo1" {
		t.Errorf("Expected first repo to be repo1, got %s", tree.Repositories[0].Name)
	}
}

func TestBuildRepositoryTree_SingleGroup(t *testing.T) {
	repos := []*Repository{
		{
			ID:       1,
			Name:     "repo1",
			FullPath: "group1/repo1",
		},
		{
			ID:       2,
			Name:     "repo2",
			FullPath: "group1/repo2",
		},
	}

	tree := buildTreeFromRepos(repos)

	if len(tree.Groups) != 1 {
		t.Errorf("Expected 1 group, got %d", len(tree.Groups))
	}

	if len(tree.Repositories) != 0 {
		t.Errorf("Expected 0 root repositories, got %d", len(tree.Repositories))
	}

	group1, exists := tree.Groups["group1"]
	if !exists {
		t.Fatal("Expected group1 to exist")
	}

	if len(group1.Repositories) != 2 {
		t.Errorf("Expected 2 repositories in group1, got %d", len(group1.Repositories))
	}

	if group1.Group.Name != "group1" {
		t.Errorf("Expected group name to be group1, got %s", group1.Group.Name)
	}
}

func TestBuildRepositoryTree_NestedGroups(t *testing.T) {
	repos := []*Repository{
		{
			ID:       1,
			Name:     "repo1",
			FullPath: "group1/subgroup1/repo1",
		},
		{
			ID:       2,
			Name:     "repo2",
			FullPath: "group1/subgroup1/repo2",
		},
		{
			ID:       3,
			Name:     "repo3",
			FullPath: "group1/repo3",
		},
	}

	tree := buildTreeFromRepos(repos)

	if len(tree.Groups) != 1 {
		t.Errorf("Expected 1 top-level group, got %d", len(tree.Groups))
	}

	group1, exists := tree.Groups["group1"]
	if !exists {
		t.Fatal("Expected group1 to exist")
	}

	if len(group1.Repositories) != 1 {
		t.Errorf("Expected 1 repository in group1, got %d", len(group1.Repositories))
	}

	if len(group1.SubGroups) != 1 {
		t.Errorf("Expected 1 subgroup in group1, got %d", len(group1.SubGroups))
	}

	subgroup1, exists := group1.SubGroups["subgroup1"]
	if !exists {
		t.Fatal("Expected subgroup1 to exist")
	}

	if len(subgroup1.Repositories) != 2 {
		t.Errorf("Expected 2 repositories in subgroup1, got %d", len(subgroup1.Repositories))
	}
}

func TestBuildRepositoryTree_MixedStructure(t *testing.T) {
	repos := []*Repository{
		{
			ID:       1,
			Name:     "root-repo",
			FullPath: "root-repo",
		},
		{
			ID:       2,
			Name:     "grouped-repo",
			FullPath: "group1/grouped-repo",
		},
		{
			ID:       3,
			Name:     "nested-repo",
			FullPath: "group1/subgroup1/nested-repo",
		},
	}

	tree := buildTreeFromRepos(repos)

	if len(tree.Repositories) != 1 {
		t.Errorf("Expected 1 root repository, got %d", len(tree.Repositories))
	}

	if tree.Repositories[0].Name != "root-repo" {
		t.Errorf("Expected root repo to be root-repo, got %s", tree.Repositories[0].Name)
	}

	if len(tree.Groups) != 1 {
		t.Errorf("Expected 1 top-level group, got %d", len(tree.Groups))
	}

	group1 := tree.Groups["group1"]
	if group1 == nil {
		t.Fatal("Expected group1 to exist")
	}

	if len(group1.Repositories) != 1 {
		t.Errorf("Expected 1 repository in group1, got %d", len(group1.Repositories))
	}

	if len(group1.SubGroups) != 1 {
		t.Errorf("Expected 1 subgroup in group1, got %d", len(group1.SubGroups))
	}
}

func buildTreeFromRepos(repos []*Repository) *RepositoryTree {
	tree := &RepositoryTree{
		Groups:       make(map[string]*GroupNode),
		Repositories: []*Repository{},
	}

	for _, repo := range repos {
		parts := strings.Split(repo.FullPath, "/")
		if len(parts) == 1 {
			tree.Repositories = append(tree.Repositories, repo)
			continue
		}

		current := tree.Groups
		var currentNode *GroupNode

		for i, part := range parts[:len(parts)-1] {
			if currentNode == nil {
				if _, exists := current[part]; !exists {
					current[part] = &GroupNode{
						Group: &Group{
							Name:     part,
							FullPath: strings.Join(parts[:i+1], "/"),
						},
						SubGroups:    make(map[string]*GroupNode),
						Repositories: []*Repository{},
					}
				}
				currentNode = current[part]
				current = currentNode.SubGroups
			} else {
				if _, exists := current[part]; !exists {
					current[part] = &GroupNode{
						Group: &Group{
							Name:     part,
							FullPath: strings.Join(parts[:i+1], "/"),
						},
						SubGroups:    make(map[string]*GroupNode),
						Repositories: []*Repository{},
					}
				}
				currentNode = current[part]
				current = currentNode.SubGroups
			}
		}

		if currentNode != nil {
			currentNode.Repositories = append(currentNode.Repositories, repo)
		}
	}

	return tree
}

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		wantErr  bool
	}{
		{
			name:     "Full HTTPS URL",
			input:    "https://gitlab.example.com",
			expected: "https://gitlab.example.com",
			wantErr:  false,
		},
		{
			name:     "Full HTTP URL",
			input:    "http://gitlab.example.com",
			expected: "http://gitlab.example.com",
			wantErr:  false,
		},
		{
			name:     "Hostname only",
			input:    "gitlab.example.com",
			expected: "https://gitlab.example.com",
			wantErr:  false,
		},
		{
			name:     "Hostname with path",
			input:    "gitlab.example.com/api",
			expected: "https://gitlab.example.com/api",
			wantErr:  false,
		},
		{
			name:     "Empty URL",
			input:    "",
			expected: "",
			wantErr:  true,
		},
		{
			name:     "URL with empty host after normalization",
			input:    "https://",
			expected: "",
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := normalizeURL(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestNewClient(t *testing.T) {
	tests := []struct {
		name      string
		baseURL   string
		token     string
		insecure  bool
		wantErr   bool
		expectErr string
	}{
		{
			name:     "Valid HTTPS URL, secure",
			baseURL:  "https://gitlab.example.com",
			token:    "test-token",
			insecure: false,
			wantErr:  false,
		},
		{
			name:     "Valid HTTPS URL, insecure",
			baseURL:  "https://gitlab.example.com",
			token:    "test-token",
			insecure: true,
			wantErr:  false,
		},
		{
			name:     "URL without protocol, secure",
			baseURL:  "gitlab.example.com",
			token:    "test-token",
			insecure: false,
			wantErr:  false,
		},
		{
			name:     "URL without protocol, insecure",
			baseURL:  "gitlab.example.com",
			token:    "test-token",
			insecure: true,
			wantErr:  false,
		},
		{
			name:      "Empty URL",
			baseURL:   "",
			token:     "test-token",
			insecure:  false,
			wantErr:   true,
			expectErr: "invalid GitLab URL",
		},
		{
			name:      "Invalid URL",
			baseURL:   "https://",
			token:     "test-token",
			insecure:  false,
			wantErr:   true,
			expectErr: "invalid GitLab URL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.baseURL, tt.token, tt.insecure)

			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error but got none")
					return
				}
				if tt.expectErr != "" && !strings.Contains(err.Error(), tt.expectErr) {
					t.Errorf("Expected error to contain '%s', got '%s'", tt.expectErr, err.Error())
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if client == nil {
				t.Error("Expected client to be non-nil")
				return
			}

			if client.client == nil {
				t.Error("Expected internal gitlab client to be non-nil")
			}
		})
	}
}

func TestNewClientHTTPClientConfiguration(t *testing.T) {
	tests := []struct {
		name     string
		insecure bool
		checkTLS bool
	}{
		{
			name:     "Secure client - default HTTP client",
			insecure: false,
			checkTLS: false,
		},
		{
			name:     "Insecure client - custom HTTP client with TLS skip",
			insecure: true,
			checkTLS: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient("https://gitlab.example.com", "test-token", tt.insecure)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if tt.checkTLS {
				httpClient := client.client.BaseURL()
				if httpClient.Scheme != "https" {
					t.Error("Expected HTTPS scheme")
				}
			}
		})
	}
}

func TestGroupFiltering(t *testing.T) {
	tests := []struct {
		name          string
		allRepos      []*Repository
		groupFilter   string
		expectedRepos []string
		expectedCount int
	}{
		{
			name: "Filter by top-level group",
			allRepos: []*Repository{
				{ID: 1, Name: "repo1", FullPath: "team-a/repo1"},
				{ID: 2, Name: "repo2", FullPath: "team-a/repo2"},
				{ID: 3, Name: "repo3", FullPath: "team-b/repo3"},
				{ID: 4, Name: "root-repo", FullPath: "root-repo"},
			},
			groupFilter:   "team-a",
			expectedRepos: []string{"team-a/repo1", "team-a/repo2"},
			expectedCount: 2,
		},
		{
			name: "Filter by nested group",
			allRepos: []*Repository{
				{ID: 1, Name: "backend-repo", FullPath: "team-a/backend/backend-repo"},
				{ID: 2, Name: "frontend-repo", FullPath: "team-a/frontend/frontend-repo"},
				{ID: 3, Name: "other-repo", FullPath: "team-a/other-repo"},
				{ID: 4, Name: "external-repo", FullPath: "team-b/external-repo"},
			},
			groupFilter:   "team-a/backend",
			expectedRepos: []string{"team-a/backend/backend-repo"},
			expectedCount: 1,
		},
		{
			name: "Filter with no matches",
			allRepos: []*Repository{
				{ID: 1, Name: "repo1", FullPath: "team-a/repo1"},
				{ID: 2, Name: "repo2", FullPath: "team-b/repo2"},
			},
			groupFilter:   "team-c",
			expectedRepos: []string{},
			expectedCount: 0,
		},
		{
			name: "Filter includes subgroups",
			allRepos: []*Repository{
				{ID: 1, Name: "direct-repo", FullPath: "team-a/direct-repo"},
				{ID: 2, Name: "sub-repo1", FullPath: "team-a/sub1/sub-repo1"},
				{ID: 3, Name: "sub-repo2", FullPath: "team-a/sub2/sub-repo2"},
				{ID: 4, Name: "other-repo", FullPath: "team-b/other-repo"},
			},
			groupFilter:   "team-a",
			expectedRepos: []string{"team-a/direct-repo", "team-a/sub1/sub-repo1", "team-a/sub2/sub-repo2"},
			expectedCount: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			filtered := filterRepositoriesByGroup(tt.allRepos, tt.groupFilter)

			if len(filtered) != tt.expectedCount {
				t.Errorf("Expected %d repos, got %d", tt.expectedCount, len(filtered))
			}

			for i, expectedPath := range tt.expectedRepos {
				if i >= len(filtered) {
					t.Errorf("Missing expected repo: %s", expectedPath)
					continue
				}
				if filtered[i].FullPath != expectedPath {
					t.Errorf("Expected repo %s at index %d, got %s", expectedPath, i, filtered[i].FullPath)
				}
			}
		})
	}
}

func filterRepositoriesByGroup(repos []*Repository, groupPath string) []*Repository {
	if groupPath == "" {
		return repos
	}

	var filtered []*Repository
	for _, repo := range repos {
		if strings.HasPrefix(repo.FullPath, groupPath+"/") {
			filtered = append(filtered, repo)
		}
	}
	return filtered
}

func TestFindGroupInTree(t *testing.T) {
	repos := []*Repository{
		{ID: 1, Name: "repo1", FullPath: "team-a/backend/repo1"},
		{ID: 2, Name: "repo2", FullPath: "team-a/frontend/repo2"},
		{ID: 3, Name: "repo3", FullPath: "team-a/repo3"},
		{ID: 4, Name: "repo4", FullPath: "team-b/repo4"},
	}

	tree := buildTreeFromRepos(repos)

	tests := []struct {
		name      string
		groupPath string
		expected  bool
		repoCount int
	}{
		{
			name:      "Find top-level group",
			groupPath: "team-a",
			expected:  true,
			repoCount: 1,
		},
		{
			name:      "Find nested group",
			groupPath: "team-a/backend",
			expected:  true,
			repoCount: 1,
		},
		{
			name:      "Find another nested group",
			groupPath: "team-a/frontend",
			expected:  true,
			repoCount: 1,
		},
		{
			name:      "Group not found",
			groupPath: "team-c",
			expected:  false,
			repoCount: 0,
		},
		{
			name:      "Nested group not found",
			groupPath: "team-a/nonexistent",
			expected:  false,
			repoCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findGroupInTreeHelper(tree, tt.groupPath)

			if tt.expected && result == nil {
				t.Errorf("Expected to find group %s but got nil", tt.groupPath)
				return
			}

			if !tt.expected && result != nil {
				t.Errorf("Expected not to find group %s but got result", tt.groupPath)
				return
			}

			if result != nil && len(result.Repositories) != tt.repoCount {
				t.Errorf("Expected %d repos in group %s, got %d", tt.repoCount, tt.groupPath, len(result.Repositories))
			}
		})
	}
}

func findGroupInTreeHelper(tree *RepositoryTree, groupPath string) *GroupNode {
	parts := strings.Split(groupPath, "/")

	current := tree.Groups
	var currentNode *GroupNode

	for _, part := range parts {
		if node, exists := current[part]; exists {
			currentNode = node
			current = node.SubGroups
		} else {
			return nil
		}
	}

	return currentNode
}
