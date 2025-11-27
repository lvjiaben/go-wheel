package websocket

import (
	"encoding/json"
	"sync"

	"go.uber.org/zap"
)

// Hub WebSocket 连接管理中心
type Hub struct {
	// Clients 所有客户端连接（key: client ID）
	Clients map[string]*Client

	// UserClients 按用户ID索引的客户端（key: user ID）
	UserClients map[int][]*Client

	// Register 注册客户端通道
	Register chan *Client

	// Unregister 注销客户端通道
	Unregister chan *Client

	// Broadcast 广播消息通道
	Broadcast chan *Message

	// Logger 日志记录器
	Logger *zap.Logger

	// mu 读写锁
	mu sync.RWMutex

	// handlers 消息处理器映射（key: 消息类型）
	handlers map[string]MessageHandler
}

// NewHub 创建新的 Hub
func NewHub(logger *zap.Logger) *Hub {
	return &Hub{
		Clients:     make(map[string]*Client),
		UserClients: make(map[int][]*Client),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Broadcast:   make(chan *Message, 256),
		Logger:      logger,
		handlers:    make(map[string]MessageHandler),
	}
}

// Run 启动 Hub，处理客户端注册、注销和消息广播
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.registerClient(client)

		case client := <-h.Unregister:
			h.unregisterClient(client)

		case message := <-h.Broadcast:
			h.broadcastMessage(message)
		}
	}
}

// RegisterHandler 注册消息处理器
func (h *Hub) RegisterHandler(msgType string, handler MessageHandler) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.handlers[msgType] = handler
	h.Logger.Info("注册 WebSocket 消息处理器", zap.String("type", msgType))
}

// HandleMessage 处理接收到的消息
func (h *Hub) HandleMessage(client *Client, data []byte) {
	var msg Message
	if err := json.Unmarshal(data, &msg); err != nil {
		h.Logger.Error("解析 WebSocket 消息失败",
			zap.String("client_id", client.ID),
			zap.Int("user_id", client.UserID),
			zap.Error(err))
		return
	}

	// 设置消息来源信息
	msg.From = client.ID
	msg.UserID = client.UserID
	msg.Username = client.Username

	h.Logger.Debug("收到 WebSocket 消息",
		zap.String("type", msg.Type),
		zap.String("from", msg.From),
		zap.Int("user_id", msg.UserID),
		zap.String("username", msg.Username))

	// 查找并执行处理器
	h.mu.RLock()
	handler, exists := h.handlers[msg.Type]
	h.mu.RUnlock()

	if exists {
		if err := handler(client, &msg); err != nil {
			h.Logger.Error("处理 WebSocket 消息失败",
				zap.String("type", msg.Type),
				zap.String("client_id", client.ID),
				zap.Int("user_id", client.UserID),
				zap.Error(err))
		}
	} else {
		h.Logger.Warn("未找到 WebSocket 消息处理器",
			zap.String("type", msg.Type),
			zap.String("client_id", client.ID))
	}
}

// SendToUser 发送消息给指定用户的所有连接
func (h *Hub) SendToUser(userID int, message *Message) {
	h.mu.RLock()
	clients := h.UserClients[userID]
	h.mu.RUnlock()

	if len(clients) == 0 {
		h.Logger.Warn("用户没有活跃的 WebSocket 连接",
			zap.Int("user_id", userID))
		return
	}

	data, err := json.Marshal(message)
	if err != nil {
		h.Logger.Error("序列化消息失败", zap.Error(err))
		return
	}

	for _, client := range clients {
		select {
		case client.Send <- data:
			h.Logger.Debug("发送消息给用户",
				zap.Int("user_id", userID),
				zap.String("client_id", client.ID))
		default:
			// 发送通道已满，关闭连接
			close(client.Send)
			h.Unregister <- client
			h.Logger.Warn("客户端发送通道已满，关闭连接",
				zap.String("client_id", client.ID),
				zap.Int("user_id", userID))
		}
	}
}

// SendToClient 发送消息给指定客户端
func (h *Hub) SendToClient(clientID string, message *Message) {
	h.mu.RLock()
	client, exists := h.Clients[clientID]
	h.mu.RUnlock()

	if !exists {
		h.Logger.Warn("客户端不存在", zap.String("client_id", clientID))
		return
	}

	data, err := json.Marshal(message)
	if err != nil {
		h.Logger.Error("序列化消息失败", zap.Error(err))
		return
	}

	select {
	case client.Send <- data:
		h.Logger.Debug("发送消息给客户端", zap.String("client_id", clientID))
	default:
		close(client.Send)
		h.Unregister <- client
		h.Logger.Warn("客户端发送通道已满，关闭连接",
			zap.String("client_id", clientID))
	}
}

