package gitlab

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strings"

	"github.com/xanzy/go-gitlab"
)

type Client struct {
	client *gitlab.Client
}

type Repository struct {
	ID            int
	Name          string
	FullPath      string
	CloneURL      string
	SSHCloneURL   string
	DefaultBranch string
	WebURL        string
}

type Group struct {
	ID       int
	Name     string
	FullPath string
}

func NewClient(baseURL, token string, insecure bool) (*Client, error) {
	normalizedURL, err := normalizeURL(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid GitLab URL: %w", err)
	}
	
	var options []gitlab.ClientOptionFunc
	options = append(options, gitlab.WithBaseURL(normalizedURL))
	
	if insecure {
		httpClient := &http.Client{
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: true,
				},
			},
		}
		options = append(options, gitlab.WithHTTPClient(httpClient))
	}
	
	client, err := gitlab.NewClient(token, options...)
	if err != nil {
		return nil, fmt.Errorf("failed to create gitlab client: %w", err)
	}
	
	return &Client{client: client}, nil
}

func normalizeURL(baseURL string) (string, error) {
	if baseURL == "" {
		return "", fmt.Errorf("URL cannot be empty")
	}
	
	if !strings.HasPrefix(baseURL, "http://") && !strings.HasPrefix(baseURL, "https://") {
		baseURL = "https://" + baseURL
	}
	
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("failed to parse URL: %w", err)
	}
	
	if parsedURL.Host == "" {
		return "", fmt.Errorf("URL must have a valid host")
	}
	
	return parsedURL.String(), nil
}

func (c *Client) ListAllRepositories() ([]*Repository, error) {
	var allRepos []*Repository
	
	opts := &gitlab.ListProjectsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
			Page:    1,
		},
		Membership: gitlab.Bool(true),
		Simple:     gitlab.Bool(false),
		OrderBy:    gitlab.String("path"),
		Sort:       gitlab.String("asc"),
	}
	
	for {
		projects, resp, err := c.client.Projects.ListProjects(opts)
		if err != nil {
			return nil, fmt.Errorf("failed to list projects: %w", err)
		}
		
		for _, project := range projects {
			repo := &Repository{
				ID:            project.ID,
				Name:          project.Name,
				FullPath:      project.PathWithNamespace,
				CloneURL:      project.HTTPURLToRepo,
				SSHCloneURL:   project.SSHURLToRepo,
				DefaultBranch: project.DefaultBranch,
				WebURL:        project.WebURL,
			}
			allRepos = append(allRepos, repo)
		}
		
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	
	sort.Slice(allRepos, func(i, j int) bool {
		return allRepos[i].FullPath < allRepos[j].FullPath
	})
	
	return allRepos, nil
}

func (c *Client) GetRepository(fullPath string) (*Repository, error) {
	project, _, err := c.client.Projects.GetProject(fullPath, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to get project %s: %w", fullPath, err)
	}
	
	return &Repository{
		ID:            project.ID,
		Name:          project.Name,
		FullPath:      project.PathWithNamespace,
		CloneURL:      project.HTTPURLToRepo,
		SSHCloneURL:   project.SSHURLToRepo,
		DefaultBranch: project.DefaultBranch,
		WebURL:        project.WebURL,
	}, nil
}

func (c *Client) ListGroups() ([]*Group, error) {
	var allGroups []*Group
	
	opts := &gitlab.ListGroupsOptions{
		ListOptions: gitlab.ListOptions{
			PerPage: 100,
			Page:    1,
		},
		AllAvailable: gitlab.Bool(true),
		OrderBy:      gitlab.String("path"),
		Sort:         gitlab.String("asc"),
	}
	
	for {
		groups, resp, err := c.client.Groups.ListGroups(opts)
		if err != nil {
			return nil, fmt.Errorf("failed to list groups: %w", err)
		}
		
		for _, group := range groups {
			g := &Group{
				ID:       group.ID,
				Name:     group.Name,
				FullPath: group.FullPath,
			}
			allGroups = append(allGroups, g)
		}
		
		if resp.NextPage == 0 {
			break
		}
		opts.Page = resp.NextPage
	}
	
	return allGroups, nil
}

type RepositoryTree struct {
	Groups       map[string]*GroupNode
	Repositories []*Repository
}

type GroupNode struct {
	Group        *Group
	SubGroups    map[string]*GroupNode
	Repositories []*Repository
}

func (c *Client) BuildRepositoryTree() (*RepositoryTree, error) {
	repos, err := c.ListAllRepositories()
	if err != nil {
		return nil, err
	}
	
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
	
	return tree, nil
}