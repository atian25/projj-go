# 提案：projj-go v1.0 初始实现

## 意图

用 Go 重写 [projj](https://github.com/popomore/projj)，产出一个独立的、跨平台的二进制文件，无需 Node.js 运行时依赖。原版的核心价值——按 `{host}/{owner}/{repo}` 结构管理 Git 仓库并支持快速目录跳转——完整保留。

## 范围

本次为初始实现，覆盖所有基础功能。

**在范围内：**
- 全部命令：`init`、`add`、`find`、`list`、`import`、`sync`、`run`、`runall`、`shell-init`
- 配置和缓存系统（JSON 格式，与原版 projj 兼容）
- Hook 系统：内置 `preadd`/`postadd` 及自定义 Hook
- Shell 集成：通过 `shell-init` 支持 bash 和 zsh
- `find` 命令多结果时的内置交互式选择器

**不在范围内：**
- Fish shell 支持
- 终端特定集成（iTerm2、Warp 等）
- Windows 支持
- 从原版 projj 的迁移工具（cache.json 格式兼容，无需迁移）

## 方案

使用 cobra 构建标准 Go CLI，子命令路由清晰。Git 操作通过 `os/exec` 委托给系统 git 二进制，继承用户的凭证配置。交互式选择器使用 bubbletea 在进程内实现，不依赖外部 fzf。

Shell 集成通过 `shell-init` 生成的 wrapper function 解决，避免与终端进程绑定。
