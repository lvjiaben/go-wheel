# WebSocket

项目内置 WebSocket 支持，使用 Hub 模式管理连接，支持用户级消息推送。

## 目录结构

```
app/websocket/
└── controller/
    └── ping.go         # Ping 控制器

pkg/websocket/
├── hub.go              # 连接管理中心
├── client.go           # 客户端连接
└── message.go          # 消息结构
```

## 连接流程

```
客户端 → WebSocket 握手 → JWT 认证 → 注册到 Hub → 消息收发
```

## 配置路由

```go
// routes/routes.go
func registerWebSocketRoutes(r *gin.Engine, c *container.Container) {
    wsHub := c.GetWebSocketHub()
    
    ws := r.Group("/ws").Use(
        middleware.ContainerMiddleware(c),
        middleware.WebSocketAuthMiddleware(c),      // JWT 认证
        middleware.WebSocketUpgradeMiddleware(wsHub), // 升级连接
    )
    {
        pingController := websocketController.NewPingController(c, wsHub)
        ws.GET("/ping", pingController.Connect)
    }
}
```

## 创建控制器

```go
// app/websocket/controller/ping.go
package controller

import (
    "github.com/gin-gonic/gin"
    "github.com/lvjiaben/go-wheel/pkg/container"
    "github.com/lvjiaben/go-wheel/pkg/websocket"
)

type PingController struct {
    container *container.Container
    hub       *websocket.Hub
}

func NewPingController(c *container.Container, hub *websocket.Hub) *PingController {
    return &PingController{
        container: c,
        hub:       hub,
    }
}

func (c *PingController) Connect(ctx *gin.Context) {
    // 连接已在中间件中建立
    // 这里可以做额外的初始化
}
```

## Hub 使用

### 发送消息给用户

```go
hub := container.GetWebSocketHub()

// 发送给指定用户的所有连接
hub.SendToUser(userID, &websocket.Message{
    Type: "notification",
    Data: map[string]interface{}{
        "title":   "新消息",
        "content": "您有一条新消息",
    },
})
```

### 发送给指定客户端

```go
hub.SendToClient(clientID, &websocket.Message{
    Type: "chat",
    Data: map[string]interface{}{
        "content": "Hello!",
    },
})
```

### 广播消息

```go
// 广播给所有连接
hub.Broadcast <- &websocket.Message{
    Type: "system",
    Data: map[string]interface{}{
        "content": "系统维护通知",
    },
}

// 或使用便捷方法
hub.BroadcastSystemMessage("announcement", "系统将于今晚维护")
```

### 发送通知

```go
// 发送通知给用户
hub.SendNotification(userID, "order_paid", map[string]interface{}{
    "order_id": 12345,
    "amount":   99.9,
})

// 广播通知
hub.BroadcastNotification("new_feature", map[string]interface{}{
    "title": "新功能上线",
})
```

## 消息处理器

### 注册处理器

```go
hub.RegisterHandler("chat", func(client *websocket.Client, msg *websocket.Message) error {
    // 处理聊天消息
    content := msg.Data["content"].(string)
    
    // 广播给所有人
    hub.Broadcast <- &websocket.Message{
        Type:     "chat",
        Data:     msg.Data,
        From:     client.ID,
        UserID:   client.UserID,
        Username: client.Username,
    }
    
    return nil
})
```

## 消息格式

```go
// pkg/websocket/message.go
type Message struct {
    Type     string                 `json:"type"`     // 消息类型
    Data     map[string]interface{} `json:"data"`     // 消息数据
    From     string                 `json:"from"`     // 发送者客户端ID
    UserID   int                    `json:"user_id"`  // 发送者用户ID
    Username string                 `json:"username"` // 发送者用户名
}
```

## 前端连接

```javascript
// 建立连接
const ws = new WebSocket('ws://localhost:8080/ws/ping?token=xxx');

ws.onopen = () => {
    console.log('WebSocket 已连接');
};

ws.onmessage = (event) => {
    const msg = JSON.parse(event.data);
    console.log('收到消息:', msg);
};

ws.onclose = () => {
    console.log('WebSocket 已断开');
};

// 发送消息
ws.send(JSON.stringify({
    type: 'chat',
    data: { content: 'Hello!' }
}));
```

## 获取连接信息

```go
// 获取总连接数
count := hub.GetClientCount()

// 获取用户连接数
userCount := hub.GetUserClientCount(userID)
```

## 最佳实践

1. **心跳检测** - 定期发送 ping 消息保持连接
2. **断线重连** - 前端实现自动重连机制
3. **消息确认** - 重要消息实现确认机制
4. **限流保护** - 限制消息发送频率

