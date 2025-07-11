package cmd

import (
	"github.com/urfave/cli/v3"
)

// GetAllCommands 返回所有可用命令
func GetAllCommands() []*cli.Command {
	return []*cli.Command{
		// Projj 核心命令
		InitCommand(),
		AddCommand(),
		FindCommand(),
		ListCommand(),
		RemoveCommand(),
		SyncCommand(),
		
		// 原有命令（保留用于演示）
		HelloCommand(),
		VersionCommand(),
		ConfigCommand(),
	}
}

// NewApp 创建并配置 CLI 应用
func NewApp() *cli.Command {
	return &cli.Command{
		Name:        "projj",
		Usage:       "Git 仓库管理工具",
		Description: "Projj 是一个强大的 Git 仓库管理工具，帮助你轻松管理多个代码仓库。\n\n" +
			"支持统一的目录结构、智能的仓库查找、灵活的别名配置等功能。",
		Version:     "1.0.0",
		Commands:    GetAllCommands(),
	}
}