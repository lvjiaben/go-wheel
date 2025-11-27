package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.uber.org/zap"
)

// Logger 日志接口
type Logger interface {
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
}

// RabbitMQConfig RabbitMQ配置
type RabbitMQConfig struct {
	Host                string
	Port                int
	VirtualHost         string
	User                string
	Pass                string
	QueueName           string
	DelayQueueName      string
	Exchange            string
	DelayExchange       string
	RetryCount          int
	ReconnectInterval   int
	HeartbeatInterval   int
	ConnectionTimeout   int
	EnableConfirmation  bool
	PrefetchCount       int
	PrefetchSize        int
	DefaultConsumerName string
}

// RabbitMQManager RabbitMQ管理器
type RabbitMQManager struct {
	config     *RabbitMQConfig
	logger     Logger
	conn       *amqp.Connection
	channel    *amqp.Channel
	mu         sync.RWMutex
	ctx        context.Context
	cancel     context.CancelFunc
	reconnectC chan struct{}
	closed     bool
}

// NewRabbitMQManager 创建RabbitMQ管理器
func NewRabbitMQManager(ctx context.Context, config *RabbitMQConfig, logger Logger) (*RabbitMQManager, error) {
	mqCtx, cancel := context.WithCancel(ctx)
	
	manager := &RabbitMQManager{
		config:     config,
		logger:     logger,
		ctx:        mqCtx,
		cancel:     cancel,
		reconnectC: make(chan struct{}, 1),
	}

	if err := manager.connect(); err != nil {
		cancel()
		return nil, err
	}

	// 启动重连监控
	go manager.monitorConnection()

	return manager, nil
}

// connect 连接到RabbitMQ
func (m *RabbitMQManager) connect() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	url := fmt.Sprintf("amqp://%s:%s@%s:%d%s",
		m.config.User,
		m.config.Pass,
		m.config.Host,
		m.config.Port,
		m.config.VirtualHost,
	)

	conn, err := amqp.Dial(url)
	if err != nil {
		m.logger.Error("连接RabbitMQ失败", zap.Error(err))
		return fmt.Errorf("连接RabbitMQ失败: %v", err)
	}

	channel, err := conn.Channel()
	if err != nil {
		conn.Close()
		m.logger.Error("创建RabbitMQ通道失败", zap.Error(err))
		return fmt.Errorf("创建RabbitMQ通道失败: %v", err)
	}

	// 设置QoS
	if err := channel.Qos(m.config.PrefetchCount, m.config.PrefetchSize, false); err != nil {
		channel.Close()
		conn.Close()
		m.logger.Error("设置QoS失败", zap.Error(err))
		return fmt.Errorf("设置QoS失败: %v", err)
	}

	// 声明普通交换机和队列
	if err := m.declareExchangeAndQueue(channel, m.config.Exchange, m.config.QueueName, false); err != nil {
		channel.Close()
		conn.Close()
		return err
	}

	// 声明延迟交换机和队列
	if err := m.declareDelayExchangeAndQueue(channel); err != nil {
		channel.Close()
		conn.Close()
		return err
	}

	m.conn = conn
	m.channel = channel

	m.logger.Info("RabbitMQ连接成功",
		zap.String("host", m.config.Host),
		zap.Int("port", m.config.Port),
	)

	return nil
}

// declareExchangeAndQueue 声明交换机和队列
func (m *RabbitMQManager) declareExchangeAndQueue(channel *amqp.Channel, exchange, queueName string, isDelay bool) error {
	// 声明交换机
	args := amqp.Table{}
	if isDelay {
		args["x-delayed-type"] = "direct"
	}

	err := channel.ExchangeDeclare(
		exchange,
		"direct",
		true,  // durable
		false, // auto-deleted
		false, // internal
		false, // no-wait
		args,
	)
	if err != nil {
		m.logger.Error("声明交换机失败",
			zap.String("exchange", exchange),
			zap.Error(err),
		)
		return fmt.Errorf("声明交换机失败: %v", err)
	}

	// 声明队列
	queueArgs := amqp.Table{}
	_, err = channel.QueueDeclare(
		queueName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		queueArgs,
	)
	if err != nil {
		m.logger.Error("声明队列失败",
			zap.String("queue", queueName),
			zap.Error(err),
		)
		return fmt.Errorf("声明队列失败: %v", err)
	}

	// 绑定队列到交换机
	err = channel.QueueBind(
		queueName,
		queueName, // routing key
		exchange,
		false,
		nil,
	)
	if err != nil {
		m.logger.Error("绑定队列失败",
			zap.String("queue", queueName),
			zap.String("exchange", exchange),
			zap.Error(err),
		)
		return fmt.Errorf("绑定队列失败: %v", err)
	}

	return nil
}

// declareDelayExchangeAndQueue 声明延迟交换机和队列
func (m *RabbitMQManager) declareDelayExchangeAndQueue(channel *amqp.Channel) error {
	// 声明延迟交换机
	err := channel.ExchangeDeclare(
		m.config.DelayExchange,
		"x-delayed-message",
		true,
		false,
		false,
		false,
		amqp.Table{
			"x-delayed-type": "direct",
		},
	)
	if err != nil {
		// 如果不支持延迟插件，使用TTL+死信队列方式
		m.logger.Warn("延迟交换机插件不可用，使用TTL+死信队列方式")
		return m.declareDelayQueueWithDLX(channel)
	}

	// 声明延迟队列
	_, err = channel.QueueDeclare(
		m.config.DelayQueueName,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		m.logger.Error("声明延迟队列失败", zap.Error(err))
		return fmt.Errorf("声明延迟队列失败: %v", err)
	}

	// 绑定延迟队列
	err = channel.QueueBind(
		m.config.DelayQueueName,
		m.config.DelayQueueName,
		m.config.DelayExchange,
		false,
		nil,
	)
	if err != nil {
		m.logger.Error("绑定延迟队列失败", zap.Error(err))
		return fmt.Errorf("绑定延迟队列失败: %v", err)
	}

	return nil
}

