package types

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

// RabbitMQConfig RabbitMQ配置
type RabbitMQConfig struct {
	Host                string
	Port                int
	User                string
	Pass                string
	VirtualHost         string
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

// 消息结构体
type RabbitMQMessage struct {
	ID        string                 `json:"id"`
	Topic     string                 `json:"topic"`
	Data      map[string]interface{} `json:"data"`
	Timestamp int64                  `json:"timestamp"`
	Retry     int                    `json:"retry"`
}

// RabbitMQManager RabbitMQ管理器
type RabbitMQManager struct {
	config          *RabbitMQConfig
	conn            *amqp.Connection
	channel         *amqp.Channel
	delayChannel    *amqp.Channel
	notifyConnClose chan *amqp.Error
	notifyChanClose chan *amqp.Error
	notifyConfirm   chan amqp.Confirmation
	isConnected     bool
	done            chan bool
}

// NewRabbitMQManager 创建RabbitMQ管理器
func NewRabbitMQManager(config *RabbitMQConfig) *RabbitMQManager {
	manager := &RabbitMQManager{
		config: config,
		done:   make(chan bool),
	}

	// 启动连接
	go manager.handleReconnect()

	return manager
}

// Push 推送消息
func (r *RabbitMQManager) Push(ctx context.Context, topic string, message interface{}) error {
	if !r.isConnected {
		return fmt.Errorf("RabbitMQ未连接")
	}

	// 转换消息
	data, err := convertToMapOrString(message)
	if err != nil {
		return fmt.Errorf("消息转换失败: %v", err)
	}

	// 创建消息
	msg := RabbitMQMessage{
		ID:        generateUniqueID(),
		Topic:     topic,
		Data:      data,
		Timestamp: time.Now().UnixNano(),
		Retry:     0,
	}

	// 序列化消息
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("消息序列化失败: %v", err)
	}

	// 发布消息
	err = r.channel.PublishWithContext(ctx,
		r.config.Exchange, // 交换机
		topic,             // 路由键
		false,             // 强制
		false,             // 立即
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // 持久化
		},
	)

	if err != nil {
		return fmt.Errorf("消息发布失败: %v", err)
	}

	return nil
}

// Pop 获取消息
func (r *RabbitMQManager) Pop(ctx context.Context, topic string) (interface{}, error) {
	if !r.isConnected {
		return nil, fmt.Errorf("RabbitMQ未连接")
	}

	// 声明队列
	queue, err := r.channel.QueueDeclare(
		topic, // 队列名
		true,  // 持久化
		false, // 自动删除
		false, // 独占
		false, // 阻塞等待
		nil,   // 参数
	)
	if err != nil {
		return nil, fmt.Errorf("队列声明失败: %v", err)
	}

	// 获取消息
	msg, ok, err := r.channel.Get(
		queue.Name, // 队列名
		true,       // 自动确认
	)
	if err != nil {
		return nil, fmt.Errorf("获取消息失败: %v", err)
	}
	if !ok {
		return nil, nil // 队列为空
	}

	// 解析消息
	var result RabbitMQMessage
	if err := json.Unmarshal(msg.Body, &result); err != nil {
		return nil, fmt.Errorf("消息解析失败: %v", err)
	}

	return result.Data, nil
}

// Close 关闭连接
func (r *RabbitMQManager) Close() error {
	if !r.isConnected {
		return nil
	}

	r.done <- true

	// 关闭通道
	if r.channel != nil {
		err := r.channel.Close()
		if err != nil {
			return fmt.Errorf("关闭通道失败: %v", err)
		}
	}

	// 关闭延迟通道
	if r.delayChannel != nil {
		err := r.delayChannel.Close()
		if err != nil {
			return fmt.Errorf("关闭延迟通道失败: %v", err)
		}
	}

	// 关闭连接
	if r.conn != nil {
		err := r.conn.Close()
		if err != nil {
			return fmt.Errorf("关闭连接失败: %v", err)
		}
	}

	r.isConnected = false
	return nil
}

