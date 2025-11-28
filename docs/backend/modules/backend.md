# Backend 后台模块

`app/backend` 目录包含后台管理 API 接口，如用户管理、角色管理、菜单管理等。

## 目录结构

```
app/backend/
├── controller/              # 控制器
│   ├── admin/              # 管理员相关
│   │   ├── admin.go       # 管理员管理
│   │   ├── menu.go        # 菜单管理
│   │   └── role.go        # 角色管理
│   ├── system/             # 系统管理
│   │   ├── attachment.go  # 附件管理
│   │   ├── config.go      # 配置管理
│   │   └── gen.go         # 代码生成
│   ├── user/               # 用户管理
│   │   └── user.go
│   ├── auth.go             # 认证控制器
│   ├── common.go           # 公共控制器
│   └── home.go             # 首页控制器
├── middleware/              # 中间件
│   └── auth.go             # 认证中间件
├── model/                   # 数据模型
│   └── admin/
│       ├── admin.go
│       ├── menu.go
│       └── role.go
├── service/                 # 服务层
│   ├── admin.go
│   ├── menu.go
│   └── role.go
└── utils/                   # 工具函数
    └── auth.go
```

## 控制器示例

### 管理员管理

```go
// app/backend/controller/admin/admin.go
type AdminController struct {
    container *container.Container
    service   *service.AdminService
}

func NewAdminController(c *container.Container) *AdminController {
    return &AdminController{
        container: c,
        service:   service.NewAdminService(c),
    }
}

// List 管理员列表
func (c *AdminController) List(ctx *gin.Context) {
    page, _ := strconv.Atoi(ctx.DefaultQuery("page", "1"))
    pageSize, _ := strconv.Atoi(ctx.DefaultQuery("page_size", "10"))
    
    list, total, err := c.service.List(page, pageSize)
    if err != nil {
        http.ErrorWithI18n(ctx, "backend.admin.list_failed", nil)
        return
    }
    
    http.SuccessWithI18n(ctx, "common.success", gin.H{
        "list":  list,
        "total": total,
    })
}

// Save 保存管理员
func (c *AdminController) Save(ctx *gin.Context) {
    var req SaveAdminRequest
    if err := ctx.ShouldBindJSON(&req); err != nil {
        http.ErrorWithI18n(ctx, "common.param_error", nil)
        return
    }
    
    if req.Id > 0 {
        err = c.service.Update(&req)
    } else {
        err = c.service.Create(&req)
    }
    
    if err != nil {
        http.ErrorWithI18n(ctx, "backend.admin.save_failed", nil)
        return
    }
    
    http.SuccessWithI18n(ctx, "backend.admin.save_success", nil)
}
```

### 认证控制器

```go
// app/backend/controller/auth.go
type AuthController struct {
    container *container.Container
}

// Login 管理员登录
func (c *AuthController) Login(ctx *gin.Context) {
    var req LoginRequest
    ctx.ShouldBindJSON(&req)
    
    // 验证验证码
    if !captcha.Verify(req.CaptchaId, req.CaptchaCode) {
        http.ErrorWithI18n(ctx, "backend.auth.captcha_error", nil)
        return
    }
    
    // 验证用户名密码
    admin, err := c.service.Login(req.Username, req.Password)
    if err != nil {
        http.ErrorWithI18n(ctx, "backend.auth.login_failed", nil)
        return
    }
    
    // 生成 Token
    token, _ := jwt.GenerateToken(admin.Id, admin.Username, secret)
    
    http.SuccessWithI18n(ctx, "backend.auth.login_success", gin.H{
        "token": token,
    })
}

// Info 获取用户信息
func (c *AuthController) Info(ctx *gin.Context) {
    adminId := ctx.GetInt("admin_id")
    admin, _ := c.service.GetById(adminId)
    
    http.SuccessWithI18n(ctx, "common.success", gin.H{
        "id":       admin.Id,
        "username": admin.Username,
        "avatar":   admin.Avatar,
        "realName": admin.Nickname,
    })
}

// Menus 获取菜单列表
func (c *AuthController) Menus(ctx *gin.Context) {
    adminId := ctx.GetInt("admin_id")
    menus, _ := c.menuService.GetMenusByAdminId(adminId)
    
    http.SuccessWithI18n(ctx, "common.success", menus)
}
```

## 权限控制

后台接口使用双重中间件保护：

```go
// 1. JWT 认证 - 验证用户身份
authMiddleware.JWTAuthCheck()

// 2. 权限检查 - 验证接口权限
authMiddleware.PermissionCheck()
```

## 路由注册

```go
// routes/routes.go
func registerBackendRoutes(r *gin.Engine, c *container.Container) {
    authMiddleware := backendMiddleware.NewAuthMiddleware(c)
    
    api := r.Group("/backend")
    {
        // 无需认证
        api.POST("/common/captcha", commonController.Captcha)
        api.POST("/auth/login", authController.Login)
        
        // 需要认证和权限
        adminGroup := api.Group("/admin", 
            authMiddleware.JWTAuthCheck(), 
            authMiddleware.PermissionCheck(),
        )
        {
            adminGroup.GET("/admin/list", adminController.List)
            adminGroup.POST("/admin/save", adminController.Save)
        }
    }
}
```

