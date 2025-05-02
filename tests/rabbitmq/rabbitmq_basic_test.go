package rabbitmq

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/lvjiaben/go-wheel/pkg/container"
	"github.com/lvjiaben/go-wheel/pkg/types"
	"github.com/streadway/amqp"
)

// RabbitMQ连接配置
type RabbitMQConfig struct {
	Host     string
	Port     int
	Username string
	Password string
	VHost    string
}

// 获取RabbitMQ连接URL
func (c *RabbitMQConfig) URL() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
		c.Username, c.Password, c.Host, c.Port, c.VHost)
}

// RabbitMQ客户端
type RabbitMQClient struct {
	config RabbitMQConfig
	conn   *amqp.Connection
	ch     *amqp.Channel
}

// 创建新的RabbitMQ客户端
func NewRabbitMQClient(config RabbitMQConfig) (*RabbitMQClient, error) {
	client := &RabbitMQClient{
		config: config,
	}

	// 连接到RabbitMQ服务器
	var err error
	client.conn, err = amqp.Dial(config.URL())
	if err != nil {
		return nil, fmt.Errorf("无法连接到RabbitMQ: %v", err)
	}

	// 创建通道
	client.ch, err = client.conn.Channel()
	if err != nil {
		client.conn.Close()
		return nil, fmt.Errorf("无法创建通道: %v", err)
	}

	return client, nil
}

// 关闭连接和通道
func (c *RabbitMQClient) Close() {
	if c.ch != nil {
		c.ch.Close()
	}
	if c.conn != nil {
		c.conn.Close()
	}
}

// 声明队列
func (c *RabbitMQClient) DeclareQueue(queueName string) (amqp.Queue, error) {
	return c.ch.QueueDeclare(
		queueName, // 队列名称
		true,      // 持久化
		false,     // 不自动删除
		false,     // 非独占
		false,     // 不等待服务器响应
		nil,       // 额外参数
	)
}

// 发送消息到队列
func (c *RabbitMQClient) PublishMessage(queueName, message string) error {
	return c.ch.Publish(
		"",        // 交换机
		queueName, // 路由键
		false,     // 强制发送
		false,     // 立即发送
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         []byte(message),
			DeliveryMode: amqp.Persistent, // 消息持久化
		},
	)
}

// 从队列消费消息
func (c *RabbitMQClient) ConsumeMessages(queueName string) (<-chan amqp.Delivery, error) {
	return c.ch.Consume(
		queueName, // 队列名称
		"",        // 消费者标识（空字符串会自动生成）
		true,      // 自动确认
		false,     // 非独占
		false,     // 不阻止其他消费者
		false,     // 不等待服务器响应
		nil,       // 额外参数
	)
}

// 扩展 MessageQueue 接口 - 用于测试
type TestMessageQueue struct {
	rabbitClient *RabbitMQClient    // 使用我们自己的客户端
	baseQueue    types.MessageQueue // 存储原始的消息队列
	config       *types.RabbitMQConfig
}

// 创建测试用消息队列
func NewTestMessageQueue(baseQueue types.MessageQueue, config *types.RabbitMQConfig) (*TestMessageQueue, error) {
	// 创建基本RabbitMQ客户端
	rabbitmqConfig := RabbitMQConfig{
		Host:     config.Host,
		Port:     config.Port,
		Username: config.User,
		Password: config.Pass,
		VHost:    config.VirtualHost,
	}

	client, err := NewRabbitMQClient(rabbitmqConfig)
	if err != nil {
		return nil, fmt.Errorf("创建RabbitMQ客户端失败: %v", err)
	}

	return &TestMessageQueue{
		rabbitClient: client,
		baseQueue:    baseQueue,
		config:       config,
	}, nil
}

// Push 转发到基础消息队列
func (t *TestMessageQueue) Push(ctx context.Context, topic string, message interface{}) error {
	return t.baseQueue.Push(ctx, topic, message)
}

