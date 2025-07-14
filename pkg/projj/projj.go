package projj

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/atian25/projj-go/internal/cache"
	"github.com/atian25/projj-go/internal/config"
	"github.com/atian25/projj-go/internal/git"
)

// Client 表示 projj 客户端
type Client struct {
	config *config.Config
	cache  *cache.Cache
}

// New 创建新的 projj 客户端
func New() (*Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, fmt.Errorf("加载配置失败: %w", err)
	}
	
	cch, err := cache.Load()
	if err != nil {
		return nil, fmt.Errorf("加载缓存失败: %w", err)
	}
	
	return &Client{
		config: cfg,
		cache:  cch,
	}, nil
}

// Init 初始化 projj 环境
func (c *Client) Init() error {
	// 创建配置目录
	configDir := config.GetConfigDir()
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %w", err)
	}
	
	// 创建基础目录
	basePath := c.config.GetBasePath()
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return fmt.Errorf("创建基础目录失败: %w", err)
	}
	
	// 保存默认配置
	if err := c.config.Save(); err != nil {
		return fmt.Errorf("保存配置失败: %w", err)
	}
	
	// 保存空缓存
	if err := c.cache.Save(); err != nil {
		return fmt.Errorf("保存缓存失败: %w", err)
	}
	
	fmt.Printf("Projj 初始化完成\n")
	fmt.Printf("配置目录: %s\n", configDir)
	fmt.Printf("基础目录: %s\n", basePath)
	
	return nil
}

// Add 添加仓库
func (c *Client) Add(repoURL string) error {
	// 解析仓库 URL
	repoInfo, err := git.ParseURL(repoURL, c.config.Alias)
	if err != nil {
		return fmt.Errorf("解析仓库 URL 失败: %w", err)
	}
	
	// 生成目标路径
	targetPath := repoInfo.GetRepoPath(c.config.GetBasePath())
	
	// 检查是否已存在
	if existingRepo := c.cache.GetByPath(targetPath); existingRepo != nil {
		return fmt.Errorf("仓库已存在: %s", targetPath)
	}
	
	fmt.Printf("正在克隆 %s 到 %s...\n", repoInfo.URL, targetPath)
	
	// 克隆仓库
	if err := git.Clone(repoInfo.URL, targetPath); err != nil {
		return fmt.Errorf("克隆仓库失败: %w", err)
	}
	
	// 添加到缓存
	repo := cache.Repository{
		Name:     repoInfo.Name,
		URL:      repoInfo.URL,
		Path:     targetPath,
		Platform: repoInfo.Platform,
	}
	c.cache.Add(repo)
	
	// 保存缓存
	if err := c.cache.Save(); err != nil {
		return fmt.Errorf("保存缓存失败: %w", err)
	}
	
	fmt.Printf("仓库添加成功: %s\n", targetPath)
	
	// 如果启用了 change_directory，输出特殊格式的路径信息供 shell 包装函数使用
	if c.config.ChangeDirectory {
		fmt.Printf("PROJJ_CHANGE_DIRECTORY=%s\n", targetPath)
	}
	
	return nil
}

// Remove 移除仓库
func (c *Client) Remove(query string, deleteFiles bool) error {
	// 查找仓库
	repos := c.cache.Find(query)
	if len(repos) == 0 {
		return fmt.Errorf("未找到匹配的仓库: %s", query)
	}
	
	if len(repos) > 1 {
		fmt.Printf("找到多个匹配的仓库:\n")
		for i, repo := range repos {
			fmt.Printf("%d. %s (%s)\n", i+1, repo.Name, repo.Path)
		}
		return fmt.Errorf("请提供更具体的查询条件")
	}
	
	repo := repos[0]
	
	// 从缓存中移除
	if !c.cache.Remove(repo.Path) {
		return fmt.Errorf("从缓存中移除仓库失败")
	}
	
	// 删除文件（如果需要）
	if deleteFiles {
		if err := os.RemoveAll(repo.Path); err != nil {
			return fmt.Errorf("删除仓库文件失败: %w", err)
		}
		fmt.Printf("已删除仓库文件: %s\n", repo.Path)
	}
	
	// 保存缓存
	if err := c.cache.Save(); err != nil {
		return fmt.Errorf("保存缓存失败: %w", err)
	}
	
	fmt.Printf("仓库移除成功: %s\n", repo.Name)
	return nil
}

