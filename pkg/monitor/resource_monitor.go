package monitor

import (
	"context"
	"runtime"
	"sync"
	"time"

	"github.com/lvjiaben/go-wheel/pkg/container"
	"go.uber.org/zap"
)

// ResourceStats 资源统计信息
type ResourceStats struct {
	Timestamp time.Time `json:"timestamp"`

	// 数据库连接池统计
	DBStats struct {
		MaxOpenConnections int `json:"max_open_connections"` // 最大打开连接数
		OpenConnections    int `json:"open_connections"`     // 当前打开连接数
		InUse              int `json:"in_use"`               // 正在使用的连接数
		Idle               int `json:"idle"`                 // 空闲连接数
		WaitCount          int64 `json:"wait_count"`          // 等待连接的总次数
		WaitDuration       int64 `json:"wait_duration_ms"`    // 等待连接的总时间（毫秒）
		MaxIdleClosed      int64 `json:"max_idle_closed"`     // 因超过最大空闲数而关闭的连接数
		MaxLifetimeClosed  int64 `json:"max_lifetime_closed"` // 因超过最大生命周期而关闭的连接数
	} `json:"db_stats"`

	// Redis 连接池统计
	RedisStats struct {
		PoolSize     int    `json:"pool_size"`      // 连接池大小
		TotalConns   uint32 `json:"total_conns"`    // 总连接数
		IdleConns    uint32 `json:"idle_conns"`     // 空闲连接数
		StaleConns   uint32 `json:"stale_conns"`    // 过期连接数
		Hits         uint64 `json:"hits"`           // 命中次数
		Misses       uint64 `json:"misses"`         // 未命中次数
		Timeouts     uint64 `json:"timeouts"`       // 超时次数
		PingResponse string `json:"ping_response"`  // Ping 响应
	} `json:"redis_stats"`

	// Goroutine 统计
	GoroutineStats struct {
		Count      int `json:"count"`       // 当前 goroutine 数量
		MaxCount   int `json:"max_count"`   // 历史最大 goroutine 数量
		CPUCount   int `json:"cpu_count"`   // CPU 核心数
		CGOCalls   int64 `json:"cgo_calls"` // CGO 调用次数
	} `json:"goroutine_stats"`

	// 内存统计
	MemoryStats struct {
		Alloc        uint64 `json:"alloc_mb"`         // 已分配内存（MB）
		TotalAlloc   uint64 `json:"total_alloc_mb"`   // 累计分配内存（MB）
		Sys          uint64 `json:"sys_mb"`           // 系统内存（MB）
		NumGC        uint32 `json:"num_gc"`           // GC 次数
		PauseTotalNs uint64 `json:"pause_total_ms"`   // GC 暂停总时间（毫秒）
		HeapAlloc    uint64 `json:"heap_alloc_mb"`    // 堆内存分配（MB）
		HeapInuse    uint64 `json:"heap_inuse_mb"`    // 堆内存使用（MB）
		StackInuse   uint64 `json:"stack_inuse_mb"`   // 栈内存使用（MB）
	} `json:"memory_stats"`
}

// ResourceMonitor 资源监控器
type ResourceMonitor struct {
	container      *container.Container
	logger         *zap.Logger
	interval       time.Duration
	ctx            context.Context
	cancel         context.CancelFunc
	mu             sync.RWMutex
	currentStats   *ResourceStats
	maxGoroutines  int
	statsHistory   []*ResourceStats // 保留最近的统计历史
	maxHistorySize int
}

// NewResourceMonitor 创建资源监控器
func NewResourceMonitor(c *container.Container, interval time.Duration) *ResourceMonitor {
	ctx, cancel := context.WithCancel(c.GetContext())
	return &ResourceMonitor{
		container:      c,
		logger:         c.GetLogger(),
		interval:       interval,
		ctx:            ctx,
		cancel:         cancel,
		maxGoroutines:  runtime.NumGoroutine(),
		statsHistory:   make([]*ResourceStats, 0, 100),
		maxHistorySize: 100, // 保留最近 100 条记录
	}
}

