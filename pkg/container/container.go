package container

import (
	"bytes"
	"context"
	"fmt"
	"io/fs"
	"math/rand"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-playground/locales/zh"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"

	"github.com/fsnotify/fsnotify"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/lvjiaben/go-wheel/pkg/config"
	"github.com/lvjiaben/go-wheel/pkg/constants"
	cronPkg "github.com/lvjiaben/go-wheel/pkg/cron"
	"github.com/lvjiaben/go-wheel/pkg/httpclient"
	queuePkg "github.com/lvjiaben/go-wheel/pkg/queue"
	"github.com/lvjiaben/go-wheel/pkg/types"
	wsPackage "github.com/lvjiaben/go-wheel/pkg/websocket"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type ComponentStatus struct {
	Name    string
	Status  string
	Message string
	Time    time.Time
}

// EmbedFS 嵌入文件系统配置
type EmbedFS struct {
	ConfigFS fs.FS // 配置文件
	I18nFS   fs.FS // 多语言文件
	ViewsFS  fs.FS // 模板文件
}

type Container struct {
	config               *config.Config
	db                   *gorm.DB
	logger               *zap.Logger
	redis                *types.RedisWrapper
	i18n                 types.I18n
	messageQ             types.MessageQueue
	delayQ               types.DelayQueue
	cron                 types.CronManager
	rabbitmq             *queuePkg.RabbitMQManager // RabbitMQ管理器
	httpClient           *httpclient.Client        // HTTP客户端
	wsHub                *wsPackage.Hub            // WebSocket Hub
	validate             *validator.Validate
	translator           ut.Translator
	mu                   sync.RWMutex
	ctx                  context.Context
	cancelFunc           context.CancelFunc
	status               map[string]ComponentStatus
	cbMonitorStarted     sync.Once // 确保熔断器监控只启动一次
	dbCBMonitorCancel    context.CancelFunc
	redisCBMonitorCancel context.CancelFunc
	customData           map[string]interface{} // 自定义数据存储（用于扩展）
	activeTransactions   sync.WaitGroup         // 活跃事务计数器（用于优雅关闭）
	goroutines           sync.WaitGroup         // Goroutine 计数器（用于等待所有 goroutine 退出）
	embedFS              *EmbedFS               // 嵌入文件系统
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

func NewContainer(embedFS *EmbedFS) *Container {
	ctx, cancel := context.WithCancel(context.Background())
	return &Container{
		ctx:        ctx,
		cancelFunc: cancel,
		status:     make(map[string]ComponentStatus),
		customData: make(map[string]interface{}),
		embedFS:    embedFS,
	}
}

// GetViewsFS 获取模板文件系统
func (c *Container) GetViewsFS() fs.FS {
	if c.embedFS != nil && c.embedFS.ViewsFS != nil {
		return c.embedFS.ViewsFS
	}
	return nil
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

	// 6. 国际化初始化
	if err := c.initializeValidate(); err != nil {
		return fmt.Errorf("验证器初始化失败: %v", err)
	}

	// 7. 消息队列初始化
	if err := c.initializeMessageQueue(); err != nil {
		return fmt.Errorf("消息队列初始化失败: %v", err)
	}

	// 8. 延迟队列初始化
	if err := c.initializeDelayQueue(); err != nil {
		return fmt.Errorf("延迟队列初始化失败: %v", err)
	}

	// 9. 定时任务初始化
	if err := c.initializeCron(); err != nil {
		return fmt.Errorf("定时任务初始化失败: %v", err)
	}

	// 10. HTTP客户端初始化
	if err := c.initializeHTTPClient(); err != nil {
		return fmt.Errorf("HTTP客户端初始化失败: %v", err)
	}

	// 11. WebSocket Hub 初始化
	if err := c.initializeWebSocketHub(); err != nil {
		return fmt.Errorf("WebSocket Hub 初始化失败: %v", err)
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

// GetRabbitMQ 获取RabbitMQ管理器
func (c *Container) GetRabbitMQ() *queuePkg.RabbitMQManager {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.rabbitmq
}

// GetQueueHelper 获取队列辅助工具
func (c *Container) GetQueueHelper() *queuePkg.QueueHelper {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.rabbitmq == nil {
		return nil
	}
	return queuePkg.NewQueueHelper(c.rabbitmq)
}

// GetHTTPClient 获取HTTP客户端
func (c *Container) GetHTTPClient() *httpclient.Client {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.httpClient
}

// GetContext 获取上下文
func (c *Container) GetContext() context.Context {
	return c.ctx
}

// SetContext 设置上下文（已废弃，不建议使用）
// Deprecated: 此方法不安全，可能导致容器生命周期管理混乱
// 请使用 GetDBWithContext 或 GetDB().WithContext() 代替
func (c *Container) SetContext(ctx context.Context) {
	// 为了向后兼容保留此方法，但不做任何操作
	// 原有的替换根 context 的行为是不安全的
}

// GetDBWithContext 获取带指定 context 的数据库连接（推荐使用）
func (c *Container) GetDBWithContext(ctx context.Context) *gorm.DB {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.db == nil {
		return nil
	}
	return c.db.WithContext(ctx)
}

// Shutdown 优雅关闭容器
func (c *Container) Shutdown() {
	// 取消上下文（通知所有 goroutine 停止）
	c.cancelFunc()

	// 等待所有 goroutine 退出（设置超时避免无限等待）
	done := make(chan struct{})
	go func() {
		c.goroutines.Wait()
		close(done)
	}()

	select {
	case <-done:
		if c.logger != nil {
			c.logger.Info("所有 goroutine 已退出")
		}
	case <-time.After(5 * time.Second):
		if c.logger != nil {
			c.logger.Warn("等待 goroutine 退出超时，强制关闭")
		}
	}

	// 停止熔断器监控 goroutine
	if c.dbCBMonitorCancel != nil {
		c.dbCBMonitorCancel()
	}
	if c.redisCBMonitorCancel != nil {
		c.redisCBMonitorCancel()
	}

	// 加锁保护资源访问
	c.mu.Lock()
	defer c.mu.Unlock()

	// 关闭定时任务
	if c.cron != nil {
		c.cron.Stop()
	}

	// 关闭消息队列
	if c.messageQ != nil {
		if err := c.messageQ.Close(); err != nil {
			if c.logger != nil {
				c.logger.Error("关闭消息队列失败", zap.Error(err))
			}
		}
	}

	// 关闭延迟队列
	if c.delayQ != nil {
		if err := c.delayQ.Close(); err != nil {
			if c.logger != nil {
				c.logger.Error("关闭延迟队列失败", zap.Error(err))
			}
		}
	}

	// 关闭Redis连接（使用新的 context 避免已取消的 context）
	if c.redis != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		if c.redis.Ping(ctx) == nil {
			if err := c.redis.Close(); err != nil {
				if c.logger != nil {
					c.logger.Error("关闭Redis连接失败", zap.Error(err))
				}
			}
		}
	}

	// 关闭数据库连接
	if c.db != nil {
		sqlDB, err := c.db.DB()
		if err != nil {
			if c.logger != nil {
				c.logger.Error("获取数据库连接失败", zap.Error(err))
			}
		} else {
			if err := sqlDB.Close(); err != nil {
				if c.logger != nil {
					c.logger.Error("关闭数据库连接失败", zap.Error(err))
				}
			}
		}
	}

	// 关闭RabbitMQ连接
	if c.rabbitmq != nil {
		if err := c.rabbitmq.Close(); err != nil {
			if c.logger != nil {
				c.logger.Error("关闭RabbitMQ连接失败", zap.Error(err))
			}
		}
	}

	if c.logger != nil {
		c.logger.Info("容器已关闭")
	}
}

// retry 重试机制（指数退避策略）
func (c *Container) retry(operation func() error, maxRetries int, initialDelay time.Duration) error {
	var err error
	delay := initialDelay
	maxDelay := constants.DefaultMaxRetryDelay // 最大延迟时间

	for i := 0; i < maxRetries; i++ {
		if err = operation(); err == nil {
			return nil
		}

		// 最后一次重试不需要等待
		if i == maxRetries-1 {
			break
		}

		if c.logger != nil {
			c.logger.Warn("操作失败，准备重试",
				zap.Int("retry", i+1),
				zap.Int("max_retries", maxRetries),
				zap.Duration("delay", delay),
				zap.Error(err))
		}

		// 等待后重试
		time.Sleep(delay)

		// 指数退避：每次延迟时间翻倍，但不超过最大延迟
		delay *= 2
		if delay > maxDelay {
			delay = maxDelay
		}
	}

	return fmt.Errorf("重试 %d 次后仍然失败: %v", maxRetries, err)
}

// retryWithJitter 带抖动的重试机制（避免惊群效应）
func (c *Container) retryWithJitter(operation func() error, maxRetries int, initialDelay time.Duration) error {
	var err error
	delay := initialDelay
	maxDelay := constants.DefaultMaxRetryDelay

	for i := 0; i < maxRetries; i++ {
		if err = operation(); err == nil {
			return nil
		}

		if i == maxRetries-1 {
			break
		}

		// 添加随机抖动（±25%）
		jitter := time.Duration(float64(delay) * (0.75 + 0.5*rand.Float64()))

		c.logger.Warn("操作失败，准备重试",
			zap.Int("retry", i+1),
			zap.Int("max_retries", maxRetries),
			zap.Duration("base_delay", delay),
			zap.Duration("jitter_delay", jitter),
			zap.Error(err))

		time.Sleep(jitter)

		// 指数退避
		delay *= 2
		if delay > maxDelay {
			delay = maxDelay
		}
	}

	return fmt.Errorf("重试 %d 次后仍然失败: %v", maxRetries, err)
}

// initializeConfig 初始化配置
func (c *Container) initializeConfig() error {
	// 创建配置实例
	v := viper.New()
	v.SetConfigType("yaml")

	// 优先从本地文件读取（开发环境热更新）
	v.SetConfigName("config")
	v.AddConfigPath(".")
	v.AddConfigPath("./configs")
	if err := v.ReadInConfig(); err != nil {
		// 本地文件不存在，从嵌入的配置读取（生产环境）
		if c.embedFS != nil && c.embedFS.ConfigFS != nil {
			embeddedConfig, err := fs.ReadFile(c.embedFS.ConfigFS, "configs/config.yaml")
			if err != nil {
				return fmt.Errorf("读取嵌入配置文件失败: %v", err)
			}
			if err := v.ReadConfig(bytes.NewReader(embeddedConfig)); err != nil {
				return fmt.Errorf("解析嵌入配置文件失败: %v", err)
			}
		} else {
			return fmt.Errorf("读取配置文件失败: %v", err)
		}
	}

	// 如果存在 .env 文件，读取并覆盖配置
	// .env 格式直接使用 . 分隔: database.host=127.0.0.1
	if _, err := os.Stat(".env"); err == nil {
		envViper := viper.New()
		envViper.SetConfigFile(".env")
		envViper.SetConfigType("env")
		if err := envViper.ReadInConfig(); err == nil {
			for _, key := range envViper.AllKeys() {
				v.Set(key, envViper.Get(key))
			}
		}
	}

	// 自动绑定环境变量（支持 APP_PORT 等格式）
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// 解析配置
	var cfg config.Config
	if err := v.Unmarshal(&cfg); err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	// 设置配置
	c.SetConfig(&cfg)

	// 如果存在外部配置文件，监听配置变化
	if _, err := os.Stat("./configs/config.yaml"); err == nil {
		v.SetConfigName("config")
		v.AddConfigPath("./configs")
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
	}

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
	if config.Database.Driver == "" {
		return fmt.Errorf("数据库驱动不能为空")
	}
	if config.Database.Driver != "mysql" && config.Database.Driver != "postgres" && config.Database.Driver != "postgresql" {
		return fmt.Errorf("不支持的数据库驱动: %s (支持: mysql, postgres)", config.Database.Driver)
	}
	if config.Database.Host == "" {
		return fmt.Errorf("数据库主机不能为空")
	}
	if config.Database.Port <= 0 {
		return fmt.Errorf("数据库端口必须大于0")
	}
	if config.Database.User == "" {
		return fmt.Errorf("数据库用户名不能为空")
	}
	if config.Database.Dbname == "" {
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

// reloadComponents 重新初始化受影响的组件（平滑切换）
func (c *Container) reloadComponents(config *config.Config) error {
	c.logger.Info("开始平滑重载组件...")

	// 保存旧的连接以便后续关闭
	c.mu.RLock()
	oldDB := c.db
	oldRedis := c.redis
	oldMessageQ := c.messageQ
	oldDelayQ := c.delayQ
	c.mu.RUnlock()

	// 重新初始化数据库连接（创建新连接）
	if err := c.initializeDB(); err != nil {
		c.logger.Error("重新初始化数据库失败", zap.Error(err))
		return fmt.Errorf("重新初始化数据库失败: %v", err)
	}

	// 重新初始化Redis连接（创建新连接）
	if err := c.initializeRedis(); err != nil {
		c.logger.Error("重新初始化Redis失败", zap.Error(err))
		return fmt.Errorf("重新初始化Redis失败: %v", err)
	}

	// 重新初始化消息队列（创建新连接）
	if err := c.initializeMessageQueue(); err != nil {
		c.logger.Error("重新初始化消息队列失败", zap.Error(err))
		return fmt.Errorf("重新初始化消息队列失败: %v", err)
	}

	// 重新初始化延迟队列（创建新连接）
	if err := c.initializeDelayQueue(); err != nil {
		c.logger.Error("重新初始化延迟队列失败", zap.Error(err))
		return fmt.Errorf("重新初始化延迟队列失败: %v", err)
	}

	// 延迟关闭旧连接，给正在使用的请求一些时间完成
	go c.gracefulCloseOldConnections(oldDB, oldRedis, oldMessageQ, oldDelayQ)

	// 重新初始化定时任务
	if err := c.initializeCron(); err != nil {
		return fmt.Errorf("重新初始化定时任务失败: %v", err)
	}

	c.logger.Info("组件平滑重载完成")
	return nil
}

// gracefulCloseOldConnections 优雅关闭旧连接
func (c *Container) gracefulCloseOldConnections(
	oldDB *gorm.DB,
	oldRedis *types.RedisWrapper,
	oldMessageQ types.MessageQueue,
	oldDelayQ types.DelayQueue,
) {
	// 等待所有活跃事务完成
	c.logger.Info("等待所有活跃事务完成...")
	done := make(chan struct{})
	go func() {
		c.activeTransactions.Wait()
		close(done)
	}()

	// 设置超时，避免无限等待
	gracePeriod := constants.DefaultGracePeriod
	select {
	case <-done:
		c.logger.Info("所有活跃事务已完成")
	case <-time.After(gracePeriod):
		c.logger.Warn(fmt.Sprintf("等待超时（%v），强制关闭旧连接", gracePeriod))
	}

	// 关闭旧的数据库连接
	if oldDB != nil {
		if sqlDB, err := oldDB.DB(); err == nil {
			if err := sqlDB.Close(); err != nil {
				c.logger.Error("关闭旧数据库连接失败", zap.Error(err))
			} else {
				c.logger.Info("旧数据库连接已关闭")
			}
		}
	}

	// 关闭旧的 Redis 连接
	if oldRedis != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if oldRedis.Ping(ctx) == nil {
			if err := oldRedis.Close(); err != nil {
				c.logger.Error("关闭旧Redis连接失败", zap.Error(err))
			} else {
				c.logger.Info("旧Redis连接已关闭")
			}
		}
	}

	// 关闭旧的消息队列
	if oldMessageQ != nil {
		if err := oldMessageQ.Close(); err != nil {
			c.logger.Error("关闭旧消息队列失败", zap.Error(err))
		} else {
			c.logger.Info("旧消息队列已关闭")
		}
	}

	// 关闭旧的延迟队列
	if oldDelayQ != nil {
		if err := oldDelayQ.Close(); err != nil {
			c.logger.Error("关闭旧延迟队列失败", zap.Error(err))
		} else {
			c.logger.Info("旧延迟队列已关闭")
		}
	}

	c.logger.Info("所有旧连接已优雅关闭")
}

// initializeLogger 初始化日志（支持日志轮转和分级）
func (c *Container) initializeLogger() error {
	logConfig := c.config.Log

	// 解析日志级别
	var level zapcore.Level
	switch logConfig.Level {
	case "debug":
		level = zapcore.DebugLevel
	case "info":
		level = zapcore.InfoLevel
	case "warn":
		level = zapcore.WarnLevel
	case "error":
		level = zapcore.ErrorLevel
	default:
		level = zapcore.InfoLevel
	}

	// 配置日志轮转
	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logConfig.Filename,   // 日志文件路径
		MaxSize:    logConfig.MaxSize,    // 单个日志文件最大大小（MB）
		MaxBackups: logConfig.MaxBackups, // 保留的旧日志文件最大数量
		MaxAge:     logConfig.MaxAge,     // 保留旧日志文件的最大天数
		Compress:   true,                 // 是否压缩旧日志文件
	})

	// 配置编码器
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder

	// 创建核心
	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderConfig),
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(writer), zapcore.AddSync(os.Stdout)), // 同时输出到文件和控制台
		level,
	)

	// 创建 logger
	logger := zap.New(core, zap.AddCaller(), zap.AddStacktrace(zapcore.ErrorLevel))

	c.SetLogger(logger)
	c.logger.Info("日志系统已初始化",
		zap.String("level", logConfig.Level),
		zap.String("filename", logConfig.Filename),
		zap.Int("max_size", logConfig.MaxSize),
		zap.Int("max_backups", logConfig.MaxBackups),
		zap.Int("max_age", logConfig.MaxAge))

	return nil
}

