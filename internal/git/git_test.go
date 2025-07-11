package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func TestParseURL(t *testing.T) {
	aliases := map[string]string{
		"github://": "git@github.com:",
		"gitlab://": "git@gitlab.com:",
		"gitee://":  "git@gitee.com:",
	}
	
	tests := []struct {
		name     string
		input    string
		expected *RepoInfo
		shouldErr bool
	}{
		{
			name:  "SSH format",
			input: "git@github.com:user/repo.git",
			expected: &RepoInfo{
				Platform: "github.com",
				Owner:    "user",
				Name:     "repo",
				URL:      "git@github.com:user/repo.git",
			},
			shouldErr: false,
		},
		{
			name:  "SSH format without .git",
			input: "git@github.com:user/repo",
			expected: &RepoInfo{
				Platform: "github.com",
				Owner:    "user",
				Name:     "repo",
				URL:      "git@github.com:user/repo.git",
			},
			shouldErr: false,
		},
		{
			name:  "HTTPS format",
			input: "https://github.com/user/repo.git",
			expected: &RepoInfo{
				Platform: "github.com",
				Owner:    "user",
				Name:     "repo",
				URL:      "https://github.com/user/repo.git",
			},
			shouldErr: false,
		},
		{
			name:  "HTTPS format without .git",
			input: "https://github.com/user/repo",
			expected: &RepoInfo{
				Platform: "github.com",
				Owner:    "user",
				Name:     "repo",
				URL:      "https://github.com/user/repo",
			},
			shouldErr: false,
		},
		{
			name:  "Platform format",
			input: "github.com/user/repo",
			expected: &RepoInfo{
				Platform: "github.com",
				Owner:    "user",
				Name:     "repo",
				URL:      "https://github.com/user/repo.git",
			},
			shouldErr: false,
		},
		{
			name:  "Short format",
			input: "user/repo",
			expected: &RepoInfo{
				Platform: "github.com",
				Owner:    "user",
				Name:     "repo",
				URL:      "https://github.com/user/repo.git",
			},
			shouldErr: false,
		},
		{
			name:  "Alias format - github",
			input: "github://user/repo",
			expected: &RepoInfo{
				Platform: "github.com",
				Owner:    "user",
				Name:     "repo",
				URL:      "git@github.com:user/repo.git",
			},
			shouldErr: false,
		},
		{
			name:  "Alias format - gitlab",
			input: "gitlab://user/repo",
			expected: &RepoInfo{
				Platform: "gitlab.com",
				Owner:    "user",
				Name:     "repo",
				URL:      "git@gitlab.com:user/repo.git",
			},
			shouldErr: false,
		},
		{
			name:      "Invalid format",
			input:     "invalid-url",
			expected:  nil,
			shouldErr: true,
		},
		{
			name:      "Empty input",
			input:     "",
			expected:  nil,
			shouldErr: true,
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseURL(tt.input, aliases)
			
			if tt.shouldErr {
				if err == nil {
					t.Errorf("Expected error for input '%s', but got none", tt.input)
				}
				return
			}
			
			if err != nil {
				t.Errorf("Unexpected error for input '%s': %v", tt.input, err)
				return
			}
			
			if result.Platform != tt.expected.Platform {
				t.Errorf("Platform mismatch for '%s': expected '%s', got '%s'", tt.input, tt.expected.Platform, result.Platform)
			}
			
			if result.Owner != tt.expected.Owner {
				t.Errorf("Owner mismatch for '%s': expected '%s', got '%s'", tt.input, tt.expected.Owner, result.Owner)
			}
			
			if result.Name != tt.expected.Name {
				t.Errorf("Name mismatch for '%s': expected '%s', got '%s'", tt.input, tt.expected.Name, result.Name)
			}
			
			if result.URL != tt.expected.URL {
				t.Errorf("URL mismatch for '%s': expected '%s', got '%s'", tt.input, tt.expected.URL, result.URL)
			}
		})
	}
}

func TestRepoInfoGetRepoPath(t *testing.T) {
	repo := &RepoInfo{
		Platform: "github.com",
		Owner:    "user",
		Name:     "repo",
		URL:      "https://github.com/user/repo.git",
	}
	
	basePath := "/home/user/projects"
	expected := filepath.Join(basePath, "github.com", "user", "repo")
	result := repo.GetRepoPath(basePath)
	
	if result != expected {
		t.Errorf("Expected path '%s', got '%s'", expected, result)
	}
}