// Start 启动监控
func (m *ResourceMonitor) Start() {
	if m.logger != nil {
		m.logger.Info("资源监控器已启动", zap.Duration("interval", m.interval))
	}

	ticker := time.NewTicker(m.interval)
	defer ticker.Stop()

	// 立即执行一次
	m.collect()

	for {
		select {
		case <-m.ctx.Done():
			if m.logger != nil {
				m.logger.Info("资源监控器已停止")
			}
			return
		case <-ticker.C:
			m.collect()
		}
	}
}

// Stop 停止监控
func (m *ResourceMonitor) Stop() {
	m.cancel()
}

// collect 收集资源统计信息
func (m *ResourceMonitor) collect() {
	stats := &ResourceStats{
		Timestamp: time.Now(),
	}

	// 收集数据库连接池统计
	m.collectDBStats(stats)

	// 收集 Redis 连接池统计
	m.collectRedisStats(stats)

	// 收集 Goroutine 统计
	m.collectGoroutineStats(stats)

	// 收集内存统计
	m.collectMemoryStats(stats)

	// 保存当前统计
	m.mu.Lock()
	m.currentStats = stats
	m.statsHistory = append(m.statsHistory, stats)
	if len(m.statsHistory) > m.maxHistorySize {
		m.statsHistory = m.statsHistory[1:]
	}
	m.mu.Unlock()

	// 记录日志
	m.logStats(stats)

	// 检查告警
	m.checkAlerts(stats)
}

// collectDBStats 收集数据库统计
func (m *ResourceMonitor) collectDBStats(stats *ResourceStats) {
	db := m.container.GetDB()
	if db == nil {
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		m.logger.Error("获取数据库连接失败", zap.Error(err))
		return
	}

	dbStats := sqlDB.Stats()
	stats.DBStats.MaxOpenConnections = dbStats.MaxOpenConnections
	stats.DBStats.OpenConnections = dbStats.OpenConnections
	stats.DBStats.InUse = dbStats.InUse
	stats.DBStats.Idle = dbStats.Idle
	stats.DBStats.WaitCount = dbStats.WaitCount
	stats.DBStats.WaitDuration = dbStats.WaitDuration.Milliseconds()
	stats.DBStats.MaxIdleClosed = dbStats.MaxIdleClosed
	stats.DBStats.MaxLifetimeClosed = dbStats.MaxLifetimeClosed
}

// collectRedisStats 收集 Redis 统计
func (m *ResourceMonitor) collectRedisStats(stats *ResourceStats) {
	redis := m.container.GetRDB()
	if redis == nil {
		return
	}

	poolStats := redis.PoolStats()
	stats.RedisStats.TotalConns = poolStats.TotalConns
	stats.RedisStats.IdleConns = poolStats.IdleConns
	stats.RedisStats.StaleConns = poolStats.StaleConns
	stats.RedisStats.Hits = uint64(poolStats.Hits)
	stats.RedisStats.Misses = uint64(poolStats.Misses)
	stats.RedisStats.Timeouts = uint64(poolStats.Timeouts)

	// 获取连接池大小（从配置中获取）
	if config := m.container.GetConfig(); config != nil {
		stats.RedisStats.PoolSize = config.Redis.PoolSize
	}

	// Ping 测试
	ctx, cancel := context.WithTimeout(m.ctx, 2*time.Second)
	defer cancel()
	if err := redis.Ping(ctx).Err(); err != nil {
		stats.RedisStats.PingResponse = "FAILED: " + err.Error()
	} else {
		stats.RedisStats.PingResponse = "PONG"
	}
}

