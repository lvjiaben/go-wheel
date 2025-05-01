#!/bin/bash

# 设置变量
APP_NAME="go-admin"
VERSION="1.0.0"
BUILD_TIME=$(date +%Y%m%d%H%M%S)
GIT_COMMIT=$(git rev-parse --short HEAD)

# 创建构建目录
BUILD_DIR="build"
mkdir -p $BUILD_DIR

# 构建应用
echo "Building $APP_NAME..."
go build -ldflags "-X main.Version=$VERSION -X main.BuildTime=$BUILD_TIME -X main.GitCommit=$GIT_COMMIT" -o $BUILD_DIR/$APP_NAME cmd/server/main.go

# 复制配置文件
echo "Copying configuration files..."
cp configs/config.yaml $BUILD_DIR/

# 创建Docker镜像
echo "Building Docker image..."
docker build -t $APP_NAME:$VERSION -f deployments/docker/Dockerfile .

# 清理
echo "Cleaning up..."
rm -rf $BUILD_DIR

echo "Build completed successfully!" 