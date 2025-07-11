# Projj-Go 产品需求文档 (PRD)

## 项目概述

### 项目背景
基于 Node.js 版本的 [projj](https://github.com/popomore/projj) <mcreference link="https://github.com/popomore/projj" index="0">0</mcreference>，使用 Go 语言重新实现一个 Git 仓库管理工具。

### 项目目标
- 提供统一的 Git 仓库管理方案
- 支持多平台（GitHub、GitLab、Gitee 等）
- 提供灵活的 Hook 机制
- 简化仓库的克隆、查找和管理流程

### 核心价值
解决开发者在管理多个 Git 仓库时遇到的目录结构混乱、重名冲突、批量操作困难等问题。

## 功能需求

### 1. 核心功能

#### 1.1 初始化 (`projj init`)
- **功能描述**: 初始化 projj 工作环境
- **行为**:
  - 创建配置目录 `~/.projj/`
  - 生成默认配置文件 `~/.projj/config.json`
  - 创建缓存文件 `~/.projj/cache.json`
  - 设置默认基础目录 `~/projj`
- **配置项**:
  ```json
  {
    "base": "~/projj",
    "change_directory": false,
    "alias": {
      "github://": "git@github.com:",
      "gitlab://": "git@gitlab.com:"
    },
    "hooks": {}
  }
  ```

#### 1.2 添加仓库 (`projj add`)
- **功能描述**: 克隆并管理 Git 仓库
- **支持格式**:
  - 完整 Git URL: `git@github.com:user/repo.git`
  - HTTPS URL: `https://github.com/user/repo.git`
  - 别名格式: `github://user/repo`
  - 简短格式: `user/repo` (默认 GitHub)
- **目录结构**:
  ```
  ~/projj/
  ├── github.com/
  │   └── user/
  │       └── repo/
  ├── gitlab.com/
  │   └── user/
  │       └── repo/
  └── gitee.com/
      └── user/
          └── repo/
  ```
- **行为**:
  - 解析 Git URL 并确定目标目录
  - 执行 `git clone`
  - 更新缓存文件
  - 执行 `preadd` 和 `postadd` Hook

#### 1.3 查找仓库 (`projj find`)
- **功能描述**: 快速定位仓库路径
- **使用方式**:
  - `projj find repo` - 查找包含 "repo" 的仓库
  - `projj find user/repo` - 精确查找
  - `projj find` - 列出所有仓库
- **输出**: 仓库的绝对路径
- **可选功能**: 支持自动切换目录 (配置 `change_directory: true`)

#### 1.4 导入仓库 (`projj import`)
- **功能描述**: 导入现有的 Git 仓库
- **支持方式**:
  - `projj import ~/code` - 从指定目录导入
  - `projj import --cache` - 从缓存文件导入
- **行为**:
  - 扫描目录下的 `.git` 文件夹
  - 解析 remote origin URL
  - 移动到标准目录结构
  - 更新缓存文件

#### 1.5 同步缓存 (`projj sync`)
- **功能描述**: 同步缓存与实际文件系统
- **行为**:
  - 检查缓存中的仓库是否存在
  - 移除不存在的仓库记录
  - 扫描并添加新的仓库

### 2. Hook 系统

#### 2.1 命令 Hook
- **preadd**: 在 `projj add` 之前执行
- **postadd**: 在 `projj add` 之后执行
- **配置示例**:
  ```json
  {
    "hooks": {
      "postadd": "echo 'Repository added successfully'",
      "preadd": "echo 'Adding repository...'"
    }
  }
  ```

#### 2.2 自定义 Hook (`projj run`)
- **功能描述**: 在当前仓库执行自定义命令
- **使用方式**: `projj run <hook_name>`
- **配置示例**:
  ```json
  {
    "hooks": {
      "install": "npm install",
      "build": "npm run build",
      "clean": "rm -rf node_modules dist"
    }
  }
  ```

#### 2.3 批量 Hook (`projj runall`)
- **功能描述**: 在所有仓库中执行命令
- **使用方式**: `projj runall <hook_name>`
- **行为**: 遍历缓存中的所有仓库并执行指定 Hook

### 3. 配置管理

#### 3.1 配置查看 (`projj config`)
- `projj config list` - 列出所有配置
- `projj config get <key>` - 获取配置值
- `projj config set <key> <value>` - 设置配置值
- `projj config path` - 显示配置文件路径

#### 3.2 别名管理 (`projj alias`)
- `projj alias list` - 列出所有别名
- `projj alias add <name> <url>` - 添加别名
- `projj alias remove <name>` - 删除别名

### 4. 仓库管理

#### 4.1 列出仓库 (`projj list`)
- 显示所有管理的仓库
- 支持过滤和搜索
- 显示仓库状态（是否有未提交更改）

#### 4.2 移除仓库 (`projj remove`)
- `projj remove <repo>` - 从管理中移除仓库
- 可选择是否删除本地文件

#### 4.3 更新仓库 (`projj update`)
- `projj update <repo>` - 更新指定仓库
- `projj update --all` - 更新所有仓库
- 执行 `git pull` 操作

## 技术需求

### 1. 技术栈
- **语言**: Go 1.21+
- **CLI 框架**: urfave/cli v3
- **配置格式**: JSON
- **Git 操作**: go-git 或调用系统 git 命令

### 2. 项目结构
```
projj-go/
├── cmd/
│   ├── root.go          # 应用配置
│   ├── init.go          # init 命令
│   ├── add.go           # add 命令
│   ├── find.go          # find 命令
│   ├── import.go        # import 命令
│   ├── sync.go          # sync 命令
│   ├── run.go           # run/runall 命令
│   ├── config.go        # config 命令
│   ├── alias.go         # alias 命令
│   ├── list.go          # list 命令
│   ├── remove.go        # remove 命令
│   └── update.go        # update 命令
├── internal/
│   ├── config/          # 配置管理
│   ├── cache/           # 缓存管理
│   ├── git/             # Git 操作
│   ├── hook/            # Hook 系统
│   └── utils/           # 工具函数
├── pkg/
│   └── projj/           # 核心逻辑
└── main.go
```

### 3. 数据结构

#### 3.1 配置文件结构
```go
type Config struct {
    Base            string            `json:"base"`
    ChangeDirectory bool              `json:"change_directory"`
    Alias           map[string]string `json:"alias"`
    Hooks           map[string]string `json:"hooks"`
}
```

#### 3.2 缓存文件结构
```go
type Repository struct {
    Name     string    `json:"name"`
    URL      string    `json:"url"`
    Path     string    `json:"path"`
    Platform string    `json:"platform"`
    AddedAt  time.Time `json:"added_at"`
}

type Cache struct {
    Repositories []Repository `json:"repositories"`
    UpdatedAt    time.Time    `json:"updated_at"`
}
```

## 用户体验需求

### 1. 命令行界面
- 提供清晰的帮助信息
- 支持命令补全
- 友好的错误提示
- 进度显示（克隆大仓库时）

### 2. 兼容性
- 与原版 projj 的配置文件兼容
- 支持从 Node.js 版本迁移
- 跨平台支持（Windows、macOS、Linux）

### 3. 性能要求
- 快速的仓库查找（支持模糊匹配）
- 并发操作支持（批量更新）
- 大量仓库的高效管理

## 实现优先级

### Phase 1 (MVP)
1. `projj init` - 初始化功能
2. `projj add` - 基础添加功能
3. `projj find` - 查找功能
4. `projj list` - 列表功能
5. 基础配置管理

### Phase 2
1. `projj import` - 导入功能
2. `projj sync` - 同步功能
3. Hook 系统基础实现
4. 别名支持

### Phase 3
1. `projj run/runall` - Hook 执行
2. `projj update` - 更新功能
3. `projj remove` - 移除功能
4. 高级配置和优化

## 测试需求

### 1. 单元测试
- 配置管理模块
- 缓存管理模块
- Git URL 解析
- Hook 执行

### 2. 集成测试
- 完整的添加流程
- 导入现有仓库
- 批量操作

### 3. 端到端测试
- 真实 Git 仓库操作
- 跨平台兼容性
- 性能测试

## 文档需求

1. **用户手册**: 详细的使用说明
2. **迁移指南**: 从 Node.js 版本迁移
3. **开发文档**: API 和架构说明
4. **配置参考**: 所有配置选项说明

---

这个 PRD 为 projj-go 的开发提供了完整的功能规划和技术指导，确保我们能够构建一个功能完整、性能优秀的 Git 仓库管理工具。