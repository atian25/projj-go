package cmd

import (
	"context"
	"fmt"

	"github.com/atian25/projj-go/pkg/projj"
	"github.com/urfave/cli/v3"
)

// ListCommand 返回 list 命令
func ListCommand() *cli.Command {
	return &cli.Command{
		Name:    "list",
		Usage:   "列出所有仓库",
		Action:  listAction,
		Aliases: []string{"ls"},
		Flags: []cli.Flag{
			&cli.BoolFlag{
				Name:  "details",
				Usage: "显示详细信息",
				Aliases: []string{"d"},
			},
			&cli.BoolFlag{
				Name:  "path-only",
				Usage: "只显示路径",
				Aliases: []string{"p"},
			},
		},
		Description: `列出所有管理的仓库。

示例:
  projj list              # 列出所有仓库
  projj list --details    # 显示详细信息
  projj list --path-only  # 只显示路径`,
	}
}

func listAction(ctx context.Context, cmd *cli.Command) error {
	client, err := projj.New()
	if err != nil {
		return fmt.Errorf("创建客户端失败: %w", err)
	}
	
	repos, err := client.List()
	if err != nil {
		return fmt.Errorf("获取仓库列表失败: %w", err)
	}
	
	if len(repos) == 0 {
		fmt.Println("未找到任何仓库，请先使用 'projj add' 添加仓库")
		return nil
	}
	
	// 根据标志决定输出格式
	if cmd.Bool("path-only") {
		for _, repo := range repos {
			fmt.Println(repo.Path)
		}
	} else {
		showDetails := cmd.Bool("details")
		output := projj.FormatRepoList(repos, showDetails)
		fmt.Println(output)
		
		fmt.Printf("\n总计: %d 个仓库\n", len(repos))
	}
	
	return nil
}