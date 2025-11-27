package cron

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
	"go.uber.org/zap"
)

// Logger 日志接口（避免循环依赖）
type Logger interface {
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
}

// RedisClient Redis客户端接口（用于分布式锁）
type RedisClient interface {
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error)
	Del(ctx context.Context, keys ...string) error
	Expire(ctx context.Context, key string, expiration time.Duration) error
}

// Container 容器接口（避免循环依赖）
// 定义任务所需的最小接口
type Container interface {
	GetLogger() Logger
	GetDB() interface{}
	GetRedis() interface{}
}

// TaskFactory 任务工厂函数类型
type TaskFactory func(Container) Task

var (
	// taskRegistry 全局任务注册表
	taskRegistry = make(map[string]TaskFactory)
	// registryMu 注册表互斥锁
	registryMu sync.RWMutex
)

// Register 注册任务工厂函数
// 在每个任务文件的 init() 中调用
func Register(name string, factory TaskFactory) {
	registryMu.Lock()
	defer registryMu.Unlock()

	if _, exists := taskRegistry[name]; exists {
		panic("任务已注册: " + name)
	}

	taskRegistry[name] = factory
}

// GetRegisteredTasks 获取所有已注册的任务名称
func GetRegisteredTasks() []string {
	registryMu.RLock()
	defer registryMu.RUnlock()

	names := make([]string, 0, len(taskRegistry))
	for name := range taskRegistry {
		names = append(names, name)
	}
	return names
}

// ContainerWithCron 带 Cron 的容器接口（用于注册任务）
type ContainerWithCron interface {
	Container
	GetCron() interface{}
}

// RegisterAllTasks 注册所有定时任务（自动从注册表加载）
func RegisterAllTasks(c ContainerWithCron) error {
	cronInterface := c.GetCron()
	cronManager, ok := cronInterface.(*CronManager)
	if !ok || cronManager == nil {
		c.GetLogger().Warn("定时任务管理器未初始化或类型不匹配")
		return nil
	}

	registryMu.RLock()
	defer registryMu.RUnlock()

	// 从注册表创建所有任务实例
	tasks := make([]Task, 0, len(taskRegistry))
	for _, factory := range taskRegistry {
		task := factory(c)
		tasks = append(tasks, task)
	}

	// 批量注册任务
	for _, task := range tasks {
		// 创建任务的闭包，避免循环变量问题
		currentTask := task

		err := cronManager.AddJobWithContext(currentTask.GetSpec(), currentTask.GetName(), func(ctx context.Context) {
			// 如果启用分布式锁
			if currentTask.UseDistributedLock() {
				if err := cronManager.runWithDistributedLock(ctx, c, currentTask); err != nil {
					c.GetLogger().Error("任务执行失败",
						zap.String("task", currentTask.GetName()),
						zap.Error(err),
					)
				}
			} else {
				// 直接执行任务
				if err := currentTask.Run(ctx); err != nil {
					c.GetLogger().Error("任务执行失败",
						zap.String("task", currentTask.GetName()),
						zap.Error(err),
					)
				}
			}
		})
		if err != nil {
			c.GetLogger().Error("注册定时任务失败",
				zap.String("task", task.GetName()),
				zap.Error(err),
			)
			return err
		}
	}

	c.GetLogger().Info("定时任务注册完成",
		zap.Int("task_count", len(tasks)),
		zap.Strings("tasks", GetRegisteredTasks()),
	)

	return nil
}

// CronManager Cron定时任务管理器
type CronManager struct {
	cron   *cron.Cron
	logger Logger
	jobs   map[string]cron.EntryID // 任务名称 -> 任务ID映射
	mu     sync.RWMutex
	ctx    context.Context
	cancel context.CancelFunc
}

// NewCronManager 创建Cron管理器
func NewCronManager(ctx context.Context, logger Logger) *CronManager {
	cronCtx, cancel := context.WithCancel(ctx)
	
	// 创建cron实例，支持秒级定时
	cronInstance := cron.New(
		cron.WithSeconds(),                          // 支持秒级定时
		cron.WithChain(cron.Recover(cron.DefaultLogger)), // 自动恢复panic
	)

	return &CronManager{
		cron:   cronInstance,
		logger: logger,
		jobs:   make(map[string]cron.EntryID),
		ctx:    cronCtx,
		cancel: cancel,
	}
}

