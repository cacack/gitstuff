package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

type Config struct {
	GitLab GitLabConfig `yaml:"gitlab"`
	Local  LocalConfig  `yaml:"local"`
}

type GitLabConfig struct {
	URL      string `yaml:"url"`
	Token    string `yaml:"token"`
	Insecure bool   `yaml:"insecure"`
}

type LocalConfig struct {
	BaseDir string `yaml:"base_dir"`
}

func Load() (*Config, error) {
	var config Config
	
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}
	
	if config.GitLab.URL == "" || config.GitLab.Token == "" {
		return nil, fmt.Errorf("gitlab url and token must be configured")
	}
	
	if config.Local.BaseDir == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("failed to get user home directory: %w", err)
		}
		config.Local.BaseDir = filepath.Join(home, "gitlab-repos")
	}
	
	return &config, nil
}

func Create(gitlabURL, token, baseDir string, insecure bool) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}
	
	if baseDir == "" {
		baseDir = filepath.Join(home, "gitlab-repos")
	}
	
	config := Config{
		GitLab: GitLabConfig{
			URL:      gitlabURL,
			Token:    token,
			Insecure: insecure,
		},
		Local: LocalConfig{
			BaseDir: baseDir,
		},
	}
	
	configPath := filepath.Join(home, ".gitstuff.yaml")
	
	data, err := yaml.Marshal(&config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}
	
	err = os.WriteFile(configPath, data, 0600)
	if err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}
	
	fmt.Printf("Configuration created at: %s\n", configPath)
	return nil
}