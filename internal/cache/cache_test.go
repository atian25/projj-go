package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestNewRepository(t *testing.T) {
	repo := &Repository{
		Name:     "test-repo",
		Path:     "/path/to/repo",
		URL:      "https://github.com/user/repo.git",
		Platform: "github.com",
		AddedAt:  time.Now(),
	}
	
	if repo.Name != "test-repo" {
		t.Errorf("Expected name 'test-repo', got '%s'", repo.Name)
	}
	
	if repo.Platform != "github.com" {
		t.Errorf("Expected platform 'github.com', got '%s'", repo.Platform)
	}
}

func TestCacheLoadSave(t *testing.T) {
	// 创建临时目录
	tempDir, err := os.MkdirTemp("", "projj-cache-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)
	
	// 设置测试环境变量
	originalConfigDir := os.Getenv("PROJJ_CONFIG_DIR")
	os.Setenv("PROJJ_CONFIG_DIR", tempDir)
	defer os.Setenv("PROJJ_CONFIG_DIR", originalConfigDir)
	
	// 创建缓存实例
	cache := &Cache{
		Repositories: []Repository{
			{
				Name:     "repo1",
				Path:     "/path/to/repo1",
				URL:      "https://github.com/user/repo1.git",
				Platform: "github.com",
				AddedAt:  time.Now(),
			},
			{
				Name:     "repo2",
				Path:     "/path/to/repo2",
				URL:      "https://gitlab.com/user/repo2.git",
				Platform: "gitlab.com",
				AddedAt:  time.Now(),
			},
		},
		UpdatedAt: time.Now(),
	}
	
	// 测试保存
	err = cache.Save()
	if err != nil {
		t.Fatalf("Failed to save cache: %v", err)
	}
	
	// 检查文件是否存在
	cachePath := GetCachePath()
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		t.Error("Cache file should exist after save")
	}
	
	// 测试加载
	loadedCache, err := Load()
	if err != nil {
		t.Fatalf("Failed to load cache: %v", err)
	}
	
	// 验证数据
	if len(loadedCache.Repositories) != 2 {
		t.Errorf("Expected 2 repositories, got %d", len(loadedCache.Repositories))
	}
	
	if loadedCache.Repositories[0].Name != "repo1" {
		t.Errorf("Expected first repo name 'repo1', got '%s'", loadedCache.Repositories[0].Name)
	}
}

func TestCacheAdd(t *testing.T) {
	cache := &Cache{
		Repositories: []Repository{},
		UpdatedAt:    time.Now(),
	}
	
	repo := Repository{
		Name:     "new-repo",
		Path:     "/path/to/new-repo",
		URL:      "https://github.com/user/new-repo.git",
		Platform: "github.com",
		AddedAt:  time.Now(),
	}
	
	cache.Add(repo)
	
	if len(cache.Repositories) != 1 {
		t.Errorf("Expected 1 repository after add, got %d", len(cache.Repositories))
	}
	
	if cache.Repositories[0].Name != "new-repo" {
		t.Errorf("Expected repo name 'new-repo', got '%s'", cache.Repositories[0].Name)
	}
	
	// 测试添加重复仓库（相同路径）
	duplicate := Repository{
		Name:     "updated-repo",
		Path:     "/path/to/new-repo", // 相同路径
		URL:      "https://github.com/user/updated-repo.git",
		Platform: "github.com",
		AddedAt:  time.Now(),
	}
	
	cache.Add(duplicate)
	
	// 应该仍然只有一个仓库（重复路径的会被更新）
	if len(cache.Repositories) != 1 {
		t.Errorf("Expected 1 repository after adding duplicate path, got %d", len(cache.Repositories))
	}
	
	// 验证仓库信息已更新
	if cache.Repositories[0].Name != "updated-repo" {
		t.Errorf("Expected updated repo name 'updated-repo', got '%s'", cache.Repositories[0].Name)
	}
}

func TestCacheRemove(t *testing.T) {
	cache := &Cache{
		Repositories: []Repository{
			{
				Name:     "repo1",
				Path:     "/path/to/repo1",
				URL:      "https://github.com/user/repo1.git",
				Platform: "github.com",
				AddedAt:  time.Now(),
			},
			{
				Name:     "repo2",
				Path:     "/path/to/repo2",
				URL:      "https://gitlab.com/user/repo2.git",
				Platform: "gitlab.com",
				AddedAt:  time.Now(),
			},
		},
		UpdatedAt: time.Now(),
	}
	
	// 测试移除存在的仓库（通过路径）
	removed := cache.Remove("/path/to/repo1")
	if !removed {
		t.Error("Should have removed repo1")
	}
	
	if len(cache.Repositories) != 1 {
		t.Errorf("Expected 1 repository after remove, got %d", len(cache.Repositories))
	}
	
	if cache.Repositories[0].Name != "repo2" {
		t.Errorf("Expected remaining repo to be 'repo2', got '%s'", cache.Repositories[0].Name)
	}
	
	// 测试移除不存在的仓库
	removed = cache.Remove("/nonexistent/path")
	if removed {
		t.Error("Should not have removed nonexistent repo")
	}
}

