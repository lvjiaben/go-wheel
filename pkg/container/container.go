package container

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"sync/atomic"
	"time"

	"github.com/fsnotify/fsnotify"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/lvjiaben/go-wheel/pkg/config"
	"github.com/lvjiaben/go-wheel/pkg/types"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type ComponentStatus struct {
	Name    string
	Status  string
	Message string
	Time    time.Time
}

type Container struct {
	config     *config.Config
	db         *gorm.DB
	logger     *zap.Logger
	redis      *types.RedisWrapper
	i18n       types.I18n
	messageQ   types.MessageQueue
	delayQ     types.DelayQueue
	cron       types.CronManager
	validate   *validator.Validate
	translator ut.Translator
	mu         sync.RWMutex
	ctx        context.Context
	cancelFunc context.CancelFunc
	status     map[string]ComponentStatus
	DB         *gorm.DB
	RDB        *redis.Client
	Config     *viper.Viper
}

// CircuitBreaker 熔断器
type CircuitBreaker struct {
	mu               sync.RWMutex
	failureCount     int64 // 改为int64类型，使用原子操作
	lastFailure      time.Time
	state            int32 // 0: closed, 1: open, 2: half-open
	resetTimeout     time.Duration
	failureThreshold int64 // 改为int64类型，与failureCount匹配
	successes        int64 // 连续成功次数
	successThreshold int64 // 半开状态下恢复所需的连续成功次数
}

const (
	stateClosed = iota
	stateOpen
	stateHalfOpen
)

// NewCircuitBreaker 创建熔断器
func NewCircuitBreaker(failureThreshold int, resetTimeout time.Duration) *CircuitBreaker {
	return &CircuitBreaker{
		failureThreshold: int64(failureThreshold),
		resetTimeout:     resetTimeout,
		state:            stateClosed,
		successThreshold: 5, // 默认5次连续成功后恢复
	}
}

// Execute 执行操作
func (cb *CircuitBreaker) Execute(operation func() error) error {
	currentState := atomic.LoadInt32(&cb.state)

	// 熔断器打开状态
	if currentState == stateOpen {
		cb.mu.RLock()
		since := time.Since(cb.lastFailure)
		cb.mu.RUnlock()

		// 检查是否达到重置时间
		if since >= cb.resetTimeout {
			// 尝试切换到半开状态
			if atomic.CompareAndSwapInt32(&cb.state, stateOpen, stateHalfOpen) {
				atomic.StoreInt64(&cb.successes, 0) // 重置成功计数
				// 继续执行请求
			} else {
				return fmt.Errorf("熔断器已打开")
			}
		} else {
			return fmt.Errorf("熔断器已打开")
		}
	}

	// 半开或关闭状态，执行操作
	err := operation()

	if err != nil {
		// 操作失败
		if currentState == stateHalfOpen {
			// 半开状态下，失败立即回到打开状态
			atomic.StoreInt32(&cb.state, stateOpen)
			cb.mu.Lock()
			cb.lastFailure = time.Now()
			cb.mu.Unlock()
		} else if currentState == stateClosed {
			// 关闭状态下，记录失败并检查是否需要打开熔断器
			newCount := atomic.AddInt64(&cb.failureCount, 1)
			if newCount >= cb.failureThreshold {
				atomic.StoreInt32(&cb.state, stateOpen)
				cb.mu.Lock()
				cb.lastFailure = time.Now()
				cb.mu.Unlock()
			}
		}
		return err
	}

	// 操作成功
	if currentState == stateHalfOpen {
		// 半开状态下，记录成功并检查是否可以关闭熔断器
		newSuccesses := atomic.AddInt64(&cb.successes, 1)
		if newSuccesses >= cb.successThreshold {
			atomic.StoreInt32(&cb.state, stateClosed)
			atomic.StoreInt64(&cb.failureCount, 0)
		}
	} else if currentState == stateClosed {
		// 关闭状态下，重置失败计数
		atomic.StoreInt64(&cb.failureCount, 0)
	}

	return nil
}

// ForceOpen 强制打开熔断器
func (cb *CircuitBreaker) ForceOpen() {
	atomic.StoreInt32(&cb.state, stateOpen)
	cb.mu.Lock()
	cb.lastFailure = time.Now()
	cb.mu.Unlock()
}