// initializeValidate 初始化验证器
func (c *Container) initializeValidate() error {
	validate := validator.New()
	zh := zh.New()
	uni := ut.New(zh, zh)
	trans, _ := uni.GetTranslator("zh")
	if err := zh_translations.RegisterDefaultTranslations(validate, trans); err != nil {
		return fmt.Errorf("注册验证器翻译失败: %v", err)
	}
	c.SetValidator(validate)
	c.SetTranslator(trans)
	return nil
}

// startCircuitBreakerMonitor 启动熔断器监控（带独立 context 管理）
func (c *Container) startCircuitBreakerMonitor(cb *CircuitBreaker, name string) context.CancelFunc {
	ctx, cancel := context.WithCancel(c.ctx)
	ticker := time.NewTicker(constants.DefaultCircuitBreakerCheckInterval)

	c.goroutines.Add(1) // 增加 goroutine 计数
	go func() {
		defer c.goroutines.Done() // goroutine 退出时减少计数
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				if c.logger != nil {
					c.logger.Info(fmt.Sprintf("%s 熔断器监控已停止", name))
				}
				return
			case <-ticker.C:
				cb.Check()
			}
		}
	}()

	return cancel
}

// initializeDB 初始化数据库（支持 MySQL 和 PostgreSQL）
func (c *Container) initializeDB() error {
	config := c.GetConfig()
	if config == nil {
		return fmt.Errorf("配置未初始化")
	}

	// 停止旧的监控 goroutine（如果存在）
	if c.dbCBMonitorCancel != nil {
		c.dbCBMonitorCancel()
	}

	// 创建熔断器并启动监控
	cb := NewCircuitBreaker(constants.DefaultCircuitBreakerThreshold, constants.DefaultCircuitBreakerTimeout)
	c.dbCBMonitorCancel = c.startCircuitBreakerMonitor(cb, "数据库")

	return c.retry(func() error {
		// 检查熔断器状态
		if cb.IsOpen() {
			return fmt.Errorf("数据库连接断路器已触发")
		}

		return cb.Execute(func() error {
			var dialector gorm.Dialector
			dbConfig := config.Database

			// 根据驱动类型选择不同的 dialector
			switch dbConfig.Driver {
			case "mysql":
				// URL 编码时区参数
				timezone := dbConfig.Timezone
				if timezone == "" {
					timezone = "Local"
				}

				dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=%s",
					dbConfig.User,
					dbConfig.Pass,
					dbConfig.Host,
					dbConfig.Port,
					dbConfig.Dbname,
					dbConfig.Charset,
					url.QueryEscape(timezone),
				)
				dialector = mysql.Open(dsn)
				c.GetLogger().Info("使用 MySQL 数据库",
					zap.String("host", dbConfig.Host),
					zap.Int("port", dbConfig.Port),
					zap.String("database", dbConfig.Dbname),
					zap.String("timezone", timezone),
				)

			case "postgres", "postgresql":
				dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
					dbConfig.Host,
					dbConfig.User,
					dbConfig.Pass,
					dbConfig.Dbname,
					dbConfig.Port,
					dbConfig.SSLMode,
					dbConfig.Timezone,
				)
				dialector = postgres.Open(dsn)
				c.GetLogger().Info("使用 PostgreSQL 数据库",
					zap.String("host", dbConfig.Host),
					zap.Int("port", dbConfig.Port),
					zap.String("database", dbConfig.Dbname),
				)

			default:
				return fmt.Errorf("不支持的数据库驱动: %s (支持: mysql, postgres)", dbConfig.Driver)
			}

			// 打开数据库连接
			db, err := gorm.Open(dialector, &gorm.Config{
				NamingStrategy: schema.NamingStrategy{
					SingularTable: true, // 使用单数表名
				},
			})
			if err != nil {
				c.GetLogger().Error("连接数据库失败",
					zap.Error(err),
					zap.String("driver", dbConfig.Driver),
				)
				return fmt.Errorf("连接数据库失败: %v", err)
			}

			// 获取底层 SQL DB
			sqlDB, err := db.DB()
			if err != nil {
				c.GetLogger().Error("获取数据库连接失败", zap.Error(err))
				return fmt.Errorf("获取数据库连接失败: %v", err)
			}

			// 设置连接池参数
			sqlDB.SetMaxIdleConns(dbConfig.MaxIdleConns)
			sqlDB.SetMaxOpenConns(dbConfig.MaxOpenConns)
			sqlDB.SetConnMaxLifetime(time.Duration(dbConfig.MaxLifetime) * time.Second)
			sqlDB.SetConnMaxIdleTime(time.Duration(dbConfig.MaxIdleTime) * time.Second)

			c.SetDB(db)
			c.GetLogger().Info("数据库连接成功",
				zap.String("driver", dbConfig.Driver),
				zap.Int("max_open_conns", dbConfig.MaxOpenConns),
				zap.Int("max_idle_conns", dbConfig.MaxIdleConns),
			)
			return nil
		})
	}, constants.DefaultMaxRetries, constants.DefaultInitialDelay)
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

	// 停止旧的监控 goroutine（如果存在）
	if c.redisCBMonitorCancel != nil {
		c.redisCBMonitorCancel()
	}

	// 创建熔断器并启动监控
	cb := NewCircuitBreaker(constants.DefaultCircuitBreakerThreshold, constants.DefaultCircuitBreakerTimeout)
	c.redisCBMonitorCancel = c.startCircuitBreakerMonitor(cb, "Redis")

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

			ctx, cancel := context.WithTimeout(c.ctx, constants.DefaultRedisTimeout)
			defer cancel()

			if err := client.Ping(ctx).Err(); err != nil {
				c.GetLogger().Error("连接Redis失败", zap.Error(err))
				return fmt.Errorf("连接Redis失败: %v", err)
			}

			c.SetRedis(&types.RedisWrapper{Client: client})
			return nil
		})
	}, constants.DefaultMaxRetries, constants.DefaultInitialDelay)
}

