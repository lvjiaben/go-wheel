package monitor

import (
	"context"
	"runtime"
	"time"

	"github.com/lvjiaben/go-wheel/pkg/container"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"go.uber.org/zap"
)

// PrometheusMetrics Prometheus 指标收集器
type PrometheusMetrics struct {
	container *container.Container
	logger    *zap.Logger

	// 数据库指标
	dbOpenConnections     prometheus.Gauge
	dbInUseConnections    prometheus.Gauge
	dbIdleConnections     prometheus.Gauge
	dbWaitCount           prometheus.Counter
	dbWaitDuration        prometheus.Counter
	dbMaxIdleClosed       prometheus.Counter
	dbMaxLifetimeClosed   prometheus.Counter

	// Redis 指标
	redisTotalConns       prometheus.Gauge
	redisIdleConns        prometheus.Gauge
	redisStaleConns       prometheus.Gauge
	redisHits             prometheus.Counter
	redisMisses           prometheus.Counter
	redisTimeouts         prometheus.Counter
	redisPingSuccess      prometheus.Gauge

	// Goroutine 指标
	goroutineCount        prometheus.Gauge
	goroutineMaxCount     prometheus.Gauge

	// 内存指标
	memoryAlloc           prometheus.Gauge
	memoryTotalAlloc      prometheus.Counter
	memorySys             prometheus.Gauge
	memoryNumGC           prometheus.Counter
	memoryPauseTotal      prometheus.Counter
	memoryHeapAlloc       prometheus.Gauge
	memoryHeapInuse       prometheus.Gauge
	memoryStackInuse      prometheus.Gauge

	// HTTP 请求指标
	httpRequestsTotal     *prometheus.CounterVec
	httpRequestDuration   *prometheus.HistogramVec
	httpRequestSize       *prometheus.HistogramVec
	httpResponseSize      *prometheus.HistogramVec

	// 业务指标
	loginAttempts         *prometheus.CounterVec
	loginFailures         *prometheus.CounterVec
	activeUsers           prometheus.Gauge
}

// NewPrometheusMetrics 创建 Prometheus 指标收集器
func NewPrometheusMetrics(c *container.Container, namespace string) *PrometheusMetrics {
	if namespace == "" {
		namespace = "goweb"
	}

	pm := &PrometheusMetrics{
		container: c,
		logger:    c.GetLogger(),

		// 数据库指标
		dbOpenConnections: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "db",
			Name:      "open_connections",
			Help:      "Number of open database connections",
		}),
		dbInUseConnections: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "db",
			Name:      "in_use_connections",
			Help:      "Number of in-use database connections",
		}),
		dbIdleConnections: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "db",
			Name:      "idle_connections",
			Help:      "Number of idle database connections",
		}),
		dbWaitCount: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "db",
			Name:      "wait_count_total",
			Help:      "Total number of connections waited for",
		}),
		dbWaitDuration: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "db",
			Name:      "wait_duration_seconds_total",
			Help:      "Total time blocked waiting for a new connection",
		}),
		dbMaxIdleClosed: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "db",
			Name:      "max_idle_closed_total",
			Help:      "Total number of connections closed due to SetMaxIdleConns",
		}),
		dbMaxLifetimeClosed: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "db",
			Name:      "max_lifetime_closed_total",
			Help:      "Total number of connections closed due to SetConnMaxLifetime",
		}),

		// Redis 指标
		redisTotalConns: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "redis",
			Name:      "total_connections",
			Help:      "Total number of Redis connections",
		}),
		redisIdleConns: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "redis",
			Name:      "idle_connections",
			Help:      "Number of idle Redis connections",
		}),
		redisStaleConns: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "redis",
			Name:      "stale_connections",
			Help:      "Number of stale Redis connections",
		}),
		redisHits: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "redis",
			Name:      "hits_total",
			Help:      "Total number of Redis pool hits",
		}),
		redisMisses: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "redis",
			Name:      "misses_total",
			Help:      "Total number of Redis pool misses",
		}),
		redisTimeouts: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "redis",
			Name:      "timeouts_total",
			Help:      "Total number of Redis timeouts",
		}),
		redisPingSuccess: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "redis",
			Name:      "ping_success",
			Help:      "Redis ping success (1) or failure (0)",
		}),

		// Goroutine 指标
		goroutineCount: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "runtime",
			Name:      "goroutines",
			Help:      "Number of goroutines",
		}),
		goroutineMaxCount: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "runtime",
			Name:      "goroutines_max",
			Help:      "Maximum number of goroutines",
		}),

		// 内存指标
		memoryAlloc: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "memory",
			Name:      "alloc_bytes",
			Help:      "Bytes of allocated heap objects",
		}),
		memoryTotalAlloc: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "memory",
			Name:      "total_alloc_bytes",
			Help:      "Cumulative bytes allocated for heap objects",
		}),
		memorySys: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "memory",
			Name:      "sys_bytes",
			Help:      "Total bytes of memory obtained from the OS",
		}),
		memoryNumGC: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "memory",
			Name:      "gc_total",
			Help:      "Number of completed GC cycles",
		}),
		memoryPauseTotal: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "memory",
			Name:      "gc_pause_seconds_total",
			Help:      "Total GC pause time",
		}),
		memoryHeapAlloc: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "memory",
			Name:      "heap_alloc_bytes",
			Help:      "Bytes of allocated heap objects",
		}),
		memoryHeapInuse: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "memory",
			Name:      "heap_inuse_bytes",
			Help:      "Bytes in in-use spans",
		}),
		memoryStackInuse: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "memory",
			Name:      "stack_inuse_bytes",
			Help:      "Bytes in stack spans",
		}),

		// HTTP 请求指标
		httpRequestsTotal: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "http",
				Name:      "requests_total",
				Help:      "Total number of HTTP requests",
			},
			[]string{"method", "path", "status"},
		),
		httpRequestDuration: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: "http",
				Name:      "request_duration_seconds",
				Help:      "HTTP request duration in seconds",
				Buckets:   prometheus.DefBuckets,
			},
			[]string{"method", "path"},
		),
		httpRequestSize: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: "http",
				Name:      "request_size_bytes",
				Help:      "HTTP request size in bytes",
				Buckets:   prometheus.ExponentialBuckets(100, 10, 8),
			},
			[]string{"method", "path"},
		),
		httpResponseSize: promauto.NewHistogramVec(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Subsystem: "http",
				Name:      "response_size_bytes",
				Help:      "HTTP response size in bytes",
				Buckets:   prometheus.ExponentialBuckets(100, 10, 8),
			},
			[]string{"method", "path"},
		),

		// 业务指标
		loginAttempts: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "business",
				Name:      "login_attempts_total",
				Help:      "Total number of login attempts",
			},
			[]string{"type"}, // type: admin, user
		),
		loginFailures: promauto.NewCounterVec(
			prometheus.CounterOpts{
				Namespace: namespace,
				Subsystem: "business",
				Name:      "login_failures_total",
				Help:      "Total number of login failures",
			},
			[]string{"type", "reason"}, // reason: invalid_credentials, rate_limited, etc.
		),
		activeUsers: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "business",
			Name:      "active_users",
			Help:      "Number of active users",
		}),
	}

	return pm
}

