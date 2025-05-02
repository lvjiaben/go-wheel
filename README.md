# Go Admin Framework | Go 后台管理框架

一个高性能、易扩展的Golang后台管理框架，提供完整的后台管理功能和丰富的组件支持。专为现代化Web应用设计，适用于企业级后台系统、API服务、微服务架构等场景。

<p align="center">
  <a href="./docs/guide/README.md">
    <img src="https://img.shields.io/badge/文档-查看文档-blue.svg" alt="查看文档">
  </a>
  <a href="./docs/guide/README_EN.md">
    <img src="https://img.shields.io/badge/Document-English-blue.svg" alt="English">
  </a>
  <a href="./LICENSE">
    <img src="https://img.shields.io/badge/License-MIT-yellow.svg" alt="License">
  </a>
  <a href="https://golang.org">
    <img src="https://img.shields.io/badge/Golang-1.23+-orange.svg" alt="Golang">
  </a>
</p>

[English](./README_EN.md) | 简体中文

## 功能特性

- **认证授权**：完整的JWT认证与RBAC权限管理
- **依赖注入**：清晰的依赖注入体系，简化组件使用
- **多语言支持**：基于i18n的国际化支持，方便全球化应用
- **消息队列**：集成RabbitMQ支持普通队列和延迟队列，高效处理异步任务
- **定时任务**：基于cron的定时任务管理器，精确控制任务执行
- **缓存支持**：Redis缓存管理，支持多种数据结构与分布式锁
- **健康检查**：服务健康状态监控与自动恢复
- **日志记录**：结构化日志记录与管理，支持多级别输出
- **参数验证**：强大的请求参数验证，确保数据安全
- **API文档**：Swagger自动生成API文档，简化接口对接
- **数据库**：集成Gorm支持多种数据库，优化数据查询性能
- **微服务支持**：内置服务发现与负载均衡，简化微服务架构
- **监控告警**：系统性能监控与故障告警机制
- **优雅关闭**：支持优雅停机，确保任务平滑过渡

## 核心优势

- **高性能**：采用Go语言高并发特性，提供卓越的性能表现
- **易扩展**：模块化设计，轻松扩展新功能
- **低耦合**：基于依赖注入，组件间松耦合
- **开箱即用**：丰富的内置组件，减少开发工作量
- **安全可靠**：内置多重安全机制，保障系统稳定运行
- **文档完善**：详尽的文档和丰富的使用示例

## 快速开始

### 1. 安装依赖
```bash
git clone https://github.com/yourusername/go-admin.git
cd go-admin
go mod tidy
```

### 2. 配置环境
修改`configs/config.yaml`配置文件，设置数据库、Redis等连接信息。

```yaml
app:
  name: Go-Admin
  mode: development
  port: 8080
  
mysql:
  host: localhost
  port: 3306
  user: root
  pass: password
  dbname: go_admin
  
redis:
  host: localhost
  port: 6379
  pass: ""
  db: 0
```

### 3. 启动服务
```bash
go run main.go
```

### 4. 访问API文档
```
http://localhost:8080/swagger/index.html
```

## 核心组件

框架集成了多种常用组件，使开发更加便捷：

### Redis缓存
用于数据缓存、计数器、分布式锁等场景，提供高效的数据存取：

```go
// 在服务中使用Redis缓存
func (s *UserService) GetUserProfile(userID string) (*User, error) {
    cacheKey := fmt.Sprintf("user:%s", userID)
    
    // 尝试从缓存获取
    user, err := s.redisCache.Get(ctx, cacheKey)
    if err == nil {
        return user, nil 
    }
    
    // 从数据库获取并缓存
    user, err = s.userRepo.FindByID(userID)
    if err != nil {
        return nil, err
    }
    
    s.redisCache.Set(ctx, cacheKey, user, 30*time.Minute)
    return user, nil
}
```

### RabbitMQ队列
用于异步处理、任务分发、事件广播等，支持普通队列和延迟队列：

```go
// 发送消息
s.messageQueue.Push(ctx, "order_created", orderData)

// 设置消费者
msgs, _ := s.messageQueue.Consume(ctx, "order_created", "order-processor")
go func() {
    for msg := range msgs {
        // 处理订单消息
        processOrder(msg.Body)
        msg.Ack(false)
    }
}()
```

### Cron定时任务
用于定期执行任务、数据统计、报表生成等：

```go
// 添加定时任务
s.cronManager.AddTask("daily-report", "0 0 0 * * *", func() error {
    return s.generateDailyReport()
})

// 管理任务
s.cronManager.RemoveTask("task-id")
```

### 健康检查
监控服务健康状态，支持自定义检查策略，确保系统可靠运行：

```go
// 注册健康检查
s.healthCheck.Register(&DatabaseChecker{
    db: s.db,
})

// 获取健康状态
status := s.healthCheck.Check(ctx)
```

<p align="center">
  <a href="./docs/guide/components.md">
    <img src="https://img.shields.io/badge/查看-组件文档-green.svg" alt="查看组件文档">
  </a>
</p>

## 目录结构

```
admin/
├── app/               # 应用代码
│   ├── api/           # API控制器
│   │   ├── v1/        # API版本1
│   │   └── v2/        # API版本2
│   ├── backend/       # 后端服务
│   │   ├── controller/ # 控制器
│   │   ├── service/    # 服务层
│   │   └── model/      # 数据模型
│   └── generate/      # 代码生成
├── configs/           # 配置文件
│   └── config.yaml    # 主配置文件
├── deployments/       # 部署相关
├── docs/              # 文档
│   ├── api/           # API文档
│   └── guide/         # 使用指南
├── pkg/               # 公共包
│   ├── auth/          # 认证
│   ├── cache/         # 缓存
│   ├── container/     # 依赖注入容器
│   ├── cron/          # 定时任务
│   ├── db/            # 数据库
│   ├── logger/        # 日志
│   ├── mq/            # 消息队列
│   ├── health/        # 健康检查
│   ├── types/         # 类型定义
│   └── validator/     # 验证器
├── routes/            # 路由定义
├── scripts/           # 脚本
└── tests/             # 测试
    ├── cron/          # 定时任务测试
    ├── health/        # 健康检查测试
    ├── rabbitmq/      # 消息队列测试
    └── redis/         # Redis缓存测试
```

