package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestAddProvider_GitLab(t *testing.T) {
	tempDir := t.TempDir()

	originalHome := os.Getenv("HOME")
	t.Cleanup(func() {
		os.Setenv("HOME", originalHome)
	})
	os.Setenv("HOME", tempDir)

	err := AddProvider("gitlab-main", "gitlab", "https://gitlab.com", "gl-token", "/custom/dir", false, "my-group")
	if err != nil {
		t.Fatalf("AddProvider failed: %v", err)
	}

	config, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(config.Providers) != 1 {
		t.Fatalf("Expected 1 provider, got %d", len(config.Providers))
	}

	provider := config.Providers[0]
	if provider.Name != "gitlab-main" {
		t.Errorf("Expected name 'gitlab-main', got '%s'", provider.Name)
	}
	if provider.Type != "gitlab" {
		t.Errorf("Expected type 'gitlab', got '%s'", provider.Type)
	}
	if provider.URL != "https://gitlab.com" {
		t.Errorf("Expected URL 'https://gitlab.com', got '%s'", provider.URL)
	}
	if provider.Token != "gl-token" {
		t.Errorf("Expected token 'gl-token', got '%s'", provider.Token)
	}
	if provider.Group != "my-group" {
		t.Errorf("Expected group 'my-group', got '%s'", provider.Group)
	}
	if config.Local.BaseDir != "/custom/dir" {
		t.Errorf("Expected base dir '/custom/dir', got '%s'", config.Local.BaseDir)
	}
}

func TestAddProvider_GitHub(t *testing.T) {
	tempDir := t.TempDir()

	originalHome := os.Getenv("HOME")
	t.Cleanup(func() {
		os.Setenv("HOME", originalHome)
	})
	os.Setenv("HOME", tempDir)

	err := AddProvider("github-main", "github", "https://github.com", "gh-token", "", false, "my-org")
	if err != nil {
		t.Fatalf("AddProvider failed: %v", err)
	}

	config, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(config.Providers) != 1 {
		t.Fatalf("Expected 1 provider, got %d", len(config.Providers))
	}

	provider := config.Providers[0]
	if provider.Name != "github-main" {
		t.Errorf("Expected name 'github-main', got '%s'", provider.Name)
	}
	if provider.Type != "github" {
		t.Errorf("Expected type 'github', got '%s'", provider.Type)
	}
	if provider.URL != "https://github.com" {
		t.Errorf("Expected URL 'https://github.com', got '%s'", provider.URL)
	}
	if provider.Token != "gh-token" {
		t.Errorf("Expected token 'gh-token', got '%s'", provider.Token)
	}
	if provider.Group != "my-org" {
		t.Errorf("Expected group 'my-org', got '%s'", provider.Group)
	}

	expectedBaseDir := filepath.Join(tempDir, "gitstuff-repos")
	if config.Local.BaseDir != expectedBaseDir {
		t.Errorf("Expected base dir '%s', got '%s'", expectedBaseDir, config.Local.BaseDir)
	}
}

