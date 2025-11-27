package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"go.uber.org/zap"
)

// Container 容器接口
type Container interface {
	GetLogger() Logger
	GetDB() interface{}
	GetRedis() interface{}
}

// Consumer 消费者接口
type Consumer interface {
	// GetTopic 获取消费的主题
	GetTopic() string
	
	// Handle 处理消息
	Handle(ctx context.Context, message []byte) error
	
	// GetDescription 获取消费者描述
	GetDescription() string
	
	// GetConcurrency 获取并发数（默认1）
	GetConcurrency() int
	
	// OnError 错误处理（可选）
	OnError(err error, message []byte)
}

// BaseConsumer 基础消费者
type BaseConsumer struct {
	Topic       string
	Description string
	Concurrency int
}

func (c *BaseConsumer) GetTopic() string {
	return c.Topic
}

func (c *BaseConsumer) GetDescription() string {
	return c.Description
}

func (c *BaseConsumer) GetConcurrency() int {
	if c.Concurrency <= 0 {
		return 1
	}
	return c.Concurrency
}

func (c *BaseConsumer) OnError(err error, message []byte) {
	// 默认不处理
}

// ConsumerFactory 消费者工厂函数
type ConsumerFactory func(Container) Consumer

var (
	// consumerRegistry 全局消费者注册表
	consumerRegistry = make(map[string]ConsumerFactory)
	// registryMu 注册表互斥锁
	registryMu sync.RWMutex
)

// Register 注册消费者（同时处理即时和延迟消息）
func Register(topic string, factory ConsumerFactory) {
	registryMu.Lock()
	defer registryMu.Unlock()

	if _, exists := consumerRegistry[topic]; exists {
		panic("消费者已注册: " + topic)
	}

	consumerRegistry[topic] = factory
}

// GetRegisteredConsumers 获取所有已注册的消费者主题
func GetRegisteredConsumers() []string {
	registryMu.RLock()
	defer registryMu.RUnlock()

	topics := make([]string, 0, len(consumerRegistry))
	for topic := range consumerRegistry {
		topics = append(topics, topic)
	}
	return topics
}

// ContainerWithQueue 带队列的容器接口
type ContainerWithQueue interface {
	Container
	GetMessageQueue() interface{}
	GetDelayQueue() interface{}
}

// StartAllConsumers 启动所有消费者
func StartAllConsumers(c ContainerWithQueue, manager *RabbitMQManager) error {
	registryMu.RLock()
	defer registryMu.RUnlock()

	// 启动所有消费者
	for topic, factory := range consumerRegistry {
		consumer := factory(c)
		if err := startConsumer(c, manager, consumer); err != nil {
			c.GetLogger().Error("启动消费者失败",
				zap.String("topic", topic),
				zap.Error(err),
			)
			return err
		}
	}

	c.GetLogger().Info("所有消费者启动完成",
		zap.Int("consumers", len(consumerRegistry)),
	)

	return nil
}

// startConsumer 启动单个消费者
func startConsumer(c Container, manager *RabbitMQManager, consumer Consumer) error {
	channel, err := manager.GetChannel()
	if err != nil {
		return err
	}

	// 为每个 topic 创建独立的队列
	topic := consumer.GetTopic()
	exchange := manager.config.Exchange

	// 声明队列
	queueName := topic
	_, err = channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		nil,
	)
	if err != nil {
		return fmt.Errorf("声明队列失败: %v", err)
	}

	// 绑定队列到交换机
	err = channel.QueueBind(
		queueName,
		topic, // routing key = topic
		exchange,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("绑定队列失败: %v", err)
	}

	// 开始消费
	msgs, err := channel.Consume(
		queueName,
		fmt.Sprintf("%s-%s", manager.config.DefaultConsumerName, topic),
		false, // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,
	)
	if err != nil {
		return fmt.Errorf("开始消费失败: %v", err)
	}

// 启动多个goroutine处理消息
concurrency := consumer.GetConcurrency()
for i := 0; i < concurrency; i++ {
	go func(workerID int) {
		c.GetLogger().Info("消费者工作线程启动",
			zap.String("topic", consumer.GetTopic()),
			zap.Int("worker_id", workerID),
		)

		for msg := range msgs {
			ctx := context.Background()

			c.GetLogger().Debug("收到消息",
				zap.String("topic", consumer.GetTopic()),
				zap.Int("worker_id", workerID),
			)

			if err := consumer.Handle(ctx, msg.Body); err != nil {
				c.GetLogger().Error("处理消息失败",
					zap.String("topic", consumer.GetTopic()),
					zap.Int("worker_id", workerID),
					zap.Error(err),
				)
				consumer.OnError(err, msg.Body)
				msg.Nack(false, true) // 重新入队
			} else {
				msg.Ack(false)
				c.GetLogger().Debug("消息处理成功",
					zap.String("topic", consumer.GetTopic()),
					zap.Int("worker_id", workerID),
				)
			}
		}

		c.GetLogger().Warn("消费者工作线程退出",
			zap.String("topic", consumer.GetTopic()),
			zap.Int("worker_id", workerID),
		)
	}(i)
}

	c.GetLogger().Info("消费者启动成功",
		zap.String("topic", consumer.GetTopic()),
		zap.String("description", consumer.GetDescription()),
		zap.Int("concurrency", concurrency),
	)

	return nil
}

// Message 通用消息结构
type Message struct {
	Data interface{} `json:"data"`
}

// UnmarshalMessage 反序列化消息
func UnmarshalMessage(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

