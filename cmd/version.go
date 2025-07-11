package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

// VersionCommand 返回 version 命令的定义
func VersionCommand() *cli.Command {
	return &cli.Command{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "显示版本信息",
		Action:  versionAction,
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:    "short",
				Aliases: []string{"s"},
				Usage:   "只显示版本号",
			},
		},
	}
}

func versionAction(ctx context.Context, cmd *cli.Command) error {
	version := "1.0.0" // 这里可以从外部传入或从配置文件读取
	
	if cmd.Bool("short") {
		fmt.Println(version)
	} else {
		fmt.Printf("projj-go version %s\n", version)
		fmt.Println("Built with urfave/cli v3")
	}
	
	return nil
}