func TestAddProvider_MultipleProviders(t *testing.T) {
	tempDir := t.TempDir()

	originalHome := os.Getenv("HOME")
	t.Cleanup(func() {
		os.Setenv("HOME", originalHome)
	})
	os.Setenv("HOME", tempDir)

	// Add first provider
	err := AddProvider("gitlab-main", "gitlab", "https://gitlab.com", "gl-token", "/shared/dir", false, "")
	if err != nil {
		t.Fatalf("First AddProvider failed: %v", err)
	}

	// Add second provider
	err = AddProvider("github-main", "github", "https://github.com", "gh-token", "", true, "my-org")
	if err != nil {
		t.Fatalf("Second AddProvider failed: %v", err)
	}

	config, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(config.Providers) != 2 {
		t.Fatalf("Expected 2 providers, got %d", len(config.Providers))
	}

	// Check first provider
	found := false
	for _, provider := range config.Providers {
		if provider.Name == "gitlab-main" {
			found = true
			if provider.Type != "gitlab" {
				t.Errorf("Expected type 'gitlab', got '%s'", provider.Type)
			}
			if provider.URL != "https://gitlab.com" {
				t.Errorf("Expected URL 'https://gitlab.com', got '%s'", provider.URL)
			}
			if provider.Insecure != false {
				t.Errorf("Expected insecure false, got %v", provider.Insecure)
			}
		}
	}
	if !found {
		t.Error("gitlab-main provider not found")
	}

	// Check second provider
	found = false
	for _, provider := range config.Providers {
		if provider.Name == "github-main" {
			found = true
			if provider.Type != "github" {
				t.Errorf("Expected type 'github', got '%s'", provider.Type)
			}
			if provider.URL != "https://github.com" {
				t.Errorf("Expected URL 'https://github.com', got '%s'", provider.URL)
			}
			if provider.Insecure != true {
				t.Errorf("Expected insecure true, got %v", provider.Insecure)
			}
			if provider.Group != "my-org" {
				t.Errorf("Expected group 'my-org', got '%s'", provider.Group)
			}
		}
	}
	if !found {
		t.Error("github-main provider not found")
	}

	// Base dir should remain from first provider
	if config.Local.BaseDir != "/shared/dir" {
		t.Errorf("Expected base dir '/shared/dir', got '%s'", config.Local.BaseDir)
	}
}

func TestAddProvider_ValidationErrors(t *testing.T) {
	tempDir := t.TempDir()

	originalHome := os.Getenv("HOME")
	t.Cleanup(func() {
		os.Setenv("HOME", originalHome)
	})
	os.Setenv("HOME", tempDir)

	tests := []struct {
		name         string
		providerName string
		providerType string
		url          string
		token        string
		wantErr      bool
		errContains  string
	}{
		{
			name:         "empty name",
			providerName: "",
			providerType: "gitlab",
			url:          "https://gitlab.com",
			token:        "token",
			wantErr:      true,
			errContains:  "provider name is required",
		},
		{
			name:         "empty type",
			providerName: "test",
			providerType: "",
			url:          "https://gitlab.com",
			token:        "token",
			wantErr:      true,
			errContains:  "provider type is required",
		},
		{
			name:         "invalid type",
			providerName: "test",
			providerType: "bitbucket",
			url:          "https://bitbucket.org",
			token:        "token",
			wantErr:      true,
			errContains:  "unsupported provider type",
		},
		{
			name:         "empty URL",
			providerName: "test",
			providerType: "gitlab",
			url:          "",
			token:        "token",
			wantErr:      true,
			errContains:  "provider URL is required",
		},
		{
			name:         "empty token",
			providerName: "test",
			providerType: "gitlab",
			url:          "https://gitlab.com",
			token:        "",
			wantErr:      true,
			errContains:  "provider token is required",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := AddProvider(tt.providerName, tt.providerType, tt.url, tt.token, "", false, "")
			if (err != nil) != tt.wantErr {
				t.Errorf("AddProvider() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && !strings.Contains(err.Error(), tt.errContains) {
				t.Errorf("AddProvider() error = %v, should contain '%s'", err, tt.errContains)
			}
		})
	}
}

