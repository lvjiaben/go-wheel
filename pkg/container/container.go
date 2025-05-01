package container

import (
	"context"
	"sync"

	"admin/pkg/types"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Container struct {
	config     *types.Config
	db         *gorm.DB
	logger     *zap.Logger
	redis      *types.RedisClient
	i18n       types.I18n
	messageQ   types.MessageQueue
	delayQ     types.DelayQueue
	cron       types.CronManager
	mu         sync.RWMutex
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewContainer() *Container {
	ctx, cancel := context.WithCancel(context.Background())
	return &Container{
		ctx:        ctx,
		cancelFunc: cancel,
	}
}

// SetConfig 设置配置
func (c *Container) SetConfig(config *types.Config) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.config = config
}

// GetConfig 获取配置
func (c *Container) GetConfig() *types.Config {
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
func (c *Container) SetRedis(redis *types.RedisClient) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.redis = redis
}

// GetRedis 获取Redis客户端
func (c *Container) GetRedis() *types.RedisClient {
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

// Shutdown 关闭容器
func (c *Container) Shutdown() {
	c.cancelFunc()
}