// Find 查找仓库
func (c *Client) Find(query string) ([]cache.Repository, error) {
	return c.cache.Find(query), nil
}

// List 列出所有仓库
func (c *Client) List() ([]cache.Repository, error) {
	return c.cache.Repositories, nil
}

// Sync 同步缓存
func (c *Client) Sync() error {
	added, removed, err := c.cache.Sync()
	if err != nil {
		return fmt.Errorf("同步缓存失败: %w", err)
	}
	
	if err := c.cache.Save(); err != nil {
		return fmt.Errorf("保存缓存失败: %w", err)
	}
	
	fmt.Printf("同步完成: 添加 %d 个，移除 %d 个仓库\n", added, removed)
	return nil
}

// GetConfig 获取配置
func (c *Client) GetConfig() *config.Config {
	return c.config
}

// SaveConfig 保存配置
func (c *Client) SaveConfig() error {
	return c.config.Save()
}

// Import 导入现有仓库
func (c *Client) Import(sourcePath string) error {
	// 检查源路径是否存在
	if _, err := os.Stat(sourcePath); os.IsNotExist(err) {
		return fmt.Errorf("源路径不存在: %s", sourcePath)
	}
	
	var imported int
	
	// 遍历源路径下的所有目录
	err := filepath.Walk(sourcePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// 检查是否是 Git 仓库
		if info.IsDir() && git.IsGitRepository(path) {
			// 获取远程 URL
			remoteURL, err := git.GetRemoteURL(path)
			if err != nil {
				fmt.Printf("警告: 无法获取 %s 的远程 URL: %v\n", path, err)
				return nil
			}
			
			// 解析仓库信息
			repoInfo, err := git.ParseURL(remoteURL, c.config.Alias)
			if err != nil {
				fmt.Printf("警告: 无法解析 %s 的 URL %s: %v\n", path, remoteURL, err)
				return nil
			}
			
			// 生成目标路径
			targetPath := repoInfo.GetRepoPath(c.config.GetBasePath())
			
			// 如果目标路径与当前路径不同，移动仓库
			if path != targetPath {
				if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
					return fmt.Errorf("创建目标目录失败: %w", err)
				}
				
				if err := os.Rename(path, targetPath); err != nil {
					return fmt.Errorf("移动仓库失败: %w", err)
				}
				
				fmt.Printf("移动仓库: %s -> %s\n", path, targetPath)
				path = targetPath
			}
			
			// 添加到缓存
			repo := cache.Repository{
				Name:     repoInfo.Name,
				URL:      repoInfo.URL,
				Path:     path,
				Platform: repoInfo.Platform,
			}
			c.cache.Add(repo)
			imported++
			
			// 跳过子目录
			return filepath.SkipDir
		}
		
		return nil
	})
	
	if err != nil {
		return fmt.Errorf("导入仓库失败: %w", err)
	}
	
	// 保存缓存
	if err := c.cache.Save(); err != nil {
		return fmt.Errorf("保存缓存失败: %w", err)
	}
	
	fmt.Printf("导入完成: 共导入 %d 个仓库\n", imported)
	return nil
}

// FormatRepoList 格式化仓库列表输出
func FormatRepoList(repos []cache.Repository, showDetails bool) string {
	if len(repos) == 0 {
		return "未找到任何仓库"
	}
	
	var lines []string
	for _, repo := range repos {
		if showDetails {
			lines = append(lines, fmt.Sprintf("%s\n  URL: %s\n  Path: %s\n  Platform: %s\n  Added: %s",
				repo.Name, repo.URL, repo.Path, repo.Platform, repo.AddedAt.Format("2006-01-02 15:04:05")))
		} else {
			lines = append(lines, fmt.Sprintf("%s (%s)", repo.Name, repo.Path))
		}
	}
	
	return strings.Join(lines, "\n")
}