package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config 表示 projj 的配置结构
type Config struct {
	Base            string                       `json:"base"`
	ChangeDirectory bool                         `json:"change_directory"`
	Alias           map[string]string            `json:"alias"`
	Hooks           map[string]string            `json:"hooks"`
	PostAdd         map[string]map[string]string `json:"postadd"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	return &Config{
		Base:            filepath.Join(homeDir, "projj"),
		ChangeDirectory: false,
		Alias: map[string]string{
			"github://": "git@github.com:",
			"gitlab://": "git@gitlab.com:",
			"gitee://":  "git@gitee.com:",
		},
		Hooks:   make(map[string]string),
		PostAdd: make(map[string]map[string]string),
	}
}

// GetConfigDir 获取配置目录路径
func GetConfigDir() string {
	// 支持测试时使用自定义配置目录
	if testDir := os.Getenv("PROJJ_CONFIG_DIR"); testDir != "" {
		return testDir
	}
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".projj")
}

// GetConfigPath 获取配置文件路径
func GetConfigPath() string {
	return filepath.Join(GetConfigDir(), "config.json")
}

// Load 加载配置文件
func Load() (*Config, error) {
	configPath := GetConfigPath()
	
	// 如果配置文件不存在，返回默认配置
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return DefaultConfig(), nil
	}
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}
	
	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}
	
	return &config, nil
}

// Save 保存配置文件
func (c *Config) Save() error {
	configDir := GetConfigDir()
	
	// 确保配置目录存在
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %w", err)
	}
	
	configPath := GetConfigPath()
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}
	
	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("写入配置文件失败: %w", err)
	}
	
	return nil
}

// ExpandPath 展开路径中的 ~ 符号
func (c *Config) ExpandPath(path string) string {
	if path == "~" {
		homeDir, _ := os.UserHomeDir()
		return homeDir
	}
	if filepath.HasPrefix(path, "~/") {
		homeDir, _ := os.UserHomeDir()
		return filepath.Join(homeDir, path[2:])
	}
	return path
}

// GetBasePath 获取展开后的基础路径
func (c *Config) GetBasePath() string {
	return c.ExpandPath(c.Base)
}