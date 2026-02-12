package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/atian25/projj-go/internal/config"
	"github.com/spf13/cobra"
)

func newInitCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "init",
		Short: "交互式初始化 projj 配置",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfgPath := config.DefaultConfigPath()

			existing, err := config.Load(cfgPath)
			if err != nil {
				return err
			}

			defaultBase := existing.Base
			if defaultBase == "" {
				home, _ := os.UserHomeDir()
				defaultBase = home + "/projj"
			}

			fmt.Fprintf(os.Stderr, "Base 目录 [%s]: ", defaultBase)
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			input = strings.TrimSpace(input)

			if input != "" {
				existing.Base = input
			} else {
				existing.Base = defaultBase
			}

			if err := config.Save(cfgPath, existing); err != nil {
				return err
			}

			if err := os.MkdirAll(existing.BaseDir(), 0o755); err != nil {
				return err
			}

			fmt.Fprintf(os.Stderr, "✓ 配置已保存到 %s\n", cfgPath)
			fmt.Fprintf(os.Stderr, "✓ Base 目录: %s\n", existing.BaseDir())
			fmt.Fprintf(os.Stderr, "\n在 .zshrc / .bashrc 中添加以下内容以启用 shell 集成：\n")
			fmt.Fprintf(os.Stderr, "  eval \"$(projj shell-init)\"\n")
			return nil
		},
	}
}
