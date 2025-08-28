package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"gopkg.in/yaml.v3"
)

func TestCreate(t *testing.T) {
	tempDir := t.TempDir()
	
	originalHome := os.Getenv("HOME")
	t.Cleanup(func() {
		os.Setenv("HOME", originalHome)
	})
	os.Setenv("HOME", tempDir)
	
	gitlabURL := "https://gitlab.example.com"
	token := "test-token"
	baseDir := "/custom/base/dir"
	
	err := Create(gitlabURL, token, baseDir, false)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	
	configPath := filepath.Join(tempDir, ".gitstuff.yaml")
	
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Fatal("Config file was not created")
	}
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}
	
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}
	
	if config.GitLab.URL != gitlabURL {
		t.Errorf("Expected URL %s, got %s", gitlabURL, config.GitLab.URL)
	}
	
	if config.GitLab.Token != token {
		t.Errorf("Expected token %s, got %s", token, config.GitLab.Token)
	}
	
	if config.Local.BaseDir != baseDir {
		t.Errorf("Expected base dir %s, got %s", baseDir, config.Local.BaseDir)
	}
	
	if config.GitLab.Insecure != false {
		t.Errorf("Expected insecure to be false, got %v", config.GitLab.Insecure)
	}
	
	info, err := os.Stat(configPath)
	if err != nil {
		t.Fatalf("Failed to stat config file: %v", err)
	}
	
	if info.Mode().Perm() != 0600 {
		t.Errorf("Expected file permissions 0600, got %o", info.Mode().Perm())
	}
}

func TestCreateWithDefaultBaseDir(t *testing.T) {
	tempDir := t.TempDir()
	
	originalHome := os.Getenv("HOME")
	t.Cleanup(func() {
		os.Setenv("HOME", originalHome)
	})
	os.Setenv("HOME", tempDir)
	
	gitlabURL := "https://gitlab.example.com"
	token := "test-token"
	
	err := Create(gitlabURL, token, "", true)
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}
	
	configPath := filepath.Join(tempDir, ".gitstuff.yaml")
	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("Failed to read config file: %v", err)
	}
	
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}
	
	expectedBaseDir := filepath.Join(tempDir, "gitlab-repos")
	if config.Local.BaseDir != expectedBaseDir {
		t.Errorf("Expected base dir %s, got %s", expectedBaseDir, config.Local.BaseDir)
	}
	
	if config.GitLab.Insecure != true {
		t.Errorf("Expected insecure to be true, got %v", config.GitLab.Insecure)
	}
}

func TestLoad(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, ".gitstuff.yaml")
	
	originalHome := os.Getenv("HOME")
	t.Cleanup(func() {
		os.Setenv("HOME", originalHome)
	})
	os.Setenv("HOME", tempDir)
	
	testConfig := Config{
		GitLab: GitLabConfig{
			URL:   "https://gitlab.example.com",
			Token: "test-token",
		},
		Local: LocalConfig{
			BaseDir: "",
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
	
	viper.Reset()
	viper.SetConfigFile(configPath)
	err = viper.ReadInConfig()
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}
	
	config, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	
	if config.GitLab.URL != testConfig.GitLab.URL {
		t.Errorf("Expected URL %s, got %s", testConfig.GitLab.URL, config.GitLab.URL)
	}
	
	if config.GitLab.Token != testConfig.GitLab.Token {
		t.Errorf("Expected token %s, got %s", testConfig.GitLab.Token, config.GitLab.Token)
	}
	
	expectedBaseDir := filepath.Join(tempDir, "gitlab-repos")
	if config.Local.BaseDir != expectedBaseDir {
		t.Errorf("Expected base dir %s, got %s", expectedBaseDir, config.Local.BaseDir)
	}
}

func TestLoadWithMissingURL(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, ".gitstuff.yaml")
	
	originalHome := os.Getenv("HOME")
	t.Cleanup(func() {
		os.Setenv("HOME", originalHome)
	})
	os.Setenv("HOME", tempDir)
	
	testConfig := Config{
		GitLab: GitLabConfig{
			Token: "test-token",
		},
		Local: LocalConfig{
			BaseDir: "/custom/base/dir",
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
	
	viper.Reset()
	viper.SetConfigFile(configPath)
	err = viper.ReadInConfig()
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}
	
	_, err = Load()
	if err == nil {
		t.Fatal("Expected error when URL is missing")
	}
}

func TestLoadWithMissingToken(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, ".gitstuff.yaml")
	
	originalHome := os.Getenv("HOME")
	t.Cleanup(func() {
		os.Setenv("HOME", originalHome)
	})
	os.Setenv("HOME", tempDir)
	
	testConfig := Config{
		GitLab: GitLabConfig{
			URL: "https://gitlab.example.com",
		},
		Local: LocalConfig{
			BaseDir: "/custom/base/dir",
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
	
	viper.Reset()
	viper.SetConfigFile(configPath)
	err = viper.ReadInConfig()
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}
	
	_, err = Load()
	if err == nil {
		t.Fatal("Expected error when token is missing")
	}
}

func TestLoadWithInsecureFlag(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, ".gitstuff.yaml")
	
	originalHome := os.Getenv("HOME")
	t.Cleanup(func() {
		os.Setenv("HOME", originalHome)
	})
	os.Setenv("HOME", tempDir)
	
	testConfig := Config{
		GitLab: GitLabConfig{
			URL:      "https://gitlab.example.com",
			Token:    "test-token",
			Insecure: true,
		},
		Local: LocalConfig{
			BaseDir: "",
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
	
	viper.Reset()
	viper.SetConfigFile(configPath)
	err = viper.ReadInConfig()
	if err != nil {
		t.Fatalf("Failed to read config: %v", err)
	}
	
	config, err := Load()
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	
	if config.GitLab.URL != testConfig.GitLab.URL {
		t.Errorf("Expected URL %s, got %s", testConfig.GitLab.URL, config.GitLab.URL)
	}
	
	if config.GitLab.Token != testConfig.GitLab.Token {
		t.Errorf("Expected token %s, got %s", testConfig.GitLab.Token, config.GitLab.Token)
	}
	
	if config.GitLab.Insecure != true {
		t.Errorf("Expected insecure to be true, got %v", config.GitLab.Insecure)
	}
	
	expectedBaseDir := filepath.Join(tempDir, "gitlab-repos")
	if config.Local.BaseDir != expectedBaseDir {
		t.Errorf("Expected base dir %s, got %s", expectedBaseDir, config.Local.BaseDir)
	}
}