// 处理重连
func (r *RabbitMQManager) handleReconnect() {
	for {
		r.isConnected = false

		// 连接
		conn, err := r.connect()
		if err != nil {
			select {
			case <-r.done:
				return
			case <-time.After(time.Duration(r.config.ReconnectInterval) * time.Second):
				continue
			}
		}

		// 设置连接
		r.conn = conn
		r.notifyConnClose = make(chan *amqp.Error)
		r.conn.NotifyClose(r.notifyConnClose)

		// 初始化通道
		if err := r.initChannel(); err != nil {
			select {
			case <-r.done:
				return
			case <-time.After(time.Duration(r.config.ReconnectInterval) * time.Second):
				continue
			}
		}

		// 初始化延迟通道
		if err := r.initDelayChannel(); err != nil {
			select {
			case <-r.done:
				return
			case <-time.After(time.Duration(r.config.ReconnectInterval) * time.Second):
				continue
			}
		}

		r.isConnected = true

		select {
		case <-r.done:
			return
		case <-r.notifyConnClose:
			// 连接断开，重新连接
		}
	}
}

// 连接
func (r *RabbitMQManager) connect() (*amqp.Connection, error) {
	// 构建连接URL
	url := fmt.Sprintf("amqp://%s:%s@%s:%d%s",
		r.config.User,
		r.config.Pass,
		r.config.Host,
		r.config.Port,
		r.config.VirtualHost,
	)

	// 设置连接参数
	config := amqp.Config{
		Heartbeat: time.Duration(r.config.HeartbeatInterval) * time.Second,
		Dial:      amqp.DefaultDial(time.Duration(r.config.ConnectionTimeout) * time.Second),
	}

	// 连接
	return amqp.DialConfig(url, config)
}

// 初始化通道
func (r *RabbitMQManager) initChannel() error {
	channel, err := r.conn.Channel()
	if err != nil {
		return fmt.Errorf("创建通道失败: %v", err)
	}

	// 设置QoS
	err = channel.Qos(
		r.config.PrefetchCount, // 预取数量
		r.config.PrefetchSize,  // 预取大小
		false,                  // 全局
	)
	if err != nil {
		return fmt.Errorf("设置QoS失败: %v", err)
	}

	// 声明交换机
	err = channel.ExchangeDeclare(
		r.config.Exchange, // 交换机名
		"direct",          // 类型
		true,              // 持久化
		false,             // 自动删除
		false,             // 内部
		false,             // 阻塞等待
		nil,               // 参数
	)
	if err != nil {
		return fmt.Errorf("交换机声明失败: %v", err)
	}

	// 声明队列
	_, err = channel.QueueDeclare(
		r.config.QueueName, // 队列名
		true,               // 持久化
		false,              // 自动删除
		false,              // 独占
		false,              // 阻塞等待
		nil,                // 参数
	)
	if err != nil {
		return fmt.Errorf("队列声明失败: %v", err)
	}

	// 绑定队列
	err = channel.QueueBind(
		r.config.QueueName, // 队列名
		"#",                // 路由键
		r.config.Exchange,  // 交换机
		false,              // 阻塞等待
		nil,                // 参数
	)
	if err != nil {
		return fmt.Errorf("队列绑定失败: %v", err)
	}

	// 设置确认模式
	if r.config.EnableConfirmation {
		err = channel.Confirm(false)
		if err != nil {
			return fmt.Errorf("设置确认模式失败: %v", err)
		}
		r.notifyConfirm = make(chan amqp.Confirmation)
		channel.NotifyPublish(r.notifyConfirm)
	}

	r.channel = channel
	r.notifyChanClose = make(chan *amqp.Error)
	r.channel.NotifyClose(r.notifyChanClose)

	return nil
}

// 初始化延迟通道
func (r *RabbitMQManager) initDelayChannel() error {
	channel, err := r.conn.Channel()
	if err != nil {
		return fmt.Errorf("创建延迟通道失败: %v", err)
	}

	// 声明延迟交换机
	err = channel.ExchangeDeclare(
		r.config.DelayExchange, // 交换机名
		"x-delayed-message",    // 类型
		true,                   // 持久化
		false,                  // 自动删除
		false,                  // 内部
		false,                  // 阻塞等待
		amqp.Table{
			"x-delayed-type": "direct",
		}, // 参数
	)
	if err != nil {
		// 如果不支持延迟交换机，使用普通交换机
		err = channel.ExchangeDeclare(
			r.config.DelayExchange, // 交换机名
			"direct",               // 类型
			true,                   // 持久化
			false,                  // 自动删除
			false,                  // 内部
			false,                  // 阻塞等待
			nil,                    // 参数
		)
		if err != nil {
			return fmt.Errorf("延迟交换机声明失败: %v", err)
		}
	}

	// 声明延迟队列
	_, err = channel.QueueDeclare(
		r.config.DelayQueueName, // 队列名
		true,                    // 持久化
		false,                   // 自动删除
		false,                   // 独占
		false,                   // 阻塞等待
		nil,                     // 参数
	)
	if err != nil {
		return fmt.Errorf("延迟队列声明失败: %v", err)
	}

	// 绑定延迟队列
	err = channel.QueueBind(
		r.config.DelayQueueName, // 队列名
		"#",                     // 路由键
		r.config.DelayExchange,  // 交换机
		false,                   // 阻塞等待
		nil,                     // 参数
	)
	if err != nil {
		return fmt.Errorf("延迟队列绑定失败: %v", err)
	}

	r.delayChannel = channel
	return nil
}

