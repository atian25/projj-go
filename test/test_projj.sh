#!/bin/bash

# Projj 功能测试脚本
# 使用独立的测试环境，不影响全局配置

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 测试配置
TEST_DIR="$(pwd)/test/fixtures"
PROJJ_BIN="$(pwd)/bin/projj"

# 清理函数
cleanup() {
    echo -e "${YELLOW}清理测试环境...${NC}"
    rm -rf "$TEST_DIR"
    unset PROJJ_CONFIG_DIR
}

# 错误处理
error_exit() {
    echo -e "${RED}错误: $1${NC}" >&2
    cleanup
    exit 1
}

# 成功信息
success() {
    echo -e "${GREEN}✓ $1${NC}"
}

# 信息输出
info() {
    echo -e "${BLUE}ℹ $1${NC}"
}

# 警告信息
warn() {
    echo -e "${YELLOW}⚠ $1${NC}"
}

echo -e "${BLUE}=== Projj 功能测试 ===${NC}"

# 检查二进制文件是否存在
if [ ! -f "$PROJJ_BIN" ]; then
    error_exit "找不到 projj 二进制文件，请先运行 'make build'"
fi

# 设置测试环境
info "设置测试环境..."
export PROJJ_CONFIG_DIR="$TEST_DIR/.projj"
mkdir -p "$TEST_DIR"

# 注册清理函数
trap cleanup EXIT

echo
echo -e "${BLUE}=== 测试 1: 初始化 ===${NC}"
info "运行 projj init..."
$PROJJ_BIN init
success "初始化完成"

# 检查配置文件是否创建
if [ -f "$PROJJ_CONFIG_DIR/config.json" ]; then
    success "配置文件已创建"
else
    error_exit "配置文件未创建"
fi

if [ -f "$PROJJ_CONFIG_DIR/cache.json" ]; then
    success "缓存文件已创建"
else
    error_exit "缓存文件未创建"
fi

echo
echo -e "${BLUE}=== 测试 2: 查看帮助 ===${NC}"
info "运行 projj --help..."
$PROJJ_BIN --help
success "帮助信息显示正常"

echo
echo -e "${BLUE}=== 测试 3: 列出仓库（空列表）===${NC}"
info "运行 projj list..."
$PROJJ_BIN list
success "空列表显示正常"

echo
echo -e "${BLUE}=== 测试 4: 查找仓库（空结果）===${NC}"
info "运行 projj find..."
$PROJJ_BIN find
success "空查找结果显示正常"

echo
echo -e "${BLUE}=== 测试 5: 添加测试仓库 ===${NC}"
# 创建一个模拟的 Git 仓库用于测试
TEST_REPO_DIR="$TEST_DIR/test-repo"
mkdir -p "$TEST_REPO_DIR"
cd "$TEST_REPO_DIR"
git init --quiet
git config user.email "test@example.com"
git config user.name "Test User"
echo "# Test Repository" > README.md
git add README.md
git commit -m "Initial commit" --quiet
cd - > /dev/null

info "添加本地测试仓库..."
# 注意：这里我们需要修改 projj 来支持本地路径，或者创建一个真实的远程仓库
# 暂时跳过这个测试，因为需要网络连接
warn "跳过添加仓库测试（需要网络连接）"

echo
echo -e "${BLUE}=== 测试 6: 配置管理 ===${NC}"
info "测试配置命令..."
$PROJJ_BIN config list
success "配置列表显示正常"

echo
echo -e "${BLUE}=== 测试 7: 同步缓存 ===${NC}"
info "运行 projj sync..."
$PROJJ_BIN sync
success "同步完成"

echo
echo -e "${GREEN}=== 所有测试完成 ===${NC}"
info "测试环境: $TEST_DIR"
info "配置目录: $PROJJ_CONFIG_DIR"

echo
echo -e "${BLUE}=== 测试环境信息 ===${NC}"
echo "配置文件内容:"
cat "$PROJJ_CONFIG_DIR/config.json" | head -20
echo
echo "缓存文件内容:"
cat "$PROJJ_CONFIG_DIR/cache.json" | head -20

echo
echo -e "${GREEN}✓ 所有基础功能测试通过！${NC}"
echo -e "${YELLOW}注意: 网络相关功能（如添加远程仓库）需要手动测试${NC}"