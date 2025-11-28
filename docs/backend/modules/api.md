# API 接口模块

`app/api/v1` 目录包含前端用户 API 接口，如登录、注册、个人中心等。

## 目录结构

```
app/api/v1/
├── controller/           # 控制器
│   ├── captcha.go       # 验证码
│   ├── index.go         # 首页
│   ├── menu.go          # 菜单
│   ├── sms.go           # 短信
│   └── user.go          # 用户
├── middleware/           # 中间件
│   └── auth.go          # 认证中间件
├── service/              # 服务层
│   └── user.go          # 用户服务
└── validator/            # 验证器
    └── user.go          # 用户验证
```

## 控制器示例

### 创建控制器

```go
// app/api/v1/controller/user.go
package controller

import (
    "github.com/gin-gonic/gin"
    "github.com/lvjiaben/go-wheel/pkg/container"
    "github.com/lvjiaben/go-wheel/pkg/utils/http"
)

type UserController struct {
    container *container.Container
    service   *service.UserService
}

func NewUserController(c *container.Container) *UserController {
    return &UserController{
        container: c,
        service:   service.NewUserService(c),
    }
}

func (c *UserController) Login(ctx *gin.Context) {
    var req validator.LoginRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        http.ErrorWithI18n(ctx, "common.param_error", nil)
        return
    }
    
    user, token, err := c.service.Login(req.Username, req.Password)
    if err != nil {
        http.ErrorWithI18n(ctx, "api.user.login_failed", nil)
        return
    }
    
    http.SuccessWithI18n(ctx, "api.user.login_success", gin.H{
        "token": token,
        "user":  user,
    })
}
```

### 服务层

```go
// app/api/v1/service/user.go
package service

type UserService struct {
    container *container.Container
}

func NewUserService(c *container.Container) *UserService {
    return &UserService{container: c}
}

func (s *UserService) Login(username, password string) (*model.User, string, error) {
    var user model.User
    if err := s.container.GetDB().Where("username = ?", username).First(&user).Error; err != nil {
        return nil, "", err
    }
    
    // 验证密码
    if !utils.CheckPassword(password, user.Password) {
        return nil, "", errors.New("密码错误")
    }
    
    // 生成 Token
    token, err := jwt.GenerateToken(user.Id, user.Username, s.container.GetConfig().Jwt.Secret)
    
    return &user, token, err
}
```

### 验证器

```go
// app/api/v1/validator/user.go
package validator

type LoginRequest struct {
    Username string `json:"username" binding:"required"`
    Password string `json:"password" binding:"required"`
}

type RegisterRequest struct {
    Username string `json:"username" binding:"required,min=3,max=20"`
    Password string `json:"password" binding:"required,min=6"`
    Mobile   string `json:"mobile" binding:"required,len=11"`
    Code     string `json:"code" binding:"required,len=6"`
}
```

## 中间件

### JWT 认证

```go
// app/api/v1/middleware/auth.go
func (m *AuthMiddleware) JWTAuthCheck() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        // 验证 Token...
        c.Set("user_id", userId)
        c.Next()
    }
}

// 可选认证（登录和未登录都可访问）
func (m *AuthMiddleware) OptionalJWTAuth() gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token != "" {
            // 尝试解析 Token
            // 成功则设置用户信息
        }
        c.Next()
    }
}
```

## 路由注册

```go
// routes/routes.go
func registerApiRoutes(r *gin.Engine, c *container.Container) {
    userController := v1.NewUserController(c)
    authMiddleware := apiMiddleware.NewAuthMiddleware(c)
    
    api := r.Group("/api")
    {
        // 无需登录
        api.POST("/user/login", userController.Login)
        api.POST("/user/reg", userController.Reg)
        
        // 需要登录
        userGroup := api.Group("/user", authMiddleware.JWTAuthCheck())
        {
            userGroup.POST("/logout", userController.Logout)
            userGroup.GET("/info", userController.Info)
        }
    }
}
```

## 响应格式

```json
{
    "code": 200,
    "message": "操作成功",
    "data": {
        "token": "eyJhbGciOiJIUzI1NiIs...",
        "user": {
            "id": 1,
            "username": "user1"
        }
    }
}
```