// collectGoroutineStats 收集 Goroutine 统计
func (m *ResourceMonitor) collectGoroutineStats(stats *ResourceStats) {
	currentCount := runtime.NumGoroutine()
	stats.GoroutineStats.Count = currentCount
	stats.GoroutineStats.CPUCount = runtime.NumCPU()
	stats.GoroutineStats.CGOCalls = runtime.NumCgoCall()

	// 更新最大值
	if currentCount > m.maxGoroutines {
		m.maxGoroutines = currentCount
	}
	stats.GoroutineStats.MaxCount = m.maxGoroutines
}

// collectMemoryStats 收集内存统计
func (m *ResourceMonitor) collectMemoryStats(stats *ResourceStats) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	stats.MemoryStats.Alloc = memStats.Alloc / 1024 / 1024
	stats.MemoryStats.TotalAlloc = memStats.TotalAlloc / 1024 / 1024
	stats.MemoryStats.Sys = memStats.Sys / 1024 / 1024
	stats.MemoryStats.NumGC = memStats.NumGC
	stats.MemoryStats.PauseTotalNs = memStats.PauseTotalNs / 1000 / 1000
	stats.MemoryStats.HeapAlloc = memStats.HeapAlloc / 1024 / 1024
	stats.MemoryStats.HeapInuse = memStats.HeapInuse / 1024 / 1024
	stats.MemoryStats.StackInuse = memStats.StackInuse / 1024 / 1024
}

// logStats 记录统计日志
func (m *ResourceMonitor) logStats(stats *ResourceStats) {
	if m.logger == nil {
		return
	}
	m.logger.Info("资源统计",
		// 数据库
		zap.Int("db_open", stats.DBStats.OpenConnections),
		zap.Int("db_in_use", stats.DBStats.InUse),
		zap.Int("db_idle", stats.DBStats.Idle),
		// Redis
		zap.Uint32("redis_total", stats.RedisStats.TotalConns),
		zap.Uint32("redis_idle", stats.RedisStats.IdleConns),
		zap.String("redis_ping", stats.RedisStats.PingResponse),
		// Goroutine
		zap.Int("goroutines", stats.GoroutineStats.Count),
		zap.Int("goroutines_max", stats.GoroutineStats.MaxCount),
		// 内存
		zap.Uint64("mem_alloc_mb", stats.MemoryStats.Alloc),
		zap.Uint64("mem_sys_mb", stats.MemoryStats.Sys),
		zap.Uint32("gc_count", stats.MemoryStats.NumGC),
	)
}

// checkAlerts 检查告警
func (m *ResourceMonitor) checkAlerts(stats *ResourceStats) {
	// 数据库连接池告警
	if stats.DBStats.MaxOpenConnections > 0 {
		usage := float64(stats.DBStats.InUse) / float64(stats.DBStats.MaxOpenConnections)
		if usage > 0.8 {
			m.logger.Warn("数据库连接池使用率过高",
				zap.Float64("usage", usage*100),
				zap.Int("in_use", stats.DBStats.InUse),
				zap.Int("max", stats.DBStats.MaxOpenConnections))
		}
	}

	// Goroutine 数量告警
	if stats.GoroutineStats.Count > 10000 {
		m.logger.Warn("Goroutine 数量过多",
			zap.Int("count", stats.GoroutineStats.Count),
			zap.Int("max", stats.GoroutineStats.MaxCount))
	}

	// 内存使用告警
	if stats.MemoryStats.Alloc > 1024 { // 超过 1GB
		m.logger.Warn("内存使用过高",
			zap.Uint64("alloc_mb", stats.MemoryStats.Alloc),
			zap.Uint64("sys_mb", stats.MemoryStats.Sys))
	}
}

// GetCurrentStats 获取当前统计信息
func (m *ResourceMonitor) GetCurrentStats() *ResourceStats {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.currentStats
}

// GetStatsHistory 获取统计历史
func (m *ResourceMonitor) GetStatsHistory() []*ResourceStats {
	m.mu.RLock()
	defer m.mu.RUnlock()
	// 返回副本
	history := make([]*ResourceStats, len(m.statsHistory))
	copy(history, m.statsHistory)
	return history
}

