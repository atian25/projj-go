package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/urfave/cli/v3"
)

// ConfigCommand 返回 config 命令的定义
func ConfigCommand() *cli.Command {
	return &cli.Command{
		Name:  "config",
		Usage: "配置管理",
		Commands: []*cli.Command{
			{
				Name:   "get",
				Usage:  "获取配置值",
				Action: configGetAction,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "key",
						Aliases:  []string{"k"},
						Usage:    "配置键名",
						Required: true,
					},
				},
			},
			{
				Name:   "set",
				Usage:  "设置配置值",
				Action: configSetAction,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "key",
						Aliases:  []string{"k"},
						Usage:    "配置键名",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "value",
						Aliases:  []string{"v"},
						Usage:    "配置值",
						Required: true,
					},
				},
			},
			{
				Name:   "list",
				Usage:  "列出所有配置",
				Action: configListAction,
			},
			{
				Name:   "path",
				Usage:  "显示配置文件路径",
				Action: configPathAction,
			},
		},
	}
}

func configGetAction(ctx context.Context, cmd *cli.Command) error {
	key := cmd.String("key")
	fmt.Printf("获取配置: %s\n", key)
	// 这里可以实现实际的配置读取逻辑
	fmt.Printf("配置值: (示例值)\n")
	return nil
}

func configSetAction(ctx context.Context, cmd *cli.Command) error {
	key := cmd.String("key")
	value := cmd.String("value")
	fmt.Printf("设置配置: %s = %s\n", key, value)
	// 这里可以实现实际的配置写入逻辑
	fmt.Println("配置已保存")
	return nil
}

func configListAction(ctx context.Context, cmd *cli.Command) error {
	fmt.Println("当前配置:")
	fmt.Println("  user.name = 示例用户")
	fmt.Println("  user.email = user@example.com")
	// 这里可以实现实际的配置列表逻辑
	return nil
}

func configPathAction(ctx context.Context, cmd *cli.Command) error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("无法获取用户主目录: %w", err)
	}
	
	configPath := filepath.Join(homeDir, ".projj-go", "config.yaml")
	fmt.Printf("配置文件路径: %s\n", configPath)
	return nil
}