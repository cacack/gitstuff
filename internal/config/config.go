package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Providers []ProviderConfig `yaml:"providers"`
	Local     LocalConfig      `yaml:"local"`
}

type ProviderConfig struct {
	Name     string `yaml:"name"`
	Type     string `yaml:"type"` // "gitlab" or "github"
	URL      string `yaml:"url"`
	Token    string `yaml:"token"`
	Insecure bool   `yaml:"insecure"`
	Group    string `yaml:"group"`
}

type LocalConfig struct {
	BaseDir string `yaml:"base_dir"`
}

// Legacy LocalConfig with different field name
type LegacyLocalConfig struct {
	BaseDir string `yaml:"basedir"`
}

// Legacy support structures
type GitLabConfig struct {
	URL      string `yaml:"url"`
	Token    string `yaml:"token"`
	Insecure bool   `yaml:"insecure"`
	Group    string `yaml:"group"`
}

type LegacyConfig struct {
	GitLab GitLabConfig      `yaml:"gitlab"`
	Local  LegacyLocalConfig `yaml:"local"`
}

func Load() (*Config, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	configPath := filepath.Join(home, ".gitstuff.yaml")

	// Read config file directly
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("config file not found at %s - run 'gitstuff config' to set up", configPath)
	}

	var config Config
	var legacyConfig LegacyConfig

	// Try to unmarshal as new format first
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	// If no providers but legacy GitLab config exists, migrate it
	if len(config.Providers) == 0 {
		if err := yaml.Unmarshal(data, &legacyConfig); err == nil && legacyConfig.GitLab.URL != "" {
			config.Providers = []ProviderConfig{
				{
					Name:     "gitlab",
					Type:     "gitlab",
					URL:      legacyConfig.GitLab.URL,
					Token:    legacyConfig.GitLab.Token,
					Insecure: legacyConfig.GitLab.Insecure,
					Group:    legacyConfig.GitLab.Group,
				},
			}
			config.Local = LocalConfig{BaseDir: legacyConfig.Local.BaseDir}

			// Save migrated config
			if saveErr := saveConfig(&config, configPath); saveErr != nil {
				return nil, fmt.Errorf("failed to save migrated config: %w", saveErr)
			}
		}
	}

	if len(config.Providers) == 0 {
		return nil, fmt.Errorf("no SCM providers configured - run 'gitstuff config' to set up")
	}

	// Validate provider configurations
	for _, provider := range config.Providers {
		if provider.URL == "" || provider.Token == "" {
			return nil, fmt.Errorf("provider %s is missing URL or token", provider.Name)
		}
		if provider.Type != "gitlab" && provider.Type != "github" {
			return nil, fmt.Errorf("provider %s has unsupported type %s", provider.Name, provider.Type)
		}
	}

	if config.Local.BaseDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user home directory: %w", err)
		}
		config.Local.BaseDir = filepath.Join(home, "gitstuff-repos")
	}

	return &config, nil
}

func AddProvider(name, providerType, url, token, baseDir string, insecure bool, group string) error {
	// Validate input parameters
	if name == "" {
		return fmt.Errorf("provider name is required")
	}
	if providerType == "" {
		return fmt.Errorf("provider type is required")
	}
	if providerType != "gitlab" && providerType != "github" {
		return fmt.Errorf("unsupported provider type: %s (supported: gitlab, github)", providerType)
	}
	if url == "" {
		return fmt.Errorf("provider URL is required")
	}
	if token == "" {
		return fmt.Errorf("provider token is required")
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	configPath := filepath.Join(home, ".gitstuff.yaml")

	// Load existing config or create new one
	var config Config
	if data, err := os.ReadFile(configPath); err == nil {
		if err := yaml.Unmarshal(data, &config); err != nil {
			return fmt.Errorf("failed to unmarshal existing config: %w", err)
		}
	}

	// Set default base directory if not set
	if config.Local.BaseDir == "" {
		if baseDir == "" {
			baseDir = filepath.Join(home, "gitstuff-repos")
		}
		config.Local.BaseDir = baseDir
	}

	// Check if provider already exists
	for i, provider := range config.Providers {
		if provider.Name == name {
			config.Providers[i] = ProviderConfig{
				Name:     name,
				Type:     providerType,
				URL:      url,
				Token:    token,
				Insecure: insecure,
				Group:    group,
			}
			return saveConfig(&config, configPath)
		}
	}

	// Add new provider
	config.Providers = append(config.Providers, ProviderConfig{
		Name:     name,
		Type:     providerType,
		URL:      url,
		Token:    token,
		Insecure: insecure,
		Group:    group,
	})

	return saveConfig(&config, configPath)
}

func saveConfig(config *Config, configPath string) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	err = os.WriteFile(configPath, data, 0600)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	fmt.Printf("Configuration updated at: %s\n", configPath)
	return nil
}

// Legacy Create function for backward compatibility
func Create(gitlabURL, token, baseDir string, insecure bool, group string) error {
	return AddProvider("gitlab", "gitlab", gitlabURL, token, baseDir, insecure, group)
}
