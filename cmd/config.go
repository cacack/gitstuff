package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"gitstuff/internal/config"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure GitLab connection settings",
	Long:  `Configure GitLab URL, access token, and local repository base directory.`,
	RunE:  runConfig,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().StringP("url", "u", "", "GitLab instance URL")
	configCmd.Flags().StringP("token", "t", "", "GitLab access token")
	configCmd.Flags().StringP("base-dir", "d", "", "Base directory for cloned repositories")
	configCmd.Flags().BoolP("insecure", "k", false, "Skip SSL certificate verification (for self-signed certificates)")
}

func runConfig(cmd *cobra.Command, args []string) error {
	url, _ := cmd.Flags().GetString("url")
	token, _ := cmd.Flags().GetString("token")
	baseDir, _ := cmd.Flags().GetString("base-dir")
	insecure, _ := cmd.Flags().GetBool("insecure")

	reader := bufio.NewReader(os.Stdin)

	if url == "" {
		fmt.Print("GitLab URL (e.g., https://gitlab.com): ")
		url, _ = reader.ReadString('\n')
		url = strings.TrimSpace(url)
	}

	if token == "" {
		fmt.Print("GitLab Access Token: ")
		tokenBytes, err := term.ReadPassword(syscall.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read token: %w", err)
		}
		token = string(tokenBytes)
		fmt.Println()
	}

	if baseDir == "" {
		fmt.Print("Base directory for repositories (default: ~/gitlab-repos): ")
		baseDir, _ = reader.ReadString('\n')
		baseDir = strings.TrimSpace(baseDir)
	}

	if !insecure && !cmd.Flags().Changed("insecure") {
		fmt.Print("Skip SSL certificate verification? (y/N): ")
		response, _ := reader.ReadString('\n')
		response = strings.ToLower(strings.TrimSpace(response))
		insecure = response == "y" || response == "yes"
	}

	return config.Create(url, token, baseDir, insecure)
}
