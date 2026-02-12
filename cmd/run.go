package cmd

import (
	"fmt"
	"os"

	"github.com/atian25/projj-go/internal/config"
	"github.com/atian25/projj-go/internal/hook"
	"github.com/spf13/cobra"
)

func newRunCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "run <hook>",
		Short: "在当前目录执行指定 hook",
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

			cwd, err := os.Getwd()
			if err != nil {
				return err
			}

			return hook.Run(hookCmd, cwd, cfg)
		},
	}
}
