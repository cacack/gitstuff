package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestGetRepositoryStatus_NonExistent(t *testing.T) {
	tempDir := t.TempDir()
	nonExistentPath := filepath.Join(tempDir, "nonexistent")
	
	status, err := GetRepositoryStatus(nonExistentPath)
	if err != nil {
		t.Fatalf("GetRepositoryStatus failed: %v", err)
	}
	
	if status.Exists {
		t.Error("Expected repository to not exist")
	}
}

func TestGetRepositoryStatus_ExistsButNotGit(t *testing.T) {
	tempDir := t.TempDir()
	repoDir := filepath.Join(tempDir, "notgit")
	
	err := os.MkdirAll(repoDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}
	
	status, err := GetRepositoryStatus(repoDir)
	if err != nil {
		t.Fatalf("GetRepositoryStatus failed: %v", err)
	}
	
	if !status.Exists {
		t.Error("Expected directory to exist")
	}
	
	if status.IsGitRepo {
		t.Error("Expected directory to not be a git repository")
	}
}

func TestGetRepositoryStatus_ValidGitRepo(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available in PATH")
	}
	
	tempDir := t.TempDir()
	repoDir := filepath.Join(tempDir, "testrepo")
	
	err := os.MkdirAll(repoDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}
	
	cmd := exec.Command("git", "-C", repoDir, "init")
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}
	
	cmd = exec.Command("git", "-C", repoDir, "config", "user.name", "Test User")
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Failed to set git user name: %v", err)
	}
	
	cmd = exec.Command("git", "-C", repoDir, "config", "user.email", "test@example.com")
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Failed to set git user email: %v", err)
	}
	
	testFile := filepath.Join(repoDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	
	cmd = exec.Command("git", "-C", repoDir, "add", "test.txt")
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Failed to add file to git: %v", err)
	}
	
	cmd = exec.Command("git", "-C", repoDir, "commit", "-m", "Initial commit")
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Failed to commit: %v", err)
	}
	
	status, err := GetRepositoryStatus(repoDir)
	if err != nil {
		t.Fatalf("GetRepositoryStatus failed: %v", err)
	}
	
	if !status.Exists {
		t.Error("Expected repository to exist")
	}
	
	if !status.IsGitRepo {
		t.Error("Expected directory to be a git repository")
	}
	
	if status.CurrentBranch == "" {
		t.Error("Expected current branch to be set")
	}
	
	if status.HasChanges {
		t.Error("Expected no uncommitted changes")
	}
}

func TestGetRepositoryStatus_WithChanges(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available in PATH")
	}
	
	tempDir := t.TempDir()
	repoDir := filepath.Join(tempDir, "testrepo")
	
	err := os.MkdirAll(repoDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}
	
	cmd := exec.Command("git", "-C", repoDir, "init")
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}
	
	cmd = exec.Command("git", "-C", repoDir, "config", "user.name", "Test User")
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Failed to set git user name: %v", err)
	}
	
	cmd = exec.Command("git", "-C", repoDir, "config", "user.email", "test@example.com")
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Failed to set git user email: %v", err)
	}
	
	testFile := filepath.Join(repoDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	
	cmd = exec.Command("git", "-C", repoDir, "add", "test.txt")
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Failed to add file to git: %v", err)
	}
	
	cmd = exec.Command("git", "-C", repoDir, "commit", "-m", "Initial commit")
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Failed to commit: %v", err)
	}
	
	err = os.WriteFile(testFile, []byte("modified content"), 0644)
	if err != nil {
		t.Fatalf("Failed to modify test file: %v", err)
	}
	
	status, err := GetRepositoryStatus(repoDir)
	if err != nil {
		t.Fatalf("GetRepositoryStatus failed: %v", err)
	}
	
	if !status.HasChanges {
		t.Error("Expected uncommitted changes")
	}
}

