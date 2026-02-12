package testutil

import (
	"os"
	"path/filepath"
	"testing"
)

// Env holds paths for an isolated test environment under .tmp/{test-name}/.
type Env struct {
	Dir        string
	ConfigPath string
	CachePath  string
}

// NewTestEnv creates a temporary directory under the project's .tmp/ folder,
// isolated per test. The directory is automatically removed when the test ends.
func NewTestEnv(t *testing.T) *Env {
	t.Helper()

	// Locate project root (walk up from this file's package dir).
	projectRoot := projectRoot()
	tmpDir := filepath.Join(projectRoot, ".tmp", t.Name())

	if err := os.MkdirAll(tmpDir, 0o755); err != nil {
		t.Fatalf("testutil: create temp dir: %v", err)
	}

	t.Cleanup(func() {
		os.RemoveAll(tmpDir)
	})

	return &Env{
		Dir:        tmpDir,
		ConfigPath: filepath.Join(tmpDir, "config.json"),
		CachePath:  filepath.Join(tmpDir, "cache.json"),
	}
}

func projectRoot() string {
	// Walk up from the current working directory until we find go.mod.
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "."
		}
		dir = parent
	}
}
