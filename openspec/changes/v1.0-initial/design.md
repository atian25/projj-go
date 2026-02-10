# 技术设计：projj-go v1.0

## 技术选型

| 关注点 | 选型 | 理由 |
|--------|------|------|
| CLI 框架 | `github.com/spf13/cobra` | Go CLI 事实标准，子命令模型契合 |
| 交互式选择器 | `github.com/charmbracelet/bubbletea` + `bubbles/list` | 自包含 TUI，不依赖外部 fzf |
| Git 操作 | `os/exec` 调用系统 git | 继承用户 SSH key、credential helper、git 配置 |
| JSON 处理 | 标准库 `encoding/json` | 无需额外依赖 |

## 项目结构

```
projj-go/
├── main.go
├── cmd/
│   ├── root.go          # cobra 根命令，持久化 flag，版本信息
│   ├── init.go
│   ├── add.go           # 含内部 _add 子命令
│   ├── find.go
│   ├── list.go
│   ├── import.go
│   ├── sync.go
│   ├── run.go
│   ├── runall.go
│   └── shell_init.go
└── internal/
    ├── config/
    │   └── config.go    # 读写 ~/.projj/config.json；默认值处理
    ├── cache/
    │   └── cache.go     # 读写 ~/.projj/cache.json；增删查
    ├── git/
    │   └── git.go       # Clone(url, dest)；IsGitRepo(path)；ScanRepos(path)
    ├── hook/
    │   └── hook.go      # Run(name, dir, config)；注入 PROJJ_HOOK_CONFIG
    ├── picker/
    │   └── picker.go    # Pick(items []string) (string, error)；bubbletea UI
    └── urlutil/
        └── urlutil.go   # Parse(url) → (host, owner, repo)；ExpandAlias
```

## 关键设计决策

### URL 解析

支持三种 URL 格式：

| 格式 | 示例 | 解析结果 |
|------|------|---------|
| HTTPS | `https://github.com/user/repo` | host=`github.com` |
| SCP (SSH) | `git@github.com:user/repo.git` | host=`github.com` |
| 别名 | `github:user/repo` | 先展开别名，再走正常解析流程 |

repo 名称末尾的 `.git` 后缀会被自动去除。

### shell-init 输出

`projj shell-init` 检测 `$SHELL`，输出对应代码：

**bash/zsh：**
```bash
projj() {
  if [ "$1" = "add" ]; then
    local dir
    dir=$(command projj _add "${@:2}") && cd "$dir"
  elif [ "$1" = "find" ]; then
    local dir
    dir=$(command projj find "${@:2}") && cd "$dir"
  else
    command projj "$@"
  fi
}
```

`add` 的实际逻辑拆分为内部 `_add` 子命令，执行克隆并将目标路径输出到 stdout；shell function 捕获该路径后执行 `cd`。

### 交互式选择器

`picker.Pick(items)` 使用 bubbletea 渲染可过滤的列表。用户按 Ctrl+C 或 Esc 取消时返回 error；选中后返回对应 item 字符串。`find` 命令在缓存查询结果超过一个时调用此函数。

### 配置与缓存的初始化

`config` 和 `cache` 包均暴露 `Load()` 函数，行为如下：
1. 展开路径中的 `~`
2. 如父目录不存在则自动创建
3. 文件不存在时返回零值结构体（不报错）

这确保全新安装时执行 `list` 等命令能优雅处理，而非 panic。

### Hook 执行

```
PROJJ_HOOK_CONFIG=<完整 config JSON> <hook 命令>
```

Hook 命令通过 `sh -c` 在目标目录中执行，stdout/stderr 继承父进程，用户可直接看到输出。
