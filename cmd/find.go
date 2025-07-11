package cmd

import (
	"context"
	"fmt"

	"github.com/atian25/projj-go/pkg/projj"
	"github.com/urfave/cli/v3"
)

// FindCommand 返回 find 命令
func FindCommand() *cli.Command {
	return &cli.Command{
		Name:      "find",
		Usage:     "查找仓库",
		Action:    findAction,
		Aliases:   []string{"f"},
		ArgsUsage: "[query]",
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
		Description: `查找管理的仓库。

如果不提供查询条件，将列出所有仓库。
查询支持仓库名、路径和 URL 的模糊匹配。

示例:
  projj find                # 列出所有仓库
  projj find golang         # 查找包含 "golang" 的仓库
  projj find golang/go      # 精确查找
  projj find --details go   # 显示详细信息
  projj find --path-only go # 只显示路径`,
	}
}

func findAction(ctx context.Context, cmd *cli.Command) error {
	client, err := projj.New()
	if err != nil {
		return fmt.Errorf("创建客户端失败: %w", err)
	}
	
	var query string
	if cmd.Args().Len() > 0 {
		query = cmd.Args().Get(0)
	}
	
	repos, err := client.Find(query)
	if err != nil {
		return fmt.Errorf("查找仓库失败: %w", err)
	}
	
	if len(repos) == 0 {
		if query == "" {
			fmt.Println("未找到任何仓库，请先使用 'projj add' 添加仓库")
		} else {
			fmt.Printf("未找到匹配 '%s' 的仓库\n", query)
		}
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
	}
	
	return nil
}