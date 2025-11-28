# JWT 认证

项目使用 [golang-jwt](https://github.com/golang-jwt/jwt) 实现 JWT 认证。

## 配置

在 `configs/config.yaml` 中配置：

```yaml
jwt:
  secret: "your-secret-key"   # JWT 密钥（生产环境请修改）
  expire_day: 7               # 过期天数
  issuer: "go-admin"          # 签发者
```

## 基本用法

### 生成 Token

```go
import "github.com/lvjiaben/go-wheel/pkg/jwt"

// 生成 Token
token, err := jwt.GenerateToken(
    userId,           // 用户ID
    username,         // 用户名
    config.Jwt.Secret, // 密钥
    config.Jwt.ExpireDay, // 过期天数
)
```

### 解析 Token

```go
// 解析 Token
claims, err := jwt.ParseToken(tokenString, config.Jwt.Secret)
if err != nil {
    // Token 无效或已过期
}

// 获取用户信息
userId := claims.Id
username := claims.Username
```

## Claims 结构

```go
// pkg/jwt/jwt.go
type CustomClaims struct {
    Id       int    `json:"id"`       // 用户ID
    Username string `json:"username"` // 用户名
    jwt.RegisteredClaims              // 标准声明
}
```

标准声明包含：
- `ExpiresAt` - 过期时间
- `IssuedAt` - 签发时间
- `NotBefore` - 生效时间

## 在中间件中使用

```go
// app/api/v1/middleware/auth.go
func (m *AuthMiddleware) JWTAuthCheck() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 获取 Token
        token := c.GetHeader("Authorization")
        if token == "" {
            http.ErrorWithI18n(c, "common.unauthorized", nil)
            c.Abort()
            return
        }
        
        // 移除 Bearer 前缀
        if strings.HasPrefix(token, "Bearer ") {
            token = token[7:]
        }
        
        // 解析 Token
        claims, err := jwt.ParseToken(token, m.container.GetConfig().Jwt.Secret)
        if err != nil {
            http.ErrorWithI18n(c, "common.token_invalid", nil)
            c.Abort()
            return
        }
        
        // 设置用户信息到上下文
        c.Set("user_id", claims.Id)
        c.Set("username", claims.Username)
        
        c.Next()
    }
}
```

## 可选认证

允许未登录用户访问，但登录用户可获取更多信息：

```go
func (m *AuthMiddleware) OptionalJWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token != "" {
            if strings.HasPrefix(token, "Bearer ") {
                token = token[7:]
            }
            
            claims, err := jwt.ParseToken(token, m.container.GetConfig().Jwt.Secret)
            if err == nil {
                c.Set("user_id", claims.Id)
                c.Set("username", claims.Username)
            }
        }
        
        c.Next()
    }
}
```

## 在控制器中获取用户

```go
func (c *UserController) Info(ctx *gin.Context) {
    userId := ctx.GetInt("user_id")
    username := ctx.GetString("username")
    
    // 使用用户信息...
}
```

## 单点登录 (SSO)

配置单点登录，同一账号只能在一处登录：

```yaml
api:
  login_sso: true
```

实现方式：
1. 登录时将 Token 存入 Redis
2. 验证时检查 Token 是否与 Redis 中一致
3. 新登录会使旧 Token 失效

```go
// 登录时存储 Token
redis.Set(ctx, fmt.Sprintf("user:token:%d", userId), token, expireTime)

// 验证时检查
storedToken, _ := redis.Get(ctx, fmt.Sprintf("user:token:%d", userId))
if storedToken != currentToken {
    // Token 已被其他登录替换
}
```

## 刷新 Token

```go
// 检查 Token 是否即将过期（如 1 天内）
if claims.ExpiresAt.Time.Sub(time.Now()) < 24*time.Hour {
    // 生成新 Token
    newToken, _ := jwt.GenerateToken(claims.Id, claims.Username, secret, expireDays)
    // 返回新 Token 给前端
}
```

## 安全建议

1. **密钥安全** - 生产环境使用强密钥，不要提交到代码库
2. **HTTPS** - 生产环境必须使用 HTTPS
3. **合理过期** - 根据安全需求设置过期时间
4. **黑名单** - 实现 Token 黑名单用于强制登出