func TestCacheFind(t *testing.T) {
	cache := &Cache{
		Repositories: []Repository{
			{
				Name:     "hello-world",
				Path:     "/path/to/hello-world",
				URL:      "https://github.com/user/hello-world.git",
				Platform: "github.com",
				AddedAt:  time.Now(),
			},
			{
				Name:     "my-project",
				Path:     "/path/to/my-project",
				URL:      "https://gitlab.com/user/my-project.git",
				Platform: "gitlab.com",
				AddedAt:  time.Now(),
			},
			{
				Name:     "another-hello",
				Path:     "/path/to/another-hello",
				URL:      "https://github.com/user/another-hello.git",
				Platform: "github.com",
				AddedAt:  time.Now(),
			},
		},
		UpdatedAt: time.Now(),
	}
	
	// 测试精确匹配
	results := cache.Find("hello-world")
	if len(results) != 1 {
		t.Errorf("Expected 1 result for exact match, got %d", len(results))
	}
	
	// 测试模糊匹配
	results = cache.Find("hello")
	if len(results) != 2 {
		t.Errorf("Expected 2 results for fuzzy match 'hello', got %d", len(results))
	}
	
	// 测试路径匹配
	results = cache.Find("/path/to/my-project")
	if len(results) != 1 {
		t.Errorf("Expected 1 result for path match, got %d", len(results))
	}
	
	// 测试 URL 匹配
	results = cache.Find("gitlab.com")
	if len(results) != 1 {
		t.Errorf("Expected 1 result for URL match, got %d", len(results))
	}
	
	// 测试无匹配
	results = cache.Find("nonexistent")
	if len(results) != 0 {
		t.Errorf("Expected 0 results for no match, got %d", len(results))
	}
}

func TestCacheGetByPath(t *testing.T) {
	cache := &Cache{
		Repositories: []Repository{
			{
				Name:     "repo1",
				Path:     "/path/to/repo1",
				URL:      "https://github.com/user/repo1.git",
				Platform: "github.com",
				AddedAt:  time.Now(),
			},
		},
		UpdatedAt: time.Now(),
	}
	
	// 测试存在的路径
	repo := cache.GetByPath("/path/to/repo1")
	if repo == nil {
		t.Error("Should find repo by path")
	}
	
	if repo.Name != "repo1" {
		t.Errorf("Expected repo name 'repo1', got '%s'", repo.Name)
	}
	
	// 测试不存在的路径
	repo = cache.GetByPath("/nonexistent/path")
	if repo != nil {
		t.Error("Should not find repo for nonexistent path")
	}
}

func TestCacheJSON(t *testing.T) {
	cache := &Cache{
		Repositories: []Repository{
			{
				Name:     "test-repo",
				Path:     "/path/to/test-repo",
				URL:      "https://github.com/user/test-repo.git",
				Platform: "github.com",
				AddedAt:  time.Now(),
			},
		},
		UpdatedAt: time.Now(),
	}
	
	// 测试 JSON 序列化
	data, err := json.Marshal(cache)
	if err != nil {
		t.Fatalf("Failed to marshal cache: %v", err)
	}
	
	// 测试 JSON 反序列化
	var newCache Cache
	err = json.Unmarshal(data, &newCache)
	if err != nil {
		t.Fatalf("Failed to unmarshal cache: %v", err)
	}
	
	// 验证数据一致性
	if len(newCache.Repositories) != 1 {
		t.Errorf("Expected 1 repository after JSON round-trip, got %d", len(newCache.Repositories))
	}
	
	if newCache.Repositories[0].Name != "test-repo" {
		t.Errorf("Repository name mismatch after JSON round-trip")
	}
}

func TestGetCachePath(t *testing.T) {
	// 测试环境变量设置
	testDir := "/test/cache/dir"
	original := os.Getenv("PROJJ_CONFIG_DIR")
	os.Setenv("PROJJ_CONFIG_DIR", testDir)
	defer os.Setenv("PROJJ_CONFIG_DIR", original)
	
	cachePath := GetCachePath()
	expected := filepath.Join(testDir, "cache.json")
	if cachePath != expected {
		t.Errorf("Expected cache path %s, got %s", expected, cachePath)
	}
}