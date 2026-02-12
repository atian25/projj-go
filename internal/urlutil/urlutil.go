package urlutil

import (
	"fmt"
	"strings"
)

// RepoInfo holds parsed repository information.
type RepoInfo struct {
	Host  string
	Owner string
	Repo  string
}

// RelPath returns "host/owner/repo".
func (r RepoInfo) RelPath() string {
	return fmt.Sprintf("%s/%s/%s", r.Host, r.Owner, r.Repo)
}

// Parse parses a git URL into a RepoInfo, expanding aliases using the provided map.
// Supported formats:
//   - https://github.com/user/repo[.git]
//   - git@github.com:user/repo[.git]
//   - alias:user/repo   (alias expanded using aliasMap)
func Parse(rawURL string, aliasMap map[string]string) (RepoInfo, error) {
	url := expandAlias(rawURL, aliasMap)
	url = strings.TrimSuffix(url, ".git")

	if strings.HasPrefix(url, "https://") || strings.HasPrefix(url, "http://") {
		return parseHTTPS(url)
	}

	if strings.Contains(url, "@") && strings.Contains(url, ":") {
		return parseSCP(url)
	}

	return RepoInfo{}, fmt.Errorf("unsupported URL format: %s", rawURL)
}

func expandAlias(url string, aliasMap map[string]string) string {
	if aliasMap == nil {
		return url
	}
	for alias, prefix := range aliasMap {
		if strings.HasPrefix(url, alias+":") {
			return prefix + url[len(alias)+1:]
		}
	}
	return url
}

func parseHTTPS(url string) (RepoInfo, error) {
	// strip scheme
	url = strings.TrimPrefix(url, "https://")
	url = strings.TrimPrefix(url, "http://")

	parts := strings.SplitN(url, "/", 3)
	if len(parts) != 3 {
		return RepoInfo{}, fmt.Errorf("invalid HTTPS URL: %s", url)
	}
	repo := strings.TrimSuffix(parts[2], ".git")
	return RepoInfo{Host: parts[0], Owner: parts[1], Repo: repo}, nil
}

func parseSCP(url string) (RepoInfo, error) {
	// git@github.com:user/repo
	atIdx := strings.Index(url, "@")
	colonIdx := strings.Index(url, ":")
	if atIdx < 0 || colonIdx < 0 || colonIdx < atIdx {
		return RepoInfo{}, fmt.Errorf("invalid SCP URL: %s", url)
	}
	host := url[atIdx+1 : colonIdx]
	rest := url[colonIdx+1:]
	parts := strings.SplitN(rest, "/", 2)
	if len(parts) != 2 {
		return RepoInfo{}, fmt.Errorf("invalid SCP URL path: %s", url)
	}
	repo := strings.TrimSuffix(parts[1], ".git")
	return RepoInfo{Host: host, Owner: parts[0], Repo: repo}, nil
}
