# 路由说明

路由定义在 `routes/routes.go` 文件中，按功能模块分组。

## 路由结构

```
/                           # 首页
/public/*                   # 静态资源
/backend/*                  # 后台管理 API
/api/*                      # 前端用户 API
/ws/*                       # WebSocket
```

## 后台路由 `/backend`

### 公共接口（无需认证）

| 路径 | 方法 | 说明 |
|------|------|------|
| `/backend/common/captcha` | POST | 获取图形验证码 |
| `/backend/auth/login` | POST | 管理员登录 |

### 认证接口（需要登录）

| 路径 | 方法 | 说明 |
|------|------|------|
| `/backend/auth/logout` | POST | 退出登录 |
| `/backend/auth/menus` | GET | 获取菜单列表 |
| `/backend/auth/permission` | GET | 获取权限列表 |
| `/backend/auth/profile` | POST | 修改个人信息 |
| `/backend/auth/password` | POST | 修改密码 |
| `/backend/auth/info` | GET | 获取用户信息 |

### 管理员管理（需要权限）

| 路径 | 方法 | 说明 |
|------|------|------|
| `/backend/admin/admin/list` | GET | 管理员列表 |
| `/backend/admin/admin/save` | POST | 保存管理员 |
| `/backend/admin/admin/delete/:id` | DELETE | 删除管理员 |

### 角色管理

| 路径 | 方法 | 说明 |
|------|------|------|
| `/backend/admin/role/list` | GET | 角色列表 |
| `/backend/admin/role/save` | POST | 保存角色 |
| `/backend/admin/role/delete/:id` | DELETE | 删除角色 |

### 菜单管理

| 路径 | 方法 | 说明 |
|------|------|------|
| `/backend/admin/menu/list` | GET | 菜单列表 |
| `/backend/admin/menu/save` | POST | 保存菜单 |
| `/backend/admin/menu/delete` | POST | 删除菜单 |

### 系统管理

| 路径 | 方法 | 说明 |
|------|------|------|
| `/backend/system/attachment/*` | - | 附件管理 |
| `/backend/system/config/*` | - | 配置管理 |
| `/backend/system/gen/*` | - | 代码生成 |

## 前端 API `/api`

### 公共接口

| 路径 | 方法 | 说明 |
|------|------|------|
| `/api/common/captcha` | POST | 获取图形验证码 |
| `/api/sms/send` | POST | 发送短信验证码 |

### 用户认证（无需登录）

| 路径 | 方法 | 说明 |
|------|------|------|
| `/api/user/login` | POST | 用户登录 |
| `/api/user/mobile-login` | POST | 手机验证码登录 |
| `/api/user/reg` | POST | 用户注册 |
| `/api/user/reset-pwd` | POST | 重置密码 |

### 用户操作（需要登录）

| 路径 | 方法 | 说明 |
|------|------|------|
| `/api/user/logout` | POST | 退出登录 |
| `/api/user/change-mobile` | POST | 修改手机号 |
| `/api/user/change-pwd` | POST | 修改密码 |
| `/api/user/info` | GET | 获取用户信息 |
| `/api/menu/all` | GET | 获取菜单列表 |

## WebSocket `/ws`

| 路径 | 说明 |
|------|------|
| `/ws/ping` | WebSocket 连接测试 |

## 中间件

### 全局中间件

```go
r.Use(middleware.I18nMiddleware(c))      // 国际化
r.Use(middleware.ContainerMiddleware(c)) // 容器注入
```

### 认证中间件

```go
// 必须登录
authMiddleware.JWTAuthCheck()

// 可选登录（登录和未登录都可访问）
authMiddleware.OptionalJWTAuth()

// 权限检查
authMiddleware.PermissionCheck()
```

## 添加新路由

在 `routes/routes.go` 中添加：

```go
// 1. 创建控制器
myController := controller.NewMyController(c)

// 2. 注册路由
api.GET("/my-route", myController.MyMethod)

// 3. 需要认证的路由
api.GET("/my-route", authMiddleware.JWTAuthCheck(), myController.MyMethod)

// 4. 需要权限的路由
api.GET("/my-route", 
    authMiddleware.JWTAuthCheck(), 
    authMiddleware.PermissionCheck(), 
    myController.MyMethod,
)
```