// initializeI18n 初始化国际化
func (c *Container) initializeI18n() error {
	// 创建 i18n 管理器
	i18nManager := types.NewI18nManager()

	// 优先从本地文件加载（开发环境热更新）
	langDir := filepath.Join("configs", "i18n")
	if err := i18nManager.LoadTranslations(langDir); err != nil {
		// 本地文件不存在，从嵌入的文件系统加载（生产环境）
		if c.embedFS != nil && c.embedFS.I18nFS != nil {
			subFS, err := fs.Sub(c.embedFS.I18nFS, "configs/i18n")
			if err != nil {
				return fmt.Errorf("获取i18n子目录失败: %v", err)
			}
			if err := i18nManager.LoadTranslationsFromFS(subFS); err != nil {
				return fmt.Errorf("加载语言文件失败: %v", err)
			}
		} else {
			return fmt.Errorf("加载语言文件失败: %v", err)
		}
	}

	// 设置到容器
	c.SetI18n(i18nManager)
	return nil
}

// initializeMessageQueue 初始化消息队列
func (c *Container) initializeMessageQueue() error {
	config := c.GetConfig()
	if config == nil {
		return fmt.Errorf("配置未初始化")
	}

	if !config.RabbitMQ.State {
		c.logger.Info("RabbitMQ未启用，跳过初始化")
		return nil
	}

	mqConfig := &queuePkg.RabbitMQConfig{
		Host:                config.RabbitMQ.Host,
		Port:                config.RabbitMQ.Port,
		VirtualHost:         config.RabbitMQ.VirtualHost,
		User:                config.RabbitMQ.User,
		Pass:                config.RabbitMQ.Pass,
		QueueName:           config.RabbitMQ.QueueName,
		DelayQueueName:      config.RabbitMQ.DelayQueueName,
		Exchange:            config.RabbitMQ.Exchange,
		DelayExchange:       config.RabbitMQ.DelayExchange,
		RetryCount:          config.RabbitMQ.RetryCount,
		ReconnectInterval:   config.RabbitMQ.ReconnectInterval,
		HeartbeatInterval:   config.RabbitMQ.HeartbeatInterval,
		ConnectionTimeout:   config.RabbitMQ.ConnectionTimeout,
		EnableConfirmation:  config.RabbitMQ.EnableConfirmation,
		PrefetchCount:       config.RabbitMQ.PrefetchCount,
		PrefetchSize:        config.RabbitMQ.PrefetchSize,
		DefaultConsumerName: config.RabbitMQ.DefaultConsumerName,
	}

	manager, err := queuePkg.NewRabbitMQManager(c.ctx, mqConfig, c.logger)
	if err != nil {
		return fmt.Errorf("初始化RabbitMQ失败: %v", err)
	}

	c.mu.Lock()
	c.rabbitmq = manager
	c.mu.Unlock()

	c.logger.Info("RabbitMQ初始化成功")
	return nil
}

