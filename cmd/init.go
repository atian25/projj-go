package cmd

import (
	"context"
	"fmt"

	"github.com/atian25/projj-go/pkg/projj"
	"github.com/urfave/cli/v3"
)

// InitCommand 返回 init 命令
func InitCommand() *cli.Command {
	return &cli.Command{
		Name:    "init",
		Usage:   "初始化 projj 环境",
		Action:  initAction,
		Aliases: []string{"i"},
	}
}

func initAction(ctx context.Context, cmd *cli.Command) error {
	client, err := projj.New()
	if err != nil {
		return fmt.Errorf("创建客户端失败: %w", err)
	}
	
	return client.Init()
}