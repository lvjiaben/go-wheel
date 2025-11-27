package middleware

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/lvjiaben/go-wheel/pkg/container"
	"github.com/lvjiaben/go-wheel/pkg/jwt"
	wsPackage "github.com/lvjiaben/go-wheel/pkg/websocket"
	"go.uber.org/zap"
)

// WebSocket 升级器配置
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	// CheckOrigin 检查请求来源
	// 生产环境需要严格检查，这里为了方便开发允许所有来源
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WebSocketUpgradeMiddleware WebSocket 升级中间件
// 这是一个公共底层中间件，用于将 HTTP 连接升级为 WebSocket 连接
// 注意：此中间件应该在认证中间件之后使用，以便从上下文中获取用户信息
func WebSocketUpgradeMiddleware(hub *wsPackage.Hub) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 升级 HTTP 连接为 WebSocket
		conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
		if err != nil {
			hub.Logger.Error("WebSocket 升级失败",
				zap.Error(err))
			ctx.JSON(500, gin.H{
				"code":    500,
				"message": "WebSocket 升级失败",
			})
			return
		}

		// 从上下文中获取用户信息（由认证中间件设置）
		// 如果没有认证信息，则使用默认值（匿名用户）
		userID, exists := ctx.Get("user_id")
		if !exists {
			userID = 0
		}

		username, exists := ctx.Get("username")
		if !exists {
			username = "anonymous"
		}

		// 创建客户端
		client := &wsPackage.Client{
			ID:       generateClientID(),
			Conn:     conn,
			Send:     make(chan []byte, 256),
			Hub:      hub,
			UserID:   userID.(int),
			Username: username.(string),
			Metadata: make(map[string]interface{}),
		}

		// 注册客户端
		hub.Register <- client

		// 启动读写协程
		go client.WritePump()
		go client.ReadPump()
	}
}

// generateClientID 生成客户端ID
func generateClientID() string {
	return fmt.Sprintf("%d-%s", time.Now().UnixNano(), uuid.New().String()[:8])
}

// WebSocketAuthMiddleware WebSocket JWT 认证中间件
// 这是一个公共底层中间件，用于验证 WebSocket 连接的 JWT token
// 支持两种方式传递 token：
// 1. 查询参数：ws://host/path?token=xxx
// 2. Authorization 头：Authorization: Bearer xxx
func WebSocketAuthMiddleware(c *container.Container) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 1. 尝试从查询参数获取 token
		token := ctx.Query("token")

		// 2. 如果查询参数没有，尝试从 Authorization 头获取
		if token == "" {
			authHeader := ctx.Request.Header.Get("Authorization")
			if strings.HasPrefix(authHeader, "Bearer ") {
				token = strings.TrimPrefix(authHeader, "Bearer ")
			}
		}

		// 3. 如果没有 token，返回未授权错误
		if token == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "未提供认证令牌",
			})
			ctx.Abort()
			return
		}

		// 4. 解析 JWT token
		secret := c.GetConfig().Jwt.Secret
		claims, err := jwt.ParseToken(token, secret)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"code":    401,
				"message": "无效的认证令牌",
			})
			ctx.Abort()
			return
		}

		// 5. 从 claims 中获取用户信息
		userID := int(claims.Id)
		username := claims.Username

		// 6. 设置用户信息到上下文（供后续中间件和控制器使用）
		ctx.Set("user_id", userID)
		ctx.Set("username", username)
		ctx.Set("token", token)

		ctx.Next()
	}
}