// ForceClose 强制关闭熔断器
func (cb *CircuitBreaker) ForceClose() {
	atomic.StoreInt32(&cb.state, stateClosed)
	atomic.StoreInt64(&cb.failureCount, 0)
}

// GetState 获取当前状态
func (cb *CircuitBreaker) GetState() int32 {
	return atomic.LoadInt32(&cb.state)
}

// IsOpen 判断熔断器是否打开
func (cb *CircuitBreaker) IsOpen() bool {
	return atomic.LoadInt32(&cb.state) == stateOpen
}

// Check 检查熔断器状态，如需要则切换到半开状态
func (cb *CircuitBreaker) Check() {
	if atomic.LoadInt32(&cb.state) != stateOpen {
		return
	}

	cb.mu.RLock()
	since := time.Since(cb.lastFailure)
	cb.mu.RUnlock()

	if since >= cb.resetTimeout {
		// 尝试切换到半开状态
		if atomic.CompareAndSwapInt32(&cb.state, stateOpen, stateHalfOpen) {
			atomic.StoreInt64(&cb.successes, 0) // 重置成功计数
		}
	}
}

func NewContainer() *Container {
	ctx, cancel := context.WithCancel(context.Background())
	return &Container{
		ctx:        ctx,
		cancelFunc: cancel,
		status:     make(map[string]ComponentStatus),
	}
}

// Initialize 按顺序初始化组件
func (c *Container) Initialize() error {
	// 1. 配置初始化
	if err := c.initializeConfig(); err != nil {
		return fmt.Errorf("配置初始化失败: %v", err)
	}

	// 2. 日志初始化
	if err := c.initializeLogger(); err != nil {
		return fmt.Errorf("日志初始化失败: %v", err)
	}

	// 3. 数据库初始化
	if err := c.initializeDB(); err != nil {
		return fmt.Errorf("数据库初始化失败: %v", err)
	}

	// 4. Redis初始化
	if err := c.initializeRedis(); err != nil {
		return fmt.Errorf("Redis初始化失败: %v", err)
	}

	// 5. 国际化初始化
	if err := c.initializeI18n(); err != nil {
		return fmt.Errorf("国际化初始化失败: %v", err)
	}

	// 6. 消息队列初始化
	if err := c.initializeMessageQueue(); err != nil {
		return fmt.Errorf("消息队列初始化失败: %v", err)
	}

	// 7. 延迟队列初始化
	if err := c.initializeDelayQueue(); err != nil {
		return fmt.Errorf("延迟队列初始化失败: %v", err)
	}

	// 8. 定时任务初始化
	if err := c.initializeCron(); err != nil {
		return fmt.Errorf("定时任务初始化失败: %v", err)
	}

	return nil
}

// HealthCheck 检查所有组件的健康状态
func (c *Container) HealthCheck() map[string]ComponentStatus {
	c.mu.RLock()
	defer c.mu.RUnlock()

	status := make(map[string]ComponentStatus)

	// 检查数据库连接
	if c.db != nil {
		sqlDB, err := c.db.DB()
		if err != nil {
			status["db"] = ComponentStatus{
				Name:    "database",
				Status:  "error",
				Message: fmt.Sprintf("获取数据库连接失败: %v", err),
				Time:    time.Now(),
			}
		} else if err := sqlDB.Ping(); err != nil {
			status["db"] = ComponentStatus{
				Name:    "database",
				Status:  "error",
				Message: fmt.Sprintf("数据库连接失败: %v", err),
				Time:    time.Now(),
			}
		} else {
			status["db"] = ComponentStatus{
				Name:    "database",
				Status:  "ok",
				Message: "连接正常",
				Time:    time.Now(),
			}
		}
	}

	// 检查Redis连接
	if c.redis != nil {
		if err := c.redis.Set(c.ctx, "health_check", "ok", time.Second); err != nil {
			status["redis"] = ComponentStatus{
				Name:    "redis",
				Status:  "error",
				Message: fmt.Sprintf("Redis连接失败: %v", err),
				Time:    time.Now(),
			}
		} else {
			status["redis"] = ComponentStatus{
				Name:    "redis",
				Status:  "ok",
				Message: "连接正常",
				Time:    time.Now(),
			}
		}
	}

	// 检查消息队列
	if c.messageQ != nil {
		status["message_queue"] = ComponentStatus{
			Name:    "message_queue",
			Status:  "ok",
			Message: "运行正常",
			Time:    time.Now(),
		}
	}

	// 检查延迟队列
	if c.delayQ != nil {
		status["delay_queue"] = ComponentStatus{
			Name:    "delay_queue",
			Status:  "ok",
			Message: "运行正常",
			Time:    time.Now(),
		}
	}

	// 检查定时任务
	if c.cron != nil {
		status["cron"] = ComponentStatus{
			Name:    "cron",
			Status:  "ok",
			Message: "运行正常",
			Time:    time.Now(),
		}
	}

	return status
}

