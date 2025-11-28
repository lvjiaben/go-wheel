# 配置管理

项目使用 [Viper](https://github.com/spf13/viper) 管理配置，支持热更新。

## 配置文件

配置文件位于 `configs/config.yaml`：

```yaml
# 应用配置
app:
  name: "FrameWork"
  mode: "debug"           # debug, release, test
  port: 8801
  version: "1.0.0"
  max_request_body: 10    # 最大请求体大小（MB）

# 日志配置
log:
  level: "debug"          # debug, info, warn, error
  filename: "runtime.log"
  max_size: 200           # 单个日志文件最大大小（MB）
  max_backups: 7          # 保留的旧日志文件数量
  max_age: 30             # 保留旧日志文件的最大天数

# 数据库配置（MySQL）
database:
  driver: "mysql"         # mysql 或 postgres
  host: "127.0.0.1"
  port: 3306
  user: "root"
  pass: ""
  dbname: "go-admin"
  charset: "utf8mb4"
  timezone: "Asia/Shanghai"
  max_open_conns: 100
  max_idle_conns: 20
  max_lifetime: 3600
  max_idle_time: 1800

# Redis 配置
redis:
  state: true             # 是否启用
  host: "127.0.0.1"
  port: 6379
  pass: ""
  db: 0
  pool_size: 100

# RabbitMQ 配置
rabbitmq:
  state: true
  host: "127.0.0.1"
  port: 5672
  user: "guest"
  pass: "guest"
  queue_name: "go-admin"
  exchange: "go-admin-exchange"

# JWT 配置
jwt:
  secret: "your-secret-key"
  expire_day: 7
  issuer: "go-admin"

# 文件上传配置
upload:
  type: "local"           # local, oss, qiniu, cos
  base_url: "http://localhost:8801"
  upload_path: "/public/uploads"
  max_size: 10485760      # 10MB
```

## 读取配置

```go
// 获取配置
config := container.GetConfig()

// 读取配置项
appName := config.App.Name
dbHost := config.Database.Host
jwtSecret := config.Jwt.Secret

// 检查功能是否启用
if config.Redis.State {
    // Redis 已启用
}
```

## 配置结构

```go
// pkg/config/config.go
type Config struct {
    App      AppConfig      `mapstructure:"app"`
    Log      LogConfig      `mapstructure:"log"`
    Database DatabaseConfig `mapstructure:"database"`
    Redis    RedisConfig    `mapstructure:"redis"`
    RabbitMQ RabbitMQConfig `mapstructure:"rabbitmq"`
    Jwt      JwtConfig      `mapstructure:"jwt"`
    Upload   UploadConfig   `mapstructure:"upload"`
    Admin    AdminConfig    `mapstructure:"admin"`
    Api      ApiConfig      `mapstructure:"api"`
}

type AppConfig struct {
    Name           string `mapstructure:"name"`
    Mode           string `mapstructure:"mode"`
    Port           int    `mapstructure:"port"`
    Version        string `mapstructure:"version"`
    MaxRequestBody int    `mapstructure:"max_request_body"`
}
```

## .env 文件配置

支持使用 `.env` 文件覆盖配置（适合生产环境部署）。

**key 格式与 `config.yaml` 路径一致，使用 `.` 分隔层级：**

```bash
# .env
app.port=8801
app.mode=release

database.host=127.0.0.1
database.port=3306
database.user=root
database.pass=your_password
database.dbname=go_admin

redis.host=127.0.0.1
redis.port=6379
redis.pass=

jwt.secret=your_jwt_secret

rabbitmq.state=false

# 多层嵌套示例
upload.type=oss
upload.oss.endpoint=oss-cn-hangzhou.aliyuncs.com
upload.oss.access_key_id=your_key
```

项目提供 `.env.example` 作为模板，部署时复制并修改：

```bash
cp .env.example .env
vim .env
```

配置优先级（从高到低）：
1. `.env` 文件
2. 本地 `configs/config.yaml`
3. 嵌入的配置文件（生产环境）

## 配置热更新

配置文件修改后自动重新加载：

```go
// pkg/container/container.go
viper.OnConfigChange(func(e fsnotify.Event) {
    logger.Info("配置文件已更新", zap.String("file", e.Name))
    // 重新加载配置
    viper.Unmarshal(&config)
})
viper.WatchConfig()
```

## 多环境配置

```bash
# 开发环境
configs/config.yaml

# 生产环境
configs/config.prod.yaml

# 测试环境
configs/config.test.yaml
```

启动时指定配置文件：

```bash
./app -config=configs/config.prod.yaml
```

## 数据库切换

支持 MySQL 和 PostgreSQL，只需修改 `driver` 字段：

```yaml
# MySQL
database:
  driver: "mysql"
  port: 3306
  charset: "utf8mb4"

# PostgreSQL
database:
  driver: "postgres"
  port: 5432
  sslmode: "disable"
```

## 嵌入配置（Embed）

项目使用 Go 的 `embed` 功能将配置文件、模板、i18n 文件编译进二进制：

```go
// embed.go（项目根目录）
package main

import "embed"

//go:embed configs/config.yaml
var ConfigFS embed.FS

//go:embed configs/i18n/*.yaml
var I18nFS embed.FS

//go:embed app/views/*.html
var ViewsFS embed.FS
```

### 加载优先级

| 文件类型 | 开发环境 | 生产环境 |
|----------|----------|----------|
| 配置文件 | 本地 `configs/config.yaml` | 嵌入的配置 + `.env` 覆盖 |
| 模板文件 | 本地 `app/views/*.html` | 嵌入的模板 |
| i18n 文件 | 本地 `configs/i18n/*.yaml` | 嵌入的 i18n |

### 生产环境部署

只需上传一个二进制文件和 `.env` 配置：

```bash
# 构建
make build-linux

# 上传到服务器
scp tmp/admin-linux user@server:/opt/app/

# 服务器上创建 .env
cat > /opt/app/.env << EOF
app.port=8801
app.mode=release
database.host=127.0.0.1
database.pass=your_password
jwt.secret=your_secret
EOF

# 运行
cd /opt/app && ./admin-linux
```

::: tip 热更新支持
开发环境下，本地文件优先于嵌入文件，修改配置/模板/i18n 后无需重新编译即可生效。
:::