// AddJob 添加定时任务
// spec: Cron表达式，如 "0 30 * * * *" 表示每小时的30分执行
// name: 任务名称（用于管理和日志）
// cmd: 任务执行函数
func (m *CronManager) AddJob(spec string, name string, cmd func()) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 包装任务函数，添加日志和错误处理
	wrappedCmd := func() {
		// 检查上下文是否已取消
		select {
		case <-m.ctx.Done():
			m.logger.Info("定时任务已取消",
				zap.String("task", name),
			)
			return
		default:
		}

		m.logger.Info("定时任务开始执行",
			zap.String("task", name),
			zap.String("spec", spec),
		)

		// 执行任务
		defer func() {
			if r := recover(); r != nil {
				m.logger.Error("定时任务执行panic",
					zap.String("task", name),
					zap.Any("error", r),
				)
			}
		}()

		cmd()

		m.logger.Info("定时任务执行完成",
			zap.String("task", name),
		)
	}

	// 添加任务到cron
	entryID, err := m.cron.AddFunc(spec, wrappedCmd)
	if err != nil {
		m.logger.Error("添加定时任务失败",
			zap.String("task", name),
			zap.String("spec", spec),
			zap.Error(err),
		)
		return err
	}

	// 如果任务名称已存在，先移除旧任务
	if oldEntryID, exists := m.jobs[name]; exists {
		m.cron.Remove(oldEntryID)
		m.logger.Info("移除旧的定时任务",
			zap.String("task", name),
		)
	}

	// 保存任务ID
	m.jobs[name] = entryID

	m.logger.Info("添加定时任务成功",
		zap.String("task", name),
		zap.String("spec", spec),
		zap.Int("entry_id", int(entryID)),
	)

	return nil
}

// AddJobWithContext 添加带上下文的定时任务
func (m *CronManager) AddJobWithContext(spec string, name string, cmd func(context.Context)) error {
	wrappedCmd := func() {
		cmd(m.ctx)
	}
	return m.AddJob(spec, name, wrappedCmd)
}

// RegisterTask 注册任务（使用 Task 接口）
func (m *CronManager) RegisterTask(task Task) error {
	return m.AddJobWithContext(task.GetSpec(), task.GetName(), func(ctx context.Context) {
		if err := task.Run(ctx); err != nil {
			m.logger.Error("任务执行失败",
				zap.String("task", task.GetName()),
				zap.Error(err),
			)
		}
	})
}

// RegisterTasks 批量注册任务
func (m *CronManager) RegisterTasks(tasks []Task) error {
	for _, task := range tasks {
		if err := m.RegisterTask(task); err != nil {
			return err
		}
	}
	return nil
}

// RemoveJob 移除指定任务
func (m *CronManager) RemoveJob(name string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if entryID, exists := m.jobs[name]; exists {
		m.cron.Remove(entryID)
		delete(m.jobs, name)
		m.logger.Info("移除定时任务",
			zap.String("task", name),
		)
	}
}

// Start 启动定时任务调度器
func (m *CronManager) Start() {
	m.cron.Start()
	m.logger.Info("定时任务调度器已启动",
		zap.Int("job_count", len(m.jobs)),
	)
}

// Stop 停止定时任务调度器
func (m *CronManager) Stop() {
	m.cancel() // 取消上下文
	ctx := m.cron.Stop()
	<-ctx.Done() // 等待所有任务完成
	m.logger.Info("定时任务调度器已停止")
}

// GetJobs 获取所有任务列表
func (m *CronManager) GetJobs() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	jobs := make([]string, 0, len(m.jobs))
	for name := range m.jobs {
		jobs = append(jobs, name)
	}
	return jobs
}

// GetJobCount 获取任务数量
func (m *CronManager) GetJobCount() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.jobs)
}

// IsRunning 检查调度器是否运行中
func (m *CronManager) IsRunning() bool {
	return len(m.cron.Entries()) > 0
}

// GetEntries 获取所有任务条目（用于调试）
func (m *CronManager) GetEntries() []cron.Entry {
	return m.cron.Entries()
}

// runWithDistributedLock 使用分布式锁执行任务
func (m *CronManager) runWithDistributedLock(ctx context.Context, container Container, task Task) error {
	redisInterface := container.GetRedis()
	if redisInterface == nil {
		m.logger.Warn("Redis未初始化，跳过分布式锁",
			zap.String("task", task.GetName()),
		)
		return task.Run(ctx)
	}

	// 尝试类型断言为 RedisClient
	redis, ok := redisInterface.(RedisClient)
	if !ok {
		// 如果不是 RedisClient 接口，尝试通过反射获取方法
		m.logger.Warn("Redis客户端不支持分布式锁接口，降级为直接执行",
			zap.String("task", task.GetName()),
		)
		return task.Run(ctx)
	}

	lockKey := fmt.Sprintf("cron:lock:%s", task.GetName())
	lockTimeout := task.GetLockTimeout()

	// 尝试获取分布式锁
	locked, err := redis.SetNX(ctx, lockKey, "1", lockTimeout)
	if err != nil {
		m.logger.Error("获取分布式锁失败",
			zap.String("task", task.GetName()),
			zap.String("lock_key", lockKey),
			zap.Error(err),
		)
		return err
	}

	if !locked {
		// 其他实例正在执行，跳过本次执行
		m.logger.Info("任务已被其他实例执行，跳过",
			zap.String("task", task.GetName()),
			zap.String("lock_key", lockKey),
		)
		return nil
	}

	// 确保释放锁
	defer func() {
		if err := redis.Del(ctx, lockKey); err != nil {
			m.logger.Error("释放分布式锁失败",
				zap.String("task", task.GetName()),
				zap.String("lock_key", lockKey),
				zap.Error(err),
			)
		}
	}()

	m.logger.Info("获取分布式锁成功，开始执行任务",
		zap.String("task", task.GetName()),
		zap.String("lock_key", lockKey),
		zap.Duration("lock_timeout", lockTimeout),
	)

	// 执行任务
	return task.Run(ctx)
}

