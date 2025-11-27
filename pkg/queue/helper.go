package queue

import (
	"context"
	"time"
)

// QueueHelper 队列辅助工具
type QueueHelper struct {
	manager *RabbitMQManager
}

// NewQueueHelper 创建队列辅助工具
func NewQueueHelper(manager *RabbitMQManager) *QueueHelper {
	return &QueueHelper{
		manager: manager,
	}
}

// Push 发送普通消息
func (h *QueueHelper) Push(ctx context.Context, topic string, message interface{}) error {
	return h.manager.Publish(ctx, topic, message)
}

// PushDelay 发送延迟消息
func (h *QueueHelper) PushDelay(ctx context.Context, topic string, message interface{}, delay time.Duration) error {
	return h.manager.PublishDelay(ctx, topic, message, delay)
}

