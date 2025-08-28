package cmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"gitstuff/internal/config"
	"gitstuff/internal/git"
	"gitstuff/internal/gitlab"

	"github.com/spf13/cobra"
)

// GitLabClientInterface defines the methods we need from the GitLab client
type GitLabClientInterface interface {
	ListAllRepositories() ([]*gitlab.Repository, error)
	ListRepositoriesInGroup(groupPath string) ([]*gitlab.Repository, error)
	BuildRepositoryTree() (*gitlab.RepositoryTree, error)
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all GitLab repositories",
	Long:  `List all GitLab repositories with their status including clone status and current branch.`,
	RunE:  runList,
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().BoolP("tree", "t", false, "Display repositories in tree structure with groups")
	listCmd.Flags().BoolP("status", "s", true, "Show local repository status")
	listCmd.Flags().BoolP("verbose", "v", false, "Show additional details like URLs")
	listCmd.Flags().StringP("group", "g", "", "Filter repositories to only those in the specified group")
}

func runList(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("failed to load config: %w (run 'gitstuff config' first)", err)
	}

	client, err := gitlab.NewClient(cfg.GitLab.URL, cfg.GitLab.Token, cfg.GitLab.Insecure)
	if err != nil {
		return err
	}

	showTree, _ := cmd.Flags().GetBool("tree")
	showStatus, _ := cmd.Flags().GetBool("status")
	showVerbose, _ := cmd.Flags().GetBool("verbose")
	groupFilter, _ := cmd.Flags().GetString("group")

	// Use group from flag first, then config, then empty string
	targetGroup := groupFilter
	if targetGroup == "" {
		targetGroup = cfg.GitLab.Group
	}

	if showTree {
		return displayRepositoryTree(client, cfg, showStatus, showVerbose, targetGroup)
	} else {
		return displayRepositoryList(client, cfg, showStatus, showVerbose, targetGroup)
	}
}

func displayRepositoryList(client GitLabClientInterface, cfg *config.Config, showStatus, showVerbose bool, groupFilter string) error {
	var repos []*gitlab.Repository
	var err error

	if groupFilter != "" {
		repos, err = client.ListRepositoriesInGroup(groupFilter)
	} else {
		repos, err = client.ListAllRepositories()
	}
	if err != nil {
		return err
	}

	fmt.Printf("Found %d repositories:\n\n", len(repos))

	for _, repo := range repos {
		fmt.Printf("📁 %s\n", repo.FullPath)

		if showVerbose {
			fmt.Printf("   URL: %s\n", repo.WebURL)
		}

		if showStatus {
			localPath := filepath.Join(cfg.Local.BaseDir, repo.FullPath)
			status, err := git.GetRepositoryStatus(localPath)
			if err != nil {
				fmt.Printf("   Status: ❌ Error checking status: %v\n", err)
			} else {
				displayStatus(status)
			}
		}

		fmt.Print("\n")
	}

	return nil
}

func displayRepositoryTree(client GitLabClientInterface, cfg *config.Config, showStatus, showVerbose bool, groupFilter string) error {
	tree, err := client.BuildRepositoryTree()
	if err != nil {
		return err
	}

	if groupFilter != "" {
		fmt.Printf("Repository tree structure (filtered by group: %s):\n", groupFilter)
		displayFilteredTree(tree, groupFilter, cfg, showStatus, showVerbose)
	} else {
		fmt.Println("Repository tree structure:")

		if len(tree.Repositories) > 0 {
			fmt.Println("Root repositories:")
			for _, repo := range tree.Repositories {
				fmt.Printf("📁 %s\n", repo.Name)

				if showVerbose {
					fmt.Printf("   URL: %s\n", repo.WebURL)
				}

				if showStatus {
					localPath := filepath.Join(cfg.Local.BaseDir, repo.FullPath)
					status, err := git.GetRepositoryStatus(localPath)
					if err != nil {
						fmt.Printf("   Status: ❌ Error: %v\n", err)
					} else {
						displayStatus(status)
					}
				}
				fmt.Print("\n")
			}
		}

		for groupName, groupNode := range tree.Groups {
			displayGroup(groupNode, 0, cfg, showStatus, showVerbose)
			_ = groupName
		}
	}

	return nil
}

func displayFilteredTree(tree *gitlab.RepositoryTree, groupFilter string, cfg *config.Config, showStatus, showVerbose bool) {
	targetGroup := findGroupInTree(tree, groupFilter)
	if targetGroup != nil {
		displayGroup(targetGroup, 0, cfg, showStatus, showVerbose)
	} else {
		fmt.Printf("Group '%s' not found\n", groupFilter)
	}
}

func findGroupInTree(tree *gitlab.RepositoryTree, groupPath string) *gitlab.GroupNode {
	parts := strings.Split(groupPath, "/")

	current := tree.Groups
	var currentNode *gitlab.GroupNode

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

func displayGroup(group *gitlab.GroupNode, indent int, cfg *config.Config, showStatus, showVerbose bool) {
	prefix := strings.Repeat("  ", indent)
	fmt.Printf("%s📂 %s/\n", prefix, group.Group.Name)

	for _, repo := range group.Repositories {
		fmt.Printf("%s  📁 %s\n", prefix, repo.Name)

		if showVerbose {
			fmt.Printf("%s     URL: %s\n", prefix, repo.WebURL)
		}

		if showStatus {
			localPath := filepath.Join(cfg.Local.BaseDir, repo.FullPath)
			status, err := git.GetRepositoryStatus(localPath)
			if err != nil {
				fmt.Printf("%s     Status: ❌ Error: %v\n", prefix, err)
			} else {
				fmt.Printf("%s     ", prefix)
				displayStatus(status)
			}
		}
		fmt.Print("\n")
	}

	for _, subGroup := range group.SubGroups {
		displayGroup(subGroup, indent+1, cfg, showStatus, showVerbose)
	}
}

func displayStatus(status *git.Status) {
	if !status.Exists {
		fmt.Print("Status: ❌ Not cloned\n")
		return
	}

	if !status.IsGitRepo {
		fmt.Print("Status: ⚠️  Directory exists but not a git repository\n")
		return
	}

	fmt.Printf("Status: ✅ Cloned")
	if status.CurrentBranch != "" {
		fmt.Printf(" (branch: %s)", status.CurrentBranch)
	}
	if status.HasChanges {
		fmt.Print(" 🔄 Has uncommitted changes")
	}
	fmt.Print("\n")
}
