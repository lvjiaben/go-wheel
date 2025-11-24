#!/bin/bash

# 测试覆盖率脚本
# 生成代码覆盖率报告

set -e

echo "========================================="
echo "Go 测试覆盖率分析"
echo "========================================="
echo ""

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 检查 Go 版本
echo -e "${YELLOW}检查 Go 版本...${NC}"
go version
echo ""

# 创建覆盖率输出目录
COVERAGE_DIR="coverage"
mkdir -p ${COVERAGE_DIR}

# 覆盖率文件
COVERAGE_FILE="${COVERAGE_DIR}/coverage.out"
COVERAGE_HTML="${COVERAGE_DIR}/coverage.html"
COVERAGE_FUNC="${COVERAGE_DIR}/coverage_func.txt"

echo -e "${YELLOW}运行测试并生成覆盖率数据...${NC}"
echo ""

# 运行测试并生成覆盖率
if [ "$1" == "--race" ]; then
    echo -e "${BLUE}启用竞态检测...${NC}"
    go test -race -coverprofile=${COVERAGE_FILE} -covermode=atomic ./...
else
    go test -coverprofile=${COVERAGE_FILE} -covermode=atomic ./...
fi

if [ $? -ne 0 ]; then
    echo -e "${RED}测试失败！${NC}"
    exit 1
fi

echo ""
echo -e "${GREEN}✓ 测试通过${NC}"
echo ""

# 生成覆盖率报告
echo -e "${YELLOW}生成覆盖率报告...${NC}"

# 1. 总体覆盖率
echo ""
echo "========================================="
echo "总体覆盖率"
echo "========================================="
go tool cover -func=${COVERAGE_FILE} | tail -1
TOTAL_COVERAGE=$(go tool cover -func=${COVERAGE_FILE} | tail -1 | awk '{print $3}')
echo ""

# 2. 按包统计覆盖率
echo "========================================="
echo "各包覆盖率详情"
echo "========================================="
go tool cover -func=${COVERAGE_FILE} > ${COVERAGE_FUNC}

# 提取每个包的覆盖率
awk '
BEGIN {
    pkg = ""
    total = 0
    covered = 0
}
{
    if ($1 ~ /\.go:/) {
        # 提取包名
        split($1, parts, "/")
        new_pkg = parts[1]
        for (i = 2; i < length(parts); i++) {
            new_pkg = new_pkg "/" parts[i]
        }
        
        if (new_pkg != pkg && pkg != "") {
            if (total > 0) {
                printf "%-50s %6.1f%%\n", pkg, (covered/total)*100
            }
            total = 0
            covered = 0
        }
        pkg = new_pkg
        
        # 统计覆盖率
        if ($3 != "") {
            total++
            if ($3 != "0.0%") {
                covered++
            }
        }
    }
}
END {
    if (total > 0 && pkg != "") {
        printf "%-50s %6.1f%%\n", pkg, (covered/total)*100
    }
}
' ${COVERAGE_FUNC}

echo ""

# 3. 生成 HTML 报告
echo -e "${YELLOW}生成 HTML 覆盖率报告...${NC}"
go tool cover -html=${COVERAGE_FILE} -o ${COVERAGE_HTML}
echo -e "${GREEN}✓ HTML 报告已生成: ${COVERAGE_HTML}${NC}"
echo ""

# 4. 查找未覆盖的代码
echo "========================================="
echo "未覆盖的函数（覆盖率 < 50%）"
echo "========================================="
awk '
{
    if ($1 ~ /\.go:/ && $3 != "" && $3 != "total:") {
        # 提取覆盖率百分比
        coverage = $3
        gsub(/%/, "", coverage)
        if (coverage + 0 < 50) {
            printf "%-60s %s\n", $2, $3
        }
    }
}
' ${COVERAGE_FUNC} | head -20

echo ""

# 5. 覆盖率评分
echo "========================================="
echo "覆盖率评分"
echo "========================================="

# 提取覆盖率数值
COVERAGE_NUM=$(echo ${TOTAL_COVERAGE} | sed 's/%//')

if (( $(echo "$COVERAGE_NUM >= 80" | bc -l) )); then
    echo -e "${GREEN}优秀 (>= 80%): ${TOTAL_COVERAGE}${NC}"
    GRADE="A"
elif (( $(echo "$COVERAGE_NUM >= 60" | bc -l) )); then
    echo -e "${BLUE}良好 (>= 60%): ${TOTAL_COVERAGE}${NC}"
    GRADE="B"
elif (( $(echo "$COVERAGE_NUM >= 40" | bc -l) )); then
    echo -e "${YELLOW}及格 (>= 40%): ${TOTAL_COVERAGE}${NC}"
    GRADE="C"
else
    echo -e "${RED}不及格 (< 40%): ${TOTAL_COVERAGE}${NC}"
    GRADE="D"
fi

echo ""

# 6. 生成覆盖率徽章（可选）
if command -v convert &> /dev/null; then
    echo -e "${YELLOW}生成覆盖率徽章...${NC}"
    # 这里可以集成 shields.io 或其他徽章生成工具
    echo "Coverage: ${TOTAL_COVERAGE}" > ${COVERAGE_DIR}/badge.txt
fi

# 7. 输出文件位置
echo "========================================="
echo "生成的文件"
echo "========================================="
echo "覆盖率数据: ${COVERAGE_FILE}"
echo "HTML 报告:  ${COVERAGE_HTML}"
echo "函数报告:  ${COVERAGE_FUNC}"
echo ""

# 8. 打开 HTML 报告（可选）
if [ "$2" == "--open" ]; then
    echo -e "${YELLOW}在浏览器中打开 HTML 报告...${NC}"
    if [[ "$OSTYPE" == "darwin"* ]]; then
        open ${COVERAGE_HTML}
    elif [[ "$OSTYPE" == "linux-gnu"* ]]; then
        xdg-open ${COVERAGE_HTML}
    else
        echo "请手动打开: ${COVERAGE_HTML}"
    fi
fi

# 9. 覆盖率阈值检查
MIN_COVERAGE=40
if (( $(echo "$COVERAGE_NUM < $MIN_COVERAGE" | bc -l) )); then
    echo -e "${RED}=========================================${NC}"
    echo -e "${RED}警告: 覆盖率 ${TOTAL_COVERAGE} 低于最低要求 ${MIN_COVERAGE}%${NC}"
    echo -e "${RED}=========================================${NC}"
    exit 1
fi

echo -e "${GREEN}=========================================${NC}"
echo -e "${GREEN}覆盖率分析完成！评分: ${GRADE}${NC}"
echo -e "${GREEN}=========================================${NC}"
exit 0