// declareDelayQueueWithDLX 使用死信队列实现延迟
func (m *RabbitMQManager) declareDelayQueueWithDLX(channel *amqp.Channel) error {
	// 声明死信交换机
	dlxExchange := m.config.DelayExchange + ".dlx"
	err := channel.ExchangeDeclare(dlxExchange, "direct", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("声明死信交换机失败: %v", err)
	}

	// 声明目标队列（延迟消息最终到达的队列）
	_, err = channel.QueueDeclare(m.config.DelayQueueName, true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("声明目标队列失败: %v", err)
	}

	// 绑定目标队列到死信交换机
	err = channel.QueueBind(m.config.DelayQueueName, m.config.DelayQueueName, dlxExchange, false, nil)
	if err != nil {
		return fmt.Errorf("绑定目标队列失败: %v", err)
	}

	return nil
}

// monitorConnection 监控连接状态
func (m *RabbitMQManager) monitorConnection() {
	for {
		select {
		case <-m.ctx.Done():
			return
		case <-m.reconnectC:
			m.reconnect()
		}
	}
}

// reconnect 重新连接
func (m *RabbitMQManager) reconnect() {
	m.mu.Lock()
	if m.closed {
		m.mu.Unlock()
		return
	}
	m.mu.Unlock()

	for i := 0; i < m.config.RetryCount; i++ {
		select {
		case <-m.ctx.Done():
			return
		default:
		}

		m.logger.Info("尝试重新连接RabbitMQ",
			zap.Int("attempt", i+1),
			zap.Int("max_retries", m.config.RetryCount),
		)

		if err := m.connect(); err != nil {
			m.logger.Error("重新连接失败", zap.Error(err))
			time.Sleep(time.Duration(m.config.ReconnectInterval) * time.Second)
			continue
		}

		m.logger.Info("重新连接成功")
		return
	}

	m.logger.Error("重新连接失败，已达到最大重试次数")
}

// Publish 发布消息到普通队列
func (m *RabbitMQManager) Publish(ctx context.Context, topic string, message interface{}) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.channel == nil {
		return fmt.Errorf("RabbitMQ通道未初始化")
	}

	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("序列化消息失败: %v", err)
	}

	err = m.channel.PublishWithContext(
		ctx,
		m.config.Exchange,
		topic,
		false, // mandatory
		false, // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
		},
	)

	if err != nil {
		m.logger.Error("发布消息失败",
			zap.String("topic", topic),
			zap.Error(err),
		)
		// 触发重连
		select {
		case m.reconnectC <- struct{}{}:
		default:
		}
		return fmt.Errorf("发布消息失败: %v", err)
	}

	m.logger.Debug("消息发布成功",
		zap.String("topic", topic),
	)

	return nil
}

// PublishDelay 发布延迟消息（使用 TTL + 死信队列方式）
func (m *RabbitMQManager) PublishDelay(ctx context.Context, topic string, message interface{}, delay time.Duration) error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.channel == nil {
		return fmt.Errorf("RabbitMQ通道未初始化")
	}

	body, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("序列化消息失败: %v", err)
	}

	// 创建临时延迟队列（使用 TTL + 死信队列）
	delayQueueName := fmt.Sprintf("%s.delay.%d", topic, delay.Milliseconds())

	// 声明临时延迟队列，设置 TTL 和死信交换机
	_, err = m.channel.QueueDeclare(
		delayQueueName,
		true,  // durable
		false, // auto-delete
		false, // exclusive
		false, // no-wait
		amqp.Table{
			"x-message-ttl":             int64(delay.Milliseconds()),
			"x-dead-letter-exchange":    m.config.Exchange,
			"x-dead-letter-routing-key": topic,
		},
	)
	if err != nil {
		m.logger.Error("声明延迟队列失败", zap.Error(err))
		return fmt.Errorf("声明延迟队列失败: %v", err)
	}

	// 发布消息到临时延迟队列
	err = m.channel.PublishWithContext(
		ctx,
		"",             // 直接发送到队列，不经过交换机
		delayQueueName, // routing key = 队列名
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent,
			Timestamp:    time.Now(),
		},
	)

	if err != nil {
		m.logger.Error("发布延迟消息失败",
			zap.String("topic", topic),
			zap.Duration("delay", delay),
			zap.Error(err),
		)
		select {
		case m.reconnectC <- struct{}{}:
		default:
		}
		return fmt.Errorf("发布延迟消息失败: %v", err)
	}

	m.logger.Debug("延迟消息发布成功",
		zap.String("topic", topic),
		zap.Duration("delay", delay),
		zap.String("delay_queue", delayQueueName),
	)

	return nil
}

// Close 关闭连接
func (m *RabbitMQManager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.closed = true
	m.cancel()

	if m.channel != nil {
		if err := m.channel.Close(); err != nil {
			m.logger.Error("关闭RabbitMQ通道失败", zap.Error(err))
		}
	}

	if m.conn != nil {
		if err := m.conn.Close(); err != nil {
			m.logger.Error("关闭RabbitMQ连接失败", zap.Error(err))
		}
	}

	m.logger.Info("RabbitMQ连接已关闭")
	return nil
}

// GetChannel 获取通道（用于消费者）
func (m *RabbitMQManager) GetChannel() (*amqp.Channel, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.channel == nil {
		return nil, fmt.Errorf("RabbitMQ通道未初始化")
	}

	return m.channel, nil
}
