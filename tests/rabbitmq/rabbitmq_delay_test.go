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

// 延迟队列管理器
type DelayQueueManager struct {
	client        *RabbitMQClient
	queueName     string
	exchangeName  string
	deadQueueName string
	deadExchange  string
	delayTime     time.Duration
}

// 创建延迟队列管理器
func NewDelayQueueManager(client *RabbitMQClient, queueName string, delayTime time.Duration) (*DelayQueueManager, error) {
	manager := &DelayQueueManager{
		client:        client,
		queueName:     queueName,
		exchangeName:  queueName + "-exchange",
		deadQueueName: queueName + "-dlx",
		deadExchange:  queueName + "-dlx-exchange",
		delayTime:     delayTime,
	}

	// 声明交换机
	err := client.ch.ExchangeDeclare(
		manager.exchangeName, // 交换机名称
		"direct",             // 类型
		true,                 // 持久化
		false,                // 不自动删除
		false,                // 非内部使用
		false,                // 不等待
		nil,                  // 额外参数
	)
	if err != nil {
		return nil, fmt.Errorf("声明交换机失败: %v", err)
	}

	// 声明死信交换机
	err = client.ch.ExchangeDeclare(
		manager.deadExchange, // 交换机名称
		"direct",             // 类型
		true,                 // 持久化
		false,                // 不自动删除
		false,                // 非内部使用
		false,                // 不等待
		nil,                  // 额外参数
	)
	if err != nil {
		return nil, fmt.Errorf("声明死信交换机失败: %v", err)
	}

	// 声明死信队列（用于接收过期消息）
	_, err = client.ch.QueueDeclare(
		manager.deadQueueName, // 队列名称
		true,                  // 持久化
		false,                 // 不自动删除
		false,                 // 非独占
		false,                 // 不等待
		nil,                   // 额外参数
	)
	if err != nil {
		return nil, fmt.Errorf("声明死信队列失败: %v", err)
	}

	// 绑定死信队列到死信交换机
	err = client.ch.QueueBind(
		manager.deadQueueName, // 队列名称
		"",                    // 路由键
		manager.deadExchange,  // 交换机名称
		false,                 // 不等待
		nil,                   // 额外参数
	)
	if err != nil {
		return nil, fmt.Errorf("绑定死信队列失败: %v", err)
	}

	// 声明延迟队列（主队列）
	args := amqp.Table{
		"x-dead-letter-exchange":    manager.deadExchange,
		"x-dead-letter-routing-key": "",
		"x-message-ttl":             int32(delayTime.Milliseconds()),
	}
	_, err = client.ch.QueueDeclare(
		manager.queueName, // 队列名称
		true,              // 持久化
		false,             // 不自动删除
		false,             // 非独占
		false,             // 不等待
		args,              // 额外参数
	)
	if err != nil {
		return nil, fmt.Errorf("声明延迟队列失败: %v", err)
	}

	// 绑定延迟队列到正常交换机
	err = client.ch.QueueBind(
		manager.queueName,    // 队列名称
		"",                   // 路由键
		manager.exchangeName, // 交换机名称
		false,                // 不等待
		nil,                  // 额外参数
	)
	if err != nil {
		return nil, fmt.Errorf("绑定延迟队列失败: %v", err)
	}

	return manager, nil
}

// 发布延迟消息
func (m *DelayQueueManager) PublishDelayMessage(message string) error {
	return m.client.ch.Publish(
		m.exchangeName, // 交换机
		"",             // 路由键
		false,          // 强制发送
		false,          // 立即发送
		amqp.Publishing{
			ContentType:  "text/plain",
			Body:         []byte(message),
			DeliveryMode: amqp.Persistent,
		},
	)
}

// 消费处理后的消息（从死信队列消费）
func (m *DelayQueueManager) ConsumeProcessedMessages() (<-chan amqp.Delivery, error) {
	return m.client.ch.Consume(
		m.deadQueueName, // 队列名称
		"",              // 消费者标识
		true,            // 自动确认
		false,           // 非独占
		false,           // 不阻止其他消费者
		false,           // 不等待
		nil,             // 额外参数
	)
}

// 扩展DelayQueue接口 - 用于测试
type TestDelayQueue struct {
	rabbitClient *RabbitMQClient  // 使用我们自己的客户端
	baseQueue    types.DelayQueue // 存储原始的延迟队列
	config       *types.RabbitMQConfig
	delayManager *DelayQueueManager
}

// 创建测试用延迟队列
func NewTestDelayQueue(baseQueue types.DelayQueue, config *types.RabbitMQConfig) (*TestDelayQueue, error) {
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

	// 创建延迟队列管理器
	delayManager, err := NewDelayQueueManager(client, config.DelayQueueName, 5*time.Second) // 5秒的默认延迟
	if err != nil {
		client.Close()
		return nil, fmt.Errorf("创建延迟队列管理器失败: %v", err)
	}

	return &TestDelayQueue{
		rabbitClient: client,
		baseQueue:    baseQueue,
		config:       config,
		delayManager: delayManager,
	}, nil
}

// Push 实现延迟消息发送
func (t *TestDelayQueue) Push(ctx context.Context, topic string, message interface{}, delay time.Duration) error {
	// 主要使用底层队列实现
	return t.baseQueue.Push(ctx, topic, message, delay)
}

// 添加消费延迟消息的方法
func (t *TestDelayQueue) ConsumeDelayedMessages(topic string) (<-chan amqp.Delivery, error) {
	// 从死信队列消费已处理的延迟消息
	return t.delayManager.ConsumeProcessedMessages()
}

