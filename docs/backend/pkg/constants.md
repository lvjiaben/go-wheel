# Constants 常量

项目在 `pkg/constants` 中定义了全局常量。

## JWT 相关

```go
const (
    // JWT 密钥
    JWTSecretKey = "jwt_secret_key"
    
    // Token 过期时间（小时）
    JWTExpireHours = 24
    
    // Refresh Token 过期时间（小时）
    JWTRefreshExpireHours = 168  // 7天
    
    // Token 类型
    TokenTypeAccess  = "access"
    TokenTypeRefresh = "refresh"
)
```

## 密码相关

```go
const (
    // 密码最小长度
    PasswordMinLength = 6
    
    // 密码最大长度
    PasswordMaxLength = 32
    
    // 密码加密成本
    PasswordBcryptCost = 10
)
```

## 验证码相关

```go
const (
    // 验证码长度
    CaptchaLength = 6
    
    // 验证码过期时间（分钟）
    CaptchaExpireMinutes = 5
    
    // 图形验证码宽度
    CaptchaImageWidth = 240
    
    // 图形验证码高度
    CaptchaImageHeight = 80
)
```

## 重试相关

```go
const (
    // 最大重试次数
    MaxRetryAttempts = 3
    
    // 重试间隔（毫秒）
    RetryIntervalMs = 1000
    
    // 重试退避倍数
    RetryBackoffMultiplier = 2.0
)
```

## 熔断器相关

```go
const (
    // 熔断器名称
    CircuitBreakerName = "default"
    
    // 熔断阈值（失败次数）
    CircuitBreakerThreshold = 5
    
    // 熔断恢复时间（秒）
    CircuitBreakerTimeout = 30
)
```

## 分页相关

```go
const (
    // 默认页码
    DefaultPage = 1
    
    // 默认每页数量
    DefaultPageSize = 20
    
    // 最大每页数量
    MaxPageSize = 100
)
```

## 缓存相关

```go
const (
    // 缓存前缀
    CachePrefix = "go_wheel:"
    
    // 用户缓存前缀
    CacheUserPrefix = "go_wheel:user:"
    
    // 权限缓存前缀
    CachePermissionPrefix = "go_wheel:permission:"
    
    // 默认缓存时间（秒）
    DefaultCacheTTL = 3600
)
```

## 使用示例

```go
import "github.com/lvjiaben/go-wheel/pkg/constants"

// 使用 JWT 过期时间
token, err := jwt.GenerateToken(userId, constants.JWTExpireHours)

// 使用分页默认值
page := constants.DefaultPage
pageSize := constants.DefaultPageSize

// 使用缓存前缀
cacheKey := constants.CacheUserPrefix + strconv.Itoa(userId)
```

## 自定义常量

如需添加新常量，在 `pkg/constants/constants.go` 中添加：

```go
const (
    // 自定义业务常量
    OrderStatusPending   = 0
    OrderStatusPaid      = 1
    OrderStatusShipped   = 2
    OrderStatusCompleted = 3
    OrderStatusCancelled = -1
)
```

## 最佳实践

1. **使用常量代替魔法数字** - 提高代码可读性
2. **按功能分组** - JWT、密码、验证码等分组定义
3. **添加注释** - 说明常量用途
4. **统一命名** - 使用大驼峰命名，前缀表示分类

