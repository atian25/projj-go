package cmd

import (
	"github.com/urfave/cli/v3"
)

// GetAllCommands 返回所有可用的命令
func GetAllCommands() []*cli.Command {
	return []*cli.Command{
		HelloCommand(),
		VersionCommand(),
		ConfigCommand(),
	}
}

// NewApp 创建并配置应用程序
func NewApp() *cli.Command {
	return &cli.Command{
		Name:        "projj-go",
		Usage:       "Go CLI 程序示例",
		Description: "一个使用 urfave/cli v3 构建的示例 CLI 程序，展示了良好的代码组织结构",
		Version:     "1.0.0",
		Commands:    GetAllCommands(),
	}
}