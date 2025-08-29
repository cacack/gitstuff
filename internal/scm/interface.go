package scm

// Repository represents a repository from any SCM provider
type Repository struct {
	ID            string
	Name          string
	FullPath      string
	CloneURL      string
	SSHCloneURL   string
	DefaultBranch string
	WebURL        string
	Provider      string // "gitlab" or "github"
}

// Group represents a group/organization from any SCM provider
type Group struct {
	ID       string
	Name     string
	FullPath string
	Provider string
}

// RepositoryTree represents the hierarchical structure of repositories
type RepositoryTree struct {
	Groups       map[string]*GroupNode
	Repositories []*Repository
}

// GroupNode represents a group and its contents in the tree
type GroupNode struct {
	Group        *Group
	SubGroups    map[string]*GroupNode
	Repositories []*Repository
}

// Client interface that both GitLab and GitHub clients must implement
type Client interface {
	// ListAllRepositories returns all repositories the user has access to
	ListAllRepositories() ([]*Repository, error)

	// ListRepositoriesInGroup returns repositories within a specific group/organization
	ListRepositoriesInGroup(groupPath string) ([]*Repository, error)

	// BuildRepositoryTree builds a hierarchical tree structure of repositories
	BuildRepositoryTree() (*RepositoryTree, error)

	// GetProviderType returns the provider type ("gitlab" or "github")
	GetProviderType() string
}
