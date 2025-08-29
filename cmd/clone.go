package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"gitstuff/internal/config"
	"gitstuff/internal/git"
	"gitstuff/internal/scm"

	"github.com/spf13/cobra"
)

var cloneCmd = &cobra.Command{
	Use:   "clone [repository-path]",
	Short: "Clone repositories from configured SCM providers",
	Long: `Clone a specific repository or all repositories from configured providers.
If no repository path is provided, all repositories will be cloned.
Repository path format: 'owner/repo' (searches all providers)`,
	RunE: runClone,
}

func init() {
	rootCmd.AddCommand(cloneCmd)
	cloneCmd.Flags().BoolP("all", "a", false, "Clone all repositories")
	cloneCmd.Flags().BoolP("ssh", "s", false, "Use SSH for cloning (default: HTTPS)")
	cloneCmd.Flags().BoolP("update", "u", false, "Pull latest changes for already cloned repositories")
}

func runClone(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w (run 'gitstuff config' first)", err)
	}

	if len(cfg.Providers) == 0 {
		return fmt.Errorf("no providers configured")
	}

	// Create clients for all providers
	clients := make([]scm.Client, 0, len(cfg.Providers))
	for _, providerConfig := range cfg.Providers {
		client, err := createClient(providerConfig)
		if err != nil {
			return fmt.Errorf("failed to create client for provider %s: %w", providerConfig.Name, err)
		}
		clients = append(clients, client)
	}

	cloneAll, _ := cmd.Flags().GetBool("all")
	useSSH, _ := cmd.Flags().GetBool("ssh")
	update, _ := cmd.Flags().GetBool("update")

	if cloneAll || len(args) == 0 {
		return cloneAllRepositories(clients, cfg, useSSH, update)
	}

	return cloneSingleRepository(clients, cfg, args[0], useSSH, update)
}

func cloneAllRepositories(clients []scm.Client, cfg *config.Config, useSSH, update bool) error {
	var allRepos []*scm.Repository

	// Collect all repositories from all providers
	for _, client := range clients {
		repos, err := client.ListAllRepositories()
		if err != nil {
			fmt.Printf("❌ Error getting repositories from %s provider: %v\n", client.GetProviderType(), err)
			continue
		}
		allRepos = append(allRepos, repos...)
	}

	fmt.Printf("Found %d repositories to clone/update\n\n", len(allRepos))

	successful := 0
	failed := 0

	for i, repo := range allRepos {
		fmt.Printf("[%d/%d] Processing %s [%s]...\n", i+1, len(allRepos), repo.FullPath, repo.Provider)

		localPath := filepath.Join(cfg.Local.BaseDir, repo.Provider, repo.FullPath)
		status, err := git.GetRepositoryStatus(localPath)
		if err != nil {
			fmt.Printf("❌ Error checking status: %v\n\n", err)
			failed++
			continue
		}

		if status.Exists && status.IsGitRepo {
			if update {
				fmt.Printf("🔄 Pulling latest changes...\n")
				if err := git.PullRepository(localPath); err != nil {
					fmt.Printf("❌ Failed to pull: %v\n\n", err)
					failed++
				} else {
					fmt.Printf("✅ Updated successfully\n\n")
					successful++
				}
			} else {
				fmt.Printf("⏭️  Already cloned (use --update to pull latest changes)\n\n")
				successful++
			}
			continue
		}

		cloneURL := repo.CloneURL
		if useSSH {
			cloneURL = repo.SSHCloneURL
		}

		fmt.Printf("📥 Cloning from %s...\n", cloneURL)
		if err := git.CloneRepository(cloneURL, localPath, useSSH); err != nil {
			fmt.Printf("❌ Failed to clone: %v\n\n", err)
			failed++
		} else {
			fmt.Printf("✅ Cloned successfully\n\n")
			successful++
		}
	}

	fmt.Printf("Summary: %d successful, %d failed\n", successful, failed)
	return nil
}

func cloneSingleRepository(clients []scm.Client, cfg *config.Config, repoPath string, useSSH, update bool) error {
	// Search for the repository across all providers
	var foundRepo *scm.Repository

	for _, client := range clients {
		// Try to find the repository in this provider
		repo, err := findRepositoryByPath(client, repoPath)
		if err == nil && repo != nil {
			foundRepo = repo
			break
		}
	}

	if foundRepo == nil {
		return fmt.Errorf("repository '%s' not found in any configured provider", repoPath)
	}

	fmt.Printf("Found repository: %s [%s]\n", foundRepo.FullPath, foundRepo.Provider)

	localPath := filepath.Join(cfg.Local.BaseDir, foundRepo.Provider, foundRepo.FullPath)
	status, err := git.GetRepositoryStatus(localPath)
	if err != nil {
		return fmt.Errorf("error checking repository status: %w", err)
	}

	if status.Exists && status.IsGitRepo {
		if update {
			fmt.Printf("🔄 Pulling latest changes...\n")
			if err := git.PullRepository(localPath); err != nil {
				return fmt.Errorf("failed to pull repository: %w", err)
			}
			fmt.Printf("✅ Repository updated successfully\n")
		} else {
			fmt.Printf("⏭️  Repository already cloned at: %s\n", localPath)
			fmt.Printf("   Use --update flag to pull latest changes\n")
		}
		return nil
	}

	if status.Exists && !status.IsGitRepo {
		return fmt.Errorf("directory %s exists but is not a git repository", localPath)
	}

	cloneURL := foundRepo.CloneURL
	if useSSH {
		cloneURL = foundRepo.SSHCloneURL
	}

	fmt.Printf("📥 Cloning from %s to %s...\n", cloneURL, localPath)
	if err := git.CloneRepository(cloneURL, localPath, useSSH); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	fmt.Printf("✅ Repository cloned successfully\n")
	return nil
}

// findRepositoryByPath searches for a repository by its path (owner/repo format)
func findRepositoryByPath(client scm.Client, repoPath string) (*scm.Repository, error) {
	// Get all repositories from this provider
	repos, err := client.ListAllRepositories()
	if err != nil {
		return nil, err
	}

	// Search for exact match or partial match
	for _, repo := range repos {
		if repo.FullPath == repoPath || strings.HasSuffix(repo.FullPath, repoPath) {
			return repo, nil
		}
	}

	return nil, fmt.Errorf("repository not found")
}
