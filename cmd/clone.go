package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/spf13/cobra"
	"gitstuff/internal/config"
	"gitstuff/internal/git"
	"gitstuff/internal/gitlab"
)

var cloneCmd = &cobra.Command{
	Use:   "clone [repository-path]",
	Short: "Clone GitLab repositories",
	Long: `Clone a specific GitLab repository or all repositories.
If no repository path is provided, all repositories will be cloned.`,
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
	
	client, err := gitlab.NewClient(cfg.GitLab.URL, cfg.GitLab.Token, cfg.GitLab.Insecure)
	if err != nil {
		return err
	}
	
	cloneAll, _ := cmd.Flags().GetBool("all")
	useSSH, _ := cmd.Flags().GetBool("ssh")
	update, _ := cmd.Flags().GetBool("update")
	
	if cloneAll || len(args) == 0 {
		return cloneAllRepositories(client, cfg, useSSH, update)
	}
	
	return cloneSingleRepository(client, cfg, args[0], useSSH, update)
}

func cloneAllRepositories(client *gitlab.Client, cfg *config.Config, useSSH, update bool) error {
	repos, err := client.ListAllRepositories()
	if err != nil {
		return err
	}
	
	fmt.Printf("Found %d repositories to clone/update\n\n", len(repos))
	
	successful := 0
	failed := 0
	
	for i, repo := range repos {
		fmt.Printf("[%d/%d] Processing %s...\n", i+1, len(repos), repo.FullPath)
		
		localPath := filepath.Join(cfg.Local.BaseDir, repo.FullPath)
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

func cloneSingleRepository(client *gitlab.Client, cfg *config.Config, repoPath string, useSSH, update bool) error {
	repo, err := client.GetRepository(repoPath)
	if err != nil {
		return err
	}
	
	fmt.Printf("Processing repository: %s\n", repo.FullPath)
	
	localPath := filepath.Join(cfg.Local.BaseDir, repo.FullPath)
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
	
	cloneURL := repo.CloneURL
	if useSSH {
		cloneURL = repo.SSHCloneURL
	}
	
	fmt.Printf("📥 Cloning from %s to %s...\n", cloneURL, localPath)
	if err := git.CloneRepository(cloneURL, localPath, useSSH); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}
	
	fmt.Printf("✅ Repository cloned successfully\n")
	return nil
}