package middleware

import (
	"strings"

	"github.com/lvjiaben/go-wheel/pkg/utils/http"

	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/app/backend/model/admin"
	"github.com/lvjiaben/go-wheel/app/backend/utils"
	"github.com/lvjiaben/go-wheel/pkg/container"
	"github.com/lvjiaben/go-wheel/pkg/jwt"
)

// AuthMiddleware 认证中间件
type AuthMiddleware struct {
	container *container.Container
	authUtils *utils.AuthUtils
}

// NewAuthMiddleware 创建认证中间件
func NewAuthMiddleware(c *container.Container) *AuthMiddleware {
	return &AuthMiddleware{
		container: c,
		authUtils: utils.NewAuthUtils(c),
	}
}

// PermissionCheck 权限检查中间件
func (m *AuthMiddleware) PermissionCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取当前用户ID
		userId := c.GetInt("admin_id")
		if userId == 0 {
			http.AuthErrorI18n(c, "backend.auth.unauthorized")
			c.Abort()
			return
		}

		// 获取用户名
		username := c.GetString("username")
		if username == "" {
			http.AuthErrorI18n(c, "backend.auth.unauthorized")
			c.Abort()
			return
		}

		// 检查路径权限
		path := c.Request.URL.Path

		// 超级管理员拥有所有权限
		if isSuper, _ := m.authUtils.IsAdminSuper(userId); isSuper {
			c.Next()
			return
		}

		// 使用自定义权限验证
		if !m.hasPermission(username, path) {
			http.ForbiddenErrorI18n(c, "backend.auth.forbidden")
			c.Abort()
			return
		}

		c.Next()
	}
}

// JWTAuthCheck JWT认证检查中间件
func (m *AuthMiddleware) JWTAuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取Authorization请求头
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			http.AuthErrorI18n(c, "backend.auth.unauthorized")
			c.Abort()
			return
		}

		// 检查Bearer前缀
		tokenParts := strings.SplitN(authHeader, " ", 2)
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			http.AuthErrorI18n(c, "backend.auth.invalid_token")
			c.Abort()
			return
		}

		tokenString := tokenParts[1]

		// 解析JWT令牌
		secret := m.container.GetConfig().Jwt.Secret
		claims, err := jwt.ParseToken(tokenString, secret)
		if err != nil {
			http.AuthErrorI18n(c, "backend.auth.invalid_token")
			c.Abort()
			return
		}

		// 从令牌中获取用户信息
		adminId := int(claims.Id)
		username := claims.Username

		// 查询用户是否存在且状态正常
		var adminUser admin.Admin
		if err := m.container.GetDB().Where("id = ? AND status = 1", adminId).First(&adminUser).Error; err != nil {
			http.AuthErrorI18n(c, "backend.auth.user_not_found")
			c.Abort()
			return
		}

		if m.container.GetConfig().GetBool("admin.login_sso") && adminUser.Token != tokenString {
			http.AuthErrorI18n(c, "backend.auth.token_expired")
			c.Abort()
			return
		}

		// 设置用户信息到上下文
		c.Set("container", m.container)
		c.Set("admin_id", adminId)
		c.Set("username", username)

		c.Next()
	}
}

// hasPermission 检查用户是否有访问指定路径和方法的权限
func (m *AuthMiddleware) hasPermission(username, path string) bool {
	var count int64

	// 使用 GORM 链式操作替代原生 SQL，避免 SQL 注入风险
	err := m.container.GetDB().Table("admin aa").
		Joins("JOIN admin_role_admin ara ON aa.id = ara.admin_id").
		Joins("JOIN admin_role_menu arm ON ara.role_id = arm.role_id").
		Joins("JOIN admin_menu am ON arm.menu_id = am.id").
		Where("aa.username = ?", username).
		Where("am.type = ?", "button").
		Where(m.container.GetDB().Where("am.route = ?", path).
			Or(m.container.GetDB().Where("am.route != ?", "").
				Where(m.container.GetDB().Where("? LIKE CONCAT(am.route, '/%')", path).
					Or("am.route = ?", path)))).
		Count(&count).Error

	return err == nil && count > 0
}
