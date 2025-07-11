package cmd

import (
	"context"
	"fmt"

	"github.com/urfave/cli/v3"
)

// HelloCommand 返回 hello 命令的定义
func HelloCommand() *cli.Command {
	return &cli.Command{
		Name:    "hello",
		Usage:   "打招呼",
		Action:  helloAction,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "name",
				Aliases: []string{"n"},
				Value:   "世界",
				Usage:   "要打招呼的名字",
			},
		},
	}
}

func helloAction(ctx context.Context, cmd *cli.Command) error {
	name := cmd.String("name")
	
	// 如果有位置参数，优先使用位置参数
	if cmd.Args().Len() > 0 {
		name = cmd.Args().Get(0)
	}
	
	fmt.Printf("你好, %s!\n", name)
	return nil
}