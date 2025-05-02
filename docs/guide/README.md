# Go Admin 框架用户指南

<p align="center">
  <a href="../../README.md">
    <img src="https://img.shields.io/badge/返回-首页-blue.svg" alt="返回首页">
  </a>
  <a href="./README_EN.md">
    <img src="https://img.shields.io/badge/Switch-English-blue.svg" alt="English">
  </a>
</p>

欢迎使用 Go Admin 框架！本指南将帮助您了解如何使用框架的各种功能和组件。

## 目录

- [快速入门](#快速入门)
- [基础概念](#基础概念)
- [组件使用](./components.md)
  - [Redis 缓存](./components/redis.md)
  - [RabbitMQ 消息队列](./components/rabbitmq.md)
  - [Cron 定时任务](./components/cron.md)
  - [健康检查](./components/health.md)
- [最佳实践](./best_practices.md)
- [配置说明](./configuration.md)
- [常见问题](./faq.md)

## 快速入门

### 安装

克隆仓库并安装依赖：

```bash
git clone https://github.com/yourusername/go-admin.git
cd go-admin
go mod tidy
```

### 配置

复制示例配置文件并修改：

```bash
cp configs/config.example.yaml configs/config.yaml
```

编辑 `config.yaml` 文件，设置数据库、Redis 等连接信息。

### 启动服务

```bash
go run main.go
```

访问 [http://localhost:8080](http://localhost:8080) 即可看到管理界面。

## 基础概念

Go Admin 框架采用了清晰的分层架构：

1. **Controller 层**：处理 HTTP 请求，进行参数验证，调用 Service 层
2. **Service 层**：实现业务逻辑，调用 Model 层和其他服务
3. **Model 层**：定义数据结构和数据库交互

### 目录结构说明

```
app/
├── api/              # API 控制器
│   ├── v1/           # API 版本 1
│   └── v2/           # API 版本 2
└── backend/          # 后端服务
    ├── controller/   # 控制器
    ├── service/      # 服务层
    └── model/        # 数据模型
```

### 核心组件

Go Admin 框架集成了多种常用组件，这些组件已经在框架启动时完成初始化，您只需要在需要的地方注入使用：

- **RedisService**：Redis 缓存服务
- **MQService**：RabbitMQ 消息队列服务
- **CronService**：定时任务服务
- **HealthService**：健康检查服务

所有这些服务都是通过依赖注入的方式使用，无需手动初始化。详细使用方法请参考[组件使用文档](./components.md)。

## 开发流程

### 1. 创建新的 API

在 `app/api/v1` 目录下创建新的控制器：

```go
// app/api/v1/user_controller.go
package v1

import (
    "github.com/gin-gonic/gin"
    "github.com/lvjiaben/go-wheel/app/backend/service"
)

type UserController struct {
    userService *service.UserService
}

func NewUserController(userService *service.UserService) *UserController {
    return &UserController{
        userService: userService,
    }
}

// GetUser 获取用户信息
// @Summary 获取用户信息
// @Description 根据用户ID获取用户信息
// @Tags 用户管理
// @Accept json
// @Produce json
// @Param id path string true "用户ID"
// @Success 200 {object} model.User
// @Router /api/v1/users/{id} [get]
func (c *UserController) GetUser(ctx *gin.Context) {
    id := ctx.Param("id")
    user, err := c.userService.GetUserByID(id)
    if err != nil {
        ctx.JSON(400, gin.H{"error": err.Error()})
        return
    }
    ctx.JSON(200, user)
}
```

### 2. 实现业务逻辑

在 `app/backend/service` 目录下创建相应的服务：

```go
// app/backend/service/user_service.go
package service

import (
    "encoding/json"
    "time"
    "github.com/lvjiaben/go-wheel/app/backend/model"
)

type UserService struct {
    redisService *RedisService
}

func NewUserService(redisService *RedisService) *UserService {
    return &UserService{
        redisService: redisService,
    }
}

func (s *UserService) GetUserByID(id string) (*model.User, error) {
    // 尝试从缓存获取
    cacheKey := "user:" + id
    userData, err := s.redisService.Get(cacheKey)
    
    // 缓存不存在，从数据库获取
    if err != nil {
        // 从数据库获取用户数据
        user := &model.User{
            ID: id,
            // ... 其他字段
        }
        
        // 将用户数据存入缓存
        jsonData, _ := json.Marshal(user)
        s.redisService.Set(cacheKey, string(jsonData), 30*time.Minute)
        
        return user, nil
    }
    
    // 解析缓存数据
    var user model.User
    json.Unmarshal([]byte(userData), &user)
    
    return &user, nil
}
```

### 3. 注册路由

在 `routes` 目录下注册 API 路由：

```go
// routes/api.go
package routes

import (
    "github.com/gin-gonic/gin"
    v1 "github.com/lvjiaben/go-wheel/app/api/v1"
)

func RegisterAPIRoutes(r *gin.Engine, userController *v1.UserController) {
    api := r.Group("/api")
    {
        v1Group := api.Group("/v1")
        {
            users := v1Group.Group("/users")
            {
                users.GET("/:id", userController.GetUser)
                // 其他用户相关路由
            }
        }
    }
}
```

## 下一步

- 查看[组件使用文档](./components.md)，了解如何使用框架提供的各种组件
- 查看[配置说明](./configuration.md)，了解如何配置框架
- 查看[最佳实践](./best_practices.md)，了解如何更好地使用框架

如果您有任何问题，请查看[常见问题](./faq.md)或提交 Issue。 