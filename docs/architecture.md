# 项目架构说明

## 目录结构

```
.
├── main.go              # 程序入口点
├── cmd/                  # 命令定义目录
│   ├── root.go          # 应用根配置
│   ├── hello.go         # hello 命令
│   ├── version.go       # version 命令
│   └── config.go        # config 命令（包含子命令）
├── docs/                # 文档目录
│   └── architecture.md  # 架构说明
├── go.mod               # Go 模块文件
├── go.sum               # 依赖锁定文件
├── Makefile             # 构建脚本
└── README.md            # 项目说明
```

## 代码组织原则

### 1. 单一职责原则
- `main.go`: 只负责应用启动
- `cmd/`: 每个文件负责一个主要命令
- 每个命令文件包含命令定义和相关的 Action 函数

### 2. 模块化设计
- 使用 `cmd` 包统一管理所有命令
- 通过 `GetAllCommands()` 函数集中注册命令
- 通过 `NewApp()` 函数创建应用配置

### 3. 可扩展性
- 添加新命令只需在 `cmd/` 目录下创建新文件
- 在 `root.go` 的 `GetAllCommands()` 中注册新命令
- 支持嵌套子命令（如 config 命令）

## 命令类型

### 简单命令
如 `hello.go` 和 `version.go`，包含：
- 命令定义函数（返回 `*cli.Command`）
- Action 处理函数
- 相关的标志参数

### 复合命令
如 `config.go`，包含：
- 主命令定义
- 多个子命令定义
- 每个子命令的 Action 函数

## 最佳实践

### 1. 命名约定
- 命令文件名：`{command}.go`
- 命令函数名：`{Command}Command()`
- Action 函数名：`{command}Action()`

### 2. 错误处理
- Action 函数返回 `error`
- 使用 `fmt.Errorf()` 包装错误信息
- 在 main 函数中统一处理错误

### 3. 配置管理
- 版本号等全局配置在 `root.go` 中定义
- 命令特定的配置在各自文件中管理

### 4. 测试友好
- 每个命令都是独立的函数，便于单元测试
- Action 函数接收 context，支持超时和取消

## 扩展示例

添加新的 `deploy` 命令：

1. 创建 `cmd/deploy.go`：
```go
package cmd

import (
    "context"
    "github.com/urfave/cli/v3"
)

func DeployCommand() *cli.Command {
    return &cli.Command{
        Name:   "deploy",
        Usage:  "部署应用",
        Action: deployAction,
        // ... 其他配置
    }
}

func deployAction(ctx context.Context, cmd *cli.Command) error {
    // 部署逻辑
    return nil
}
```

2. 在 `cmd/root.go` 中注册：
```go
func GetAllCommands() []*cli.Command {
    return []*cli.Command{
        HelloCommand(),
        VersionCommand(),
        ConfigCommand(),
        DeployCommand(), // 新增
    }
}
```

这种组织方式使得代码结构清晰，易于维护和扩展。