// SetConfig 设置配置
func (c *Container) SetConfig(cfg *config.Config) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.config = cfg
}

// GetConfig 获取配置
func (c *Container) GetConfig() *config.Config {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.config
}

// SetDB 设置数据库连接
func (c *Container) SetDB(db *gorm.DB) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.db = db
}

// GetDB 获取数据库连接
func (c *Container) GetDB() *gorm.DB {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.db
}

// SetLogger 设置日志
func (c *Container) SetLogger(logger *zap.Logger) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.logger = logger
}

// GetLogger 获取日志
func (c *Container) GetLogger() *zap.Logger {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.logger
}

// SetRedis 设置Redis客户端
func (c *Container) SetRedis(redis *types.RedisWrapper) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.redis = redis
}

// GetRedis 获取Redis客户端
func (c *Container) GetRedis() *types.RedisWrapper {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.redis
}

// SetI18n 设置多语言
func (c *Container) SetI18n(i18n types.I18n) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.i18n = i18n
}

// GetI18n 获取多语言
func (c *Container) GetI18n() types.I18n {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.i18n
}

// SetMessageQueue 设置消息队列
func (c *Container) SetMessageQueue(mq types.MessageQueue) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.messageQ = mq
}

// GetMessageQueue 获取消息队列
func (c *Container) GetMessageQueue() types.MessageQueue {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.messageQ
}

// SetDelayQueue 设置延迟队列
func (c *Container) SetDelayQueue(dq types.DelayQueue) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.delayQ = dq
}

// GetDelayQueue 获取延迟队列
func (c *Container) GetDelayQueue() types.DelayQueue {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.delayQ
}

// SetCron 设置定时任务
func (c *Container) SetCron(cron types.CronManager) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.cron = cron
}

// GetCron 获取定时任务
func (c *Container) GetCron() types.CronManager {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.cron
}

// GetContext 获取上下文
func (c *Container) GetContext() context.Context {
	return c.ctx
}

// SetContext 设置上下文
func (c *Container) SetContext(ctx context.Context) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ctx = ctx
}

// Shutdown 优雅关闭容器
func (c *Container) Shutdown() error {
	// 取消上下文
	c.cancelFunc()

	// 关闭定时任务
	if c.cron != nil {
		c.cron.Stop()
	}

	// 关闭消息队列
	if c.messageQ != nil {
		if err := c.messageQ.Close(); err != nil {
			c.logger.Error("关闭消息队列失败", zap.Error(err))
		}
	}

	// 关闭延迟队列
	if c.delayQ != nil {
		if err := c.delayQ.Close(); err != nil {
			c.logger.Error("关闭延迟队列失败", zap.Error(err))
		}
	}

	// 关闭Redis连接
	if c.redis != nil && c.redis.Ping(c.ctx) == nil {
		if err := c.redis.Close(); err != nil {
			c.logger.Error("关闭Redis连接失败", zap.Error(err))
		}
	}

	// 关闭数据库连接
	if c.db != nil {
		sqlDB, err := c.db.DB()
		if err != nil {
			c.logger.Error("获取数据库连接失败", zap.Error(err))
		} else {
			if err := sqlDB.Close(); err != nil {
				c.logger.Error("关闭数据库连接失败", zap.Error(err))
			}
		}
	}

	return nil
}

