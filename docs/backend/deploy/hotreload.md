# 热更新

项目使用 [Air](https://github.com/cosmtrek/air) 实现开发环境热更新。

## 安装 Air

```bash
# 使用 go install
go install github.com/cosmtrek/air@latest

# 或使用 curl
curl -sSfL https://raw.githubusercontent.com/cosmtrek/air/master/install.sh | sh -s -- -b $(go env GOPATH)/bin
```

## 配置文件

项目根目录的 `.air.toml` 配置：

```toml
root = "."
testdata_dir = "testdata"
tmp_dir = "tmp"

[build]
  # 构建命令
  cmd = "go build -o ./tmp/main ."
  
  # 输出二进制文件
  bin = "./tmp/main"
  
  # 延迟时间（毫秒）
  delay = 1000
  
  # 排除目录
  exclude_dir = [
    "assets", 
    "tmp", 
    "vendor", 
    "testdata", 
    "node_modules", 
    "vben-admin", 
    "docs", 
    "public"
  ]
  
  # 排除文件
  exclude_regex = ["_test.go"]
  
  # 监听文件类型
  include_ext = ["go", "tpl", "tmpl", "html", "yaml", "yml", "json"]
  
  # 构建错误日志
  log = "build-errors.log"
  
  # 遇到错误停止
  stop_on_error = true
  
  # 自动重启
  rerun = true
  rerun_delay = 500

[color]
  build = "yellow"
  main = "magenta"
  runner = "green"
  watcher = "cyan"

[log]
  main_only = false
  time = true

[misc]
  clean_on_exit = false

[screen]
  clear_on_rebuild = false
  keep_scroll = true
```

## 使用方法

### 启动热更新

```bash
# 在项目根目录执行
air

# 或指定配置文件
air -c .air.toml
```

### 使用 Makefile

```bash
make dev
```

## 工作原理

1. Air 监听 `include_ext` 指定的文件类型
2. 文件变化后等待 `delay` 毫秒
3. 执行 `cmd` 构建命令
4. 重启 `bin` 指定的二进制文件

## 常见问题

### 排除不需要监听的目录

修改 `exclude_dir` 配置：

```toml
exclude_dir = ["assets", "tmp", "vendor", "node_modules", "vben-admin"]
```

### 监听配置文件变化

确保 `include_ext` 包含配置文件类型：

```toml
include_ext = ["go", "yaml", "yml", "json"]
```

### 构建失败不重启

设置 `stop_on_error = true`：

```toml
stop_on_error = true
```

### 清理临时文件

```bash
rm -rf tmp/
```

## 生产环境

生产环境不使用热更新，直接编译运行：

```bash
# 编译
go build -o main .

# 运行
./main

# 或使用 Makefile
make build
./main
```

## 与 Docker 配合

开发环境 Docker Compose 配置：

```yaml
services:
  app:
    build: .
    volumes:
      - .:/app
    command: air -c .air.toml
    ports:
      - "8801:8801"
```

