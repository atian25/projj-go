package cmd

import (
	"fmt"
	"os"

	"github.com/atian25/projj-go/internal/cache"
	"github.com/atian25/projj-go/internal/config"
	"github.com/atian25/projj-go/internal/hook"
	"github.com/spf13/cobra"
)

func newRunAllCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "runall <hook>",
		Short: "在所有缓存仓库中执行指定 hook",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			hookName := args[0]

			cfg, err := config.Load(config.DefaultConfigPath())
			if err != nil {
				return err
			}

			hookCmd, ok := cfg.Hooks[hookName]
			if !ok {
				fmt.Fprintf(os.Stderr, "hook %q 未在配置中定义\n", hookName)
				os.Exit(1)
			}

			c, err := cache.Load(cache.DefaultCachePath())
			if err != nil {
				return err
			}

			success, failed := 0, 0
			for key, path := range c {
				fmt.Printf("[%s]\n", key)
				if err := hook.Run(hookCmd, path, cfg); err != nil {
					fmt.Fprintf(os.Stderr, "  ✗ 失败: %v\n", err)
					failed++
				} else {
					success++
				}
			}

			fmt.Printf("\n汇总: 成功 %d，失败 %d\n", success, failed)
			if failed > 0 {
				os.Exit(1)
			}
			return nil
		},
	}
}
