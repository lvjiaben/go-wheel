# Queue 消息队列

项目使用 RabbitMQ 作为消息队列，支持普通队列和延迟队列。

## 配置

在 `configs/config.yaml` 中配置：

```yaml
rabbitmq:
  host: "127.0.0.1"
  port: 5672
  virtual_host: "/"
  user: "guest"
  pass: "guest"
  queue_name: "default_queue"
  delay_queue_name: "delay_queue"
  exchange: "default_exchange"
  delay_exchange: "delay_exchange"
  retry_count: 3
  reconnect_interval: 5
  prefetch_count: 10
```

## 目录结构

```
app/queue/
└── consumers/
    ├── example.go      # 示例消费者
    └── order.go        # 订单消费者
```

## 发送消息

### 普通消息

```go
// 获取消息队列
mq := container.GetMessageQueue()

// 发送消息
err := mq.Publish("order.created", map[string]interface{}{
    "order_id": 12345,
    "user_id":  1,
    "amount":   99.9,
})
```

### 延迟消息

```go
// 获取延迟队列
delayQ := container.GetDelayQueue()

// 发送延迟消息（30分钟后执行）
err := delayQ.PublishDelay("order.timeout", map[string]interface{}{
    "order_id": 12345,
}, 30*time.Minute)
```

## 创建消费者

### 1. 定义消费者

```go
// app/queue/consumers/order.go
package consumers

import (
    "context"
    "encoding/json"
    
    "github.com/lvjiaben/go-wheel/pkg/container"
)

type OrderConsumer struct {
    container *container.Container
}

func NewOrderConsumer(c *container.Container) *OrderConsumer {
    return &OrderConsumer{container: c}
}

// GetQueueName 队列名称
func (c *OrderConsumer) GetQueueName() string {
    return "order_queue"
}

// GetRoutingKey 路由键
func (c *OrderConsumer) GetRoutingKey() string {
    return "order.*"
}

// Handle 处理消息
func (c *OrderConsumer) Handle(ctx context.Context, body []byte) error {
    var msg map[string]interface{}
    if err := json.Unmarshal(body, &msg); err != nil {
        return err
    }
    
    logger := c.container.GetLogger()
    logger.Info("处理订单消息", zap.Any("data", msg))
    
    // 业务逻辑...
    
    return nil
}
```

### 2. 注册消费者

```go
// main.go
func main() {
    c := container.NewContainer()
    
    // 注册消费者
    orderConsumer := consumers.NewOrderConsumer(c)
    c.GetMessageQueue().RegisterConsumer(orderConsumer)
    
    // 启动消费
    c.GetMessageQueue().StartConsume()
}
```

## 延迟队列原理

使用 TTL + 死信队列实现延迟消息：

```
生产者 → 延迟交换机 → 延迟队列(TTL) → 死信交换机 → 目标队列 → 消费者
```

1. 消息发送到延迟队列，设置 TTL
2. 消息过期后转发到死信交换机
3. 死信交换机路由到目标队列
4. 消费者从目标队列消费

## 消息确认

```go
// 手动确认模式
func (c *OrderConsumer) Handle(ctx context.Context, body []byte, ack func(), nack func()) error {
    // 处理消息...
    
    if success {
        ack()  // 确认消息
    } else {
        nack() // 拒绝消息，重新入队
    }
    
    return nil
}
```

## 重试机制

消息处理失败时自动重试：

```yaml
rabbitmq:
  retry_count: 3  # 最大重试次数
```

## 自动重连

连接断开时自动重连：

```go
// pkg/queue/rabbitmq.go
func (m *RabbitMQManager) monitorConnection() {
    for {
        select {
        case <-m.ctx.Done():
            return
        case <-m.reconnectC:
            // 重连逻辑
            for i := 0; i < m.config.RetryCount; i++ {
                if err := m.connect(); err == nil {
                    break
                }
                time.Sleep(time.Duration(m.config.ReconnectInterval) * time.Second)
            }
        }
    }
}
```

## 使用场景

| 场景 | 队列类型 | 说明 |
|------|----------|------|
| 订单创建通知 | 普通队列 | 实时处理 |
| 订单超时取消 | 延迟队列 | 30分钟后检查 |
| 邮件发送 | 普通队列 | 异步发送 |
| 定时提醒 | 延迟队列 | 指定时间触发 |

## 最佳实践

1. **消息幂等** - 消费者应该是幂等的
2. **消息持久化** - 重要消息开启持久化
3. **死信处理** - 处理失败的消息进入死信队列
4. **监控告警** - 监控队列积压情况

