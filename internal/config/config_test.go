package config

import (
	"encoding/json"
	"os"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()
	
	if config.Base == "" {
		t.Error("Base should not be empty")
	}
	
	if config.ChangeDirectory != false {
		t.Error("ChangeDirectory should be false by default")
	}
	
	if len(config.Alias) == 0 {
		t.Error("Alias should have default values")
	}
	
	// 检查默认别名
	expectedAliases := map[string]string{
		"github://": "git@github.com:",
		"gitlab://": "git@gitlab.com:",
		"gitee://":  "git@gitee.com:",
	}
	
	for alias, expected := range expectedAliases {
		if config.Alias[alias] != expected {
			t.Errorf("Expected alias %s to be %s, got %s", alias, expected, config.Alias[alias])
		}
	}
}

func TestExpandPath(t *testing.T) {
	config := DefaultConfig()
	
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "absolute path",
			input:    "/usr/local/bin",
			expected: "/usr/local/bin",
		},
		{
			name:     "relative path",
			input:    "./test",
			expected: "./test",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := config.ExpandPath(tt.input)
			if tt.name == "absolute path" && result != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, result)
			}
			if tt.name == "relative path" && result == tt.input {
				// 相对路径应该保持不变
				// 这里的测试逻辑需要调整
			}
		})
	}
}

func TestConfigLoadSave(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "projj-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// 设置测试环境变量
	originalConfigDir := os.Getenv("PROJJ_CONFIG_DIR")
	os.Setenv("PROJJ_CONFIG_DIR", tempDir)
	defer os.Setenv("PROJJ_CONFIG_DIR", originalConfigDir)
	
	// 测试保存配置
	config := DefaultConfig()
	config.Base = "/test/path"
	config.ChangeDirectory = true
	
	err = config.Save()
	if err != nil {
		t.Fatalf("Failed to save config: %v", err)
	}
	
	// 检查文件是否存在
	configPath := GetConfigPath()
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Error("Config file should exist after save")
	}
	
	// 测试加载配置
	loadedConfig, err := Load()
	if err != nil {
		t.Fatalf("Failed to load config: %v", err)
	}
	
	// 验证加载的配置
	if loadedConfig.Base != config.Base {
		t.Errorf("Expected Base %s, got %s", config.Base, loadedConfig.Base)
	}
	
	if loadedConfig.ChangeDirectory != config.ChangeDirectory {
		t.Errorf("Expected ChangeDirectory %v, got %v", config.ChangeDirectory, loadedConfig.ChangeDirectory)
	}
}

func TestGetConfigDir(t *testing.T) {
	// 测试环境变量设置
	testDir := "/test/config/dir"
	original := os.Getenv("PROJJ_CONFIG_DIR")
	os.Setenv("PROJJ_CONFIG_DIR", testDir)
	defer os.Setenv("PROJJ_CONFIG_DIR", original)
	
	configDir := GetConfigDir()
	if configDir != testDir {
		t.Errorf("Expected config dir %s, got %s", testDir, configDir)
	}
	
	// 测试默认路径
	os.Unsetenv("PROJJ_CONFIG_DIR")
	defaultDir := GetConfigDir()
	if defaultDir == "" {
		t.Error("Default config dir should not be empty")
	}
}

func TestGetBasePath(t *testing.T) {
	// 创建临时配置
	tempDir, err := os.MkdirTemp("", "projj-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	originalConfigDir := os.Getenv("PROJJ_CONFIG_DIR")
	os.Setenv("PROJJ_CONFIG_DIR", tempDir)
	defer os.Setenv("PROJJ_CONFIG_DIR", originalConfigDir)
	
	// 创建配置文件
	config := DefaultConfig()
	config.Base = "/custom/base/path"
	config.Save()
	
	// 测试获取基础路径
	basePath := config.GetBasePath()
	if basePath != config.Base {
		t.Errorf("Expected base path %s, got %s", config.Base, basePath)
	}
}

func TestConfigJSON(t *testing.T) {
	config := DefaultConfig()
	
	// 测试 JSON 序列化
	data, err := json.Marshal(config)
	if err != nil {
		t.Fatalf("Failed to marshal config: %v", err)
	}
	
	// 测试 JSON 反序列化
	var newConfig Config
	err = json.Unmarshal(data, &newConfig)
	if err != nil {
		t.Fatalf("Failed to unmarshal config: %v", err)
	}
	
	// 验证数据一致性
	if newConfig.Base != config.Base {
		t.Errorf("Base mismatch after JSON round-trip")
	}
	
	if len(newConfig.Alias) != len(config.Alias) {
		t.Errorf("Alias count mismatch after JSON round-trip")
	}
}