## 开发指南

### 添加新的API

1. 在`app/api/v1`目录下定义API控制器
```go
// UserController 用户控制器
type UserController struct {
    userService *service.UserService
}

// NewUserController 创建用户控制器
func NewUserController(userService *service.UserService) *UserController {
    return &UserController{
        userService: userService,
    }
}

// GetProfile 获取用户资料
// @Summary 获取用户资料
// @Description 获取当前登录用户的详细资料
// @Tags 用户
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} response.Response{data=model.User}
// @Router /api/v1/user/profile [get]
func (c *UserController) GetProfile(ctx *gin.Context) {
    userID := ctx.GetString("user_id")
    user, err := c.userService.GetUserProfile(userID)
    // 处理响应...
}
```

2. 在`app/backend/service`中实现业务逻辑
```go
// UserService 用户服务
type UserService struct {
    db         *gorm.DB
    redisCache *service.RedisService
}

// NewUserService 创建用户服务
func NewUserService(db *gorm.DB, redisCache *service.RedisService) *UserService {
    return &UserService{
        db:         db,
        redisCache: redisCache,
    }
}

// GetUserProfile 获取用户资料
func (s *UserService) GetUserProfile(userID string) (*model.User, error) {
    // 实现业务逻辑...
}
```

3. 在`app/backend/model`中定义数据模型
```go
// User 用户模型
type User struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    Username  string    `json:"username" gorm:"unique;not null"`
    Email     string    `json:"email" gorm:"unique;not null"`
    Password  string    `json:"-" gorm:"not null"`
    Status    int       `json:"status" gorm:"default:1"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

4. 在`routes`目录中注册路由
```go
// 注册用户相关路由
func registerUserRoutes(r *gin.RouterGroup, authMiddleware gin.HandlerFunc, userController *controller.UserController) {
    userGroup := r.Group("/user").Use(authMiddleware)
    {
        userGroup.GET("/profile", userController.GetProfile)
        userGroup.PUT("/profile", userController.UpdateProfile)
        // 其他路由...
    }
}
```

### 使用框架组件

所有核心组件都是通过依赖注入的方式使用，无需手动初始化：

```go
// 服务中使用Redis
type UserService struct {
    redisService *service.RedisService
    db           *gorm.DB
    messageQueue types.MessageQueue
}

func NewUserService(redisService *service.RedisService, db *gorm.DB, messageQueue types.MessageQueue) *UserService {
    return &UserService{
        redisService: redisService,
        db:           db,
        messageQueue: messageQueue,
    }
}

// 使用Redis缓存
func (s *UserService) GetUserByID(id string) (*User, error) {
    cacheKey := "user:" + id
    // 尝试从缓存获取
    data, err := s.redisService.Get(cacheKey)
    // ...
}
```

## 项目实例

1. **用户认证与授权**：完整的JWT认证实现
2. **商品管理系统**：包含分类、商品、库存管理
3. **订单处理流程**：订单创建、支付、发货、完成
4. **统计报表生成**：销售统计、用户分析、趋势图表

## 部署

### Docker部署
```bash
# 构建镜像
docker build -t go-admin .

# 运行容器
docker run -p 8080:8080 \
  -e MYSQL_HOST=mysql_host \
  -e REDIS_HOST=redis_host \
  -e APP_MODE=production \
  go-admin
```

### Kubernetes部署
```bash
# 应用部署配置
kubectl apply -f deployments/k8s/

# 查看部署状态
kubectl get pods -l app=go-admin
```

### 生产环境配置建议
- 使用环境变量覆盖敏感配置
- 启用TLS加密和请求限流
- 配置日志收集与监控
- 设置定期数据备份

## 性能优化

框架内置了多种性能优化措施：

- **连接池**：数据库和Redis连接池管理
- **缓存策略**：多级缓存设计，减少数据库压力
- **请求合并**：合并重复请求，降低下游服务压力
- **异步处理**：非核心流程异步处理
- **资源监控**：实时监控资源使用情况，及时发现瓶颈

## 文档中心

<p align="center">
  <a href="./docs/guide/README.md">
    <img src="https://img.shields.io/badge/用户指南-User_Guide-blue.svg" alt="用户指南">
  </a>
  <a href="./docs/api/README.md">
    <img src="https://img.shields.io/badge/API文档-API_Document-blue.svg" alt="API文档">
  </a>
  <a href="./docs/guide/components.md">
    <img src="https://img.shields.io/badge/组件文档-Components-blue.svg" alt="组件文档">
  </a>
  <a href="./docs/guide/configuration.md">
    <img src="https://img.shields.io/badge/配置说明-Configuration-blue.svg" alt="配置说明">
  </a>
</p>

## 社区与支持

- **问题反馈**：通过 [GitHub Issues](https://github.com/yourusername/go-admin/issues) 提交问题
- **功能请求**：通过 [GitHub Discussions](https://github.com/yourusername/go-admin/discussions) 讨论新功能
- **贡献代码**：欢迎提交 [Pull Request](https://github.com/yourusername/go-admin/pulls)

## 许可证

[MIT License](./LICENSE) - 您可以自由使用、修改和分发本项目的代码
