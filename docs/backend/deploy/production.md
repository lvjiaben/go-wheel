# 生产环境部署

本文介绍如何将项目部署到生产服务器。

## 构建

```bash
# 构建 Linux amd64 版本
make build-linux

# 输出文件：tmp/admin-linux
```

## 部署文件

只需要上传以下文件：

| 文件 | 说明 |
|------|------|
| `admin-linux` | 可执行文件（包含嵌入的配置、模板、i18n） |
| `.env` | 环境配置文件 |

## 配置 .env

```bash
# 复制示例文件
cp .env.example .env

# 编辑配置
vim .env
```

`.env` 文件内容：

```bash
# 应用配置
APP_PORT=8801
APP_MODE=release

# 数据库配置
DATABASE_HOST=127.0.0.1
DATABASE_PORT=3306
DATABASE_USER=root
DATABASE_PASS=your_password
DATABASE_DBNAME=go_admin

# Redis配置
REDIS_HOST=127.0.0.1
REDIS_PORT=6379
REDIS_PASS=

# JWT配置
JWT_SECRET=your_jwt_secret_change_this

# RabbitMQ配置（可选）
RABBITMQ_STATE=false
```

## Nginx 反向代理

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://127.0.0.1:8801;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        
        # WebSocket 支持
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_read_timeout 86400s;
    }
}
```

::: tip CORS 说明
项目后端已内置 CORS 中间件，Nginx 无需额外配置跨域头。
:::

## Systemd 服务

创建服务文件：

```bash
sudo vim /etc/systemd/system/go-admin.service
```

内容：

```ini
[Unit]
Description=Go Admin Service
After=network.target mysql.service redis.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/go-admin
ExecStart=/opt/go-admin/admin-linux
Restart=always
RestartSec=5
StandardOutput=append:/var/log/go-admin/stdout.log
StandardError=append:/var/log/go-admin/stderr.log

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
# 创建日志目录
sudo mkdir -p /var/log/go-admin
sudo chown www-data:www-data /var/log/go-admin

# 启动服务
sudo systemctl daemon-reload
sudo systemctl enable go-admin
sudo systemctl start go-admin

# 查看状态
sudo systemctl status go-admin
```

## 部署脚本

```bash
#!/bin/bash
# deploy.sh

SERVER="user@your-server"
REMOTE_PATH="/opt/go-admin"

# 构建
make build-linux

# 上传
scp tmp/admin-linux $SERVER:$REMOTE_PATH/

# 重启服务
ssh $SERVER "sudo systemctl restart go-admin"

echo "部署完成"
```

## 健康检查

```bash
# 检查服务状态
curl http://127.0.0.1:8801/api/common/captcha -X POST

# 检查进程
ps aux | grep admin-linux
```

