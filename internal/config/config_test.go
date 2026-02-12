package config

import (
	"path/filepath"
	"testing"

	"github.com/atian25/projj-go/internal/testutil"
)

func TestLoadNonExistent(t *testing.T) {
	env := testutil.NewTestEnv(t)
	cfg, err := Load(env.ConfigPath)
	if err != nil {
		t.Fatalf("Load non-existent: %v", err)
	}
	if cfg.Base != "" {
		t.Errorf("expected empty base, got %q", cfg.Base)
	}
}

func TestSaveAndLoad(t *testing.T) {
	env := testutil.NewTestEnv(t)
	cfg := Config{
		Base: "/tmp/projj",
		Hooks: map[string]string{
			"preadd": "echo pre",
		},
		Alias: map[string]string{
			"github": "git@github.com:",
		},
	}

	if err := Save(env.ConfigPath, cfg); err != nil {
		t.Fatalf("Save: %v", err)
	}

	loaded, err := Load(env.ConfigPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Base != cfg.Base {
		t.Errorf("Base = %q, want %q", loaded.Base, cfg.Base)
	}
	if loaded.Hooks["preadd"] != "echo pre" {
		t.Errorf("Hooks[preadd] = %q", loaded.Hooks["preadd"])
	}
	if loaded.Alias["github"] != "git@github.com:" {
		t.Errorf("Alias[github] = %q", loaded.Alias["github"])
	}
}

func TestLoadCreatesParentDir(t *testing.T) {
	env := testutil.NewTestEnv(t)
	deepPath := filepath.Join(env.Dir, "a", "b", "config.json")
	_, err := Load(deepPath)
	if err != nil {
		t.Fatalf("Load with deep path: %v", err)
	}
}

func TestDefaultValues(t *testing.T) {
	env := testutil.NewTestEnv(t)
	cfg := Config{Base: "~/projj"}
	if err := Save(env.ConfigPath, cfg); err != nil {
		t.Fatalf("Save: %v", err)
	}
	loaded, err := Load(env.ConfigPath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if loaded.Hooks == nil {
		t.Error("Hooks should not be nil after Load")
	}
	if loaded.Alias == nil {
		t.Error("Alias should not be nil after Load")
	}
}
