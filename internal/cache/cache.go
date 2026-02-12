package cache

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// Cache maps "host/owner/repo" -> absolute path.
type Cache map[string]string

// DefaultCachePath returns the default path to cache.json.
func DefaultCachePath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".projj", "cache.json")
}

// Load reads cache from the given path. Missing file returns an empty Cache.
func Load(path string) (Cache, error) {
	path = expandHome(path)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return Cache{}, nil
	}
	if err != nil {
		return nil, err
	}

	var c Cache
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, err
	}
	if c == nil {
		return Cache{}, nil
	}
	return c, nil
}

// Save writes the cache to the given path.
func Save(path string, c Cache) error {
	path = expandHome(path)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// Add inserts or updates a key→path entry and saves.
func Add(path string, key string, dir string) error {
	c, err := Load(path)
	if err != nil {
		return err
	}
	c[key] = dir
	return Save(path, c)
}

// Remove deletes a key and saves.
func Remove(path string, key string) error {
	c, err := Load(path)
	if err != nil {
		return err
	}
	delete(c, key)
	return Save(path, c)
}

// Find returns all entries whose key contains the query string.
func (c Cache) Find(query string) []string {
	var results []string
	for key := range c {
		if strings.Contains(key, query) {
			results = append(results, c[key])
		}
	}
	return results
}

// FindKeys returns all keys that contain the query string.
func (c Cache) FindKeys(query string) []string {
	var results []string
	for key := range c {
		if strings.Contains(key, query) {
			results = append(results, key)
		}
	}
	return results
}

func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return path
}