// Pop 获取延迟消息
func (t *TestDelayQueue) Pop(ctx context.Context, topic string) (interface{}, error) {
	return t.baseQueue.Pop(ctx, topic)
}

// Close 清理资源
func (t *TestDelayQueue) Close() error {
	if t.rabbitClient != nil {
		t.rabbitClient.Close()
	}
	return t.baseQueue.Close()
}

// 订单延迟服务示例
type OrderDelayService struct {
	container  *container.Container
	delayQueue *TestDelayQueue // 使用测试延迟队列
}

// 创建新的订单延迟服务
func NewOrderDelayService(c *container.Container, delayQueue *TestDelayQueue) *OrderDelayService {
	return &OrderDelayService{
		container:  c,
		delayQueue: delayQueue,
	}
}

// 延迟订单数据结构
type DelayOrder struct {
	ID        string    `json:"id"`
	Amount    float64   `json:"amount"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	CheckAt   time.Time `json:"check_at"`
}

// 创建延迟订单检查
func (s *OrderDelayService) CreateOrderPaymentCheck(orderID string, delay time.Duration) error {
	// 创建订单检查数据
	checkData := map[string]interface{}{
		"order_id":   orderID,
		"check_time": time.Now().Add(delay).Format(time.RFC3339),
		"action":     "payment_check",
	}

	// 将订单检查任务加入延迟队列
	return s.delayQueue.Push(
		context.Background(),
		"order_payment_check",
		checkData,
		delay,
	)
}

// 设置订单检查消费者
func (s *OrderDelayService) SetupConsumers() error {
	// 消费延迟消息
	messages, err := s.delayQueue.ConsumeDelayedMessages("order_payment_check")
	if err != nil {
		return err
	}

	go func() {
		for msg := range messages {
			// 解析消息
			var data map[string]interface{}
			if err := json.Unmarshal(msg.Body, &data); err != nil {
				log.Printf("解析延迟订单消息失败: %v", err)
				continue
			}

			log.Printf("处理订单付款检查: 订单ID=%v, 检查时间=%v",
				data["order_id"], data["check_time"])

			// 在实际应用中，这里会检查订单支付状态并进行相应处理
			// ...

			// 确认消息
			msg.Ack(false)
		}
	}()

	return nil
}

// 测试延迟队列基本功能
func TestDelayQueueWithDI(t *testing.T) {
	// 如果没有RabbitMQ环境，则跳过测试
	t.Skip("需要RabbitMQ环境才能运行此测试")

	// 创建容器
	c := container.NewContainer()

	// 初始化延迟队列配置
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

	// 创建延迟队列管理器
	delayQueueManager := types.NewDelayQueueManager(mqManager)

	// 设置到容器
	c.SetMessageQueue(mqManager)
	c.SetDelayQueue(delayQueueManager)

	// 创建测试延迟队列
	testDelayQueue, err := NewTestDelayQueue(delayQueueManager, config)
	if err != nil {
		t.Fatalf("创建测试延迟队列失败: %v", err)
	}

	// 创建订单延迟服务
	orderDelayService := NewOrderDelayService(c, testDelayQueue)

	// 设置消费者
	if err := orderDelayService.SetupConsumers(); err != nil {
		t.Fatalf("设置订单延迟消费者失败: %v", err)
	}

	// 发送延迟订单检查
	orderID := "delay-order-123"
	delay := 3 * time.Second

	startTime := time.Now()
	if err := orderDelayService.CreateOrderPaymentCheck(orderID, delay); err != nil {
		t.Fatalf("创建订单延迟检查失败: %v", err)
	}

	t.Logf("订单检查已加入延迟队列，将在约 %v 秒后执行", delay.Seconds())

	// 等待消息处理
	time.Sleep(delay + 1*time.Second)

	elapsed := time.Since(startTime)
	if elapsed < delay {
		t.Errorf("消息处理时间过短，期望至少 %v，实际 %v", delay, elapsed)
	} else {
		t.Logf("延迟消息处理成功，延迟时间: %v", elapsed)
	}
}

// 实际使用延迟队列的示例
func ExampleDelayQueueUsage() {
	// 在实际应用中，容器会在应用启动时初始化
	c := container.NewContainer()

	// 创建RabbitMQ管理器和延迟队列管理器
	config := &types.RabbitMQConfig{
		Host:           "localhost",
		Port:           5672,
		User:           "guest",
		Pass:           "guest",
		VirtualHost:    "/",
		QueueName:      "go-admin",
		DelayQueueName: "go-admin-delay",
		Exchange:       "go-admin-exchange",
		DelayExchange:  "go-admin-delay-exchange",
	}
	mqManager := types.NewRabbitMQManager(config)
	delayQueueManager := types.NewDelayQueueManager(mqManager)

	// 创建测试延迟队列
	testDelayQueue, err := NewTestDelayQueue(delayQueueManager, config)
	if err != nil {
		log.Fatalf("创建测试延迟队列失败: %v", err)
	}

	// 创建订单延迟服务
	orderDelayService := NewOrderDelayService(c, testDelayQueue)

	// 设置消费者
	orderDelayService.SetupConsumers()

	// 创建延迟订单检查
	orderID := "example-delay-order"
	delay := 10 * time.Second

	if err := orderDelayService.CreateOrderPaymentCheck(orderID, delay); err != nil {
		log.Fatalf("创建订单延迟检查失败: %v", err)
	}

	fmt.Printf("订单 %s 的支付检查将在 %v 秒后执行\n", orderID, delay.Seconds())
}