// DelayQueueManager RabbitMQ延迟队列管理器
type DelayQueueManager struct {
	rabbit *RabbitMQManager
}

// NewDelayQueueManager 创建RabbitMQ延迟队列管理器
func NewDelayQueueManager(rabbit *RabbitMQManager) *DelayQueueManager {
	return &DelayQueueManager{
		rabbit: rabbit,
	}
}

// Push 推送延迟消息
func (d *DelayQueueManager) Push(ctx context.Context, topic string, message interface{}, delay time.Duration) error {
	if !d.rabbit.isConnected {
		return fmt.Errorf("RabbitMQ未连接")
	}

	// 转换消息
	data, err := convertToMapOrString(message)
	if err != nil {
		return fmt.Errorf("消息转换失败: %v", err)
	}

	// 创建消息
	msg := RabbitMQMessage{
		ID:        generateUniqueID(),
		Topic:     topic,
		Data:      data,
		Timestamp: time.Now().UnixNano(),
		Retry:     0,
	}

	// 序列化消息
	body, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("消息序列化失败: %v", err)
	}

	// 发布延迟消息
	headers := amqp.Table{}
	if delay > 0 {
		headers["x-delay"] = int(delay.Milliseconds())
	}

	err = d.rabbit.delayChannel.PublishWithContext(ctx,
		d.rabbit.config.DelayExchange, // 交换机
		topic,                         // 路由键
		false,                         // 强制
		false,                         // 立即
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // 持久化
			Headers:      headers,
		},
	)

	if err != nil {
		return fmt.Errorf("延迟消息发布失败: %v", err)
	}

	return nil
}

// Pop 获取延迟消息
func (d *DelayQueueManager) Pop(ctx context.Context, topic string) (interface{}, error) {
	if !d.rabbit.isConnected {
		return nil, fmt.Errorf("RabbitMQ未连接")
	}

	// 声明队列
	queue, err := d.rabbit.delayChannel.QueueDeclare(
		topic, // 队列名
		true,  // 持久化
		false, // 自动删除
		false, // 独占
		false, // 阻塞等待
		nil,   // 参数
	)
	if err != nil {
		return nil, fmt.Errorf("延迟队列声明失败: %v", err)
	}

	// 获取消息
	msg, ok, err := d.rabbit.delayChannel.Get(
		queue.Name, // 队列名
		true,       // 自动确认
	)
	if err != nil {
		return nil, fmt.Errorf("获取延迟消息失败: %v", err)
	}
	if !ok {
		return nil, nil // 队列为空
	}

	// 解析消息
	var result RabbitMQMessage
	if err := json.Unmarshal(msg.Body, &result); err != nil {
		return nil, fmt.Errorf("延迟消息解析失败: %v", err)
	}

	return result.Data, nil
}

// Close 关闭连接
func (d *DelayQueueManager) Close() error {
	// 使用RabbitMQ管理器关闭连接
	return nil
}

// 生成唯一ID
func generateUniqueID() string {
	return fmt.Sprintf("%d%d", time.Now().UnixNano(), time.Now().Unix())
}

// 转换消息为map或string
func convertToMapOrString(message interface{}) (map[string]interface{}, error) {
	var data map[string]interface{}

	// 如果是字符串，直接包装
	if str, ok := message.(string); ok {
		data = map[string]interface{}{
			"message": str,
		}
		return data, nil
	}

	// 如果是map，直接使用
	if mapData, ok := message.(map[string]interface{}); ok {
		return mapData, nil
	}

	// 其他类型，尝试JSON序列化
	bytes, err := json.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("消息序列化失败: %v", err)
	}

	// 反序列化为map
	err = json.Unmarshal(bytes, &data)
	if err != nil {
		// 如果反序列化失败，包装为字符串
		data = map[string]interface{}{
			"message": string(bytes),
		}
	}

	return data, nil
}
