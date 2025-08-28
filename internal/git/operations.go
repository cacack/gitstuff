package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Status struct {
	Exists        bool
	CurrentBranch string
	IsGitRepo     bool
	HasChanges    bool
}

func GetRepositoryStatus(repoPath string) (*Status, error) {
	status := &Status{}
	
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		status.Exists = false
		return status, nil
	}
	
	status.Exists = true
	
	gitDir := filepath.Join(repoPath, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		status.IsGitRepo = false
		return status, nil
	}
	
	status.IsGitRepo = true
	
	cmd := exec.Command("git", "-C", repoPath, "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to get current branch: %w", err)
	}
	
	status.CurrentBranch = strings.TrimSpace(string(output))
	
	cmd = exec.Command("git", "-C", repoPath, "status", "--porcelain")
	output, err = cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to check git status: %w", err)
	}
	
	status.HasChanges = len(strings.TrimSpace(string(output))) > 0
	
	return status, nil
}

func CloneRepository(cloneURL, targetPath string, useSSH bool) error {
	if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
		return fmt.Errorf("failed to create target directory: %w", err)
	}
	
	var cmd *exec.Cmd
	if useSSH {
		cmd = exec.Command("git", "clone", cloneURL, targetPath)
	} else {
		cmd = exec.Command("git", "clone", cloneURL, targetPath)
	}
	
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}
	
	return nil
}

func PullRepository(repoPath string) error {
	cmd := exec.Command("git", "-C", repoPath, "pull")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to pull repository: %w", err)
	}
	
	return nil
}