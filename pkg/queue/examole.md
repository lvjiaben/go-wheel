package queue

import (
	"context"
	"fmt"

	queuePkg "github.com/lvjiaben/go-wheel/pkg/queue"
	"go.uber.org/zap"
)

// 自动注册消费者
func init() {
	queuePkg.Register("test.queue", NewTestConsumer)
}

// TestConsumer 测试消费者（同时处理即时和延迟消息）
type TestConsumer struct {
	queuePkg.BaseConsumer
	container queuePkg.Container
}

// NewTestConsumer 创建测试消费者
func NewTestConsumer(c queuePkg.Container) queuePkg.Consumer {
	return &TestConsumer{
		BaseConsumer: queuePkg.BaseConsumer{
			Topic:       "test.queue",
			Description: "测试队列消费者（支持即时和延迟消息）",
			Concurrency: 2, // 2个并发工作线程
		},
		container: c,
	}
}

// Handle 处理消息（即时消息和延迟消息都会到这里）
func (c *TestConsumer) Handle(ctx context.Context, message []byte) error {
	logger := c.container.GetLogger()

	logger.Info("收到消息",
		zap.String("topic", c.GetTopic()),
		zap.ByteString("message", message),
	)

	// 解析消息
	var msg map[string]interface{}
	if err := queuePkg.UnmarshalMessage(message, &msg); err != nil {
		logger.Error("解析消息失败", zap.Error(err))
		return err
	}

	// 处理业务逻辑
	fmt.Printf("✅ 处理消息: %+v\n", msg)

	return nil
}

// OnError 错误处理
func (c *TestConsumer) OnError(err error, message []byte) {
	c.container.GetLogger().Error("消息处理失败",
		zap.String("topic", c.GetTopic()),
		zap.Error(err),
		zap.ByteString("message", message),
	)
}


// ExampleUsage 演示如何使用队列
func ExampleUsage(c *container.Container) {
	// 获取队列辅助工具
	queueHelper := c.GetQueueHelper()
	if queueHelper == nil {
		return
	}

	ctx := context.Background()

	// 1. 发送即时消息到 test.queue
	message := map[string]interface{}{
		"user_id": 123,
		"action":  "login",
		"time":    time.Now(),
	}
	queueHelper.Push(ctx, "test.queue", message)

	// 2. 发送延迟消息到 test.queue（5分钟后执行）
	// 注意：延迟消息最终也会到达 test.queue，由同一个消费者处理
	delayMessage := map[string]interface{}{
		"order_id": 456,
		"action":   "cancel_order",
		"time":     time.Now(),
	}
	queueHelper.PushDelay(ctx, "test.queue", delayMessage, 5*time.Minute)

	// 3. 发送延迟消息（10秒后执行，用于测试）
	testDelayMessage := map[string]interface{}{
		"test": "delay_10s",
		"time": time.Now(),
	}
	queueHelper.PushDelay(ctx, "test.queue", testDelayMessage, 10*time.Second)
}

