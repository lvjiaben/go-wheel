# 组件使用文档

<p align="center">
  <a href="./README.md">
    <img src="https://img.shields.io/badge/返回-指南首页-blue.svg" alt="返回指南首页">
  </a>
  <a href="./components_en.md">
    <img src="https://img.shields.io/badge/Switch-English-blue.svg" alt="English">
  </a>
</p>

Go Admin框架集成了多种常用组件，这些组件已在框架启动时完成初始化，您只需在需要的地方通过依赖注入使用。本文档将详细介绍各个组件的使用方法。

## 目录

- [Redis缓存](#redis缓存)
  - [基本操作](#基本操作)
  - [实际应用场景](#实际应用场景)
  - [性能优化建议](#性能优化建议)
- [RabbitMQ消息队列](#rabbitmq消息队列)
  - [基本操作](#基本操作-1)
  - [实际应用场景](#实际应用场景-1)
  - [延迟队列](#延迟队列)
  - [消息可靠性](#消息可靠性)
- [Cron定时任务](#cron定时任务)
  - [基本操作](#基本操作-2)
  - [实际应用场景](#实际应用场景-2)
  - [管理与监控](#管理与监控)
- [健康检查](#健康检查)
  - [基本操作](#基本操作-3)
  - [实际应用场景](#实际应用场景-3)
  - [自定义检查器](#自定义检查器)

## Redis缓存

Redis缓存服务提供了对Redis的完整操作，包括字符串、哈希表、列表、集合等数据结构的增删改查。

### 使用方法

在服务或控制器中，通过依赖注入使用Redis服务：

```go
type UserService struct {
    redisService *service.RedisService
}

func NewUserService(redisService *service.RedisService) *UserService {
    return &UserService{
        redisService: redisService,
    }
}
```

### 基本操作

#### 字符串操作

```go
// 设置值
err := s.redisService.Set("key", "value", 30*time.Minute)

// 获取值
value, err := s.redisService.Get("key")

// 检查键是否存在
exists, err := s.redisService.Exists("key")

// 删除键
err := s.redisService.Delete("key")

// 设置过期时间
err := s.redisService.Expire("key", 1*time.Hour)
```

#### 哈希表操作

```go
// 设置哈希字段
err := s.redisService.HSet("user:1", "name", "张三")
err = s.redisService.HSet("user:1", "age", "25")

// 获取哈希字段
name, err := s.redisService.HGet("user:1", "name")

// 获取所有哈希字段
userData, err := s.redisService.HGetAll("user:1")
```

#### 列表操作

```go
// 从左侧推入元素
err := s.redisService.LPush("list", "item1", "item2")

// 从右侧推入元素
err := s.redisService.RPush("list", "item3")

// 从左侧弹出元素
item, err := s.redisService.LPop("list")

// 获取列表范围
items, err := s.redisService.LRange("list", 0, -1)
```

#### 集合操作

```go
// 添加集合成员
err := s.redisService.SAdd("set", "member1", "member2")

// 获取所有成员
members, err := s.redisService.SMembers("set")

// 判断成员是否存在
exists, err := s.redisService.SIsMember("set", "member1")
```

#### 有序集合操作

```go
// 添加有序集合成员
err := s.redisService.ZAdd("zset", 1.0, "member1")
err = s.redisService.ZAdd("zset", 2.0, "member2")

// 获取范围（按分数升序）
members, err := s.redisService.ZRange("zset", 0, -1)
```

### 实际应用场景

#### 缓存场景

```go
func (s *UserService) GetUserByID(id string) (*model.User, error) {
    // 缓存键
    cacheKey := "user:" + id
    
    // 尝试从缓存获取
    userData, err := s.redisService.Get(cacheKey)
    if err == nil {
        // 缓存命中，解析数据
        var user model.User
        json.Unmarshal([]byte(userData), &user)
        return &user, nil
    }
    
    // 缓存未命中，从数据库获取
    user, err := s.userRepo.GetByID(id)
    if err != nil {
        return nil, err
    }
    
    // 存入缓存
    jsonData, _ := json.Marshal(user)
    s.redisService.Set(cacheKey, string(jsonData), 30*time.Minute)
    
    return user, nil
}
```

#### 计数器场景

```go
func (s *ArticleService) IncrementViewCount(articleID string) (int64, error) {
    // 计数器键
    counterKey := "article:view:" + articleID
    
    // 自增计数
    count, err := s.redisService.Incr(counterKey)
    if err != nil {
        return 0, err
    }
    
    // 设置过期时间（如果需要）
    s.redisService.Expire(counterKey, 24*time.Hour)
    
    return count, nil
}
```

#### 分布式锁

```go
func (s *PaymentService) ProcessPayment(orderID string, amount float64) error {
    // 锁的键名
    lockKey := "lock:order:" + orderID
    // 获取锁，设置10秒过期时间，防止死锁
    acquired, err := s.redisService.SetNX(lockKey, "1", 10*time.Second)
    if err != nil {
        return err
    }
    
    if !acquired {
        return errors.New("订单正在处理中，请稍后再试")
    }
    
    // 确保最后释放锁
    defer s.redisService.Delete(lockKey)
    
    // 执行支付处理逻辑...
    
    return nil
}
```

### 性能优化建议

1. **合理使用过期时间**：为缓存数据设置合理的过期时间，避免缓存过度增长
2. **批量操作**：尽可能使用批量命令（如MGET, MSET, HMSET）减少网络开销
3. **使用管道技术**：对于多个无依赖关系的操作，使用Pipeline并行执行
4. **合理的键命名**：采用规范的键命名方案，便于管理和查询
5. **考虑序列化性能**：大数据量时注意JSON序列化与反序列化的性能影响

```go
// 使用Pipeline批量操作示例
pipe := s.redisService.Pipeline()
pipe.Set("key1", "value1", 1*time.Hour)
pipe.Set("key2", "value2", 1*time.Hour)
pipe.Incr("counter")
_, err := pipe.Exec()
```

## RabbitMQ消息队列

RabbitMQ消息队列服务提供了对RabbitMQ的完整操作，包括普通队列和延迟队列的发送和接收。

### 使用方法

在服务或控制器中，通过依赖注入使用MQ服务：

```go
type NotificationService struct {
    mqService *service.MQService
}

func NewNotificationService(mqService *service.MQService) *NotificationService {
    return &NotificationService{
        mqService: mqService,
    }
}
```

### 基本操作

#### 发送消息

```go
// 发送普通消息
err := s.mqService.PublishMessage("exchange", "routingKey", []byte("消息内容"))

// 发送延迟消息（使用延迟队列）
err := s.mqService.PublishDelayMessage("delayQueue", []byte("延迟消息"), 5*time.Second)
```

#### 消费消息

通常在服务初始化时设置消费者：

```go
// 消费普通消息
msgs, err := s.mqService.ConsumeMessages("queueName")
go func() {
    for msg := range msgs {
        // 处理消息
        log.Printf("接收到消息: %s", string(msg.Body))
        msg.Ack(false) // 确认消息
    }
}()

// 消费延迟队列的消息
msgs, err := s.mqService.ConsumeProcessedMessages("delayQueue")
go func() {
    for msg := range msgs {
        // 处理延迟消息
        log.Printf("接收到延迟消息: %s", string(msg.Body))
        msg.Ack(false) // 确认消息
    }
}()
```

### 实际应用场景

#### 发送通知

```go
func (s *NotificationService) SendNotification(userID string, message string) error {
    // 创建通知数据
    notification := map[string]interface{}{
        "user_id":    userID,
        "message":    message,
        "created_at": time.Now(),
    }
    
    // 序列化为JSON
    data, err := json.Marshal(notification)
    if err != nil {
        return err
    }
    
    // 发送到通知队列
    return s.mqService.PublishMessage("notifications", userID, data)
}
```

#### 延迟任务

```go
func (s *OrderService) CreateOrder(order *model.Order) error {
    // 保存订单
    if err := s.orderRepo.Save(order); err != nil {
        return err
    }
    
    // 发送延迟消息，15分钟后检查支付状态
    checkPaymentData := map[string]string{
        "order_id": order.ID,
    }
    data, _ := json.Marshal(checkPaymentData)
    
    // 发送15分钟后处理的延迟消息
    return s.mqService.PublishDelayMessage("order_payment_check", data, 15*time.Minute)
}

// 处理延迟消息的消费者（在服务初始化时设置）
func (s *OrderService) setupConsumers() {
    msgs, _ := s.mqService.ConsumeProcessedMessages("order_payment_check")
    go func() {
        for msg := range msgs {
            var data map[string]string
            json.Unmarshal(msg.Body, &data)
            
            // 检查订单支付状态
            s.checkOrderPayment(data["order_id"])
            
            msg.Ack(false)
        }
    }()
}
```

### 延迟队列

框架提供了基于RabbitMQ的延迟队列实现，适用于需要延时处理的业务场景：

```go
// 延迟队列配置示例
type DelayQueueConfig struct {
    Exchange     string        // 交换机名称
    Queue        string        // 队列名称 
    DeadExchange string        // 死信交换机名称
    DeadQueue    string        // 死信队列名称
    TTL          time.Duration // 延迟时间
}

// 接收者处理延迟消息示例
func (s *ScheduleService) HandleExpiredTask() {
    msgs, _ := s.mqService.ConsumeProcessedMessages("expired_tasks")
    go func() {
        for msg := range msgs {
            // 处理过期任务
            var task model.Task
            json.Unmarshal(msg.Body, &task)
            s.processExpiredTask(task)
            msg.Ack(false)
        }
    }()
}
```

### 消息可靠性

框架提供了多种机制确保消息可靠性：

```go
// 1. 发送确认
err := s.mqService.PublishWithConfirm(exchange, routingKey, data)

// 2. 消费端手动确认
msgs, _ := s.mqService.ConsumeWithManualAck(queueName)
for msg := range msgs {
    // 处理消息
    err := processMessage(msg)
    
    if err != nil {
        // 消息处理失败，拒绝并重新入队
        msg.Nack(false, true)
    } else {
        // 消息处理成功，确认
        msg.Ack(false)
    }
}

// 3. 持久化设置
queueConfig := &QueueConfig{
    Durable:    true,  // 队列持久化
    Persistent: true,  // 消息持久化
}
```

## Cron定时任务

Cron定时任务服务提供了对定时任务的完整管理，基于robfig/cron/v3包实现。

### 使用方法

在服务或控制器中，通过依赖注入使用Cron服务：

```go
type TaskService struct {
    cronService *service.CronService
}

func NewTaskService(cronService *service.CronService) *TaskService {
    return &TaskService{
        cronService: cronService,
    }
}
```

### 基本操作

#### 添加定时任务

```go
// 添加每分钟执行一次的任务
jobID, err := s.cronService.AddJob("0 * * * * *", func() {
    log.Println("每分钟执行一次")
})

// 添加每天凌晨执行的任务
jobID, err := s.cronService.AddJob("0 0 0 * * *", func() {
    log.Println("每天凌晨执行")
})

// 添加工作日上午9点执行的任务
jobID, err := s.cronService.AddJob("0 0 9 * * 1-5", func() {
    log.Println("工作日上午9点执行")
})

// 使用秒级定时器
jobID, err := s.cronService.AddJob("*/10 * * * * *", func() {
    log.Println("每10秒执行一次")
})
```

#### 管理定时任务

```go
// 移除任务
s.cronService.Remove(jobID)

// 暂停任务
s.cronService.Pause(jobID)

// 恢复任务
s.cronService.Resume(jobID)

// 获取下次执行时间
nextRun := s.cronService.GetNextRun(jobID)
```

### 实际应用场景

#### 定时统计

```go
func (s *StatisticsService) setupDailyStats() {
    // 每天凌晨1点执行统计
    s.cronService.AddJob("0 0 1 * * *", func() {
        // 计算昨天的统计数据
        yesterday := time.Now().AddDate(0, 0, -1)
        s.calculateDailyStats(yesterday)
    })
}

func (s *StatisticsService) calculateDailyStats(date time.Time) {
    // 统计逻辑
    // ...
    
    // 保存统计结果
    // ...
}
```

#### 定时清理

```go
func (s *CleanupService) setupCleanupTasks() {
    // 每周日凌晨3点清理临时文件
    s.cronService.AddJob("0 0 3 * * 0", func() {
        s.cleanupTempFiles()
    })
    
    // 每月1号清理过期数据
    s.cronService.AddJob("0 0 2 1 * *", func() {
        s.cleanupExpiredData()
    })
}
```

#### 定时数据同步

```go
func (s *SyncService) setupSyncTasks() {
    // 每小时同步一次数据
    s.cronService.AddJob("0 0 * * * *", func() {
        s.syncDataFromRemote()
    })
}

func (s *SyncService) syncDataFromRemote() {
    // 连接远程数据源
    client, err := s.connectToDataSource()
    if err != nil {
        log.Printf("同步失败: %v", err)
        return
    }
    
    // 获取新数据
    newData, err := client.FetchData()
    if err != nil {
        log.Printf("获取数据失败: %v", err)
        return
    }
    
    // 处理并保存数据
    s.processAndSaveData(newData)
}
```

### 管理与监控

框架提供了对定时任务的管理和监控功能：

```go
// 获取所有任务
tasks := s.cronService.ListTasks()

// 获取任务执行历史
history := s.cronService.GetTaskHistory(jobID)

// 获取任务统计信息
stats := s.cronService.GetTaskStats(jobID)
fmt.Printf("任务执行次数: %d\n", stats.ExecutionCount)
fmt.Printf("平均执行时间: %s\n", stats.AverageExecutionTime)
fmt.Printf("最后执行时间: %s\n", stats.LastExecutionTime)
```

## 健康检查

健康检查服务提供了对系统各组件健康状态的监控，可以自定义检查策略和报警规则。

### 使用方法

在服务或控制器中，通过依赖注入使用健康检查服务：

```go
type MonitorService struct {
    healthService *service.HealthService
}

func NewMonitorService(healthService *service.HealthService) *MonitorService {
    return &MonitorService{
        healthService: healthService,
    }
}
```

### 基本操作

#### 添加健康检查

```go
// 添加数据库健康检查
s.healthService.AddCheck("database", func() (bool, error) {
    err := s.db.Ping()
    return err == nil, err
})

// 添加Redis健康检查
s.healthService.AddCheck("redis", func() (bool, error) {
    _, err := s.redisClient.Ping().Result()
    return err == nil, err
})

// 添加外部API健康检查
s.healthService.AddCheck("api", func() (bool, error) {
    resp, err := http.Get("https://api.example.com/health")
    if err != nil {
        return false, err
    }
    defer resp.Body.Close()
    return resp.StatusCode == 200, nil
})
```

#### 获取健康状态

```go
// 获取所有检查项的状态
status := s.healthService.GetStatus()

// 获取特定检查项的状态
dbStatus, err := s.healthService.GetCheckStatus("database")

// 检查所有服务是否健康
isHealthy := s.healthService.IsHealthy()
```

### 实际应用场景

#### 健康检查API

```go
func (c *HealthController) GetHealthStatus(ctx *gin.Context) {
    // 获取所有健康检查状态
    status := c.healthService.GetStatus()
    
    // 检查是否所有组件都健康
    allHealthy := true
    for _, check := range status {
        if !check.Healthy {
            allHealthy = false
            break
        }
    }
    
    if allHealthy {
        ctx.JSON(http.StatusOK, gin.H{
            "status": "healthy",
            "checks": status,
        })
    } else {
        ctx.JSON(http.StatusServiceUnavailable, gin.H{
            "status": "unhealthy",
            "checks": status,
        })
    }
}
```

#### 健康检查中间件

```go
func HealthCheckMiddleware(healthService *service.HealthService) gin.HandlerFunc {
    return func(c *gin.Context) {
        if !healthService.IsHealthy() {
            c.AbortWithStatusJSON(http.StatusServiceUnavailable, gin.H{
                "status": "服务不可用，请稍后再试",
            })
            return
        }
        c.Next()
    }
}
```

### 自定义检查器

框架支持自定义健康检查器，实现更复杂的健康检查策略：

```go
// 自定义健康检查器
type DatabaseChecker struct {
    db *gorm.DB
    threshold int // 连接池阈值
}

func (c *DatabaseChecker) Name() string {
    return "database"
}

func (c *DatabaseChecker) Check() (bool, error) {
    // 检查数据库连接
    err := c.db.Raw("SELECT 1").Error
    if err != nil {
        return false, err
    }
    
    // 检查连接池状态
    stats := c.db.DB().Stats()
    if stats.OpenConnections > c.threshold {
        return false, fmt.Errorf("数据库连接池使用率过高: %d/%d", 
            stats.OpenConnections, c.threshold)
    }
    
    return true, nil
}

// 注册自定义检查器
s.healthService.RegisterChecker(&DatabaseChecker{
    db: s.db,
    threshold: 100,
}) 