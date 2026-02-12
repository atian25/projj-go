package cache

import (
	"testing"

	"github.com/atian25/projj-go/internal/testutil"
)

func TestLoadNonExistent(t *testing.T) {
	env := testutil.NewTestEnv(t)
	c, err := Load(env.CachePath)
	if err != nil {
		t.Fatalf("Load non-existent: %v", err)
	}
	if len(c) != 0 {
		t.Errorf("expected empty cache, got %v", c)
	}
}

func TestAddAndLoad(t *testing.T) {
	env := testutil.NewTestEnv(t)
	if err := Add(env.CachePath, "github.com/user/repo", "/tmp/projj/github.com/user/repo"); err != nil {
		t.Fatalf("Add: %v", err)
	}
	c, err := Load(env.CachePath)
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if c["github.com/user/repo"] != "/tmp/projj/github.com/user/repo" {
		t.Errorf("unexpected value: %q", c["github.com/user/repo"])
	}
}

func TestRemove(t *testing.T) {
	env := testutil.NewTestEnv(t)
	_ = Add(env.CachePath, "github.com/user/repo", "/tmp/path")
	if err := Remove(env.CachePath, "github.com/user/repo"); err != nil {
		t.Fatalf("Remove: %v", err)
	}
	c, _ := Load(env.CachePath)
	if _, ok := c["github.com/user/repo"]; ok {
		t.Error("key should have been removed")
	}
}

func TestFind(t *testing.T) {
	c := Cache{
		"github.com/user/repo":    "/path/to/repo",
		"github.com/user/other":   "/path/to/other",
		"gitlab.com/user/project": "/path/to/project",
	}
	results := c.Find("user/repo")
	if len(results) != 1 || results[0] != "/path/to/repo" {
		t.Errorf("Find('user/repo') = %v, want [/path/to/repo]", results)
	}

	results = c.Find("github.com")
	if len(results) != 2 {
		t.Errorf("Find('github.com') = %v, want 2 results", results)
	}
}

func TestFindKeys(t *testing.T) {
	c := Cache{
		"github.com/user/repo":  "/path/to/repo",
		"github.com/user/other": "/path/to/other",
	}
	keys := c.FindKeys("repo")
	if len(keys) != 1 || keys[0] != "github.com/user/repo" {
		t.Errorf("FindKeys('repo') = %v", keys)
	}
}