func TestCloneRepository(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available in PATH")
	}

	tempDir := t.TempDir()
	
	sourceRepo := filepath.Join(tempDir, "source")
	targetRepo := filepath.Join(tempDir, "target")
	
	err := os.MkdirAll(sourceRepo, 0755)
	if err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}
	
	cmd := exec.Command("git", "-C", sourceRepo, "init", "--bare")
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Failed to init bare git repo: %v", err)
	}
	
	err = CloneRepository(sourceRepo, targetRepo, false)
	if err != nil {
		t.Fatalf("Failed to clone repository: %v", err)
	}
	
	if _, err := os.Stat(filepath.Join(targetRepo, ".git")); os.IsNotExist(err) {
		t.Error("Expected .git directory in cloned repository")
	}
}

func TestCloneRepository_InvalidURL(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available in PATH")
	}

	tempDir := t.TempDir()
	targetRepo := filepath.Join(tempDir, "target")
	
	err := CloneRepository("https://invalid.nonexistent.url/repo.git", targetRepo, false)
	if err == nil {
		t.Error("Expected error when cloning from invalid URL")
	}
}

func TestCloneRepository_CreateTargetDirectory(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available in PATH")
	}

	tempDir := t.TempDir()
	
	sourceRepo := filepath.Join(tempDir, "source")
	targetRepo := filepath.Join(tempDir, "nested", "deep", "target")
	
	err := os.MkdirAll(sourceRepo, 0755)
	if err != nil {
		t.Fatalf("Failed to create source directory: %v", err)
	}
	
	cmd := exec.Command("git", "-C", sourceRepo, "init", "--bare")
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Failed to init bare git repo: %v", err)
	}
	
	err = CloneRepository(sourceRepo, targetRepo, false)
	if err != nil {
		t.Fatalf("Failed to clone repository: %v", err)
	}
	
	if _, err := os.Stat(filepath.Join(targetRepo, ".git")); os.IsNotExist(err) {
		t.Error("Expected .git directory in cloned repository")
	}
	
	if _, err := os.Stat(filepath.Dir(targetRepo)); os.IsNotExist(err) {
		t.Error("Expected parent directories to be created")
	}
}

func TestPullRepository(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git not available in PATH")
	}

	tempDir := t.TempDir()
	
	bareRepo := filepath.Join(tempDir, "bare.git")
	workingRepo := filepath.Join(tempDir, "working")
	
	cmd := exec.Command("git", "-C", tempDir, "init", "--bare", bareRepo)
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Failed to init bare git repo: %v", err)
	}
	
	cmd = exec.Command("git", "clone", bareRepo, workingRepo)
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Failed to clone repo: %v", err)
	}
	
	cmd = exec.Command("git", "-C", workingRepo, "config", "user.name", "Test User")
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Failed to set git user name: %v", err)
	}
	
	cmd = exec.Command("git", "-C", workingRepo, "config", "user.email", "test@example.com")
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Failed to set git user email: %v", err)
	}
	
	testFile := filepath.Join(workingRepo, "test.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}
	
	cmd = exec.Command("git", "-C", workingRepo, "add", "test.txt")
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Failed to add file to git: %v", err)
	}
	
	cmd = exec.Command("git", "-C", workingRepo, "commit", "-m", "Initial commit")
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Failed to commit: %v", err)
	}
	
	cmd = exec.Command("git", "-C", workingRepo, "push")
	err = cmd.Run()
	if err != nil {
		t.Fatalf("Failed to push: %v", err)
	}
	
	err = PullRepository(workingRepo)
	if err != nil {
		t.Fatalf("Failed to pull repository: %v", err)
	}
}

func TestPullRepository_NonGitDirectory(t *testing.T) {
	tempDir := t.TempDir()
	nonGitDir := filepath.Join(tempDir, "notgit")
	
	err := os.MkdirAll(nonGitDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create directory: %v", err)
	}
	
	err = PullRepository(nonGitDir)
	if err == nil {
		t.Error("Expected error when pulling from non-git directory")
	}
}

func TestPullRepository_NonExistentDirectory(t *testing.T) {
	tempDir := t.TempDir()
	nonExistentDir := filepath.Join(tempDir, "nonexistent")
	
	err := PullRepository(nonExistentDir)
	if err == nil {
		t.Error("Expected error when pulling from non-existent directory")
	}
}