// 为测试消息队列添加 Consume 方法
func (t *TestMessageQueue) Consume(ctx context.Context, topic string, consumerName string) (<-chan amqp.Delivery, error) {
	// 声明队列
	queue, err := t.rabbitClient.ch.QueueDeclare(
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

	// 绑定队列到交换机
	err = t.rabbitClient.ch.QueueBind(
		queue.Name,        // 队列名
		topic,             // 路由键
		t.config.Exchange, // 交换机
		false,             // 阻塞等待
		nil,               // 参数
	)
	if err != nil {
		return nil, fmt.Errorf("队列绑定失败: %v", err)
	}

	// 消费消息
	return t.rabbitClient.ch.Consume(
		queue.Name,   // 队列名
		consumerName, // 消费者名称
		false,        // 手动确认
		false,        // 非独占
		false,        // 不阻止其他消费者
		false,        // 不等待
		nil,          // 参数
	)
}

// 清理资源
func (t *TestMessageQueue) Close() error {
	if t.rabbitClient != nil {
		t.rabbitClient.Close()
	}
	return t.baseQueue.Close()
}

// 模拟消息服务
type NotificationService struct {
	container *container.Container
	mqService *TestMessageQueue // 使用测试消息队列
}

// 创建新的通知服务
func NewNotificationService(c *container.Container, mqService *TestMessageQueue) *NotificationService {
	return &NotificationService{
		container: c,
		mqService: mqService,
	}
}

// 发送通知
func (s *NotificationService) SendNotification(userID string, message string) error {
	// 创建通知数据
	notification := map[string]interface{}{
		"user_id":    userID,
		"message":    message,
		"created_at": time.Now().Format(time.RFC3339),
	}

	// 序列化为JSON
	data, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	// 通过消息队列发送消息
	return s.mqService.Push(
		context.Background(),
		"notifications", // 主题
		data,            // 消息内容
	)
}

// 设置通知消费者
func (s *NotificationService) SetupConsumers() error {
	// 从消息队列获取通知消息
	messages, err := s.mqService.Consume(
		context.Background(),
		"notifications",         // 主题
		"notification-consumer", // 消费者名称
	)
	if err != nil {
		return err
	}

	// 处理接收到的消息
	go func() {
		for msg := range messages {
			// 处理通知消息
			var notification map[string]interface{}
			if err := json.Unmarshal(msg.Body, &notification); err != nil {
				log.Printf("解析通知消息失败: %v", err)
				continue
			}

			log.Printf("收到通知: 用户ID=%s, 消息=%s",
				notification["user_id"],
				notification["message"])

			// 确认消息已处理
			msg.Ack(false)
		}
	}()

	return nil
}

// 使用依赖注入方式测试RabbitMQ基本功能
func TestRabbitMQBasicWithDI(t *testing.T) {
	// 如果没有RabbitMQ环境，则跳过测试
	t.Skip("需要RabbitMQ环境才能运行此测试")

	// 创建容器
	c := container.NewContainer()

	// 初始化消息队列配置
	config := &types.RabbitMQConfig{
		Host:                "localhost",
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

	// 创建RabbitMQ管理器
	mqManager := types.NewRabbitMQManager(config)

	// 创建测试消息队列 - 扩展了Consume方法
	testMQService, err := NewTestMessageQueue(mqManager, config)
	if err != nil {
		t.Fatalf("创建测试消息队列失败: %v", err)
	}

	// 设置消息队列到容器
	c.SetMessageQueue(mqManager)

	// 创建通知服务，使用测试消息队列
	notificationService := NewNotificationService(c, testMQService)

	// 设置消费者
	if err := notificationService.SetupConsumers(); err != nil {
		t.Fatalf("设置消费者失败: %v", err)
	}

	// 发送测试通知
	userID := "user123"
	message := "这是一条测试通知"

	if err := notificationService.SendNotification(userID, message); err != nil {
		t.Fatalf("发送通知失败: %v", err)
	}
	t.Logf("成功发送通知给用户 %s: %s", userID, message)

	// 等待消息被处理
	time.Sleep(1 * time.Second)
}

// 订单服务示例
type OrderService struct {
	container *container.Container
	mqService *TestMessageQueue // 使用测试消息队列
}

// 创建新的订单服务
func NewOrderService(c *container.Container, mqService *TestMessageQueue) *OrderService {
	return &OrderService{
		container: c,
		mqService: mqService,
	}
}

// 订单数据结构
type Order struct {
	ID        string    `json:"id"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

// 创建订单
func (s *OrderService) CreateOrder(order *Order) error {
	// 在实际应用中，这里会保存订单到数据库
	order.CreatedAt = time.Now()
	order.Status = "pending"

	// 记录日志
	log.Printf("创建订单: ID=%s, 金额=%.2f", order.ID, order.Amount)

	// 发送订单创建通知
	orderData, err := json.Marshal(order)
	if err != nil {
		return err
	}

	// 使用消息队列发送订单创建事件
	return s.mqService.Push(
		context.Background(),
		"order_created",
		orderData,
	)
}

// 设置订单处理消费者
func (s *OrderService) SetupConsumers() error {
	// 订单创建事件消费者
	messages, err := s.mqService.Consume(
		context.Background(),
		"order_created",
		"order-processor",
	)
	if err != nil {
		return err
	}

	go func() {
		for msg := range messages {
			var order Order
			if err := json.Unmarshal(msg.Body, &order); err != nil {
				log.Printf("解析订单数据失败: %v", err)
				continue
			}

			log.Printf("处理新订单: ID=%s, 金额=%.2f", order.ID, order.Amount)

			// 在实际应用中，这里会进行订单处理逻辑
			// ...

			// 确认消息
			msg.Ack(false)
		}
	}()

	return nil
}

// 测试多个消费者场景
func TestMultipleConsumers(t *testing.T) {
	// 如果没有RabbitMQ环境，则跳过测试
	t.Skip("需要RabbitMQ环境才能运行此测试")

	// 创建容器
	c := container.NewContainer()

	// 初始化消息队列
	config := &types.RabbitMQConfig{
		Host:                "localhost",
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

	// 创建RabbitMQ管理器
	mqManager := types.NewRabbitMQManager(config)

	// 创建测试消息队列
	testMQService, err := NewTestMessageQueue(mqManager, config)
	if err != nil {
		t.Fatalf("创建测试消息队列失败: %v", err)
	}

	// 设置消息队列到容器
	c.SetMessageQueue(mqManager)

	// 创建订单服务和通知服务
	orderService := NewOrderService(c, testMQService)
	notificationService := NewNotificationService(c, testMQService)

	// 设置消费者
	if err := orderService.SetupConsumers(); err != nil {
		t.Fatalf("设置订单消费者失败: %v", err)
	}

	if err := notificationService.SetupConsumers(); err != nil {
		t.Fatalf("设置通知消费者失败: %v", err)
	}

	// 创建测试订单
	order := &Order{
		ID:     "order-123",
		Amount: 199.99,
	}

	// 创建订单并发送事件
	if err := orderService.CreateOrder(order); err != nil {
		t.Fatalf("创建订单失败: %v", err)
	}

	// 发送订单通知
	if err := notificationService.SendNotification("customer-456",
		fmt.Sprintf("您的订单 %s 已创建, 金额: ¥%.2f", order.ID, order.Amount)); err != nil {
		t.Fatalf("发送订单通知失败: %v", err)
	}

	// 等待消息处理
	time.Sleep(1 * time.Second)

	t.Logf("成功创建订单并发送通知")
}

// 实际使用RabbitMQ的示例
func ExampleRabbitMQUsage() {
	// 在实际应用中，容器会在应用启动时初始化
	c := container.NewContainer()

	// 创建RabbitMQ管理器
	config := &types.RabbitMQConfig{
		Host:        "localhost",
		Port:        5672,
		User:        "guest",
		Pass:        "guest",
		VirtualHost: "/",
		QueueName:   "go-admin",
		Exchange:    "go-admin-exchange",
	}
	mqManager := types.NewRabbitMQManager(config)
	testMQService, err := NewTestMessageQueue(mqManager, config)
	if err != nil {
		log.Fatalf("创建测试消息队列失败: %v", err)
	}

	// 服务会通过依赖注入获取容器
	orderService := NewOrderService(c, testMQService)

	// 设置消费者
	orderService.SetupConsumers()

	// 创建订单
	order := &Order{
		ID:     "example-order",
		Amount: 99.99,
	}

	// 处理订单
	if err := orderService.CreateOrder(order); err != nil {
		log.Printf("创建订单失败: %v", err)
		return
	}

	fmt.Println("订单已创建并发送到消息队列")
}
