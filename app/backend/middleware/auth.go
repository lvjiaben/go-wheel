package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/app/backend/model"
	"github.com/lvjiaben/go-wheel/app/backend/service"
	"github.com/lvjiaben/go-wheel/pkg/container"
	"github.com/lvjiaben/go-wheel/pkg/jwt"
)

// AuthMiddleware 认证中间件
type AuthMiddleware struct {
	container   *container.Container
	authService *service.AuthService
}

// NewAuthMiddleware 创建认证中间件
func NewAuthMiddleware(c *container.Container) *AuthMiddleware {
	return &AuthMiddleware{
		container:   c,
		authService: service.NewAuthService(c),
	}
}

// PermissionCheck 权限检查中间件
func (m *AuthMiddleware) PermissionCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		userId := c.GetInt("admin_id")
		if userId <= 0 {
			c.JSON(http.StatusOK, gin.H{
				"code":    401,
				"message": "请先登陆",
			})
			c.Abort()
			return
		}
		path := c.Request.URL.Path
		if !m.authService.CheckPermission(userId, path, c.Request.Method) {
			c.JSON(http.StatusOK, gin.H{
				"code":    401,
				"message": "无权访问",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

// JWTAuthCheck JWT认证检查中间件
func (m *AuthMiddleware) JWTAuthCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{
				"code":    401,
				"message": "请先登陆",
			})
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusOK, gin.H{
				"code":    401,
				"message": "请先登陆",
			})
			c.Abort()
			return
		}
		mc, err := jwt.ParseToken(parts[1], m.container.GetConfig().GetString("jwt.secret"))
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    401,
				"message": "请先登陆",
			})
			c.Abort()
			return
		}
		var user model.Admin
		if res := m.container.GetDB().First(&user, "id = ?", mc.Id); res.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    401,
				"message": "请先登陆",
			})
			c.Abort()
			return
		}
		if m.container.GetConfig().GetBool("admin.login_sso") && user.Token != parts[1] {
			c.JSON(http.StatusOK, gin.H{
				"code":    401,
				"message": "登陆失效",
			})
			c.Abort()
			return
		}
		c.Set("admin_id", user.Id)
		c.Next()
	}
}
