.PHONY: build run clean install test

# 应用名称
APP_NAME=projj-go

# 构建目录
BUILD_DIR=./bin

# 默认目标
all: build

# 构建应用
build:
	@echo "构建 $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build -o $(BUILD_DIR)/$(APP_NAME) .
	@echo "构建完成: $(BUILD_DIR)/$(APP_NAME)"

# 运行应用
run:
	go run .

# 运行测试
test:
	go test ./...

# 清理构建文件
clean:
	@echo "清理构建文件..."
	rm -rf $(BUILD_DIR)
	@echo "清理完成"

# 安装到系统
install: build
	@echo "安装 $(APP_NAME) 到系统..."
	cp $(BUILD_DIR)/$(APP_NAME) /usr/local/bin/
	@echo "安装完成"

# 格式化代码
fmt:
	go fmt ./...

# 代码检查
vet:
	go vet ./...

# 显示帮助
help:
	@echo "可用命令:"
	@echo "  build    - 构建应用"
	@echo "  run      - 运行应用"
	@echo "  test     - 运行测试"
	@echo "  clean    - 清理构建文件"
	@echo "  install  - 安装到系统"
	@echo "  fmt      - 格式化代码"
	@echo "  vet      - 代码检查"
	@echo "  help     - 显示此帮助信息"