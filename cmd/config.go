package cmd

import (
	"context"
	"fmt"
	"strings"

	"github.com/atian25/projj-go/internal/config"
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
	
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}
	
	switch key {
	case "base":
		fmt.Printf("%s\n", cfg.Base)
	case "change_directory":
		fmt.Printf("%t\n", cfg.ChangeDirectory)
	default:
		return fmt.Errorf("未知的配置键: %s", key)
	}
	
	return nil
}

func configSetAction(ctx context.Context, cmd *cli.Command) error {
	key := cmd.String("key")
	value := cmd.String("value")
	
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}
	
	switch key {
	case "base":
		cfg.Base = value
	case "change_directory":
		cfg.ChangeDirectory = strings.ToLower(value) == "true"
	default:
		return fmt.Errorf("未知的配置键: %s", key)
	}
	
	if err := cfg.Save(); err != nil {
		return fmt.Errorf("保存配置失败: %w", err)
	}
	
	fmt.Printf("配置已保存: %s = %s\n", key, value)
	return nil
}

func configListAction(ctx context.Context, cmd *cli.Command) error {
	cfg, err := config.Load()
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}
	
	fmt.Println("当前配置:")
	fmt.Printf("  base = %s\n", cfg.Base)
	fmt.Printf("  change_directory = %t\n", cfg.ChangeDirectory)
	
	if len(cfg.Hooks) > 0 {
		fmt.Println("  hooks:")
		for k, v := range cfg.Hooks {
			fmt.Printf("    %s = %s\n", k, v)
		}
	}
	
	if len(cfg.PostAdd) > 0 {
		fmt.Println("  postadd:")
		for platform, userInfo := range cfg.PostAdd {
			fmt.Printf("    %s:\n", platform)
			for k, v := range userInfo {
				fmt.Printf("      %s = %s\n", k, v)
			}
		}
	}
	
	return nil
}

func configPathAction(ctx context.Context, cmd *cli.Command) error {
	configPath := config.GetConfigPath()
	fmt.Printf("配置文件路径: %s\n", configPath)
	return nil
}