package cmd

import (
	"context"
	"fmt"

	"github.com/atian25/projj-go/pkg/projj"
	"github.com/urfave/cli/v3"
)

// SyncCommand 返回 sync 命令
func SyncCommand() *cli.Command {
	return &cli.Command{
		Name:    "sync",
		Usage:   "同步缓存与文件系统",
		Action:  syncAction,
		Aliases: []string{"s"},
		Description: `同步缓存与文件系统状态。

检查缓存中的仓库是否还存在，移除不存在的记录。
扫描并添加新发现的仓库到缓存中。

示例:
  projj sync`,
	}
}

func syncAction(ctx context.Context, cmd *cli.Command) error {
	client, err := projj.New()
	if err != nil {
		return fmt.Errorf("创建客户端失败: %w", err)
	}
	
	return client.Sync()
}