# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project

`projj-go` 是 [projj](https://github.com/popomore/projj) 的 Go 重写版，按 `$BASE/{host}/{owner}/{repo}` 层级结构管理本地 Git 仓库。

规格和技术设计文档在 `openspec/` 目录中。

## Commands

```bash
# 构建
go build ./...

# 运行所有测试
go test ./...

# 运行单个测试
go test ./internal/urlutil -run TestParseHTTPS

# 格式化
gofmt -w .
```

## Architecture

```
main.go                  # 入口，调用 cmd.Execute()
cmd/                     # cobra 子命令，每个命令一个文件
internal/
  config/                # 读写 ~/.projj/config.json
  cache/                 # 读写 ~/.projj/cache.json
  urlutil/               # Git URL 解析（HTTPS / SCP / 别名）
  git/                   # 调用系统 git：Clone、IsGitRepo、ScanRepos
  hook/                  # 执行 hook（sh -c，注入 PROJJ_HOOK_CONFIG）
  picker/                # bubbletea 交互式仓库选择器
  testutil/              # 测试辅助：NewTestEnv(t) 在 .tmp/ 下创建隔离环境
```

## Key Conventions

- `config.Load(path)` 和 `cache.Load(path)` 接受路径参数（而非硬编码 `~/.projj/`），方便测试注入 `.tmp/` 路径。
- 文件不存在时返回零值结构体，不报错。
- 所有测试通过 `testutil.NewTestEnv(t)` 隔离，不读写 `~/.projj/`，临时文件放在 `.tmp/` 下。
- `add` 命令的实际逻辑在 `runAdd()`，`_add` 子命令（供 shell wrapper 调用）和 `add` 共用同一函数，区别在于 `_add` 会将路径输出到 stdout。
