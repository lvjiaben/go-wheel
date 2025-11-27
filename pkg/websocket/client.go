package websocket

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

// 配置常量
const (
	// WriteWait 写入超时时间
	WriteWait = 10 * time.Second

	// PongWait 等待 pong 消息的超时时间
	PongWait = 60 * time.Second

	// PingPeriod 发送 ping 消息的间隔（必须小于 pongWait）
	PingPeriod = (PongWait * 9) / 10

	// MaxMessageSize 最大消息大小
	MaxMessageSize = 512 * 1024 // 512KB
)

// Client WebSocket 客户端
type Client struct {
	ID       string                 // 客户端唯一ID
	Conn     *websocket.Conn        // WebSocket 连接
	Send     chan []byte            // 发送消息通道
	Hub      *Hub                   // 所属的 Hub
	UserID   int                    // 用户ID（从JWT中获取）
	Username string                 // 用户名（从JWT中获取）
	Metadata map[string]interface{} // 额外元数据
	mu       sync.RWMutex           // 读写锁
}

// ReadPump 从 WebSocket 连接读取消息
// 应用程序在每个连接的 goroutine 中运行 ReadPump
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadLimit(MaxMessageSize)
	c.Conn.SetReadDeadline(time.Now().Add(PongWait))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(PongWait))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				c.Hub.Logger.Error("WebSocket 读取错误",
					zap.String("client_id", c.ID),
					zap.Int("user_id", c.UserID),
					zap.Error(err))
			}
			break
		}

		// 处理接收到的消息
		c.Hub.HandleMessage(c, message)
	}
}

// WritePump 向 WebSocket 连接写入消息
// 应用程序在每个连接的 goroutine 中运行 WritePump
func (c *Client) WritePump() {
	ticker := time.NewTicker(PingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(WriteWait))
			if !ok {
				// Hub 关闭了通道
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// 将队列中的消息一起发送
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(WriteWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// SetMetadata 设置元数据
func (c *Client) SetMetadata(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.Metadata[key] = value
}

// GetMetadata 获取元数据
func (c *Client) GetMetadata(key string) (interface{}, bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	val, ok := c.Metadata[key]
	return val, ok
}