func TestIsGitRepository(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "git-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// 测试非 Git 目录
	if IsGitRepository(tempDir) {
		t.Error("Empty directory should not be recognized as Git repository")
	}
	
	// 创建 .git 目录
	gitDir := filepath.Join(tempDir, ".git")
	err = os.Mkdir(gitDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create .git dir: %v", err)
	}
	
	// 测试 Git 目录
	if !IsGitRepository(tempDir) {
		t.Error("Directory with .git should be recognized as Git repository")
	}
	
	// 测试不存在的目录
	if IsGitRepository("/nonexistent/path") {
		t.Error("Nonexistent directory should not be recognized as Git repository")
	}
}

func TestClone(t *testing.T) {
	// 跳过需要网络连接的测试
	if testing.Short() {
		t.Skip("Skipping clone test in short mode")
	}
	
	// 检查 git 命令是否可用
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("Git command not available")
	}
	
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "clone-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	targetPath := filepath.Join(tempDir, "test-repo")
	
	// 测试克隆一个小的公开仓库
	repoURL := "https://github.com/octocat/Hello-World.git"
	err = Clone(repoURL, targetPath)
	if err != nil {
		t.Fatalf("Failed to clone repository: %v", err)
	}
	
	// 验证克隆结果
	if !IsGitRepository(targetPath) {
		t.Error("Cloned directory should be a Git repository")
	}
	
	// 测试克隆到已存在的目录
	err = Clone(repoURL, targetPath)
	if err == nil {
		t.Error("Should fail when cloning to existing directory")
	}
}

func TestGetRemoteURL(t *testing.T) {
	// 跳过需要 Git 仓库的测试
	if testing.Short() {
		t.Skip("Skipping remote URL test in short mode")
	}
	
	// 检查 git 命令是否可用
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("Git command not available")
	}
	
	// 创建临时 Git 仓库
	tempDir, err := os.MkdirTemp("", "git-remote-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// 初始化 Git 仓库
	cmd := exec.Command("git", "init")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}
	
	// 添加远程 URL
	remoteURL := "https://github.com/user/repo.git"
	cmd = exec.Command("git", "remote", "add", "origin", remoteURL)
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to add remote: %v", err)
	}
	
	// 测试获取远程 URL
	result, err := GetRemoteURL(tempDir)
	if err != nil {
		t.Fatalf("Failed to get remote URL: %v", err)
	}
	
	if result != remoteURL {
		t.Errorf("Expected remote URL '%s', got '%s'", remoteURL, result)
	}
	
	// 测试非 Git 目录
	nonGitDir, err := os.MkdirTemp("", "non-git-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(nonGitDir)
	
	_, err = GetRemoteURL(nonGitDir)
	if err == nil {
		t.Error("Should fail for non-Git directory")
	}
}

func TestGetStatus(t *testing.T) {
	// 跳过需要 Git 仓库的测试
	if testing.Short() {
		t.Skip("Skipping status test in short mode")
	}
	
	// 检查 git 命令是否可用
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("Git command not available")
	}
	
	// 创建临时 Git 仓库
	tempDir, err := os.MkdirTemp("", "git-status-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// 初始化 Git 仓库
	cmd := exec.Command("git", "init")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to init git repo: %v", err)
	}
	
	// 配置用户信息
	cmd = exec.Command("git", "config", "user.email", "test@example.com")
	cmd.Dir = tempDir
	cmd.Run()
	
	cmd = exec.Command("git", "config", "user.name", "Test User")
	cmd.Dir = tempDir
	cmd.Run()
	
	// 测试空仓库状态（应该是干净的）
	isClean, err := GetStatus(tempDir)
	if err != nil {
		t.Fatalf("Failed to get status: %v", err)
	}
	
	if !isClean {
		t.Error("Empty repository should be clean")
	}
	
	// 创建文件
	testFile := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}
	
	// 测试有未跟踪文件的状态
	isClean, err = GetStatus(tempDir)
	if err != nil {
		t.Fatalf("Failed to get status: %v", err)
	}
	
	if isClean {
		t.Error("Repository with untracked files should not be clean")
	}
	
	// 添加并提交文件
	cmd = exec.Command("git", "add", "test.txt")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to add file: %v", err)
	}
	
	cmd = exec.Command("git", "commit", "-m", "Initial commit")
	cmd.Dir = tempDir
	if err := cmd.Run(); err != nil {
		t.Fatalf("Failed to commit: %v", err)
	}
	
	// 测试提交后的状态（应该是干净的）
	isClean, err = GetStatus(tempDir)
	if err != nil {
		t.Fatalf("Failed to get status: %v", err)
	}
	
	if !isClean {
		t.Error("Repository after commit should be clean")
	}
}