// retry 重试机制
func (c *Container) retry(operation func() error, maxRetries int, delay time.Duration) error {
	var err error
	for i := 0; i < maxRetries; i++ {
		if err = operation(); err == nil {
			return nil
		}
		c.logger.Warn("操作失败，准备重试",
			zap.Int("retry", i+1),
			zap.Int("max_retries", maxRetries),
			zap.Error(err))
		time.Sleep(delay)
	}
	return fmt.Errorf("重试 %d 次后仍然失败: %v", maxRetries, err)
}

// initializeConfig 初始化配置
func (c *Container) initializeConfig() error {
	// 创建配置实例
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./configs")

	// 读取配置
	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析配置
	var cfg config.Config
	if err := v.Unmarshal(&cfg); err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 设置配置
	c.SetConfig(&cfg)

	// 监听配置变化
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		c.logger.Info("配置文件发生变化",
			zap.String("name", e.Name),
			zap.String("op", e.Op.String()))

		// 重新解析配置
		var newCfg config.Config
		if err := v.Unmarshal(&newCfg); err != nil {
			c.logger.Error("重新解析配置文件失败", zap.Error(err))
			return
		}

		// 验证配置
		if err := c.validateConfig(&newCfg); err != nil {
			c.logger.Error("配置验证失败", zap.Error(err))
			return
		}

		// 更新配置
		c.SetConfig(&newCfg)

		// 重新初始化受影响的组件
		if err := c.reloadComponents(&newCfg); err != nil {
			c.logger.Error("重新初始化组件失败", zap.Error(err))
		}
	})

	return nil
}

// validateConfig 验证配置
func (c *Container) validateConfig(config *config.Config) error {
	// 验证应用配置
	if config.App.Name == "" {
		return fmt.Errorf("应用名称不能为空")
	}
	if config.App.Port <= 0 {
		return fmt.Errorf("应用端口必须大于0")
	}

	// 验证数据库配置
	if config.Mysql.Host == "" {
		return fmt.Errorf("数据库主机不能为空")
	}
	if config.Mysql.Port <= 0 {
		return fmt.Errorf("数据库端口必须大于0")
	}
	if config.Mysql.User == "" {
		return fmt.Errorf("数据库用户名不能为空")
	}
	if config.Mysql.Dbname == "" {
		return fmt.Errorf("数据库名称不能为空")
	}

	// 验证Redis配置
	if config.Redis.State {
		if config.Redis.Host == "" {
			return fmt.Errorf("Redis主机不能为空")
		}
		if config.Redis.Port <= 0 {
			return fmt.Errorf("Redis端口必须大于0")
		}
	}

	return nil
}

// reloadComponents 重新初始化受影响的组件
func (c *Container) reloadComponents(config *config.Config) error {
	// 重新初始化数据库连接
	if err := c.initializeDB(); err != nil {
		return fmt.Errorf("重新初始化数据库失败: %v", err)
	}

	// 重新初始化Redis连接
	if err := c.initializeRedis(); err != nil {
		return fmt.Errorf("重新初始化Redis失败: %v", err)
	}

	// 重新初始化消息队列
	if err := c.initializeMessageQueue(); err != nil {
		return fmt.Errorf("重新初始化消息队列失败: %v", err)
	}

	// 重新初始化延迟队列
	if err := c.initializeDelayQueue(); err != nil {
		return fmt.Errorf("重新初始化延迟队列失败: %v", err)
	}

	// 重新初始化定时任务
	if err := c.initializeCron(); err != nil {
		return fmt.Errorf("重新初始化定时任务失败: %v", err)
	}

	return nil
}

// initializeLogger 初始化日志
func (c *Container) initializeLogger() error {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout"}
	config.ErrorOutputPaths = []string{"stderr"}

	logger, err := config.Build()
	if err != nil {
		return fmt.Errorf("创建日志器失败: %v", err)
	}

	c.SetLogger(logger)
	return nil
}

