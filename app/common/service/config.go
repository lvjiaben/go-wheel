package service

import (
	"sync"
	"time"

	"github.com/lvjiaben/go-wheel/app/backend/model/system"
	"github.com/lvjiaben/go-wheel/pkg/container"
	"go.uber.org/zap"
)

// ConfigCacheService 全局配置缓存服务
type ConfigCacheService struct {
	container   *container.Container
	cache       sync.Map // 配置缓存 key: string, value: string
	ticker      *time.Ticker
	stopChan    chan struct{}
	pubsubChan  chan struct{} // Redis Pub/Sub 停止通道
	initialized bool
	mu          sync.RWMutex
}

var (
	configCacheInstance *ConfigCacheService
	configCacheOnce     sync.Once
)

// NewConfigCacheService 创建配置缓存服务（单例）
func NewConfigCacheService(c *container.Container) *ConfigCacheService {
	configCacheOnce.Do(func() {
		configCacheInstance = &ConfigCacheService{
			container:  c,
			stopChan:   make(chan struct{}),
			pubsubChan: make(chan struct{}),
		}
	})
	return configCacheInstance
}

// Start 启动配置缓存服务
func (s *ConfigCacheService) Start(interval time.Duration) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.initialized {
		return
	}

	// 首次加载配置
	s.loadConfigFromDB()

	// 启动定时轮询
	s.ticker = time.NewTicker(interval)
	s.initialized = true

	// 启动定时轮询协程
	go func() {
		for {
			select {
			case <-s.ticker.C:
				s.loadConfigFromDB()
			case <-s.stopChan:
				return
			}
		}
	}()

	// 启动 Redis Pub/Sub 订阅（如果 Redis 可用）
	if s.container.GetRedis() != nil {
		go s.subscribeConfigUpdate()
	}

	s.container.GetLogger().Info("配置缓存服务已启动",
		zap.Duration("轮询间隔", interval),
		zap.Bool("Redis Pub/Sub", s.container.GetRedis() != nil))
}

// Stop 停止配置缓存服务
func (s *ConfigCacheService) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.initialized {
		return
	}

	if s.ticker != nil {
		s.ticker.Stop()
	}

	close(s.stopChan)
	close(s.pubsubChan)
	s.initialized = false

	s.container.GetLogger().Info("配置缓存服务已停止")
}

// loadConfigFromDB 从数据库加载配置到缓存
// 注意：sync.Map 本身是并发安全的，Store 操作是原子的
// 这里不需要额外的锁，因为每个 Store 操作都是独立的原子操作
func (s *ConfigCacheService) loadConfigFromDB() {
	var configs []system.Config

	if err := s.container.GetDB().Find(&configs).Error; err != nil {
		s.container.GetLogger().Error("加载配置失败", zap.Error(err))
		return
	}

	// 更新缓存（sync.Map.Store 是并发安全的）
	count := 0
	for _, config := range configs {
		s.cache.Store(config.Key, config.Value)
		count++
	}

	s.container.GetLogger().Debug("配置已重新加载",
		zap.Int("配置数量", count))
}

// Get 获取配置值
func (s *ConfigCacheService) Get(key string) string {
	// 从缓存读取
	if value, ok := s.cache.Load(key); ok {
		return value.(string)
	}

	// 缓存未命中，从数据库读取
	var config system.Config
	if err := s.container.GetDB().Where("`key` = ?", key).First(&config).Error; err != nil {
		s.container.GetLogger().Warn("配置不存在",
			zap.String("key", key),
			zap.Error(err))
		return ""
	}

	// 更新缓存
	s.cache.Store(key, config.Value)
	return config.Value
}

// GetWithDefault 获取配置值，如果不存在则返回默认值
func (s *ConfigCacheService) GetWithDefault(key string, defaultValue string) string {
	value := s.Get(key)
	if value == "" {
		return defaultValue
	}
	return value
}

// Set 设置配置值（同时更新数据库和缓存）
func (s *ConfigCacheService) Set(key string, value string) error {
	// 更新数据库
	var config system.Config
	if err := s.container.GetDB().Where("`key` = ?", key).First(&config).Error; err != nil {
		return err
	}

	if err := s.container.GetDB().Model(&config).Update("value", value).Error; err != nil {
		return err
	}

	// 更新缓存
	s.cache.Store(key, value)

	s.container.GetLogger().Info("配置已更新",
		zap.String("key", key),
		zap.String("value", value))

	return nil
}

// Refresh 手动刷新配置缓存（并通知其他服务器）
func (s *ConfigCacheService) Refresh() {
	s.loadConfigFromDB()
	s.publishConfigUpdate() // 发布更新消息，通知其他服务器
	s.container.GetLogger().Info("配置缓存已手动刷新并通知其他服务器")
}

// GetAll 获取所有配置
func (s *ConfigCacheService) GetAll() map[string]string {
	result := make(map[string]string)
	s.cache.Range(func(key, value interface{}) bool {
		result[key.(string)] = value.(string)
		return true
	})
	return result
}

// Exists 检查配置是否存在
func (s *ConfigCacheService) Exists(key string) bool {
	_, ok := s.cache.Load(key)
	return ok
}

// Delete 删除配置（同时删除数据库和缓存）
func (s *ConfigCacheService) Delete(key string) error {
	// 删除数据库记录
	if err := s.container.GetDB().Where("`key` = ?", key).Delete(&system.Config{}).Error; err != nil {
		return err
	}

	// 删除缓存
	s.cache.Delete(key)

	s.container.GetLogger().Info("配置已删除",
		zap.String("key", key))

	return nil
}

// subscribeConfigUpdate 订阅配置更新消息
func (s *ConfigCacheService) subscribeConfigUpdate() {
	redisWrapper := s.container.GetRedis()
	if redisWrapper == nil {
		return
	}

	ctx := s.container.GetContext()
	pubsub := redisWrapper.Client.Subscribe(ctx, "config:update")
	defer pubsub.Close()

	s.container.GetLogger().Info("开始订阅配置更新消息")

	for {
		select {
		case <-s.pubsubChan:
			s.container.GetLogger().Info("停止订阅配置更新消息")
			return
		default:
			msg, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				s.container.GetLogger().Error("接收配置更新消息失败", zap.Error(err))
				time.Sleep(time.Second) // 出错后等待1秒再重试
				continue
			}

			s.container.GetLogger().Info("收到配置更新消息",
				zap.String("channel", msg.Channel),
				zap.String("payload", msg.Payload))

			// 重新加载配置
			s.loadConfigFromDB()
		}
	}
}

// publishConfigUpdate 发布配置更新消息
func (s *ConfigCacheService) publishConfigUpdate() {
	redisWrapper := s.container.GetRedis()
	if redisWrapper == nil {
		return
	}

	ctx := s.container.GetContext()
	if err := redisWrapper.Client.Publish(ctx, "config:update", "refresh").Err(); err != nil {
		s.container.GetLogger().Error("发布配置更新消息失败", zap.Error(err))
		return
	}

	s.container.GetLogger().Debug("已发布配置更新消息")
}