// GetClientCount 获取当前连接的客户端数量
func (h *Hub) GetClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.Clients)
}

// GetUserClientCount 获取指定用户的连接数量
func (h *Hub) GetUserClientCount(userID int) int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.UserClients[userID])
}

// registerClient 注册客户端
func (h *Hub) registerClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.Clients[client.ID] = client
	h.UserClients[client.UserID] = append(h.UserClients[client.UserID], client)

	h.Logger.Info("WebSocket 客户端已连接",
		zap.String("client_id", client.ID),
		zap.Int("user_id", client.UserID),
		zap.String("username", client.Username),
		zap.Int("total_clients", len(h.Clients)),
		zap.Int("user_clients", len(h.UserClients[client.UserID])))
}

// unregisterClient 注销客户端
func (h *Hub) unregisterClient(client *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if _, ok := h.Clients[client.ID]; ok {
		delete(h.Clients, client.ID)
		close(client.Send)

		// 从用户索引中移除
		clients := h.UserClients[client.UserID]
		for i, c := range clients {
			if c.ID == client.ID {
				h.UserClients[client.UserID] = append(clients[:i], clients[i+1:]...)
				break
			}
		}

		// 如果用户没有连接了，删除用户索引
		if len(h.UserClients[client.UserID]) == 0 {
			delete(h.UserClients, client.UserID)
		}

		h.Logger.Info("WebSocket 客户端已断开",
			zap.String("client_id", client.ID),
			zap.Int("user_id", client.UserID),
			zap.String("username", client.Username),
			zap.Int("total_clients", len(h.Clients)))
	}
}

// broadcastMessage 广播消息给所有客户端
func (h *Hub) broadcastMessage(message *Message) {
	data, err := json.Marshal(message)
	if err != nil {
		h.Logger.Error("序列化广播消息失败", zap.Error(err))
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	h.Logger.Info("广播消息",
		zap.String("type", message.Type),
		zap.Int("client_count", len(h.Clients)))

	for _, client := range h.Clients {
		select {
		case client.Send <- data:
		default:
			close(client.Send)
			h.Unregister <- client
			h.Logger.Warn("客户端发送通道已满，关闭连接",
				zap.String("client_id", client.ID))
		}
	}
}

// SendNotification 发送通知给指定用户
// 这是一个通用方法，可以在任何地方调用
func (h *Hub) SendNotification(userID int, notificationType string, data interface{}) {
	message := &Message{
		Type: "notification",
		Data: map[string]interface{}{
			"type": notificationType,
			"data": data,
		},
		From:     "server",
		UserID:   0,
		Username: "system",
	}

	h.SendToUser(userID, message)

	h.Logger.Info("发送通知给用户",
		zap.Int("user_id", userID),
		zap.String("notification_type", notificationType))
}

// BroadcastNotification 广播通知给所有用户
// 这是一个通用方法，可以在任何地方调用
func (h *Hub) BroadcastNotification(notificationType string, data interface{}) {
	message := &Message{
		Type: "notification",
		Data: map[string]interface{}{
			"type": notificationType,
			"data": data,
		},
		From:     "server",
		UserID:   0,
		Username: "system",
	}

	h.Broadcast <- message

	h.Logger.Info("广播通知给所有用户",
		zap.String("notification_type", notificationType))
}

// SendChatMessage 发送聊天消息（通用方法）
func (h *Hub) SendChatMessage(fromUserID int, fromUsername string, toClientID string, content string) {
	message := &Message{
		Type: "chat",
		Data: map[string]interface{}{
			"content": content,
		},
		From:     toClientID,
		UserID:   fromUserID,
		Username: fromUsername,
	}

	if toClientID != "" {
		h.SendToClient(toClientID, message)
	} else {
		h.Broadcast <- message
	}
}

// BroadcastSystemMessage 广播系统消息
func (h *Hub) BroadcastSystemMessage(messageType string, content string) {
	message := &Message{
		Type: messageType,
		Data: map[string]interface{}{
			"content": content,
		},
		From:     "server",
		UserID:   0,
		Username: "system",
	}

	h.Broadcast <- message

	h.Logger.Info("广播系统消息",
		zap.String("type", messageType),
		zap.String("content", content))
}

