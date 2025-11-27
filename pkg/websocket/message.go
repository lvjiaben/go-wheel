package websocket

// Message WebSocket 消息结构
type Message struct {
	Type     string                 `json:"type"`     // 消息类型
	Data     interface{}            `json:"data"`     // 消息数据
	From     string                 `json:"from"`     // 发送者ID
	To       string                 `json:"to"`       // 接收者ID（可选，为空则广播）
	UserID   int                    `json:"user_id"`  // 用户ID
	Username string                 `json:"username"` // 用户名
	Metadata map[string]interface{} `json:"metadata"` // 元数据
}

// MessageHandler 消息处理器接口
type MessageHandler func(*Client, *Message) error

