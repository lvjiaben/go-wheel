# 安装部署

## 环境要求

### 后端环境

| 依赖 | 版本 | 说明 |
|------|------|------|
| Go | >= 1.21 | 编程语言 |
| MySQL | >= 5.7 / 8.0 | 主数据库 |
| PostgreSQL | >= 12 | 可选数据库 |
| Redis | >= 6.0 | 缓存、会话 |
| RabbitMQ | >= 3.8 | 消息队列（可选） |

### 前端环境

| 依赖 | 版本 | 说明 |
|------|------|------|
| Node.js | >= 20.10.0 | 运行环境 |
| pnpm | >= 9.0.0 | 包管理器 |

## 安装步骤

### 1. 克隆项目

```bash
git clone <repository-url>
cd go-wheel
```

### 2. 配置后端

#### 复制配置文件

```bash
cp configs/config.example.yaml configs/config.yaml
```

#### 编辑配置文件

```yaml
# configs/config.yaml
app:
  name: "go-wheel"
  mode: "debug"  # debug, release, test
  port: 8080

database:
  driver: "mysql"  # mysql, postgres
  host: "127.0.0.1"
  port: 3306
  database: "go_wheel"
  username: "root"
  password: "your_password"

redis:
  host: "127.0.0.1"
  port: 6379
  password: ""
  db: 0

jwt:
  secret: "your-jwt-secret-key"
  expire: 7200  # 秒
```

#### 导入数据库

```bash
# MySQL
mysql -u root -p go_wheel < go-admin.sql

# 或使用 Makefile
make db-import
```

### 3. 启动后端

```bash
# 安装依赖
go mod download

# 开发模式（热重载）
make dev

# 或直接运行
go run main.go

# 编译生产版本
make build
./bin/app
```

### 4. 配置前端

```bash
cd vben-admin

# 安装依赖
pnpm install
```

#### 配置 API 地址

编辑 `apps/admin/.env.development`：

```env
VITE_GLOB_API_URL=/api
```

### 5. 启动前端

```bash
# 管理后台
pnpm dev:admin

# 用户端
pnpm dev:user

# 构建生产版本
pnpm build:admin
pnpm build:user
```

## 生产部署

### 后端部署

#### 编译

```bash
# Linux
make build-linux

# 或手动编译
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o bin/app main.go
```

#### Systemd 服务

创建 `/etc/systemd/system/go-wheel.service`：

```ini
[Unit]
Description=Go-Wheel Service
After=network.target

[Service]
Type=simple
User=www
WorkingDirectory=/opt/go-wheel
ExecStart=/opt/go-wheel/bin/app
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
systemctl daemon-reload
systemctl enable go-wheel
systemctl start go-wheel
```

### 前端部署

#### Nginx 配置

```nginx
server {
    listen 80;
    server_name admin.example.com;
    root /var/www/admin;
    index index.html;

    location / {
        try_files $uri $uri/ /index.html;
    }

    location /api {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    location /ws {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
```

## Docker 部署

```bash
# 构建镜像
docker build -t go-wheel .

# 运行容器
docker run -d \
  --name go-wheel \
  -p 8080:8080 \
  -v ./configs:/app/configs \
  go-wheel
```

## 常见问题

### 数据库连接失败

检查配置文件中的数据库连接信息是否正确，确保数据库服务已启动。

### Redis 连接失败

确保 Redis 服务已启动，检查密码配置是否正确。

### 端口被占用

修改配置文件中的端口号，或停止占用端口的进程。

