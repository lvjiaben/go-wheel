# Makefile 命令

项目提供 Makefile 简化常用操作。

## 命令列表

| 命令 | 说明 |
|------|------|
| `make help` | 显示帮助信息 |
| `make dev` | 启动热更新开发服务器 |
| `make run` | 直接运行程序 |
| `make build` | 构建生产版本 |
| `make build-prod` | 构建优化的生产版本 |
| `make build-linux` | 构建 Linux amd64 版本（服务器部署） |
| `make test` | 运行测试 |
| `make clean` | 清理构建文件 |
| `make install` | 安装依赖 |
| `make fmt` | 格式化代码 |
| `make vet` | 代码检查 |
| `make mod` | 更新 Go 模块 |
| `make air-install` | 安装 Air 热更新工具 |

## 开发命令

### 启动开发服务器

```bash
make dev
```

使用 Air 实现热更新，修改代码后自动重新编译运行。

### 直接运行

```bash
make run
```

不使用热更新，直接运行 `go run main.go`。

### 安装 Air

```bash
make air-install
```

首次使用热更新前需要安装 Air 工具。

## 构建命令

### 普通构建

```bash
make build
```

构建可执行文件到 `./tmp/admin`。

### 生产构建

```bash
make build-prod
```

构建优化的静态链接版本：

```bash
CGO_ENABLED=0 go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o ./tmp/admin .
```

### Linux 服务器构建

```bash
make build-linux
```

专门用于构建 Linux amd64 服务器部署版本：

```bash
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o ./tmp/admin-linux .
```

::: tip 提示
`-a` 参数强制重新编译所有包，确保嵌入的文件（配置、模板、i18n）更新到二进制中。
:::

## 代码质量

### 格式化代码

```bash
make fmt
```

使用 `go fmt` 格式化所有 Go 代码。

### 代码检查

```bash
make vet
```

使用 `go vet` 检查代码问题。

### 运行测试

```bash
make test
```

运行所有测试用例。

## 依赖管理

### 安装依赖

```bash
make install
```

执行 `go mod tidy` 和 `go mod download`。

### 更新模块

```bash
make mod
```

执行 `go mod tidy` 和 `go mod vendor`。

## 清理

```bash
make clean
```

删除 `./tmp` 目录和构建错误日志。

## Air 配置

Air 配置文件 `.air.toml`：

```toml
root = "."
tmp_dir = "tmp"

[build]
cmd = "go build -o ./tmp/admin ./main.go"
bin = "./tmp/admin"
include_ext = ["go", "tpl", "tmpl", "html"]
exclude_dir = ["assets", "tmp", "vendor", "vben-admin"]
delay = 1000

[log]
time = true

[color]
main = "magenta"
watcher = "cyan"
build = "yellow"
runner = "green"
```

## 自定义命令

可以在 Makefile 中添加自定义命令：

```makefile
## 生成 Swagger 文档
swagger:
	@swag init

## 数据库迁移
migrate:
	@go run cmd/migrate/main.go

## 生成代码
gen:
	@go run cmd/gen/main.go
```

## 跨平台构建

```bash
# Linux
GOOS=linux GOARCH=amd64 go build -o admin-linux main.go

# Windows
GOOS=windows GOARCH=amd64 go build -o admin.exe main.go

# macOS
GOOS=darwin GOARCH=amd64 go build -o admin-mac main.go

# macOS ARM (M1/M2)
GOOS=darwin GOARCH=arm64 go build -o admin-mac-arm main.go
```

