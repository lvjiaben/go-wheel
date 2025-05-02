package initialize

import (
	"github.com/lvjiaben/go-wheel/pkg/container"
	"github.com/lvjiaben/go-wheel/pkg/types"
	"go.uber.org/zap"
)

// RabbitMQLoad 加载RabbitMQ
func RabbitMQLoad(c *container.Container) error {
	cfg := c.GetConfig()
	if cfg == nil {
		return nil
	}

	// 检查是否启用
	if !cfg.RabbitMQ.State {
		c.GetLogger().Info("RabbitMQ未启用")
		return nil
	}

	// 配置RabbitMQ
	config := &types.RabbitMQConfig{
		Host:                cfg.RabbitMQ.Host,
		Port:                cfg.RabbitMQ.Port,
		User:                cfg.RabbitMQ.User,
		Pass:                cfg.RabbitMQ.Pass,
		VirtualHost:         cfg.RabbitMQ.VirtualHost,
		QueueName:           cfg.RabbitMQ.QueueName,
		DelayQueueName:      cfg.RabbitMQ.DelayQueueName,
		Exchange:            cfg.RabbitMQ.Exchange,
		DelayExchange:       cfg.RabbitMQ.DelayExchange,
		RetryCount:          cfg.RabbitMQ.RetryCount,
		ReconnectInterval:   cfg.RabbitMQ.ReconnectInterval,
		HeartbeatInterval:   cfg.RabbitMQ.HeartbeatInterval,
		ConnectionTimeout:   cfg.RabbitMQ.ConnectionTimeout,
		EnableConfirmation:  cfg.RabbitMQ.EnableConfirmation,
		PrefetchCount:       cfg.RabbitMQ.PrefetchCount,
		PrefetchSize:        cfg.RabbitMQ.PrefetchSize,
		DefaultConsumerName: cfg.RabbitMQ.DefaultConsumerName,
	}

	// 创建RabbitMQ管理器
	rabbit := types.NewRabbitMQManager(config)

	// 设置消息队列
	c.SetMessageQueue(rabbit)

	// 设置延迟队列
	c.SetDelayQueue(types.NewDelayQueueManager(rabbit))

	c.GetLogger().Info("RabbitMQ初始化成功",
		zap.String("host", cfg.RabbitMQ.Host),
		zap.Int("port", cfg.RabbitMQ.Port),
		zap.String("queue", cfg.RabbitMQ.QueueName),
		zap.String("delayQueue", cfg.RabbitMQ.DelayQueueName))

	return nil
}

// RabbitMQMessageQueue 创建新的消息队列
func RabbitMQMessageQueue() types.MessageQueue {
	cfg := &types.RabbitMQConfig{
		Host:                "127.0.0.1",
		Port:                5672,
		User:                "guest",
		Pass:                "guest",
		VirtualHost:         "/",
		QueueName:           "go-admin",
		DelayQueueName:      "go-admin-delay",
		Exchange:            "go-admin-exchange",
		DelayExchange:       "go-admin-delay-exchange",
		RetryCount:          3,
		ReconnectInterval:   5,
		HeartbeatInterval:   30,
		ConnectionTimeout:   10,
		EnableConfirmation:  true,
		PrefetchCount:       10,
		PrefetchSize:        0,
		DefaultConsumerName: "go-admin-consumer",
	}

	return types.NewRabbitMQManager(cfg)
}

// RabbitMQDelayQueue 创建新的延迟队列
func RabbitMQDelayQueue() types.DelayQueue {
	rabbit := RabbitMQMessageQueue().(*types.RabbitMQManager)
	return types.NewDelayQueueManager(rabbit)
}
