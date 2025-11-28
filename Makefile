# Go项目 Makefile
# 支持开发、构建、部署等常用操作

BINARY_NAME=admin
MAIN_PATH=.
BUILD_DIR=./tmp

# 默认目标
.DEFAULT_GOAL := help

## 显示帮助信息
help:
	@echo "可用命令："
	@echo "  dev        - 启动热更新开发服务器 (使用Air)"
	@echo "  build      - 构建生产版本"
	@echo "  run        - 直接运行程序"
	@echo "  test       - 运行测试"
	@echo "  clean      - 清理构建文件"
	@echo "  install    - 安装依赖"
	@echo "  air-install- 安装Air工具"
	@echo "  fmt        - 格式化代码"
	@echo "  vet        - 代码检查"
	@echo "  mod        - 更新go模块"
	@echo "  build-linux    - 构建Linux AMD64版本"
	@echo "  build-linux-arm- 构建Linux ARM64版本"

## 启动开发服务器（热更新）
dev:
	@echo "🔥 启动Air热更新开发服务器..."
	@./air-dev.sh

## 直接启动Air（无额外脚本）
air:
	@echo "🔥 直接启动Air..."
	@export PATH=$$PATH:$$(go env GOPATH)/bin && air

## 安装Air工具
air-install:
	@echo "📦 安装Air热更新工具..."
	@echo "🔄 Air项目已迁移到新仓库，使用新地址安装..."
	@go install github.com/air-verse/air@latest || go install github.com/cosmtrek/air@v1.52.3
	@echo "✅ Air安装完成"

## 构建生产版本
build:
	@echo "🔨 构建生产版本..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✅ 构建完成: $(BUILD_DIR)/$(BINARY_NAME)"

## 直接运行程序
run:
	@echo "🚀 运行程序..."
	@go run $(MAIN_PATH)

## 运行测试
test:
	@echo "🧪 运行测试..."
	@go test -v ./...

## 清理构建文件
clean:
	@echo "🧹 清理构建文件..."
	@rm -rf $(BUILD_DIR)
	@rm -f build-errors.log
	@echo "✅ 清理完成"

## 安装依赖
install:
	@echo "📦 安装依赖..."
	@go mod tidy
	@go mod download
	@echo "✅ 依赖安装完成"

## 格式化代码
fmt:
	@echo "🎨 格式化代码..."
	@go fmt ./...
	@echo "✅ 代码格式化完成"

## 代码检查
vet:
	@echo "🔍 检查代码..."
	@go vet ./...
	@echo "✅ 代码检查完成"

## 更新go模块
mod:
	@echo "📄 更新go模块..."
	@go mod tidy
	@go mod vendor
	@echo "✅ 模块更新完成"

## 生产环境构建（优化）
build-prod:
	@echo "🏭 构建生产版本（优化）..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PATH)
	@echo "✅ 生产版本构建完成"

## 构建Linux AMD64版本
build-linux:
	@echo "🐧 构建Linux AMD64版本..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-linux $(MAIN_PATH)
	@echo "✅ Linux版本构建完成: $(BUILD_DIR)/$(BINARY_NAME)_linux"

## 构建Linux ARM64版本
build-linux-arm:
	@echo "🐧 构建Linux ARM64版本..."
	@mkdir -p $(BUILD_DIR)
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -a -ldflags="-s -w" -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 $(MAIN_PATH)
	@echo "✅ Linux ARM64版本构建完成: $(BUILD_DIR)/$(BINARY_NAME)_linux_arm64"

## 检查Air配置
check-air:
	@echo "🔧 检查Air配置..."
	@if [ -f ".air.toml" ]; then \
		echo "✅ Air配置文件存在"; \
		air -c .air.toml -d; \
	else \
		echo "❌ Air配置文件不存在，请运行 make dev 自动创建"; \
	fi

.PHONY: help dev air-install build run test clean install fmt vet mod build-prod build-linux build-linux-arm check-air
