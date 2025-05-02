package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/app/backend/service"
	"github.com/lvjiaben/go-wheel/pkg/container"
	"github.com/lvjiaben/go-wheel/pkg/utils"
	"go.uber.org/zap"
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

// Auth 认证中间件
func (m *AuthMiddleware) Auth() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取token
		token := c.GetHeader("Authorization")
		if token == "" {
			c.JSON(401, gin.H{"code": 401, "msg": "未登录"})
			c.Abort()
			return
		}

		// 解析token
		token = strings.TrimPrefix(token, "Bearer ")
		claims, err := utils.ParseToken(token, m.container.GetConfig().GetString("jwt.secret"))
		if err != nil {
			m.container.GetLogger().Error("解析token失败", zap.Error(err))
			c.JSON(401, gin.H{"code": 401, "msg": "未登录"})
			c.Abort()
			return
		}

		// 验证token是否有效
		var admin struct {
			Id    int    `json:"id"`
			Token string `json:"token"`
		}
		if err := m.container.GetDB().Table("admin").Where("id = ?", claims.Id).First(&admin).Error; err != nil {
			m.container.GetLogger().Error("获取用户信息失败", zap.Error(err))
			c.JSON(401, gin.H{"code": 401, "msg": "未登录"})
			c.Abort()
			return
		}

		if admin.Token != token {
			c.JSON(401, gin.H{"code": 401, "msg": "未登录"})
			c.Abort()
			return
		}

		// 验证权限
		if !m.checkPermission(claims.Id, c.Request.URL.Path, c.Request.Method) {
			c.JSON(403, gin.H{"code": 403, "msg": "无权限"})
			c.Abort()
			return
		}

		// 设置用户信息
		c.Set("admin_id", claims.Id)
		c.Set("username", claims.Username)
		c.Next()
	}
}

// checkPermission 检查权限
func (m *AuthMiddleware) checkPermission(userId int, path string, method string) bool {
	return m.authService.CheckPermission(userId, path, method)
}
