# Cron 定时任务

项目使用 [robfig/cron](https://github.com/robfig/cron) 实现定时任务，支持秒级调度和分布式锁。

## 目录结构

```
app/cron/
└── tasks/
    ├── example.go      # 示例任务
    └── cleanup.go      # 清理任务
```

## 创建任务

### 1. 定义任务

```go
// app/cron/tasks/example.go
package tasks

import (
    "context"
    "time"
    
    "github.com/lvjiaben/go-wheel/pkg/cron"
)

// ExampleTask 示例任务
type ExampleTask struct {
    container cron.Container
}

// 注册任务（在 init 中自动注册）
func init() {
    cron.Register("example", func(c cron.Container) cron.Task {
        return &ExampleTask{container: c}
    })
}

// GetName 任务名称
func (t *ExampleTask) GetName() string {
    return "example"
}

// GetSpec Cron 表达式
func (t *ExampleTask) GetSpec() string {
    return "0 */5 * * * *"  // 每5分钟执行
}

// UseDistributedLock 是否使用分布式锁
func (t *ExampleTask) UseDistributedLock() bool {
    return true
}

// GetLockTimeout 分布式锁超时时间
func (t *ExampleTask) GetLockTimeout() time.Duration {
    return 5 * time.Minute
}

// Run 执行任务
func (t *ExampleTask) Run(ctx context.Context) error {
    logger := t.container.GetLogger()
    logger.Info("执行示例任务")
    
    // 任务逻辑...
    
    return nil
}
```

### 2. Task 接口

```go
// pkg/cron/task.go
type Task interface {
    GetName() string                    // 任务名称
    GetSpec() string                    // Cron 表达式
    Run(ctx context.Context) error      // 执行任务
    UseDistributedLock() bool           // 是否使用分布式锁
    GetLockTimeout() time.Duration      // 锁超时时间
}
```

## Cron 表达式

支持 6 位 Cron 表达式（秒级）：

```
┌──────────── 秒 (0-59)
│ ┌────────── 分 (0-59)
│ │ ┌──────── 时 (0-23)
│ │ │ ┌────── 日 (1-31)
│ │ │ │ ┌──── 月 (1-12)
│ │ │ │ │ ┌── 周 (0-6, 0=周日)
│ │ │ │ │ │
* * * * * *
```

### 常用表达式

| 表达式 | 说明 |
|--------|------|
| `0 * * * * *` | 每分钟 |
| `0 */5 * * * *` | 每5分钟 |
| `0 0 * * * *` | 每小时 |
| `0 0 0 * * *` | 每天0点 |
| `0 0 8 * * *` | 每天8点 |
| `0 0 0 * * 1` | 每周一0点 |
| `0 0 0 1 * *` | 每月1号0点 |

## 分布式锁

多实例部署时，使用 Redis 分布式锁确保任务只执行一次：

```go
func (t *MyTask) UseDistributedLock() bool {
    return true  // 启用分布式锁
}

func (t *MyTask) GetLockTimeout() time.Duration {
    return 10 * time.Minute  // 锁超时时间
}
```

### 锁机制

1. 任务执行前尝试获取锁 `cron:lock:{task_name}`
2. 获取成功则执行任务
3. 获取失败则跳过本次执行
4. 任务完成后释放锁

## 手动管理任务

```go
// 获取 Cron 管理器
cronManager := container.GetCron().(*cron.CronManager)

// 添加任务
cronManager.AddJob("0 * * * * *", "my-task", func() {
    // 任务逻辑
})

// 添加带上下文的任务
cronManager.AddJobWithContext("0 * * * * *", "my-task", func(ctx context.Context) {
    // 任务逻辑
})

// 移除任务
cronManager.RemoveJob("my-task")

// 获取所有任务
jobs := cronManager.GetJobs()

// 获取任务数量
count := cronManager.GetJobCount()
```

## 启动和停止

```go
// main.go
func main() {
    // 初始化容器
    c := container.NewContainer()
    
    // 注册所有任务
    cron.RegisterAllTasks(c)
    
    // 启动调度器
    c.GetCron().(*cron.CronManager).Start()
    
    // 优雅关闭
    defer c.GetCron().(*cron.CronManager).Stop()
}
```

## 日志

任务执行会自动记录日志：

```
INFO  定时任务开始执行  {"task": "example", "spec": "0 */5 * * * *"}
INFO  定时任务执行完成  {"task": "example"}
```

## 最佳实践

1. **任务幂等** - 任务应该是幂等的，重复执行不会产生副作用
2. **超时控制** - 使用 context 控制任务超时
3. **错误处理** - 任务中的错误应该被捕获和记录
4. **分布式锁** - 多实例部署时启用分布式锁

