# 构建阶段
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 安装必要的构建工具
RUN apk add --no-cache git make

# 复制go.mod和go.sum
COPY go.mod go.sum ./

# 下载依赖
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o main cmd/server/main.go

# 运行阶段
FROM alpine:latest

WORKDIR /app

# 安装必要的运行时依赖
RUN apk add --no-cache ca-certificates tzdata

# 从构建阶段复制二进制文件
COPY --from=builder /app/main .
# 复制配置文件
COPY configs/config.yaml ./configs/
# 复制语言文件
COPY configs/i18n ./configs/i18n/

# 创建日志目录
RUN mkdir -p logs

# 设置时区
ENV TZ=Asia/Shanghai

# 暴露端口
EXPOSE 8081

# 设置健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8081/health || exit 1

# 运行应用
CMD ["./main"] 