// initializeDelayQueue 初始化延迟队列
func (c *Container) initializeDelayQueue() error {
	// 延迟队列使用同一个RabbitMQ连接，不需要单独初始化
	return nil
}

// initializeCron 初始化定时任务
func (c *Container) initializeCron() error {
	cronManager := cronPkg.NewCronManager(c.ctx, c.logger)
	c.SetCron(cronManager)
	c.logger.Info("定时任务管理器已初始化")
	return nil
}

// initializeHTTPClient 初始化HTTP客户端
func (c *Container) initializeHTTPClient() error {
	// 创建适配器，将 zap.Logger 转换为 httpclient.Logger
	loggerAdapter := &httpClientLoggerAdapter{logger: c.logger}

	// 创建HTTP客户端
	client := httpclient.NewClient(
		httpclient.WithTimeout(30*time.Second),
		httpclient.WithLogger(loggerAdapter),
		httpclient.WithHeaders(map[string]string{
			"User-Agent": "Go-Wheel/1.0",
		}),
	)

	c.mu.Lock()
	c.httpClient = client
	c.mu.Unlock()

	c.logger.Info("HTTP客户端已初始化")
	return nil
}

// httpClientLoggerAdapter zap.Logger 适配器
type httpClientLoggerAdapter struct {
	logger *zap.Logger
}

