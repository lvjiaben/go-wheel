# 目录结构

## 项目结构

```
.
├── app/                        # 应用代码
│   ├── api/                   # 前端 API 接口
│   │   └── v1/               # API v1 版本
│   │       ├── controller/   # 控制器
│   │       ├── middleware/   # 中间件
│   │       ├── service/      # 服务层
│   │       └── validator/    # 验证器
│   │
│   ├── backend/              # 后台管理接口
│   │   ├── controller/       # 控制器
│   │   │   ├── admin/       # 管理员相关
│   │   │   ├── system/      # 系统管理
│   │   │   └── user/        # 用户管理
│   │   ├── middleware/       # 后台中间件
│   │   └── service/          # 服务层
│   │
│   ├── common/               # 公共代码
│   │   ├── model/           # 数据模型
│   │   └── service/         # 公共服务
│   │
│   ├── cron/                 # 定时任务
│   │   └── tasks/           # 任务定义
│   │
│   ├── queue/                # 消息队列消费者
│   │   └── consumers/       # 消费者定义
│   │
│   ├── views/                # HTML 模板
│   │
│   └── websocket/            # WebSocket
│       └── controller/       # WebSocket 控制器
│
├── configs/                   # 配置文件
│   ├── config.yaml           # 主配置文件
│   └── i18n/                 # 多语言文件
│       ├── zh-CN.yaml       # 中文
│       └── en-US.yaml       # 英文
│
├── pkg/                       # 公共包
│   ├── captcha/              # 验证码服务
│   ├── config/               # 配置结构体
│   ├── container/            # 依赖注入容器
│   ├── cron/                 # 定时任务管理器
│   ├── httpclient/           # HTTP 客户端
│   ├── jwt/                  # JWT 工具
│   ├── middleware/           # Gin 中间件
│   ├── queue/                # 消息队列管理器
│   ├── utils/                # 工具函数
│   └── websocket/            # WebSocket Hub
│
├── routes/                    # 路由定义
│   └── routes.go             # 路由注册
│
├── storage/                   # 存储目录
│   ├── logs/                 # 日志文件
│   └── uploads/              # 上传文件
│
├── vben-admin/               # 前端项目
│
├── .air.toml                 # Air 热重载配置
├── go.mod                    # Go 模块定义
├── go.sum                    # 依赖版本锁定
├── main.go                   # 入口文件
├── Makefile                  # Make 命令
└── go-admin.sql              # 数据库初始化脚本
```

## 目录说明

### app 目录

应用业务代码，按功能模块划分：

| 目录 | 说明 |
|------|------|
| `api/v1` | 前端用户 API，如登录、注册、个人中心 |
| `backend` | 后台管理 API，如用户管理、角色管理 |
| `common` | 公共代码，如数据模型、公共服务 |
| `cron` | 定时任务定义 |
| `queue` | 消息队列消费者 |
| `views` | HTML 模板文件 |
| `websocket` | WebSocket 控制器 |

### pkg 目录

可复用的公共包：

| 目录 | 说明 |
|------|------|
| `captcha` | 图形验证码生成 |
| `config` | 配置结构体定义 |
| `container` | 依赖注入容器，管理所有服务 |
| `cron` | 定时任务管理器 |
| `httpclient` | HTTP 客户端封装 |
| `jwt` | JWT Token 工具 |
| `middleware` | Gin 中间件 |
| `queue` | RabbitMQ 消息队列 |
| `utils` | 通用工具函数 |
| `websocket` | WebSocket Hub |

### configs 目录

配置文件：

| 文件 | 说明 |
|------|------|
| `config.yaml` | 主配置文件 |
| `i18n/*.yaml` | 多语言翻译文件 |

## 代码分层

项目采用 MVC + Service 分层架构：

```
Controller → Service → Model
     ↓          ↓
 Validator   Container
```

- **Controller**: 处理 HTTP 请求，参数验证，调用 Service
- **Service**: 业务逻辑处理
- **Model**: 数据模型定义
- **Validator**: 请求参数验证规则
- **Container**: 依赖注入，管理数据库、Redis、日志等服务

