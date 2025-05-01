package initialize

import (
	"admin/pkg/queue"
)

func NewMessageQueue() *queue.MessageQueue {
	return queue.NewMessageQueue()
}

func NewDelayQueue() *queue.DelayQueue {
	return queue.NewDelayQueue()
}