func (a *httpClientLoggerAdapter) Debug(msg string, fields ...interface{}) {
	a.logger.Sugar().Debugw(msg, fields...)
}

func (a *httpClientLoggerAdapter) Info(msg string, fields ...interface{}) {
	a.logger.Sugar().Infow(msg, fields...)
}

func (a *httpClientLoggerAdapter) Error(msg string, fields ...interface{}) {
	a.logger.Sugar().Errorw(msg, fields...)
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

// GetRDB 获取 Redis 客户端（线程安全）
func (c *Container) GetRDB() *redis.Client {
	c.mu.RLock()
	defer c.mu.RUnlock()
	if c.redis == nil {
		return nil
	}
	return c.redis.Client
}

// Set 设置自定义数据
func (c *Container) Set(key string, value interface{}) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.customData[key] = value
}

// Get 获取自定义数据
func (c *Container) Get(key string) interface{} {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.customData[key]
}

// BeginTransaction 开始一个事务（增加活跃事务计数）
// 使用示例：
//
//	defer c.EndTransaction()
//	c.BeginTransaction()
//	// 执行数据库操作...
func (c *Container) BeginTransaction() {
	c.activeTransactions.Add(1)
}

// EndTransaction 结束一个事务（减少活跃事务计数）
func (c *Container) EndTransaction() {
	c.activeTransactions.Done()
}

// TrackGoroutine 跟踪一个 goroutine（在启动 goroutine 前调用）
// 使用示例：
//
//	c.TrackGoroutine()
//	go func() {
//	    defer c.UntrackGoroutine()
//	    // goroutine 逻辑...
//	}()
func (c *Container) TrackGoroutine() {
	c.goroutines.Add(1)
}

// UntrackGoroutine 取消跟踪一个 goroutine（在 goroutine 退出时调用）
func (c *Container) UntrackGoroutine() {
	c.goroutines.Done()
}

// initializeWebSocketHub 初始化 WebSocket Hub
func (c *Container) initializeWebSocketHub() error {
	c.logger.Info("开始初始化 WebSocket Hub")

	// 创建 WebSocket Hub
	c.wsHub = wsPackage.NewHub(c.logger)

	// 启动 Hub（在独立的 goroutine 中运行）
	c.TrackGoroutine()
	go func() {
		defer c.UntrackGoroutine()
		c.wsHub.Run()
	}()

	c.logger.Info("WebSocket Hub 初始化成功")
	return nil
}

// GetWebSocketHub 获取 WebSocket Hub
func (c *Container) GetWebSocketHub() *wsPackage.Hub {
	return c.wsHub
}
