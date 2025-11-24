package monitor

import (
	"sync"
	"testing"

	"github.com/lvjiaben/go-wheel/pkg/container"
)

// TestNewPrometheusMetrics 测试创建 Prometheus 指标收集器
func TestNewPrometheusMetrics(t *testing.T) {
	c := container.NewContainer()
	metrics := NewPrometheusMetrics(c, "test_new")

	if metrics == nil {
		t.Error("NewPrometheusMetrics 应该返回非 nil")
	}

	if metrics.container != c {
		t.Error("container 引用不正确")
	}
}

// TestRecordHTTPRequest 测试记录 HTTP 请求
func TestRecordHTTPRequest(t *testing.T) {
	c := container.NewContainer()
	metrics := NewPrometheusMetrics(c, "test_http")

	// 记录几个请求（不会panic就说明成功）
	metrics.RecordHTTPRequest("GET", "/api/users", "200", 0.123, 1024, 2048)
	metrics.RecordHTTPRequest("POST", "/api/users", "201", 0.456, 2048, 1024)
	metrics.RecordHTTPRequest("GET", "/api/users", "404", 0.789, 512, 256)
}

// TestRecordLoginAttempt 测试记录登录尝试
func TestRecordLoginAttempt(t *testing.T) {
	c := container.NewContainer()
	metrics := NewPrometheusMetrics(c, "test_login")

	// 记录几次登录尝试（不会panic就说明成功）
	metrics.RecordLoginAttempt("admin")
	metrics.RecordLoginAttempt("user")
	metrics.RecordLoginAttempt("admin")
}

// TestRecordLoginFailure 测试记录登录失败
func TestRecordLoginFailure(t *testing.T) {
	c := container.NewContainer()
	metrics := NewPrometheusMetrics(c, "test_failure")

	// 记录几次登录失败（不会panic就说明成功）
	metrics.RecordLoginFailure("admin", "invalid_password")
	metrics.RecordLoginFailure("user", "invalid_captcha")
	metrics.RecordLoginFailure("admin", "rate_limited")
}

// TestSetActiveUsers 测试设置活跃用户数
func TestSetActiveUsers(t *testing.T) {
	c := container.NewContainer()
	metrics := NewPrometheusMetrics(c, "test_users")

	// 设置活跃用户数（不会panic就说明成功）
	metrics.SetActiveUsers(100)
	metrics.SetActiveUsers(200)
	metrics.SetActiveUsers(150)
}

// TestConcurrentRecording 测试并发记录
func TestConcurrentRecording(t *testing.T) {
	c := container.NewContainer()
	metrics := NewPrometheusMetrics(c, "test_concurrent")

	var wg sync.WaitGroup
	iterations := 100

	// 并发记录 HTTP 请求
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			metrics.RecordHTTPRequest("GET", "/api/test", "200", 0.1, 100, 200)
		}(i)
	}

	// 并发记录登录尝试
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			metrics.RecordLoginAttempt("admin")
		}(i)
	}

	// 并发记录登录失败
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			metrics.RecordLoginFailure("admin", "test")
		}(i)
	}

	// 并发设置活跃用户
	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			metrics.SetActiveUsers(index)
		}(i)
	}

	wg.Wait()
}

// TestCollect 测试指标收集
func TestCollect(t *testing.T) {
	c := container.NewContainer()
	metrics := NewPrometheusMetrics(c, "test_collect")

	// 记录一些数据
	metrics.RecordHTTPRequest("GET", "/test", "200", 0.1, 100, 200)
	metrics.RecordLoginAttempt("admin")
	metrics.RecordLoginFailure("admin", "invalid_password")
	metrics.SetActiveUsers(100)

	// 调用 Collect 方法（不会panic就说明成功）
	metrics.Collect()
}

// BenchmarkRecordHTTPRequest 基准测试：记录 HTTP 请求
func BenchmarkRecordHTTPRequest(b *testing.B) {
	c := container.NewContainer()
	metrics := NewPrometheusMetrics(c, "bench_http")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		metrics.RecordHTTPRequest("GET", "/test", "200", 0.1, 100, 200)
	}
}

// BenchmarkRecordLoginAttempt 基准测试：记录登录尝试
func BenchmarkRecordLoginAttempt(b *testing.B) {
	c := container.NewContainer()
	metrics := NewPrometheusMetrics(c, "bench_login")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		metrics.RecordLoginAttempt("admin")
	}
}

// BenchmarkConcurrentRecording 基准测试：并发记录
func BenchmarkConcurrentRecording(b *testing.B) {
	c := container.NewContainer()
	metrics := NewPrometheusMetrics(c, "bench_concurrent")

	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			metrics.RecordHTTPRequest("GET", "/test", "200", 0.1, 100, 200)
		}
	})
}

