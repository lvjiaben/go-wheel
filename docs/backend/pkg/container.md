# Container 容器

Container 是项目的核心依赖注入容器，统一管理所有组件的生命周期。

## 概述

Container 管理以下组件：

| 组件 | 说明 | 获取方法 |
|------|------|----------|
| Config | 配置管理 | `GetConfig()` |
| Logger | 日志记录 | `GetLogger()` |
| DB | 数据库连接 | `GetDB()` |
| Redis | Redis 客户端 | `GetRedis()` |
| Cron | 定时任务 | `GetCron()` |
| MessageQueue | 消息队列 | `GetMessageQueue()` |
| DelayQueue | 延迟队列 | `GetDelayQueue()` |
| HTTPClient | HTTP 客户端 | `GetHTTPClient()` |
| WebSocketHub | WebSocket 管理 | `GetWebSocketHub()` |
| I18n | 多语言 | `GetI18n()` |
| Validator | 验证器 | `GetValidator()` |

## 初始化

```go
// main.go
func main() {
    // 创建容器
    c := container.NewContainer()
    
    // 优雅关闭
    defer c.Shutdown()
    
    // 注册路由
    r := gin.Default()
    routes.RegisterRoutes(r, c)
    
    // 启动服务
    r.Run(fmt.Sprintf(":%d", c.GetConfig().App.Port))
}
```

## 使用组件

### 数据库

```go
db := container.GetDB()

// 查询
var user model.User
db.Where("id = ?", 1).First(&user)

// 事务
db.Transaction(func(tx *gorm.DB) error {
    // 事务操作
    return nil
})
```

### Redis

```go
redis := container.GetRedis()
ctx := context.Background()

// 设置
redis.Set(ctx, "key", "value", time.Hour)

// 获取
val, _ := redis.Get(ctx, "key")

// 删除
redis.Del(ctx, "key")
```

### 日志

```go
logger := container.GetLogger()

logger.Info("操作成功", zap.Int("user_id", 1))
logger.Error("操作失败", zap.Error(err))
logger.Debug("调试信息", zap.Any("data", data))
```

### HTTP 客户端

```go
client := container.GetHTTPClient()

// GET 请求
resp, err := client.Get("https://api.example.com/users")

// POST 请求
resp, err := client.PostJSON("https://api.example.com/users", map[string]interface{}{
    "name": "张三",
})
```

### 多语言

```go
i18n := container.GetI18n()

// 获取翻译
msg := i18n.T("zh-CN", "common.success")

// 带参数
msg := i18n.T("zh-CN", "user.welcome", map[string]interface{}{
    "name": "张三",
})
```

## 熔断器

Container 内置数据库和 Redis 熔断器，防止级联故障：

```go
// 熔断器配置
cb := container.NewCircuitBreaker(
    5,              // 失败阈值
    30*time.Second, // 重置超时
)

// 执行操作
err := cb.Execute(func() error {
    return db.Ping()
})
```

### 熔断状态

| 状态 | 说明 |
|------|------|
| Closed | 正常状态，请求正常通过 |
| Open | 熔断状态，请求直接失败 |
| Half-Open | 半开状态，允许部分请求测试 |

## 组件状态

```go
// 获取组件状态
status := container.GetStatus()

for name, s := range status {
    fmt.Printf("%s: %s - %s\n", name, s.Status, s.Message)
}
```

## 优雅关闭

```go
// 关闭容器，释放所有资源
container.Shutdown()
```

关闭顺序：
1. 停止定时任务
2. 关闭消息队列
3. 关闭 WebSocket 连接
4. 关闭数据库连接
5. 关闭 Redis 连接

## 自定义数据

```go
// 存储自定义数据
container.SetCustomData("myKey", myValue)

// 获取自定义数据
value := container.GetCustomData("myKey")
```

## 在中间件中使用

```go
// 注入容器到上下文
func ContainerMiddleware(c *container.Container) gin.HandlerFunc {
    return func(ctx *gin.Context) {
        ctx.Set("container", c)
        ctx.Next()
    }
}

// 在控制器中获取
func (c *Controller) Handler(ctx *gin.Context) {
    container := ctx.MustGet("container").(*container.Container)
    db := container.GetDB()
}
```

