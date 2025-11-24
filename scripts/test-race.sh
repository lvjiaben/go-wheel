#!/bin/bash

# 竞态检测测试脚本
# 使用 Go 的 race detector 检测并发问题

set -e

echo "========================================="
echo "Go 竞态检测测试"
echo "========================================="
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# 检查 Go 版本
echo -e "${YELLOW}检查 Go 版本...${NC}"
go version
echo ""

# 设置测试参数
RACE_FLAGS="-race"
TIMEOUT="10m"
VERBOSE="-v"

# 测试计数器
TOTAL_TESTS=0
PASSED_TESTS=0
FAILED_TESTS=0

# 测试函数
run_test() {
    local package=$1
    local name=$2
    
    echo -e "${YELLOW}测试: ${name}${NC}"
    echo "包: ${package}"
    
    TOTAL_TESTS=$((TOTAL_TESTS + 1))
    
    if go test ${RACE_FLAGS} ${VERBOSE} -timeout ${TIMEOUT} ${package}; then
        echo -e "${GREEN}✓ 通过${NC}"
        PASSED_TESTS=$((PASSED_TESTS + 1))
    else
        echo -e "${RED}✗ 失败${NC}"
        FAILED_TESTS=$((FAILED_TESTS + 1))
    fi
    echo ""
}

# 1. 测试 Container
echo "========================================="
echo "1. Container 竞态检测"
echo "========================================="
run_test "./pkg/container" "Container 并发安全测试"

# 2. 测试资源监控
echo "========================================="
echo "2. 资源监控竞态检测"
echo "========================================="
run_test "./pkg/monitor" "资源监控并发安全测试"

# 3. 测试中间件
echo "========================================="
echo "3. 中间件竞态检测"
echo "========================================="
run_test "./pkg/middleware" "中间件并发安全测试"

# 4. 测试后端服务
echo "========================================="
echo "4. 后端服务竞态检测"
echo "========================================="
run_test "./app/backend/service/..." "后端服务并发安全测试"

# 5. 测试后端控制器
echo "========================================="
echo "5. 后端控制器竞态检测"
echo "========================================="
run_test "./app/backend/controller/..." "后端控制器并发安全测试"

# 6. 测试工具包
echo "========================================="
echo "6. 工具包竞态检测"
echo "========================================="
run_test "./pkg/utils/..." "工具包并发安全测试"

# 7. 测试公共服务
echo "========================================="
echo "7. 公共服务竞态检测"
echo "========================================="
run_test "./app/common/service/..." "公共服务并发安全测试"

# 8. 全局竞态检测（可选，耗时较长）
if [ "$1" == "--full" ]; then
    echo "========================================="
    echo "8. 全局竞态检测（完整扫描）"
    echo "========================================="
    run_test "./..." "全局并发安全测试"
fi

# 输出测试结果
echo "========================================="
echo "测试结果汇总"
echo "========================================="
echo -e "总测试数: ${TOTAL_TESTS}"
echo -e "${GREEN}通过: ${PASSED_TESTS}${NC}"
echo -e "${RED}失败: ${FAILED_TESTS}${NC}"
echo ""

if [ ${FAILED_TESTS} -eq 0 ]; then
    echo -e "${GREEN}=========================================${NC}"
    echo -e "${GREEN}所有竞态检测测试通过！ ✓${NC}"
    echo -e "${GREEN}=========================================${NC}"
    exit 0
else
    echo -e "${RED}=========================================${NC}"
    echo -e "${RED}发现 ${FAILED_TESTS} 个竞态问题！ ✗${NC}"
    echo -e "${RED}=========================================${NC}"
    exit 1
fi

