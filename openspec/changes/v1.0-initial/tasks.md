# 任务清单：projj-go v1.0

## 1. 项目初始化

- [ ] `go mod init github.com/atian25/projj-go`
- [ ] 添加依赖：cobra、bubbletea、bubbles
- [ ] `main.go` 入口
- [ ] `cmd/root.go` cobra 根命令，含版本 flag

## 2. 内部包

- [ ] `internal/config` — 读写 config.json，默认值，`~` 路径展开
- [ ] `internal/cache` — 读写 cache.json，增删查条目
- [ ] `internal/urlutil` — URL 解析（HTTPS、SCP/SSH、别名展开）
- [ ] `internal/git` — Clone、IsGitRepo、ScanRepos
- [ ] `internal/hook` — 按名称执行 hook，注入 PROJJ_HOOK_CONFIG
- [ ] `internal/picker` — bubbletea 交互式选择器

## 3. 命令

- [ ] `cmd/init.go` — 交互式初始化，提示输入 base 目录，写入配置
- [ ] `cmd/add.go` — 别名展开 → 解析 URL → preadd hook → clone → 写缓存 → postadd hook → 输出路径
- [ ] `cmd/add.go` — 内部 `_add` 子命令（供 shell wrapper 调用）
- [ ] `cmd/find.go` — 查询缓存，多结果时调用 picker，输出路径
- [ ] `cmd/list.go` — 打印所有缓存条目
- [ ] `cmd/import.go` — 支持 `--cache` flag，扫描路径或从缓存恢复
- [ ] `cmd/sync.go` — 校验缓存条目有效性，清理失效项
- [ ] `cmd/run.go` — 在当前目录执行指定 hook
- [ ] `cmd/runall.go` — 在所有缓存仓库执行指定 hook，输出汇总
- [ ] `cmd/shell_init.go` — 检测 $SHELL，输出 bash/zsh wrapper function

## 4. 测试

> 所有测试使用项目根目录 `.tmp/{test-name}/` 作为临时目录，不读写 `~/.projj/`。

- [ ] 在 `.gitignore` 中添加 `.tmp/`
- [ ] 封装测试辅助函数 `testutil.NewTestEnv(t)`，返回隔离的 config/cache 路径，并在测试结束时自动清理
- [ ] `urlutil` 单元测试（HTTPS、SCP/SSH、别名展开、`.git` 后缀去除）
- [ ] `config`、`cache` 单元测试（读写、默认值、路径不存在时的处理）
- [ ] `add` → `find` → `list` 集成测试流程（使用 `.tmp/` 隔离环境）

## 5. 收尾

- [ ] 更新 `CLAUDE.md`，补充 go.mod 创建后的构建/测试命令
- [ ] 验证 `go build ./...` 和 `go test ./...` 通过
