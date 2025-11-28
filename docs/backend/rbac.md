# RBAC 权限系统

项目采用 RBAC（基于角色的访问控制）权限模型，实现细粒度的权限管理。

## 数据模型

### 表结构

```sql
-- 管理员表
CREATE TABLE admin (
    id INT PRIMARY KEY AUTO_INCREMENT,
    username VARCHAR(50) NOT NULL,
    password VARCHAR(255) NOT NULL,
    status TINYINT DEFAULT 1,
    token VARCHAR(500)
);

-- 角色表
CREATE TABLE admin_role (
    id INT PRIMARY KEY AUTO_INCREMENT,
    name VARCHAR(50) NOT NULL,
    code VARCHAR(50) NOT NULL,
    status TINYINT DEFAULT 1
);

-- 菜单表
CREATE TABLE admin_menu (
    id INT PRIMARY KEY AUTO_INCREMENT,
    parent_id INT DEFAULT 0,
    name VARCHAR(50) NOT NULL,
    route VARCHAR(200),
    type ENUM('menu', 'button') DEFAULT 'menu',
    sort INT DEFAULT 0
);

-- 管理员-角色关联表
CREATE TABLE admin_role_admin (
    admin_id INT,
    role_id INT,
    PRIMARY KEY (admin_id, role_id)
);

-- 角色-菜单关联表
CREATE TABLE admin_role_menu (
    role_id INT,
    menu_id INT,
    PRIMARY KEY (role_id, menu_id)
);
```

### 关系图

```
Admin ──┬── AdminRoleAdmin ──┬── AdminRole
        │                    │
        │                    └── AdminRoleMenu ──── AdminMenu
        │
        └── 直接权限检查
```

## 权限检查流程

### 1. JWT 认证

```go
// app/backend/middleware/auth.go
func (m *AuthMiddleware) JWTAuthCheck() gin.HandlerFunc {
    return func(c *gin.Context) {
        // 1. 获取 Authorization 请求头
        authHeader := c.Request.Header.Get("Authorization")
        
        // 2. 解析 Bearer Token
        tokenString := strings.TrimPrefix(authHeader, "Bearer ")
        
        // 3. 验证 JWT
        claims, err := jwt.ParseToken(tokenString, secret)
        
        // 4. 查询用户状态
        var admin Admin
        db.Where("id = ? AND status = 1", claims.Id).First(&admin)
        
        // 5. 设置用户信息到上下文
        c.Set("admin_id", adminId)
        c.Set("username", username)
    }
}
```

### 2. 权限检查

```go
func (m *AuthMiddleware) PermissionCheck() gin.HandlerFunc {
    return func(c *gin.Context) {
        userId := c.GetInt("admin_id")
        username := c.GetString("username")
        path := c.Request.URL.Path
        
        // 超级管理员跳过权限检查
        if isSuper, _ := m.authUtils.IsAdminSuper(userId); isSuper {
            c.Next()
            return
        }
        
        // 检查路径权限
        if !m.hasPermission(username, path) {
            http.ForbiddenErrorI18n(c, "backend.auth.forbidden")
            c.Abort()
            return
        }
    }
}
```

### 3. 权限查询

```go
func (m *AuthMiddleware) hasPermission(username, path string) bool {
    var count int64
    
    m.container.GetDB().Table("admin aa").
        Joins("JOIN admin_role_admin ara ON aa.id = ara.admin_id").
        Joins("JOIN admin_role_menu arm ON ara.role_id = arm.role_id").
        Joins("JOIN admin_menu am ON arm.menu_id = am.id").
        Where("aa.username = ?", username).
        Where("am.type = ?", "button").
        Where("am.route = ? OR ? LIKE CONCAT(am.route, '/%')", path, path).
        Count(&count)
    
    return count > 0
}
```

## 使用方法

### 路由配置

```go
// routes/routes.go
adminGroup := api.Group("/admin", 
    authMiddleware.JWTAuthCheck(),      // JWT 认证
    authMiddleware.PermissionCheck(),   // 权限检查
)
{
    adminGroup.GET("/list", adminController.List)
    adminGroup.POST("/save", adminController.Save)
}
```

### 菜单类型

| 类型 | 说明 | 用途 |
|------|------|------|
| `menu` | 菜单 | 前端导航菜单 |
| `button` | 按钮 | 接口权限控制 |

### 配置菜单权限

在后台菜单管理中添加按钮类型的菜单：

```
名称: 用户列表
路由: /backend/user/list
类型: button
```

## 超级管理员

超级管理员拥有所有权限，不受权限检查限制：

```go
// 判断是否超级管理员
func (u *AuthUtils) IsAdminSuper(adminId int) (bool, error) {
    var count int64
    err := u.container.GetDB().Table("admin_role_admin ara").
        Joins("JOIN admin_role ar ON ara.role_id = ar.id").
        Where("ara.admin_id = ?", adminId).
        Where("ar.code = ?", "super_admin").
        Count(&count).Error
    return count > 0, err
}
```

## 前端配合

### 获取菜单

```typescript
// 登录后获取菜单
const menus = await getMenusApi();
// 根据菜单生成路由
```

### 按钮权限

```vue
<template>
  <Button v-if="hasPermission('user:create')">新增</Button>
</template>
```

