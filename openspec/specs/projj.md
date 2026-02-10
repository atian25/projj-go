# 规格文档：projj-go

> 本文档为 projj-go 的持续维护规格，随需求演进持续更新。

## 概述

projj-go 是一个用 Go 实现的本地代码仓库管理 CLI。核心理念是将所有 Git 仓库按 `$BASE/{host}/{owner}/{repo}` 的层级结构统一管理，并通过 Hook 机制支持自动化操作。

---

## 需求：目录结构

工具 **必须** 将所有克隆的仓库按 `$BASE/{host}/{owner}/{repo}` 的格式存放在可配置的根目录下。

#### 场景：克隆 GitHub 仓库
- **给定** base 为 `~/projj`
- **当** 用户执行 `projj add https://github.com/user/repo`
- **则** 仓库被克隆到 `~/projj/github.com/user/repo`

---

## 需求：配置文件

工具 **必须** 从 `~/.projj/config.json` 读写配置。

配置 **必须** 支持以下字段：

| 字段 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `base` | string | `~/projj` | 仓库根目录 |
| `hooks` | object | `{}` | Hook 名称 → Shell 命令映射 |
| `alias` | object | `{}` | URL 前缀别名映射 |

#### 场景：别名展开
- **给定** 配置包含 `"alias": { "github": "git@github.com:" }`
- **当** 用户执行 `projj add github:user/repo`
- **则** 别名先被展开，再进行正常的 URL 解析和克隆

---

## 需求：缓存文件

工具 **必须** 在 `~/.projj/cache.json` 维护一个缓存文件，记录仓库 key 到本地路径的映射。

缓存 key 格式 **必须** 为 `{host}/{owner}/{repo}`（不含前导斜杠）。

```json
{
  "github.com/user/repo": "/Users/tz/projj/github.com/user/repo"
}
```

---

## 需求：命令 — init

工具 **必须** 提供 `init` 命令，以交互式方式初始化 `~/.projj/config.json`。

#### 场景：首次初始化
- **给定** 配置文件不存在
- **当** 用户执行 `projj init`
- **则** 工具提示输入 base 目录（默认 `~/projj`），创建目录和配置文件

---

## 需求：命令 — add

工具 **必须** 提供 `add <url>` 命令，克隆仓库并写入缓存。

执行顺序 **必须** 为：
1. 展开 URL 中的别名
2. 解析 URL，提取 `{host}/{owner}/{repo}`
3. 执行 `preadd` hook（如已配置）
4. 将仓库克隆到 `$BASE/{host}/{owner}/{repo}`
5. 写入缓存
6. 执行 `postadd` hook（如已配置）
7. 将克隆目录路径输出到 stdout

#### 场景：目标目录已存在
- **给定** 目标目录已存在
- **当** 用户执行 `projj add <url>`
- **则** 跳过克隆，更新缓存条目，输出路径

---

## 需求：命令 — find

工具 **必须** 提供 `find [query]` 命令，在缓存中搜索匹配的仓库并将路径输出到 stdout。

- 恰好一个匹配：直接输出路径
- 多个匹配：展示交互式选择器供用户选择
- 无匹配：以非零状态退出，向 stderr 打印错误

#### 场景：模糊匹配
- **给定** 缓存中存在 `github.com/user/my-project`
- **当** 用户执行 `projj find my-proj`
- **则** 该条目被匹配并输出其路径

---

## 需求：命令 — list

工具 **必须** 提供 `list` 命令，将所有缓存的仓库以 `{host}/{owner}/{repo}` 格式逐行输出到 stdout。

---

## 需求：命令 — import

工具 **必须** 提供 `import` 命令，支持两种模式：

**模式一 — 从路径导入：** `projj import <path>`
- 递归扫描 `<path>` 中的 Git 仓库（含 `.git` 目录的文件夹）
- 将找到的每个仓库写入缓存

**模式二 — 从缓存恢复：** `projj import --cache`
- 读取 `cache.json` 中的所有条目
- 克隆本地路径不存在的仓库

---

## 需求：命令 — sync

工具 **必须** 提供 `sync` 命令，验证缓存完整性。

- 检查每个缓存条目对应的本地路径是否存在
- 删除路径已不存在的条目
- 将删除的条目输出到 stdout

---

## 需求：Hook 系统

工具 **必须** 支持在配置的 `hooks` 字段中定义 Shell 命令作为 Hook。

内置 Hook 名称：
- `preadd`：`add` 命令克隆前执行
- `postadd`：`add` 命令克隆后执行

自定义 Hook：`hooks` 中非内置名称的任意 key，通过 `run` / `runall` 命令触发。

Hook 执行时 **必须** 注入环境变量 `PROJJ_HOOK_CONFIG`，其值为完整的 config JSON 字符串。

---

## 需求：命令 — run

工具 **必须** 提供 `run <hook>` 命令，在当前工作目录执行指定 Hook。

#### 场景：Hook 不存在
- **给定** 配置中未定义该 Hook 名称
- **当** 用户执行 `projj run <hook>`
- **则** 以非零状态退出并向 stderr 打印错误

---

## 需求：命令 — runall

工具 **必须** 提供 `runall <hook>` 命令，在缓存中所有仓库目录依次执行指定 Hook。

- 单个仓库执行失败 **必须** 被记录，但 **不得** 中断其余仓库的执行
- 所有执行完成后 **必须** 打印成功/失败的汇总信息

---

## 需求：命令 — shell-init

工具 **必须** 提供 `shell-init` 命令，将 Shell 集成代码输出到 stdout。

输出代码 **必须** 兼容 bash 和 zsh（通过 `$SHELL` 自动检测）。

生成的 wrapper function **必须**：
1. 包装 `projj add`：克隆成功后在当前 Shell 中 `cd` 进入目标目录
2. 包装 `projj find`：定位仓库后在当前 Shell 中 `cd` 进入目标目录
3. 其他子命令透传给真实的 binary，行为不变

#### 场景：Shell 集成配置
- **给定** 用户在 `.zshrc` 中添加了 `eval "$(projj shell-init)"`
- **当** 用户执行 `projj add https://github.com/user/repo`
- **则** 仓库被克隆，当前 Shell 自动切换到该目录

---

## 需求：测试隔离

所有单元测试和集成测试 **必须** 使用项目根目录下的 `.tmp/` 作为临时工作目录，**不得** 读写用户全局配置（`~/.projj/config.json`、`~/.projj/cache.json`）。

- 每个测试用例 **必须** 在 `.tmp/{test-name}/` 下创建独立的 config 和 cache 文件
- 测试结束后 **应当** 清理自己创建的临时目录
- `.tmp/` 目录 **必须** 加入 `.gitignore`

---

## 非目标

- Fish shell 支持
- 终端特定集成（iTerm2、Warp 等）
- Windows 支持
