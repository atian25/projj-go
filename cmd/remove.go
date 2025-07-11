package cmd

import (
	"context"
	"fmt"

	"github.com/atian25/projj-go/pkg/projj"
	"github.com/urfave/cli/v3"
)

// RemoveCommand 返回 remove 命令
func RemoveCommand() *cli.Command {
	return &cli.Command{
		Name:      "remove",
		Usage:     "移除仓库",
		Action:    removeAction,
		Aliases:   []string{"rm", "delete"},
		ArgsUsage: "<query>",
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "delete-files",
				Usage: "同时删除本地文件",
				Aliases: []string{"d"},
			},
		},
		Description: `从 projj 管理中移除仓库。

默认情况下只从缓存中移除，不删除本地文件。
使用 --delete-files 标志可以同时删除本地文件。

示例:
  projj remove golang/go           # 只从管理中移除
  projj remove --delete-files go   # 移除并删除文件`,
	}
}

func removeAction(ctx context.Context, cmd *cli.Command) error {
	if cmd.Args().Len() == 0 {
		return fmt.Errorf("请提供要移除的仓库查询条件")
	}
	
	query := cmd.Args().Get(0)
	deleteFiles := cmd.Bool("delete-files")
	
	client, err := projj.New()
	if err != nil {
		return fmt.Errorf("创建客户端失败: %w", err)
	}
	
	return client.Remove(query, deleteFiles)
}