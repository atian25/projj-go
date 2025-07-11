# projj-go

一个用 Go 语言编写的 CLI 程序示例，使用 [urfave/cli](https://github.com/urfave/cli) 库构建。

## 功能特性

- 使用 urfave/cli v3 库进行命令行参数解析
- 内置帮助和版本信息
- 支持子命令和标志参数
- 可扩展的命令系统
- 使用 Makefile 简化构建流程

## 快速开始

### 构建程序

```bash
make build
```

### 运行程序

```bash
# 显示帮助信息
./bin/projj-go --help

# 显示版本信息
./bin/projj-go --version

# 运行 hello 命令
./bin/projj-go hello
./bin/projj-go hello 张三

# 使用标志参数
./bin/projj-go hello --name 开发者
./bin/projj-go hello -n 李四

# 查看子命令帮助
./bin/projj-go hello --help
```

### 开发模式运行

```bash
make run
```

## 可用命令

### make 命令

- `make build` - 构建应用
- `make run` - 运行应用
- `make test` - 运行测试
- `make clean` - 清理构建文件
- `make install` - 安装到系统
- `make fmt` - 格式化代码
- `make vet` - 代码检查
- `make help` - 显示帮助信息

### CLI 命令

- `hello` - 打招呼命令
  - `--name, -n` - 指定要打招呼的名字（默认：世界）
  - 支持位置参数：`hello 张三`

- `version, v` - 版本信息命令
  - `--short, -s` - 只显示版本号

- `config` - 配置管理命令
  - `get` - 获取配置值
    - `--key, -k` - 配置键名（必需）
  - `set` - 设置配置值
    - `--key, -k` - 配置键名（必需）
    - `--value, -v` - 配置值（必需）
  - `list` - 列出所有配置
  - `path` - 显示配置文件路径

## 项目结构

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
├── bin/                 # 构建输出目录
├── .gitignore           # Git 忽略文件
├── LICENSE              # 许可证
└── README.md            # 项目说明
```

## 代码组织

### 多命令架构

本项目采用模块化的命令组织方式：

- **`main.go`**: 程序入口点，只负责应用启动
- **`cmd/`**: 命令定义目录，每个文件负责一个主要命令
  - `root.go`: 应用根配置和命令注册
  - `hello.go`: hello 命令实现
  - `version.go`: version 命令实现
  - `config.go`: config 命令及其子命令实现

### 添加新命令

1. 在 `cmd/` 目录下创建新的命令文件：

```go
// cmd/deploy.go
package cmd

import (
    "context"
    "fmt"
    "github.com/urfave/cli/v3"
)

func DeployCommand() *cli.Command {
    return &cli.Command{
        Name:    "deploy",
        Aliases: []string{"d"},
        Usage:   "部署应用",
        Action:  deployAction,
        Flags: []cli.Flag{
            &cli.StringFlag{
                Name:  "env",
                Usage: "部署环境",
                Value: "production",
            },
        },
    }
}

func deployAction(ctx context.Context, cmd *cli.Command) error {
    env := cmd.String("env")
    fmt.Printf("部署到 %s 环境\n", env)
    return nil
}
```

2. 在 `cmd/root.go` 中注册新命令：

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

### 子命令组织

对于复杂的命令（如 `config`），可以包含多个子命令：

```go
func ConfigCommand() *cli.Command {
    return &cli.Command{
        Name:  "config",
        Usage: "配置管理",
        Commands: []*cli.Command{
            {
                Name:   "get",
                Usage:  "获取配置值",
                Action: configGetAction,
                // ...
            },
            {
                Name:   "set",
                Usage:  "设置配置值",
                Action: configSetAction,
                // ...
            },
        },
    }
}
```

### 最佳实践

1. **单一职责**: 每个命令文件只负责一个主要功能
2. **命名约定**: 
   - 文件名：`{command}.go`
   - 函数名：`{Command}Command()` 和 `{command}Action()`
3. **错误处理**: Action 函数应返回 `error`
4. **文档**: 为每个命令和标志提供清晰的 `Usage` 说明

详细的架构说明请参考 [`docs/architecture.md`](docs/architecture.md)。

## 扩展开发

要添加新的命令，请在 `main.go` 中的 `Commands` 切片中添加新的命令定义。

例如：

```go
{
    Name:    "newcommand",
    Usage:   "新命令的描述",
    Action:  newCommandAction,
    Flags: []cli.Flag{
        &cli.StringFlag{
            Name:    "option",
            Aliases: []string{"o"},
            Usage:   "选项描述",
        },
    },
},
```

然后实现对应的 Action 函数（注意 v3 的函数签名变化）：

```go
func newCommandAction(ctx context.Context, cmd *cli.Command) error {
    // 命令逻辑
    // 获取标志参数：cmd.String("option")
    // 获取位置参数：cmd.Args().Get(0)
    return nil
}
```

## 依赖

- [urfave/cli v3](https://github.com/urfave/cli) - 强大的 CLI 应用框架
