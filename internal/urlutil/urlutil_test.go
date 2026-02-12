package urlutil

import (
	"testing"
)

func TestParseHTTPS(t *testing.T) {
	tests := []struct {
		input string
		host  string
		owner string
		repo  string
	}{
		{"https://github.com/user/repo", "github.com", "user", "repo"},
		{"https://github.com/user/repo.git", "github.com", "user", "repo"},
		{"http://github.com/user/repo", "github.com", "user", "repo"},
	}
	for _, tt := range tests {
		info, err := Parse(tt.input, nil)
		if err != nil {
			t.Errorf("Parse(%q) error: %v", tt.input, err)
			continue
		}
		if info.Host != tt.host || info.Owner != tt.owner || info.Repo != tt.repo {
			t.Errorf("Parse(%q) = {%s %s %s}, want {%s %s %s}",
				tt.input, info.Host, info.Owner, info.Repo,
				tt.host, tt.owner, tt.repo)
		}
	}
}

func TestParseSCP(t *testing.T) {
	tests := []struct {
		input string
		host  string
		owner string
		repo  string
	}{
		{"git@github.com:user/repo", "github.com", "user", "repo"},
		{"git@github.com:user/repo.git", "github.com", "user", "repo"},
		{"git@gitlab.com:org/project.git", "gitlab.com", "org", "project"},
	}
	for _, tt := range tests {
		info, err := Parse(tt.input, nil)
		if err != nil {
			t.Errorf("Parse(%q) error: %v", tt.input, err)
			continue
		}
		if info.Host != tt.host || info.Owner != tt.owner || info.Repo != tt.repo {
			t.Errorf("Parse(%q) = {%s %s %s}, want {%s %s %s}",
				tt.input, info.Host, info.Owner, info.Repo,
				tt.host, tt.owner, tt.repo)
		}
	}
}

func TestParseAlias(t *testing.T) {
	aliases := map[string]string{
		"github": "git@github.com:",
		"gl":     "https://gitlab.com/",
	}

	tests := []struct {
		input string
		host  string
		owner string
		repo  string
	}{
		{"github:user/repo", "github.com", "user", "repo"},
		{"github:user/repo.git", "github.com", "user", "repo"},
		{"gl:user/repo", "gitlab.com", "user", "repo"},
	}
	for _, tt := range tests {
		info, err := Parse(tt.input, aliases)
		if err != nil {
			t.Errorf("Parse(%q) error: %v", tt.input, err)
			continue
		}
		if info.Host != tt.host || info.Owner != tt.owner || info.Repo != tt.repo {
			t.Errorf("Parse(%q) = {%s %s %s}, want {%s %s %s}",
				tt.input, info.Host, info.Owner, info.Repo,
				tt.host, tt.owner, tt.repo)
		}
	}
}

func TestParseInvalid(t *testing.T) {
	_, err := Parse("not-a-url", nil)
	if err == nil {
		t.Error("expected error for invalid URL")
	}
}

func TestRelPath(t *testing.T) {
	info := RepoInfo{Host: "github.com", Owner: "user", Repo: "repo"}
	if got := info.RelPath(); got != "github.com/user/repo" {
		t.Errorf("RelPath() = %q, want %q", got, "github.com/user/repo")
	}
}
