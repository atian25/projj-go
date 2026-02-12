package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/atian25/projj-go/internal/cache"
	"github.com/atian25/projj-go/internal/config"
	"github.com/atian25/projj-go/internal/git"
	"github.com/atian25/projj-go/internal/hook"
	"github.com/atian25/projj-go/internal/urlutil"
	"github.com/spf13/cobra"
)

func newAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "add <url>",
		Short: "克隆仓库到对应目录并写入缓存",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := runAdd(args[0])
			return err
		},
	}
}

// newInternalAddCmd is the internal _add subcommand used by the shell wrapper.
// It prints the destination path to stdout so the shell function can cd to it.
func newInternalAddCmd() *cobra.Command {
	return &cobra.Command{
		Use:    "_add <url>",
		Hidden: true,
		Args:   cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			dest, err := runAdd(args[0])
			if err != nil {
				return err
			}
			fmt.Println(dest)
			return nil
		},
	}
}

func runAdd(rawURL string) (string, error) {
	cfgPath := config.DefaultConfigPath()
	cfg, err := config.Load(cfgPath)
	if err != nil {
		return "", err
	}

	info, err := urlutil.Parse(rawURL, cfg.Alias)
	if err != nil {
		return "", err
	}

	dest := filepath.Join(cfg.BaseDir(), info.RelPath())
	cacheKey := info.RelPath()
	cachePath := cache.DefaultCachePath()

	// preadd hook
	if cmd, ok := cfg.Hooks["preadd"]; ok {
		if err := hook.Run(cmd, dest, cfg); err != nil {
			return "", fmt.Errorf("preadd hook: %w", err)
		}
	}

	if err := git.Clone(rawURL, dest); err != nil {
		return "", err
	}

	if err := cache.Add(cachePath, cacheKey, dest); err != nil {
		return "", err
	}

	// postadd hook
	if cmd, ok := cfg.Hooks["postadd"]; ok {
		if err := hook.Run(cmd, dest, cfg); err != nil {
			return "", fmt.Errorf("postadd hook: %w", err)
		}
	}

	return dest, nil
}