func TestLoadWithLegacyMigration(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, ".gitstuff.yaml")

	originalHome := os.Getenv("HOME")
	t.Cleanup(func() {
		os.Setenv("HOME", originalHome)
	})
	os.Setenv("HOME", tempDir)

	// Create legacy config format
	legacyConfig := `
gitlab:
  url: https://gitlab.example.com
  token: legacy-token
  insecure: true
  group: legacy-group
local:
  basedir: /legacy/dir
`

	err := os.WriteFile(configPath, []byte(legacyConfig), 0600)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	config, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Should have migrated to new format
	if len(config.Providers) != 1 {
		t.Fatalf("Expected 1 provider after migration, got %d", len(config.Providers))
	}

	provider := config.Providers[0]
	if provider.Name != "gitlab" {
		t.Errorf("Expected name 'gitlab', got '%s'", provider.Name)
	}
	if provider.Type != "gitlab" {
		t.Errorf("Expected type 'gitlab', got '%s'", provider.Type)
	}
	if provider.URL != "https://gitlab.example.com" {
		t.Errorf("Expected URL 'https://gitlab.example.com', got '%s'", provider.URL)
	}
	if provider.Token != "legacy-token" {
		t.Errorf("Expected token 'legacy-token', got '%s'", provider.Token)
	}
	if provider.Insecure != true {
		t.Errorf("Expected insecure true, got %v", provider.Insecure)
	}
	if provider.Group != "legacy-group" {
		t.Errorf("Expected group 'legacy-group', got '%s'", provider.Group)
	}
	if config.Local.BaseDir != "/legacy/dir" {
		t.Errorf("Expected base dir '/legacy/dir', got '%s'", config.Local.BaseDir)
	}
}

func TestLoad_MultiProvider(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, ".gitstuff.yaml")

	originalHome := os.Getenv("HOME")
	t.Cleanup(func() {
		os.Setenv("HOME", originalHome)
	})
	os.Setenv("HOME", tempDir)

	testConfig := Config{
		Providers: []ProviderConfig{
			{
				Name:     "gitlab-main",
				Type:     "gitlab",
				URL:      "https://gitlab.com",
				Token:    "gl-token",
				Insecure: false,
				Group:    "my-group",
			},
			{
				Name:     "github-enterprise",
				Type:     "github",
				URL:      "https://github.enterprise.com",
				Token:    "gh-token",
				Insecure: true,
				Group:    "enterprise-org",
			},
		},
		Local: LocalConfig{
			BaseDir: "/multi/provider/dir",
		},
	}

	data, err := yaml.Marshal(&testConfig)
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}

	err = os.WriteFile(configPath, data, 0600)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	config, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	if len(config.Providers) != 2 {
		t.Fatalf("Expected 2 providers, got %d", len(config.Providers))
	}

	// Check all provider details
	gitlabFound := false
	githubFound := false
	for _, provider := range config.Providers {
		if provider.Name == "gitlab-main" {
			gitlabFound = true
			if provider.Type != "gitlab" {
				t.Errorf("Expected type 'gitlab', got '%s'", provider.Type)
			}
			if provider.URL != "https://gitlab.com" {
				t.Errorf("Expected URL 'https://gitlab.com', got '%s'", provider.URL)
			}
			if provider.Token != "gl-token" {
				t.Errorf("Expected token 'gl-token', got '%s'", provider.Token)
			}
			if provider.Insecure != false {
				t.Errorf("Expected insecure false, got %v", provider.Insecure)
			}
			if provider.Group != "my-group" {
				t.Errorf("Expected group 'my-group', got '%s'", provider.Group)
			}
		}
		if provider.Name == "github-enterprise" {
			githubFound = true
			if provider.Type != "github" {
				t.Errorf("Expected type 'github', got '%s'", provider.Type)
			}
			if provider.URL != "https://github.enterprise.com" {
				t.Errorf("Expected URL 'https://github.enterprise.com', got '%s'", provider.URL)
			}
			if provider.Token != "gh-token" {
				t.Errorf("Expected token 'gh-token', got '%s'", provider.Token)
			}
			if provider.Insecure != true {
				t.Errorf("Expected insecure true, got %v", provider.Insecure)
			}
			if provider.Group != "enterprise-org" {
				t.Errorf("Expected group 'enterprise-org', got '%s'", provider.Group)
			}
		}
	}

	if !gitlabFound {
		t.Error("gitlab-main provider not found")
	}
	if !githubFound {
		t.Error("github-enterprise provider not found")
	}

	if config.Local.BaseDir != "/multi/provider/dir" {
		t.Errorf("Expected base dir '/multi/provider/dir', got '%s'", config.Local.BaseDir)
	}
}
