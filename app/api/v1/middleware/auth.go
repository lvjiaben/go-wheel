package middleware

import (
	"strings"

	"github.com/lvjiaben/go-wheel/pkg/utils/http"

	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/app/backend/utils"
	"github.com/lvjiaben/go-wheel/app/common/model"
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
		userId := int(claims.Id)
		username := claims.Username

		// 查询用户是否存在且状态正常
		var user model.User
		if err := m.container.GetDB().Where("id = ? AND status = 1", userId).First(&user).Error; err != nil {
			http.AuthErrorI18n(c, "backend.auth.user_not_found")
			c.Abort()
			return
		}

		if m.container.GetConfig().GetBool("api.login_sso") && user.Token != tokenString {
			http.AuthErrorI18n(c, "backend.auth.token_expired")
			c.Abort()
			return
		}

		// 设置用户信息到上下文
		c.Set("container", m.container)
		c.Set("user_id", userId)
		c.Set("username", username)

		c.Next()
	}
}

// OptionalJWTAuth 可选JWT认证中间件（登录和未登录都可访问）
// 有token则解析设置用户信息，无token或token无效则继续执行（不返回401）
func (m *AuthMiddleware) OptionalJWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 设置container到上下文
		c.Set("container", m.container)

		// 获取Authorization请求头
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.Next()
			return
		}

		// 检查Bearer前缀
		tokenParts := strings.SplitN(authHeader, " ", 2)
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			c.Next()
			return
		}

		tokenString := tokenParts[1]

		// 解析JWT令牌
		secret := m.container.GetConfig().Jwt.Secret
		claims, err := jwt.ParseToken(tokenString, secret)
		if err != nil {
			c.Next()
			return
		}

		// 从令牌中获取用户信息
		userId := int(claims.Id)
		username := claims.Username

		// 查询用户是否存在且状态正常
		var user model.User
		if err := m.container.GetDB().Where("id = ? AND status = 1", userId).First(&user).Error; err != nil {
			c.Next()
			return
		}

		// SSO检查
		if m.container.GetConfig().GetBool("api.login_sso") && user.Token != tokenString {
			c.Next()
			return
		}

		// 设置用户信息到上下文
		c.Set("user_id", userId)
		c.Set("username", username)
		c.Set("user_mobile", user.Mobile)

		c.Next()
	}
}
