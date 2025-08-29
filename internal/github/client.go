package github

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"

	"github.com/google/go-github/v67/github"
	"golang.org/x/oauth2"

	"gitstuff/internal/scm"
)

type Client struct {
	client *github.Client
	ctx    context.Context
}

func NewClient(baseURL, token string, insecure bool) (*Client, error) {
	ctx := context.Background()

	// Validate required parameters
	if token == "" {
		return nil, fmt.Errorf("GitHub access token is required")
	}
	if baseURL == "" {
		return nil, fmt.Errorf("GitHub base URL is required")
	}

	// Create HTTP client
	var httpClient *http.Client
	if insecure {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		httpClient = &http.Client{Transport: tr}
	}

	// Set up OAuth2 token source
	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})

	var tc *http.Client
	if httpClient != nil {
		// Combine OAuth2 with custom transport
		tc = &http.Client{
			Transport: &oauth2.Transport{
				Source: ts,
				Base:   httpClient.Transport,
			},
		}
	} else {
		tc = oauth2.NewClient(ctx, ts)
	}

	client := github.NewClient(tc)

	// Handle custom GitHub Enterprise URLs
	if baseURL != "https://github.com" && baseURL != "github.com" {
		normalizedURL, err := normalizeURL(baseURL)
		if err != nil {
			return nil, fmt.Errorf("invalid GitHub URL: %w", err)
		}

		// Ensure trailing slash for GitHub client
		if !strings.HasSuffix(normalizedURL, "/") {
			normalizedURL += "/"
		}

		baseURLParsed, err := url.Parse(normalizedURL)
		if err != nil {
			return nil, fmt.Errorf("failed to parse GitHub URL: %w", err)
		}

		client.BaseURL = baseURLParsed
	}

	return &Client{client: client, ctx: ctx}, nil
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

	// Ensure proper API endpoint for GitHub Enterprise
	if !strings.Contains(parsedURL.Path, "/api/v3") && parsedURL.Host != "github.com" {
		parsedURL.Path = "/api/v3"
	}

	return parsedURL.String(), nil
}

func (c *Client) GetProviderType() string {
	return "github"
}

func (c *Client) ListAllRepositories() ([]*scm.Repository, error) {
	var allRepos []*scm.Repository

	opts := &github.RepositoryListOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
		Sort:      "full_name",
		Direction: "asc",
	}

	for {
		repos, resp, err := c.client.Repositories.List(c.ctx, "", opts)
		if err != nil {
			return nil, fmt.Errorf("failed to list repositories: %w", err)
		}

		for _, repo := range repos {
			if repo.GetFullName() == "" || repo.GetPrivate() && !repo.GetPermissions()["pull"] {
				continue // Skip repos we don't have access to
			}

			scmRepo := &scm.Repository{
				ID:            strconv.FormatInt(repo.GetID(), 10),
				Name:          repo.GetName(),
				FullPath:      repo.GetFullName(),
				CloneURL:      repo.GetCloneURL(),
				SSHCloneURL:   repo.GetSSHURL(),
				DefaultBranch: repo.GetDefaultBranch(),
				WebURL:        repo.GetHTMLURL(),
				Provider:      "github",
			}
			allRepos = append(allRepos, scmRepo)
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

func (c *Client) ListRepositoriesInGroup(orgName string) ([]*scm.Repository, error) {
	var allRepos []*scm.Repository

	opts := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{
			PerPage: 100,
		},
		Sort:      "full_name",
		Direction: "asc",
	}

	for {
		repos, resp, err := c.client.Repositories.ListByOrg(c.ctx, orgName, opts)
		if err != nil {
			return nil, fmt.Errorf("failed to list repositories for organization %s: %w", orgName, err)
		}

		for _, repo := range repos {
			scmRepo := &scm.Repository{
				ID:            strconv.FormatInt(repo.GetID(), 10),
				Name:          repo.GetName(),
				FullPath:      repo.GetFullName(),
				CloneURL:      repo.GetCloneURL(),
				SSHCloneURL:   repo.GetSSHURL(),
				DefaultBranch: repo.GetDefaultBranch(),
				WebURL:        repo.GetHTMLURL(),
				Provider:      "github",
			}
			allRepos = append(allRepos, scmRepo)
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

func (c *Client) BuildRepositoryTree() (*scm.RepositoryTree, error) {
	repos, err := c.ListAllRepositories()
	if err != nil {
		return nil, err
	}

	return buildTreeFromRepos(repos), nil
}

func buildTreeFromRepos(repos []*scm.Repository) *scm.RepositoryTree {
	tree := &scm.RepositoryTree{
		Groups:       make(map[string]*scm.GroupNode),
		Repositories: []*scm.Repository{},
	}

	for _, repo := range repos {
		parts := strings.Split(repo.FullPath, "/")
		if len(parts) == 1 {
			// Root repository (shouldn't happen in GitHub but handle it)
			tree.Repositories = append(tree.Repositories, repo)
			continue
		}

		// In GitHub, the first part is always the organization/user
		orgName := parts[0]

		if _, exists := tree.Groups[orgName]; !exists {
			tree.Groups[orgName] = &scm.GroupNode{
				Group: &scm.Group{
					ID:       orgName,
					Name:     orgName,
					FullPath: orgName,
					Provider: "github",
				},
				SubGroups:    make(map[string]*scm.GroupNode),
				Repositories: []*scm.Repository{},
			}
		}

		tree.Groups[orgName].Repositories = append(tree.Groups[orgName].Repositories, repo)
	}

	return tree
}
