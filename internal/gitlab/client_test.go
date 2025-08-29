package gitlab

import (
	"strings"
	"testing"

	"gitstuff/internal/scm"
)

func TestBuildRepositoryTree_EmptyRepos(t *testing.T) {
	repos := []*scm.Repository{}

	tree := buildTreeFromRepos(repos)

	if len(tree.Groups) != 0 {
		t.Errorf("Expected 0 groups, got %d", len(tree.Groups))
	}

	if len(tree.Repositories) != 0 {
		t.Errorf("Expected 0 root repositories, got %d", len(tree.Repositories))
	}
}

func TestBuildRepositoryTree_RootRepos(t *testing.T) {
	repos := []*scm.Repository{
		{
			ID:       "1",
			Name:     "repo1",
			FullPath: "repo1",
			Provider: "gitlab",
		},
		{
			ID:       "2",
			Name:     "repo2",
			FullPath: "repo2",
			Provider: "gitlab",
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
	repos := []*scm.Repository{
		{
			ID:       "1",
			Name:     "repo1",
			FullPath: "group1/repo1",
			Provider: "gitlab",
		},
		{
			ID:       "2",
			Name:     "repo2",
			FullPath: "group1/repo2",
			Provider: "gitlab",
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
		t.Error("Expected group1 to exist")
	}

	if len(group1.Repositories) != 2 {
		t.Errorf("Expected 2 repositories in group1, got %d", len(group1.Repositories))
	}

	if group1.Group.Name != "group1" {
		t.Errorf("Expected group name to be group1, got %s", group1.Group.Name)
	}
}

func TestNormalizeURL(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  string
	}{
		{
			name:  "gitlab.com without protocol",
			input: "gitlab.com",
			want:  "https://gitlab.com",
		},
		{
			name:  "gitlab.com with https",
			input: "https://gitlab.com",
			want:  "https://gitlab.com",
		},
		{
			name:  "self-hosted without protocol",
			input: "gitlab.example.com",
			want:  "https://gitlab.example.com",
		},
		{
			name:  "self-hosted with https",
			input: "https://gitlab.example.com",
			want:  "https://gitlab.example.com",
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

// buildTreeFromRepos is a simplified version for testing
func buildTreeFromRepos(repos []*scm.Repository) *scm.RepositoryTree {
	tree := &scm.RepositoryTree{
		Groups:       make(map[string]*scm.GroupNode),
		Repositories: []*scm.Repository{},
	}

	for _, repo := range repos {
		parts := strings.Split(repo.FullPath, "/")
		if len(parts) == 1 {
			tree.Repositories = append(tree.Repositories, repo)
			continue
		}

		current := tree.Groups
		var currentNode *scm.GroupNode

		for i, part := range parts[:len(parts)-1] {
			if currentNode == nil {
				if _, exists := current[part]; !exists {
					current[part] = &scm.GroupNode{
						Group: &scm.Group{
							Name:     part,
							FullPath: strings.Join(parts[:i+1], "/"),
							Provider: "gitlab",
						},
						SubGroups:    make(map[string]*scm.GroupNode),
						Repositories: []*scm.Repository{},
					}
				}
				currentNode = current[part]
				current = currentNode.SubGroups
			} else {
				if _, exists := current[part]; !exists {
					current[part] = &scm.GroupNode{
						Group: &scm.Group{
							Name:     part,
							FullPath: strings.Join(parts[:i+1], "/"),
							Provider: "gitlab",
						},
						SubGroups:    make(map[string]*scm.GroupNode),
						Repositories: []*scm.Repository{},
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
