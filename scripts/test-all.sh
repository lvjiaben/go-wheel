#!/bin/bash

# 综合测试脚本
# 运行所有测试：单元测试、竞态检测、覆盖率分析

set -e

echo "========================================="
echo "Go Web 框架 - 综合测试套件"
echo "========================================="
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 获取脚本所在目录
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" && pwd )"
PROJECT_ROOT="$( cd "${SCRIPT_DIR}/.." && pwd )"

cd ${PROJECT_ROOT}

# 测试结果
UNIT_TEST_RESULT=0
RACE_TEST_RESULT=0
COVERAGE_RESULT=0

# 1. 运行单元测试
echo -e "${BLUE}=========================================${NC}"
echo -e "${BLUE}步骤 1/3: 运行单元测试${NC}"
echo -e "${BLUE}=========================================${NC}"
echo ""

if go test -v ./...; then
    echo -e "${GREEN}✓ 单元测试通过${NC}"
    UNIT_TEST_RESULT=0
else
    echo -e "${RED}✗ 单元测试失败${NC}"
    UNIT_TEST_RESULT=1
fi

echo ""

# 2. 运行竞态检测
echo -e "${BLUE}=========================================${NC}"
echo -e "${BLUE}步骤 2/3: 运行竞态检测${NC}"
echo -e "${BLUE}=========================================${NC}"
echo ""

if [ -f "${SCRIPT_DIR}/test-race.sh" ]; then
    if bash "${SCRIPT_DIR}/test-race.sh"; then
        echo -e "${GREEN}✓ 竞态检测通过${NC}"
        RACE_TEST_RESULT=0
    else
        echo -e "${RED}✗ 竞态检测失败${NC}"
        RACE_TEST_RESULT=1
    fi
else
    echo -e "${YELLOW}⚠ 竞态检测脚本不存在，跳过${NC}"
    RACE_TEST_RESULT=0
fi

echo ""

# 3. 生成覆盖率报告
echo -e "${BLUE}=========================================${NC}"
echo -e "${BLUE}步骤 3/3: 生成覆盖率报告${NC}"
echo -e "${BLUE}=========================================${NC}"
echo ""

if [ -f "${SCRIPT_DIR}/test-coverage.sh" ]; then
    if bash "${SCRIPT_DIR}/test-coverage.sh"; then
        echo -e "${GREEN}✓ 覆盖率分析完成${NC}"
        COVERAGE_RESULT=0
    else
        echo -e "${RED}✗ 覆盖率分析失败${NC}"
        COVERAGE_RESULT=1
    fi
else
    echo -e "${YELLOW}⚠ 覆盖率脚本不存在，跳过${NC}"
    COVERAGE_RESULT=0
fi

echo ""

# 4. 汇总结果
echo -e "${BLUE}=========================================${NC}"
echo -e "${BLUE}测试结果汇总${NC}"
echo -e "${BLUE}=========================================${NC}"
echo ""

if [ ${UNIT_TEST_RESULT} -eq 0 ]; then
    echo -e "单元测试:   ${GREEN}✓ 通过${NC}"
else
    echo -e "单元测试:   ${RED}✗ 失败${NC}"
fi

if [ ${RACE_TEST_RESULT} -eq 0 ]; then
    echo -e "竞态检测:   ${GREEN}✓ 通过${NC}"
else
    echo -e "竞态检测:   ${RED}✗ 失败${NC}"
fi

if [ ${COVERAGE_RESULT} -eq 0 ]; then
    echo -e "覆盖率分析: ${GREEN}✓ 通过${NC}"
else
    echo -e "覆盖率分析: ${RED}✗ 失败${NC}"
fi

echo ""

# 5. 最终结果
TOTAL_FAILURES=$((UNIT_TEST_RESULT + RACE_TEST_RESULT + COVERAGE_RESULT))

if [ ${TOTAL_FAILURES} -eq 0 ]; then
    echo -e "${GREEN}=========================================${NC}"
    echo -e "${GREEN}所有测试通过！ ✓${NC}"
    echo -e "${GREEN}=========================================${NC}"
    exit 0
else
    echo -e "${RED}=========================================${NC}"
    echo -e "${RED}发现 ${TOTAL_FAILURES} 个测试失败！ ✗${NC}"
    echo -e "${RED}=========================================${NC}"
    exit 1
fi

