package config

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

// Config represents ~/.projj/config.json.
type Config struct {
	Base  string            `json:"base"`
	Hooks map[string]string `json:"hooks,omitempty"`
	Alias map[string]string `json:"alias,omitempty"`
}

// DefaultConfigPath returns the default path to config.json.
func DefaultConfigPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".projj", "config.json")
}

// Load reads config from the given path. If the file does not exist, returns a
// zero-value Config (not an error). Missing parent directories are created.
func Load(path string) (Config, error) {
	path = expandHome(path)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		return Config{}, nil
	}
	if err != nil {
		return Config{}, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}
	if cfg.Hooks == nil {
		cfg.Hooks = make(map[string]string)
	}
	if cfg.Alias == nil {
		cfg.Alias = make(map[string]string)
	}
	return cfg, nil
}

// Save writes the config to the given path.
func Save(path string, cfg Config) error {
	path = expandHome(path)
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

// BaseDir returns the expanded base directory from config.
func (c Config) BaseDir() string {
	return expandHome(c.Base)
}

func expandHome(path string) string {
	if strings.HasPrefix(path, "~/") {
		home, _ := os.UserHomeDir()
		return filepath.Join(home, path[2:])
	}
	return path
}
