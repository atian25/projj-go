package cmd

import (
	"context"
	"fmt"

	"github.com/atian25/projj-go/pkg/projj"
	"github.com/urfave/cli/v3"
)

// AddCommand 返回 add 命令
func AddCommand() *cli.Command {
	return &cli.Command{
		Name:      "add",
		Usage:     "添加 Git 仓库",
		Action:    addAction,
		Aliases:   []string{"a"},
		ArgsUsage: "<repository-url>",
		Description: `添加 Git 仓库到 projj 管理。

支持的 URL 格式:
  - SSH: git@github.com:user/repo.git
  - HTTPS: https://github.com/user/repo.git
  - 别名: github://user/repo
  - 简短: user/repo (默认 GitHub)

示例:
  projj add git@github.com:golang/go.git
  projj add https://github.com/golang/go.git
  projj add github://golang/go
  projj add golang/go`,
	}
}

func addAction(ctx context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() == 0 {
		return fmt.Errorf("请提供仓库 URL")
	}
	
	repoURL := cmd.Args().Get(0)
	
	client, err := projj.New()
	if err != nil {
		return fmt.Errorf("创建客户端失败: %w", err)
	}
	
	return client.Add(repoURL)
}