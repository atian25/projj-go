package git

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// Clone clones the given URL into dest. If dest already exists, it is a no-op.
// Returns the destination path.
func Clone(url, dest string) error {
	if _, err := os.Stat(dest); err == nil {
		return nil
	}

	if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
		return fmt.Errorf("create parent dir: %w", err)
	}

	cmd := exec.Command("git", "clone", url, dest)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git clone %s: %w", url, err)
	}
	return nil
}

// IsGitRepo reports whether dir is a git repository (contains a .git entry).
func IsGitRepo(dir string) bool {
	_, err := os.Stat(filepath.Join(dir, ".git"))
	return err == nil
}

// ScanRepos recursively walks root and returns all git repository root paths.
func ScanRepos(root string) ([]string, error) {
	var repos []string

	err := filepath.WalkDir(root, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return nil // skip unreadable entries
		}
		if !d.IsDir() {
			return nil
		}
		if IsGitRepo(path) {
			repos = append(repos, path)
			return filepath.SkipDir // don't descend into a repo
		}
		return nil
	})
	return repos, err
}
