package websocket

import (
    "github.com/gin-gonic/gin"
    "github.com/lvjiaben/go-wheel/pkg/container"
    wsPackage "github.com/lvjiaben/go-wheel/pkg/websocket"
    "go.uber.org/zap"
)

type PingController struct {
    container *container.Container
    hub       *wsPackage.Hub
}

func NewPingController(c *container.Container, hub *wsPackage.Hub) *PingController {
    ctrl := &PingController{
        container: c,
        hub:       hub,
    }
    
    // 注册消息处理器
    ctrl.registerHandlers()
    
    return ctrl
}

func (ctrl *PingController) Connect(ctx *gin.Context) {
    // 连接建立后的逻辑
    ctrl.container.GetLogger().Info("WebSocket 连接已建立",
        zap.Int("user_id", ctx.GetInt("user_id")),
        zap.String("username", ctx.GetString("username")))
}

func (ctrl *PingController) registerHandlers() {
    // 注册消息处理器
    ctrl.hub.RegisterHandler("ping", ctrl.handle)
}

func (ctrl *PingController) handle(client *wsPackage.Client, msg *wsPackage.Message) error {

	// 发送响应
	ctrl.hub.SendNotification(client.UserID, "pong", map[string]interface{}{
		"message": msg.Data,
		"from":    "来自测试Ping的输入返回",
	})

	return nil
}