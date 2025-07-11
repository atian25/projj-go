package cache

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/atian25/projj-go/internal/config"
)

// Repository 表示一个仓库的信息
type Repository struct {
	Name     string    `json:"name"`
	URL      string    `json:"url"`
	Path     string    `json:"path"`
	Platform string    `json:"platform"`
	AddedAt  time.Time `json:"added_at"`
}

// Cache 表示缓存结构
type Cache struct {
	Repositories []Repository `json:"repositories"`
	UpdatedAt    time.Time    `json:"updated_at"`
}

// GetCachePath 获取缓存文件路径
func GetCachePath() string {
	return filepath.Join(config.GetConfigDir(), "cache.json")
}

// Load 加载缓存文件
func Load() (*Cache, error) {
	cachePath := GetCachePath()
	
	// 如果缓存文件不存在，返回空缓存
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		return &Cache{
			Repositories: make([]Repository, 0),
			UpdatedAt:    time.Now(),
		}, nil
	}
	
	data, err := os.ReadFile(cachePath)
	if err != nil {
		return nil, fmt.Errorf("读取缓存文件失败: %w", err)
	}
	
	// 尝试解析为新格式
	var cache Cache
	if err := json.Unmarshal(data, &cache); err == nil && len(cache.Repositories) > 0 {
		return &cache, nil
	}
	
	// 尝试解析为旧格式（projj 原始格式）
	var oldFormat map[string]interface{}
	if err := json.Unmarshal(data, &oldFormat); err != nil {
		return nil, fmt.Errorf("解析缓存文件失败: %w", err)
	}
	
	// 转换旧格式到新格式
	cache = Cache{
		Repositories: make([]Repository, 0),
		UpdatedAt:    time.Now(),
	}
	
	for path, repoData := range oldFormat {
		// 跳过 version 字段
		if path == "version" {
			continue
		}
		
		if repoMap, ok := repoData.(map[string]interface{}); ok {
			if repoURL, exists := repoMap["repo"]; exists {
				if urlStr, ok := repoURL.(string); ok {
					repo := Repository{
						Path:     path,
						URL:      urlStr,
						Name:     extractRepoName(path),
						Platform: extractPlatform(urlStr),
						AddedAt:  time.Now(),
					}
					cache.Repositories = append(cache.Repositories, repo)
				}
			}
		}
	}
	
	return &cache, nil
}

// Save 保存缓存文件
func (c *Cache) Save() error {
	configDir := config.GetConfigDir()
	
	// 确保配置目录存在
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %w", err)
	}
	
	c.UpdatedAt = time.Now()
	cachePath := GetCachePath()
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化缓存失败: %w", err)
	}
	
	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		return fmt.Errorf("写入缓存文件失败: %w", err)
	}
	
	return nil
}

// Add 添加仓库到缓存
func (c *Cache) Add(repo Repository) {
	// 检查是否已存在
	for i, existing := range c.Repositories {
		if existing.Path == repo.Path {
			// 更新现有记录
			c.Repositories[i] = repo
			return
		}
	}
	
	// 添加新记录
	repo.AddedAt = time.Now()
	c.Repositories = append(c.Repositories, repo)
}

// Remove 从缓存中移除仓库
func (c *Cache) Remove(path string) bool {
	for i, repo := range c.Repositories {
		if repo.Path == path {
			// 移除元素
			c.Repositories = append(c.Repositories[:i], c.Repositories[i+1:]...)
			return true
		}
	}
	return false
}

// Find 查找仓库
func (c *Cache) Find(query string) []Repository {
	if query == "" {
		return c.Repositories
	}
	
	var results []Repository
	query = strings.ToLower(query)
	
	for _, repo := range c.Repositories {
		// 支持多种匹配方式
		if strings.Contains(strings.ToLower(repo.Name), query) ||
			strings.Contains(strings.ToLower(repo.Path), query) ||
			strings.Contains(strings.ToLower(repo.URL), query) {
			results = append(results, repo)
		}
	}
	
	return results
}

// GetByPath 根据路径获取仓库
func (c *Cache) GetByPath(path string) *Repository {
	for _, repo := range c.Repositories {
		if repo.Path == path {
			return &repo
		}
	}
	return nil
}

// extractRepoName 从路径中提取仓库名称
func extractRepoName(path string) string {
	parts := strings.Split(path, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return ""
}

// extractPlatform 从 URL 中提取平台信息
func extractPlatform(url string) string {
	if strings.Contains(url, "github.com") {
		return "github"
	} else if strings.Contains(url, "code.byted.org") {
		return "byted"
	} else if strings.Contains(url, "git.byted.org") {
		return "byted"
	}
	return "unknown"
}

// Sync 同步缓存与文件系统
func (c *Cache) Sync() (int, int, error) {
	var removed, added int
	
	// 检查缓存中的仓库是否还存在
	var validRepos []Repository
	for _, repo := range c.Repositories {
		if _, err := os.Stat(filepath.Join(repo.Path, ".git")); err == nil {
			validRepos = append(validRepos, repo)
		} else {
			removed++
		}
	}
	c.Repositories = validRepos
	
	// TODO: 扫描文件系统中的新仓库
	// 这里可以实现扫描基础目录下的所有 .git 目录
	
	return added, removed, nil
}