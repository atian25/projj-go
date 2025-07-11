.PHONY: build run clean install test unit-test fmt vet deps dev help

# 项目名称
PROJECT_NAME=projj

# 构建目录
BUILD_DIR=./bin

# 测试目录
TEST_DIR=test

# 默认目标
all: build

# 构建应用
build:
	@echo "构建 $(PROJECT_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(PROJECT_NAME) .
	@echo "构建完成: $(BUILD_DIR)/$(PROJECT_NAME)"

# 运行应用
run: build
	./$(BUILD_DIR)/$(PROJECT_NAME)

# 运行功能测试
test: build
	@echo "运行功能测试..."
	chmod +x $(TEST_DIR)/test_projj.sh
	./$(TEST_DIR)/test_projj.sh

# 运行单元测试
unit-test:
	@echo "运行单元测试..."
	go test -v ./...

# 清理构建文件
clean:
	@echo "清理构建文件..."
	rm -rf $(BUILD_DIR)
	rm -rf $(TEST_DIR)/fixtures
	@echo "清理完成"

# 安装到系统
install: build
	@echo "安装 $(PROJECT_NAME) 到系统..."
	cp $(BUILD_DIR)/$(PROJECT_NAME) /usr/local/bin/
	@echo "安装完成"

# 格式化代码
fmt:
	@echo "格式化代码..."
	go fmt ./...

# 代码检查
vet:
	@echo "代码检查..."
	go vet ./...

# 安装依赖
deps:
	@echo "安装依赖..."
	go mod tidy
	go mod download

# 开发环境设置
dev: deps fmt vet build test
	@echo "开发环境准备完成"

# 显示帮助
help:
	@echo "可用命令:"
	@echo "  build      - 构建应用"
	@echo "  run        - 构建并运行应用"
	@echo "  test       - 运行功能测试"
	@echo "  unit-test  - 运行单元测试"
	@echo "  clean      - 清理构建文件和测试文件"
	@echo "  install    - 安装到系统"
	@echo "  fmt        - 格式化代码"
	@echo "  vet        - 代码检查"
	@echo "  deps       - 安装依赖"
	@echo "  dev        - 开发环境完整设置"
	@echo "  help       - 显示此帮助信息"