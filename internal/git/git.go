package git

import (
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// RepoInfo 表示解析后的仓库信息
type RepoInfo struct {
	Platform string // github.com, gitlab.com, etc.
	Owner    string // 用户名或组织名
	Name     string // 仓库名
	URL      string // 完整的 Git URL
}

// ParseURL 解析各种格式的 Git URL
func ParseURL(input string, aliases map[string]string) (*RepoInfo, error) {
	// 处理别名
	for alias, replacement := range aliases {
		if strings.HasPrefix(input, alias) {
			input = strings.Replace(input, alias, replacement, 1)
			break
		}
	}
	
	// SSH 格式: git@github.com:user/repo.git
	if sshRegex := regexp.MustCompile(`^git@([^:]+):([^/]+)/([^/]+?)(?:\.git)?$`); sshRegex.MatchString(input) {
		matches := sshRegex.FindStringSubmatch(input)
		return &RepoInfo{
			Platform: matches[1],
			Owner:    matches[2],
			Name:     matches[3],
			URL:      fmt.Sprintf("git@%s:%s/%s.git", matches[1], matches[2], matches[3]),
		}, nil
	}
	
	// HTTPS 格式: https://github.com/user/repo.git
	if httpsURL, err := url.Parse(input); err == nil && httpsURL.Scheme != "" {
		pathParts := strings.Split(strings.Trim(httpsURL.Path, "/"), "/")
		if len(pathParts) >= 2 {
			repoName := strings.TrimSuffix(pathParts[1], ".git")
			return &RepoInfo{
				Platform: httpsURL.Host,
				Owner:    pathParts[0],
				Name:     repoName,
				URL:      input,
			}, nil
		}
	}
	
	// 平台格式: github.com/user/repo
	if platformRegex := regexp.MustCompile(`^([^/]+\.[^/]+)/([^/]+)/([^/]+?)(?:\.git)?$`); platformRegex.MatchString(input) {
		matches := platformRegex.FindStringSubmatch(input)
		return &RepoInfo{
			Platform: matches[1],
			Owner:    matches[2],
			Name:     matches[3],
			URL:      fmt.Sprintf("https://%s/%s/%s.git", matches[1], matches[2], matches[3]),
		}, nil
	}
	
	// 简短格式: user/repo (默认 GitHub)
	if shortRegex := regexp.MustCompile(`^([^/]+)/([^/]+)$`); shortRegex.MatchString(input) {
		matches := shortRegex.FindStringSubmatch(input)
		return &RepoInfo{
			Platform: "github.com",
			Owner:    matches[1],
			Name:     matches[2],
			URL:      fmt.Sprintf("https://github.com/%s/%s.git", matches[1], matches[2]),
		}, nil
	}
	
	return nil, fmt.Errorf("无法解析 Git URL: %s", input)
}

// GetRepoPath 根据仓库信息生成本地路径
func (r *RepoInfo) GetRepoPath(basePath string) string {
	return filepath.Join(basePath, r.Platform, r.Owner, r.Name)
}

// Clone 克隆仓库到指定路径
func Clone(repoURL, targetPath string) error {
	// 确保目标目录的父目录存在
	parentDir := filepath.Dir(targetPath)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return fmt.Errorf("创建目录失败: %w", err)
	}
	
	// 检查目标目录是否已存在
	if _, err := os.Stat(targetPath); err == nil {
		return fmt.Errorf("目标目录已存在: %s", targetPath)
	}
	
	// 执行 git clone
	cmd := exec.Command("git", "clone", repoURL, targetPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("克隆仓库失败: %w", err)
	}
	
	return nil
}

// IsGitRepository 检查目录是否是 Git 仓库
func IsGitRepository(path string) bool {
	gitDir := filepath.Join(path, ".git")
	if stat, err := os.Stat(gitDir); err == nil {
		return stat.IsDir() || stat.Mode().IsRegular() // 支持 .git 文件（submodule）
	}
	return false
}

// GetRemoteURL 获取仓库的远程 URL
func GetRemoteURL(repoPath string) (string, error) {
	cmd := exec.Command("git", "config", "--get", "remote.origin.url")
	cmd.Dir = repoPath
	
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("获取远程 URL 失败: %w", err)
	}
	
	return strings.TrimSpace(string(output)), nil
}

// Pull 拉取仓库更新
func Pull(repoPath string) error {
	cmd := exec.Command("git", "pull")
	cmd.Dir = repoPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("拉取更新失败: %w", err)
	}
	
	return nil
}

// GetStatus 获取仓库状态
func GetStatus(repoPath string) (bool, error) {
	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = repoPath
	
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("获取仓库状态失败: %w", err)
	}
	
	// 如果输出为空，说明工作区是干净的
	return len(strings.TrimSpace(string(output))) == 0, nil
}