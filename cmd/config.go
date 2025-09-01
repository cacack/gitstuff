package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"gitstuff/internal/config"
	"gitstuff/internal/verbosity"

	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configure SCM provider settings",
	Long:  `Configure GitLab or GitHub connection settings interactively.`,
	RunE:  runConfig,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().StringP("provider", "p", "", "Provider type (gitlab or github)")
	configCmd.Flags().StringP("name", "n", "", "Provider name (identifier)")
	configCmd.Flags().StringP("url", "u", "", "Provider instance URL")
	configCmd.Flags().StringP("token", "t", "", "Access token")
	configCmd.Flags().StringP("base-dir", "d", "", "Base directory for cloned repositories")
	configCmd.Flags().BoolP("insecure", "k", false, "Skip SSL certificate verification (for self-signed certificates)")
	configCmd.Flags().StringP("group", "g", "", "Default group/organization to filter repositories (optional)")
}

func runConfig(cmd *cobra.Command, args []string) error {
	// Check if flags were provided for non-interactive setup
	providerType, _ := cmd.Flags().GetString("provider")
	name, _ := cmd.Flags().GetString("name")
	url, _ := cmd.Flags().GetString("url")
	token, _ := cmd.Flags().GetString("token")
	baseDir, _ := cmd.Flags().GetString("base-dir")
	insecure, _ := cmd.Flags().GetBool("insecure")
	group, _ := cmd.Flags().GetString("group")

	if providerType != "" {
		verbosity.Debug("Running config in non-interactive mode for provider: %s", providerType)
	} else {
		verbosity.Debug("Running config in interactive mode")
	}

	reader := bufio.NewReader(os.Stdin)

	// Interactive mode if no provider type specified
	if providerType == "" {
		fmt.Println("Available SCM providers:")
		fmt.Println("1. GitLab")
		fmt.Println("2. GitHub")
		fmt.Print("Select a provider (1-2): ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			providerType = "gitlab"
		case "2":
			providerType = "github"
		default:
			return fmt.Errorf("invalid selection: %s", choice)
		}
	}

	// Validate provider type
	if providerType != "gitlab" && providerType != "github" {
		return fmt.Errorf("unsupported provider type: %s", providerType)
	}

	// Get provider name
	if name == "" {
		fmt.Printf("Provider name (identifier for this %s instance): ", providerType)
		name, _ = reader.ReadString('\n')
		name = strings.TrimSpace(name)
		if name == "" {
			name = providerType // Default to provider type
		}
	}

	// Get URL
	if url == "" {
		if providerType == "gitlab" {
			fmt.Print("GitLab URL (e.g., https://gitlab.com or gitlab.example.com): ")
		} else {
			fmt.Print("GitHub URL (leave blank for github.com or enter GitHub Enterprise URL): ")
		}
		url, _ = reader.ReadString('\n')
		url = strings.TrimSpace(url)

		if url == "" && providerType == "github" {
			url = "https://github.com"
		}
	}

	// Get token
	if token == "" {
		if providerType == "gitlab" {
			fmt.Print("GitLab Access Token: ")
		} else {
			fmt.Print("GitHub Personal Access Token: ")
		}
		tokenBytes, err := term.ReadPassword(syscall.Stdin)
		if err != nil {
			return fmt.Errorf("failed to read token: %w", err)
		}
		token = string(tokenBytes)
		fmt.Println()
	}

	// Get base directory
	if baseDir == "" && !cmd.Flags().Changed("base-dir") {
		fmt.Print("Base directory for repositories (default: ~/gitstuff-repos): ")
		baseDir, _ = reader.ReadString('\n')
		baseDir = strings.TrimSpace(baseDir)
	}

	// Get insecure setting (mainly for GitLab)
	if !insecure && !cmd.Flags().Changed("insecure") && providerType == "gitlab" {
		fmt.Print("Skip SSL certificate verification? (y/N): ")
		response, _ := reader.ReadString('\n')
		response = strings.ToLower(strings.TrimSpace(response))
		insecure = response == "y" || response == "yes"
	}

	// Get group/organization filter
	if group == "" && !cmd.Flags().Changed("group") {
		if providerType == "gitlab" {
			fmt.Print("Default GitLab group to filter repositories (optional, leave blank for all): ")
		} else {
			fmt.Print("Default GitHub organization to filter repositories (optional, leave blank for all): ")
		}
		group, _ = reader.ReadString('\n')
		group = strings.TrimSpace(group)
	}

	// Add the provider
	err := config.AddProvider(name, providerType, url, token, baseDir, insecure, group)
	if err != nil {
		return err
	}

	// Ask if user wants to add another provider
	fmt.Print("Would you like to add another provider? (y/N): ")
	response, _ := reader.ReadString('\n')
	response = strings.ToLower(strings.TrimSpace(response))

	if response == "y" || response == "yes" {
		// Recursively run the interactive config for the next provider
		// but skip base-dir since it's already set
		return runConfig(cmd, args)
	}

	fmt.Println("Configuration complete!")
	return nil
}
