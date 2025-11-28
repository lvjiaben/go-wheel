# Go-Wheel 企业级后台管理系统

一个基于 Go + Gin + Vben Admin 的现代化企业级后台管理系统，提供完整的前后端分离解决方案。

## 🌐 演示站点

| 站点 | 地址 | 账号 | 密码 |
|------|------|------|------|
| 管理后台 | http://120.26.94.219:8803/ | admin | 123456 |
| 用户端 | http://120.26.94.219:8804/ | - | - |

> ⚠️ 演示站点仅供功能预览，请勿修改系统配置或删除数据。

## ✨ 特性

- 🚀 **现代技术栈** - Go 1.21+、Gin、GORM、Vue 3、Vite 5、TypeScript
- 🔐 **完善的权限** - RBAC 权限模型、JWT 认证、接口级权限控制
- 📦 **开箱即用** - 内置用户管理、角色管理、菜单管理、配置管理
- 🔧 **代码生成** - 可视化代码生成器，快速生成 CRUD 代码
- 🌐 **国际化** - 前后端完整的多语言支持
- 🔌 **消息队列** - RabbitMQ 集成，支持延迟消息
- ⏰ **定时任务** - 分布式定时任务，支持 Cron 表达式
- 🌐 **WebSocket** - 实时通信支持

## 📦 内置功能

### 后端功能

| 模块 | 功能 |
|------|------|
| 用户认证 | JWT 登录/注销、短信登录、密码重置 |
| 权限管理 | RBAC 模型、角色管理、菜单管理 |
| 系统管理 | 配置管理、附件管理、代码生成 |
| 消息队列 | RabbitMQ 普通队列、延迟队列 |
| 定时任务 | Cron 表达式、分布式锁 |
| WebSocket | 实时通信、消息推送 |
| 日志系统 | Zap 日志、日志轮转 |
| 多语言 | i18n 国际化支持 |
| HTTP 客户端 | 链式调用、中间件支持 |

### 前端功能

| 模块 | 功能 |
|------|------|
| 登录注册 | 账号密码、手机验证码登录 |
| 布局系统 | 多种布局模式、主题切换 |
| 权限控制 | 路由权限、按钮权限 |
| 自定义组件 | 富文本编辑器、附件选择器、表格选择器 |

## 🔧 环境要求

### 后端环境

| 依赖 | 版本 | 说明 |
|------|------|------|
| Go | >= 1.21 | 编程语言 |
| MySQL | >= 5.7 / 8.0 | 主数据库（或 PostgreSQL >= 12） |
| Redis | >= 6.0 | 缓存、队列 |
| RabbitMQ | >= 3.8 | 消息队列（可选） |

### 前端环境

| 依赖 | 版本 | 说明 |
|------|------|------|
| Node.js | >= 20.10.0 | 运行环境 |
| pnpm | >= 9.0.0 | 包管理器 |

## 🚀 快速开始

### 1. 克隆项目

```bash
git clone <repository-url>
cd go-wheel
```

### 2. 配置后端

```bash
# 复制配置文件
cp configs/config.example.yaml configs/config.yaml

# 编辑配置文件，修改数据库、Redis 等配置
vim configs/config.yaml

# 导入数据库
mysql -u root -p your_database < go-admin.sql
```

### 3. 启动后端

```bash
# 安装依赖
go mod download

# 开发模式（热重载）
make dev

# 或者直接运行
go run main.go
```

### 4. 启动前端

```bash
cd vben-admin

# 安装依赖
pnpm install

# 启动开发服务器
pnpm dev:user    # 用户端
pnpm dev:admin   # 管理后台
```

## 📁 目录结构

```
.
├── app/                    # 应用代码
│   ├── api/               # API 接口（用户端）
│   ├── backend/           # 后台管理接口
│   ├── common/            # 公共代码
│   ├── cron/              # 定时任务
│   ├── queue/             # 消息队列消费者
│   ├── views/             # 视图模板
│   └── websocket/         # WebSocket 控制器
├── configs/               # 配置文件
│   ├── config.yaml        # 主配置
│   └── i18n/              # 多语言文件
├── pkg/                   # 公共包
│   ├── captcha/           # 验证码
│   ├── container/         # 依赖注入容器
│   ├── cron/              # 定时任务管理
│   ├── httpclient/        # HTTP 客户端
│   ├── jwt/               # JWT 工具
│   ├── middleware/        # 中间件
│   ├── queue/             # 队列管理
│   ├── utils/             # 工具函数
│   └── websocket/         # WebSocket
├── routes/                # 路由定义
├── vben-admin/            # 前端项目
├── Makefile               # Make 命令
└── main.go                # 入口文件
```

## 📖 文档

详细文档请查看 [docs](./docs) 目录。

## 📄 许可证

[MIT License](./LICENSE)

