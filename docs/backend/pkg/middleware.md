# Middleware 中间件

项目提供多个通用中间件，位于 `pkg/middleware` 目录。

## 中间件列表

| 中间件 | 文件 | 说明 |
|--------|------|------|
| CorsMiddleware | cors.go | CORS 跨域支持 |
| GinLogger | gin.go | 请求日志记录 |
| GinRecovery | gin.go | Panic 恢复 |
| I18nMiddleware | i18n.go | 多语言支持 |
| ContainerMiddleware | container.go | 容器注入 |
| RequestBodyLimitMiddleware | request_limit.go | 请求体大小限制 |
| WebSocketAuthMiddleware | websocket.go | WebSocket 认证 |

## CORS 跨域中间件

处理跨域请求，支持预检请求（OPTIONS）：

```go
// 使用（应放在最前面）
r.Use(middleware.CorsMiddleware())
```

支持的 HTTP 头：
- `Access-Control-Allow-Origin`: 动态匹配请求的 Origin
- `Access-Control-Allow-Methods`: GET, POST, PUT, DELETE, PATCH, OPTIONS
- `Access-Control-Allow-Headers`: Content-Type, Authorization, X-Token 等
- `Access-Control-Allow-Credentials`: true
- `Access-Control-Max-Age`: 12 小时

::: tip 提示
CORS 中间件在后端处理跨域，无需在 Nginx 中重复配置。
:::

## GinLogger 日志中间件

记录每个请求的详细信息：

```go
// 使用
r.Use(middleware.GinLogger(container))
```

日志输出：
```
INFO  [GIN] method=GET status=200 cost=1.2ms ip=127.0.0.1 path=/api/users
ERROR [GIN] method=POST status=500 cost=5ms ip=127.0.0.1 path=/api/login error=...
```

## GinRecovery 恢复中间件

捕获 Panic 并返回 500 错误：

```go
r.Use(middleware.GinRecovery(container))
```

## I18n 多语言中间件

根据 `Accept-Language` 请求头设置语言：

```go
r.Use(middleware.I18nMiddleware(container))
```

在控制器中获取：
```go
lang := ctx.GetString("lang")     // "zh-CN" 或 "en-US"
isCn := ctx.GetBool("isCn")       // true 或 false
```

## Container 容器中间件

将容器注入到请求上下文：

```go
r.Use(middleware.ContainerMiddleware(container))

// 在控制器中获取
c := ctx.MustGet("container").(*container.Container)
```

## RequestBodyLimit 请求体限制

限制请求体大小，防止大文件攻击：

```go
r.Use(middleware.RequestBodyLimitMiddleware(container))
```

配置：
```yaml
app:
  max_request_body: 10  # 最大 10MB
```

## WebSocket 认证中间件

WebSocket 连接的 JWT 认证：

```go
ws := r.Group("/ws").Use(
    middleware.WebSocketAuthMiddleware(container),
    middleware.WebSocketUpgradeMiddleware(wsHub),
)
```

## 自定义中间件

### 创建中间件

```go
// app/backend/middleware/custom.go
func CustomMiddleware(c *container.Container) gin.HandlerFunc {
    return func(ctx *gin.Context) {
        // 请求前处理
        start := time.Now()
        
        ctx.Next()
        
        // 请求后处理
        duration := time.Since(start)
        c.GetLogger().Info("请求耗时", zap.Duration("duration", duration))
    }
}
```

### 认证中间件示例

```go
func AuthMiddleware(c *container.Container) gin.HandlerFunc {
    return func(ctx *gin.Context) {
        token := ctx.GetHeader("Authorization")
        if token == "" {
            ctx.JSON(401, gin.H{"error": "未授权"})
            ctx.Abort()
            return
        }
        
        // 验证 Token...
        
        ctx.Next()
    }
}
```

### 权限中间件示例

```go
func PermissionMiddleware(c *container.Container) gin.HandlerFunc {
    return func(ctx *gin.Context) {
        adminId := ctx.GetInt("admin_id")
        path := ctx.Request.URL.Path
        method := ctx.Request.Method
        
        // 检查权限
        if !hasPermission(adminId, path, method) {
            ctx.JSON(403, gin.H{"error": "无权限"})
            ctx.Abort()
            return
        }
        
        ctx.Next()
    }
}
```

## 中间件顺序

推荐的中间件顺序：

```go
r.Use(
    middleware.CorsMiddleware(),              // 1. CORS（最外层，处理预检请求）
    middleware.GinRecovery(c),                // 2. 恢复
    middleware.GinLogger(c),                  // 3. 日志
    middleware.ContainerMiddleware(c),        // 4. 容器注入
    middleware.I18nMiddleware(c),             // 5. 多语言
    middleware.RequestBodyLimitMiddleware(c), // 6. 请求限制
)
```

::: warning 注意
CORS 中间件必须放在最前面，否则预检请求（OPTIONS）可能被其他中间件拦截。
:::

## 跳过中间件

某些路由跳过中间件：

```go
// 方式1：分组
public := r.Group("/public")  // 无认证
private := r.Group("/api", authMiddleware)  // 需认证

// 方式2：路径判断
func AuthMiddleware() gin.HandlerFunc {
    return func(ctx *gin.Context) {
        // 跳过公开路由
        if strings.HasPrefix(ctx.Request.URL.Path, "/public") {
            ctx.Next()
            return
        }
        // 认证逻辑...
    }
}
```