// Collect 收集所有指标
func (pm *PrometheusMetrics) Collect() {
	pm.collectDBMetrics()
	pm.collectRedisMetrics()
	pm.collectRuntimeMetrics()
	pm.collectMemoryMetrics()
}

// collectDBMetrics 收集数据库指标
func (pm *PrometheusMetrics) collectDBMetrics() {
	db := pm.container.GetDB()
	if db == nil {
		return
	}

	sqlDB, err := db.DB()
	if err != nil {
		return
	}

	stats := sqlDB.Stats()
	pm.dbOpenConnections.Set(float64(stats.OpenConnections))
	pm.dbInUseConnections.Set(float64(stats.InUse))
	pm.dbIdleConnections.Set(float64(stats.Idle))
	pm.dbWaitCount.Add(float64(stats.WaitCount))
	pm.dbWaitDuration.Add(stats.WaitDuration.Seconds())
	pm.dbMaxIdleClosed.Add(float64(stats.MaxIdleClosed))
	pm.dbMaxLifetimeClosed.Add(float64(stats.MaxLifetimeClosed))
}

// collectRedisMetrics 收集 Redis 指标
func (pm *PrometheusMetrics) collectRedisMetrics() {
	redis := pm.container.GetRDB()
	if redis == nil {
		return
	}

	poolStats := redis.PoolStats()
	pm.redisTotalConns.Set(float64(poolStats.TotalConns))
	pm.redisIdleConns.Set(float64(poolStats.IdleConns))
	pm.redisStaleConns.Set(float64(poolStats.StaleConns))
	pm.redisHits.Add(float64(poolStats.Hits))
	pm.redisMisses.Add(float64(poolStats.Misses))
	pm.redisTimeouts.Add(float64(poolStats.Timeouts))

	// Ping 测试
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := redis.Ping(ctx).Err(); err != nil {
		pm.redisPingSuccess.Set(0)
	} else {
		pm.redisPingSuccess.Set(1)
	}
}

// collectRuntimeMetrics 收集运行时指标
func (pm *PrometheusMetrics) collectRuntimeMetrics() {
	pm.goroutineCount.Set(float64(runtime.NumGoroutine()))
}

// collectMemoryMetrics 收集内存指标
func (pm *PrometheusMetrics) collectMemoryMetrics() {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	pm.memoryAlloc.Set(float64(memStats.Alloc))
	pm.memoryTotalAlloc.Add(float64(memStats.TotalAlloc))
	pm.memorySys.Set(float64(memStats.Sys))
	pm.memoryNumGC.Add(float64(memStats.NumGC))
	pm.memoryPauseTotal.Add(float64(memStats.PauseTotalNs) / 1e9)
	pm.memoryHeapAlloc.Set(float64(memStats.HeapAlloc))
	pm.memoryHeapInuse.Set(float64(memStats.HeapInuse))
	pm.memoryStackInuse.Set(float64(memStats.StackInuse))
}

// RecordHTTPRequest 记录 HTTP 请求
func (pm *PrometheusMetrics) RecordHTTPRequest(method, path, status string, duration float64, reqSize, respSize int) {
	pm.httpRequestsTotal.WithLabelValues(method, path, status).Inc()
	pm.httpRequestDuration.WithLabelValues(method, path).Observe(duration)
	pm.httpRequestSize.WithLabelValues(method, path).Observe(float64(reqSize))
	pm.httpResponseSize.WithLabelValues(method, path).Observe(float64(respSize))
}

// RecordLoginAttempt 记录登录尝试
func (pm *PrometheusMetrics) RecordLoginAttempt(loginType string) {
	pm.loginAttempts.WithLabelValues(loginType).Inc()
}

// RecordLoginFailure 记录登录失败
func (pm *PrometheusMetrics) RecordLoginFailure(loginType, reason string) {
	pm.loginFailures.WithLabelValues(loginType, reason).Inc()
}

// SetActiveUsers 设置活跃用户数
func (pm *PrometheusMetrics) SetActiveUsers(count int) {
	pm.activeUsers.Set(float64(count))
}