// initializeDB 初始化数据库
func (c *Container) initializeDB() error {
	config := c.GetConfig()
	if config == nil {
		return fmt.Errorf("配置未初始化")
	}

	// 创建熔断器
	cb := NewCircuitBreaker(3, 30*time.Second)

	// 定期检查熔断器状态
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-c.ctx.Done():
				return
			case <-ticker.C:
				cb.Check()
			}
		}
	}()

	return c.retry(func() error {
		// 检查熔断器状态
		if cb.IsOpen() {
			return fmt.Errorf("数据库连接断路器已触发")
		}

		return cb.Execute(func() error {
			dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
				config.Mysql.User,
				config.Mysql.Pass,
				config.Mysql.Host,
				config.Mysql.Port,
				config.Mysql.Dbname,
				config.Mysql.Charset,
			)

			db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
			if err != nil {
				c.GetLogger().Error("连接数据库失败", zap.Error(err))
				return fmt.Errorf("连接数据库失败: %v", err)
			}

			sqlDB, err := db.DB()
			if err != nil {
				c.GetLogger().Error("获取数据库连接失败", zap.Error(err))
				return fmt.Errorf("获取数据库连接失败: %v", err)
			}

			sqlDB.SetMaxIdleConns(config.Mysql.MaxIdleConns)
			sqlDB.SetMaxOpenConns(config.Mysql.MaxOpenConns)
			sqlDB.SetConnMaxLifetime(time.Hour)

			c.SetDB(db)
			return nil
		})
	}, 3, 5*time.Second)
}

// initializeRedis 初始化Redis
func (c *Container) initializeRedis() error {
	config := c.GetConfig()
	if config == nil {
		return fmt.Errorf("配置未初始化")
	}

	if !config.Redis.State {
		return nil
	}

	// 创建熔断器
	cb := NewCircuitBreaker(3, 30*time.Second)

	// 定期检查熔断器状态
	ticker := time.NewTicker(5 * time.Second)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-c.ctx.Done():
				return
			case <-ticker.C:
				cb.Check()
			}
		}
	}()

	return c.retry(func() error {
		// 检查熔断器状态
		if cb.IsOpen() {
			return fmt.Errorf("Redis连接断路器已触发")
		}

		return cb.Execute(func() error {
			client := redis.NewClient(&redis.Options{
				Addr:     fmt.Sprintf("%s:%d", config.Redis.Host, config.Redis.Port),
				Password: config.Redis.Pass,
				DB:       config.Redis.Db,
				PoolSize: config.Redis.PoolSize,
			})

			ctx, cancel := context.WithTimeout(c.ctx, 5*time.Second)
			defer cancel()

			if err := client.Ping(ctx).Err(); err != nil {
				c.GetLogger().Error("连接Redis失败", zap.Error(err))
				return fmt.Errorf("连接Redis失败: %v", err)
			}

			c.SetRedis(&types.RedisWrapper{Client: client})
			return nil
		})
	}, 3, 5*time.Second)
}

// initializeI18n 初始化国际化
func (c *Container) initializeI18n() error {
	// 创建 i18n 管理器
	i18nManager := types.NewI18nManager()

	// 加载语言文件
	langDir := filepath.Join("configs", "i18n")
	if err := i18nManager.LoadTranslations(langDir); err != nil {
		return fmt.Errorf("加载语言文件失败: %v", err)
	}

	// 设置到容器
	c.SetI18n(i18nManager)
	return nil
}

// initializeMessageQueue 初始化消息队列
func (c *Container) initializeMessageQueue() error {
	// TODO: 实现消息队列初始化
	return nil
}

// initializeDelayQueue 初始化延迟队列
func (c *Container) initializeDelayQueue() error {
	// TODO: 实现延迟队列初始化
	return nil
}

// initializeCron 初始化定时任务
func (c *Container) initializeCron() error {
	// TODO: 实现定时任务初始化
	return nil
}

// SetValidator 设置验证器
func (c *Container) SetValidator(validate *validator.Validate) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.validate = validate
}

// GetValidator 获取验证器
func (c *Container) GetValidator() *validator.Validate {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.validate
}

// SetTranslator 设置翻译器
func (c *Container) SetTranslator(translator ut.Translator) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.translator = translator
}

// GetTranslator 获取翻译器
func (c *Container) GetTranslator() ut.Translator {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.translator
}

func (c *Container) GetRDB() *redis.Client {
	return c.RDB
}
