package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/atian25/projj-go/internal/cache"
	"github.com/atian25/projj-go/internal/config"
	"github.com/atian25/projj-go/internal/git"
	"github.com/spf13/cobra"
)

func newImportCmd() *cobra.Command {
	var fromCache bool

	cmd := &cobra.Command{
		Use:   "import [path]",
		Short: "从本地目录导入仓库，或从缓存恢复",
		RunE: func(cmd *cobra.Command, args []string) error {
			if fromCache {
				return importFromCache()
			}
			if len(args) == 0 {
				return fmt.Errorf("请提供目录路径，或使用 --cache 从缓存恢复")
			}
			return importFromPath(args[0])
		},
	}
	cmd.Flags().BoolVar(&fromCache, "cache", false, "从 cache.json 恢复所有仓库")
	return cmd
}

func importFromPath(root string) error {
	cachePath := cache.DefaultCachePath()
	cfg, err := config.Load(config.DefaultConfigPath())
	if err != nil {
		return err
	}

	repos, err := git.ScanRepos(root)
	if err != nil {
		return err
	}

	base := cfg.BaseDir()
	imported := 0
	for _, repoPath := range repos {
		rel, err := filepath.Rel(base, repoPath)
		if err != nil || strings.HasPrefix(rel, "..") {
			// repo is outside base; use path components relative to root
			rel, err = filepath.Rel(root, repoPath)
			if err != nil {
				continue
			}
		}
		key := filepath.ToSlash(rel)
		if err := cache.Add(cachePath, key, repoPath); err != nil {
			fmt.Fprintf(os.Stderr, "警告: 无法写入缓存 %s: %v\n", key, err)
			continue
		}
		fmt.Printf("导入: %s\n", key)
		imported++
	}
	fmt.Printf("共导入 %d 个仓库\n", imported)
	return nil
}

func importFromCache() error {
	cachePath := cache.DefaultCachePath()
	c, err := cache.Load(cachePath)
	if err != nil {
		return err
	}

	cfg, err := config.Load(config.DefaultConfigPath())
	if err != nil {
		return err
	}

	cloned := 0
	for key, localPath := range c {
		if _, err := os.Stat(localPath); err == nil {
			continue // already exists
		}
		// Reconstruct a clone URL from the key (best-effort HTTPS).
		cloneURL := "https://" + key
		fmt.Printf("克隆: %s -> %s\n", cloneURL, localPath)
		destPath := filepath.Join(cfg.BaseDir(), key)
		if err := git.Clone(cloneURL, destPath); err != nil {
			fmt.Fprintf(os.Stderr, "错误: 克隆 %s 失败: %v\n", cloneURL, err)
		} else {
			cloned++
		}
	}
	fmt.Printf("共恢复 %d 个仓库\n", cloned)
	return nil
}
