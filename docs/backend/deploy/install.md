# 安装部署

## 环境要求

| 组件 | 版本要求 |
|------|----------|
| Go | >= 1.21 |
| MySQL | >= 5.7 或 PostgreSQL >= 10 |
| Redis | >= 6.0 |
| RabbitMQ | >= 3.8（可选） |
| Node.js | >= 18（前端） |

## 后端安装

### 1. 克隆项目

```bash
git clone https://github.com/your-repo/go-admin.git
cd go-admin
```

### 2. 安装依赖

```bash
# 安装 Go 依赖
go mod tidy
go mod download

# 或使用 Makefile
make install
```

### 3. 配置文件

```bash
# 复制配置文件
cp configs/config.example.yaml configs/config.yaml

# 编辑配置
vim configs/config.yaml
```

### 4. 导入数据库

```bash
# MySQL
mysql -u root -p go-admin < database/go-admin.sql

# PostgreSQL
psql -U postgres -d go-admin -f database/go-admin.sql
```

### 5. 启动服务

```bash
# 开发模式（热更新）
make dev

# 直接运行
make run

# 或
go run main.go
```

## 前端安装

### 1. 进入前端目录

```bash
cd vben-admin
```

### 2. 安装依赖

```bash
# 使用 pnpm（推荐）
pnpm install

# 或使用 npm
npm install
```

### 3. 启动开发服务器

```bash
# 后台管理
pnpm dev:backend

# 用户端
pnpm dev:user
```

### 4. 构建生产版本

```bash
pnpm build
```

## 生产部署

### 后端构建

```bash
# 构建 Linux 版本
make build-prod

# 或手动构建
CGO_ENABLED=0 GOOS=linux go build -o admin main.go
```

### 使用 Systemd

创建服务文件 `/etc/systemd/system/go-admin.service`：

```ini
[Unit]
Description=Go Admin Service
After=network.target mysql.service redis.service

[Service]
Type=simple
User=www
WorkingDirectory=/var/www/go-admin
ExecStart=/var/www/go-admin/admin
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
sudo systemctl daemon-reload
sudo systemctl enable go-admin
sudo systemctl start go-admin
```

### 使用 Supervisor

创建配置文件 `/etc/supervisor/conf.d/go-admin.conf`：

```ini
[program:go-admin]
directory=/var/www/go-admin
command=/var/www/go-admin/admin
autostart=true
autorestart=true
user=www
stdout_logfile=/var/log/go-admin/stdout.log
stderr_logfile=/var/log/go-admin/stderr.log
```

### Nginx 反向代理

```nginx
server {
    listen 80;
    server_name api.example.com;

    location / {
        proxy_pass http://127.0.0.1:8801;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
    }

    # WebSocket
    location /ws {
        proxy_pass http://127.0.0.1:8801;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
```

## Docker 部署

### Dockerfile

```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o admin main.go

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/admin .
COPY --from=builder /app/configs ./configs
EXPOSE 8801
CMD ["./admin"]
```

### docker-compose.yml

```yaml
version: '3'
services:
  app:
    build: .
    ports:
      - "8801:8801"
    depends_on:
      - mysql
      - redis
    volumes:
      - ./configs:/app/configs

  mysql:
    image: mysql:8.0
    environment:
      MYSQL_ROOT_PASSWORD: root
      MYSQL_DATABASE: go-admin

  redis:
    image: redis:7-alpine
```

启动：

```bash
docker-compose up -d
```

