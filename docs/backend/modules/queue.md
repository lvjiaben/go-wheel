# Queue 消息队列

项目使用 RabbitMQ 作为消息队列，支持普通队列和延迟队列。消费者采用 `init()` 自动注册机制，与 Cron 定时任务保持一致的注册方式。

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
├── queue.go           # 包入口（确保 init 被调用）
└── order_consumer.go  # 订单消费者
```

## 发送消息

### 在控制器中发送消息

```go
type OrderController struct {
    container *container.Container
}

func (ctrl *OrderController) Create(ctx *gin.Context) {
    // 获取队列辅助工具
    queueHelper := ctrl.container.GetQueueHelper()
    if queueHelper == nil {
        // RabbitMQ 未启用
        return
    }

    // 发送普通消息
    err := queueHelper.Push(ctx, "order.created", map[string]interface{}{
        "order_id": 12345,
        "user_id":  1,
        "amount":   99.9,
    })

    // 发送延迟消息（30分钟后执行）
    err = queueHelper.PushDelay(ctx, "order.timeout", map[string]interface{}{
        "order_id": 12345,
    }, 30*time.Minute)
}
```

### 在服务层发送消息

```go
type OrderService struct {
    container *container.Container
}

func (s *OrderService) CreateOrder(ctx context.Context, data *OrderData) error {
    // 创建订单逻辑...

    // 发送消息通知
    queueHelper := s.container.GetQueueHelper()
    if queueHelper != nil {
        // 发送即时消息
        queueHelper.Push(ctx, "order.created", map[string]interface{}{
            "order_id": order.ID,
            "user_id":  order.UserID,
        })

        // 发送延迟消息（30分钟后检查订单是否支付）
        queueHelper.PushDelay(ctx, "order.timeout", map[string]interface{}{
            "order_id": order.ID,
        }, 30*time.Minute)
    }

    return nil
}
```

## 创建消费者

### 1. 定义消费者（使用 init 自动注册）

```go
// app/queue/order_consumer.go
package queue

import (
    "context"
    "encoding/json"

    queuePkg "github.com/lvjiaben/go-wheel/pkg/queue"
    "go.uber.org/zap"
)

// 自动注册消费者（在 init 中注册，无需手动调用）
func init() {
    queuePkg.Register("order.created", NewOrderConsumer)
}

// OrderConsumer 订单消费者
type OrderConsumer struct {
    queuePkg.BaseConsumer
    container queuePkg.Container
}

// NewOrderConsumer 创建订单消费者
func NewOrderConsumer(c queuePkg.Container) queuePkg.Consumer {
    return &OrderConsumer{
        BaseConsumer: queuePkg.BaseConsumer{
            Topic:       "order.created",
            Description: "订单创建消费者",
            Concurrency: 2, // 并发工作线程数
        },
        container: c,
    }
}

// Handle 处理消息（同时处理即时和延迟消息）
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

### 2. 自动启动（无需手动注册）

消费者在 `main.go` 中通过导入包自动注册和启动：

```go
// main.go
package main

import (
    _ "github.com/lvjiaben/go-wheel/app/queue" // 导入队列消费者包以触发 init()
    queuePkg "github.com/lvjiaben/go-wheel/pkg/queue"
)

func main() {
    c := container.NewContainer()
    defer c.Shutdown()

    // 启动消息队列消费者（自动启动所有已注册的消费者）
    if c.GetRabbitMQ() != nil {
        if err := queuePkg.StartAllConsumers(c.AsQueueContainer(), c.GetRabbitMQ()); err != nil {
            log.Fatalf("启动消息队列消费者失败: %v", err)
        }
    }

    // ...
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

消息处理成功返回 `nil`，失败返回 `error`：

```go
func (c *OrderConsumer) Handle(ctx context.Context, body []byte) error {
    // 处理消息...

    if success {
        return nil  // 确认消息
    }
    return errors.New("处理失败") // 消息会重新入队
}

// 可选：实现 OnError 方法处理错误
func (c *OrderConsumer) OnError(err error, message []byte) {
    c.container.GetLogger().Error("消息处理失败",
        zap.Error(err),
        zap.ByteString("message", message),
    )
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

