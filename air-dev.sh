#!/bin/bash

# Air 开发脚本 - 解决代码错误时热更新问题
# 使用方法: ./air-dev.sh

# 创建临时目录
mkdir -p tmp

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}🔥 启动 Air 热更新开发服务器...${NC}"

# 确保GOPATH/bin在PATH中
export PATH=$PATH:$(go env GOPATH)/bin

# 检查Air是否安装
if ! command -v air &> /dev/null; then
    echo -e "${YELLOW}⚠️  Air未安装，正在安装...${NC}"
    # Air项目已迁移到新仓库
    go install github.com/air-verse/air@latest
    if [ $? -ne 0 ]; then
        echo -e "${RED}❌ Air安装失败，尝试使用旧仓库...${NC}"
        # 如果新仓库失败，尝试旧仓库的特定版本
        go install github.com/cosmtrek/air@v1.52.3
        if [ $? -ne 0 ]; then
            echo -e "${RED}❌ Air安装失败，请检查网络连接${NC}"
            echo -e "${BLUE}💡 你可以手动安装：${NC}"
            echo -e "${BLUE}   go install github.com/air-verse/air@latest${NC}"
            echo -e "${BLUE}   或者${NC}"
            echo -e "${BLUE}   brew install air${NC}"
            exit 1
        fi
    fi
    echo -e "${GREEN}✅ Air安装成功${NC}"
fi

# 函数：检查代码编译
check_compile() {
    echo -e "${BLUE}🔍 检查代码编译状态...${NC}"
    go build -o ./tmp/main . 2>./tmp/compile_error.log
    if [ $? -ne 0 ]; then
        echo -e "${RED}❌ 代码编译失败，请修复以下错误：${NC}"
        cat ./tmp/compile_error.log
        return 1
    else
        echo -e "${GREEN}✅ 代码编译成功${NC}"
        return 0
    fi
}

# 函数：等待代码修复
wait_for_fix() {
    echo -e "${YELLOW}⏳ 等待代码修复...${NC}"
    while true; do
        sleep 2
        if check_compile; then
            echo -e "${GREEN}🎉 代码已修复，重新启动热更新...${NC}"
            break
        fi
    done
}

# 初始编译检查
if ! check_compile; then
    wait_for_fix
fi

# 创建Air配置（如果不存在）
if [ ! -f ".air.toml" ]; then
    echo -e "${YELLOW}⚠️  .air.toml 不存在，使用默认配置${NC}"
fi

# 启动Air并处理错误
echo -e "${GREEN}🚀 启动Air热更新...${NC}"
echo -e "${BLUE}🔧 构建错误日志: build-errors.log${NC}"
echo -e "${BLUE}💡 按 Ctrl+C 停止服务${NC}"
echo "----------------------------------------"

# 设置信号处理
trap 'echo -e "\n${YELLOW}🛑 正在停止服务...${NC}"; exit 0' INT TERM

# 循环启动Air，如果因为编译错误退出则自动重试
while true; do
    echo -e "${GREEN}🚀 启动Air热更新...${NC}"
    
    # 启动Air
    air
    
    # 如果Air退出，检查是否是编译错误
    EXIT_CODE=$?
    if [ $EXIT_CODE -ne 0 ]; then
        echo -e "${RED}❌ Air退出，退出码: $EXIT_CODE${NC}"
        
        # 检查编译错误
        if [ -f "./tmp/compile_error.log" ] && [ -s "./tmp/compile_error.log" ]; then
            echo -e "${RED}📋 发现编译错误：${NC}"
            cat ./tmp/compile_error.log
            echo -e "${YELLOW}⏳ 等待修复代码后自动重启...${NC}"
            wait_for_fix
        else
            echo -e "${YELLOW}⏳ 3秒后重试...${NC}"
            sleep 3
        fi
    else
        # 正常退出
        break
    fi
done

echo -e "${BLUE}👋 开发服务器已停止${NC}"
