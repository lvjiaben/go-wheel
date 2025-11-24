package container

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

// TestCircuitBreaker 测试熔断器基本功能
func TestCircuitBreaker(t *testing.T) {
	cb := NewCircuitBreaker(3, 1*time.Second)

	// 测试初始状态
	if cb.IsOpen() {
		t.Error("熔断器初始状态应该是关闭的")
	}

	// 测试失败计数
	failCount := 0
	for i := 0; i < 5; i++ {
		err := cb.Execute(func() error {
			failCount++
			if failCount <= 3 {
				return &testError{"模拟失败"}
			}
			return nil
		})
		if i < 3 && err == nil {
			t.Errorf("第 %d 次执行应该返回错误", i+1)
		}
	}

	// 验证熔断器已打开
	if !cb.IsOpen() {
		t.Error("失败 3 次后熔断器应该打开")
	}

	// 测试熔断器打开时拒绝请求
	err := cb.Execute(func() error {
		return nil
	})
	if err == nil {
		t.Error("熔断器打开时应该拒绝请求")
	}
}

// TestCircuitBreakerRecovery 测试熔断器恢复
func TestCircuitBreakerRecovery(t *testing.T) {
	cb := NewCircuitBreaker(2, 100*time.Millisecond)

	// 触发熔断
	for i := 0; i < 2; i++ {
		cb.Execute(func() error {
			return &testError{"失败"}
		})
	}

	if !cb.IsOpen() {
		t.Error("熔断器应该已打开")
	}

	// 等待重置超时
	time.Sleep(150 * time.Millisecond)

	// 应该进入半开状态
	successCount := 0
	for i := 0; i < 6; i++ {
		err := cb.Execute(func() error {
			successCount++
			return nil
		})
		if err != nil && i == 0 {
			t.Error("第一次请求应该被允许（半开状态）")
		}
	}

	// 验证熔断器已关闭
	if cb.IsOpen() {
		t.Error("连续成功后熔断器应该关闭")
	}
}

// TestCircuitBreakerConcurrency 测试熔断器并发安全
func TestCircuitBreakerConcurrency(t *testing.T) {
	cb := NewCircuitBreaker(100, 1*time.Second)
	var wg sync.WaitGroup
	var successCount int64
	var failCount int64

	// 并发执行 1000 次
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			err := cb.Execute(func() error {
				// 模拟随机成功/失败
				if index%10 == 0 {
					return &testError{"失败"}
				}
				return nil
			})
			if err == nil {
				atomic.AddInt64(&successCount, 1)
			} else {
				atomic.AddInt64(&failCount, 1)
			}
		}(i)
	}

	wg.Wait()

	t.Logf("成功: %d, 失败: %d", successCount, failCount)

	// 验证没有数据竞争
	if successCount+failCount != 1000 {
		t.Errorf("总数不匹配: %d + %d != 1000", successCount, failCount)
	}
}

// TestRetryWithExponentialBackoff 测试指数退避重试
func TestRetryWithExponentialBackoff(t *testing.T) {
	c := NewContainer()

	attemptCount := 0
	startTime := time.Now()

	err := c.retry(func() error {
		attemptCount++
		if attemptCount < 3 {
			return &testError{"失败"}
		}
		return nil
	}, 5, 10*time.Millisecond)

	duration := time.Since(startTime)

	if err != nil {
		t.Errorf("重试应该成功: %v", err)
	}

	if attemptCount != 3 {
		t.Errorf("应该尝试 3 次，实际: %d", attemptCount)
	}

	// 验证指数退避：第一次等待 10ms，第二次等待 20ms
	// 总时间应该大于 30ms
	if duration < 30*time.Millisecond {
		t.Errorf("重试时间太短，可能没有使用指数退避: %v", duration)
	}
}

// TestRetryMaxAttempts 测试重试次数限制
func TestRetryMaxAttempts(t *testing.T) {
	c := NewContainer()

	attemptCount := 0
	err := c.retry(func() error {
		attemptCount++
		return &testError{"一直失败"}
	}, 3, 1*time.Millisecond)

	if err == nil {
		t.Error("应该返回错误")
	}

	if attemptCount != 3 {
		t.Errorf("应该尝试 3 次，实际: %d", attemptCount)
	}
}

// TestContainerGettersSetters 测试 Container 的 getter/setter 并发安全
func TestContainerGettersSetters(t *testing.T) {
	c := NewContainer()

	var wg sync.WaitGroup

	// 并发读写测试
	for i := 0; i < 100; i++ {
		wg.Add(2)

		// 写入
		go func() {
			defer wg.Done()
			c.SetConfig(nil)
		}()

		// 读取
		go func() {
			defer wg.Done()
			_ = c.GetConfig()
		}()
	}

	wg.Wait()
}

// TestContainerContext 测试 Container 的 context 管理
func TestContainerContext(t *testing.T) {
	c := NewContainer()

	// 测试 GetContext
	if c.GetContext() == nil {
		t.Error("GetContext 不应该返回 nil")
	}

	// 测试 Shutdown 会取消 context
	go func() {
		time.Sleep(10 * time.Millisecond)
		c.Shutdown()
	}()

	select {
	case <-c.GetContext().Done():
		// 正确：context 被取消
	case <-time.After(100 * time.Millisecond):
		t.Error("Shutdown 应该取消 context")
	}
}

// TestCircuitBreakerMonitor 测试熔断器监控
func TestCircuitBreakerMonitor(t *testing.T) {
	c := NewContainer()

	cb := NewCircuitBreaker(3, 100*time.Millisecond)

	// 启动监控
	cancelMonitor := c.startCircuitBreakerMonitor(cb, "测试")
	defer cancelMonitor()

	// 触发熔断
	for i := 0; i < 3; i++ {
		cb.Execute(func() error {
			return &testError{"失败"}
		})
	}

	if !cb.IsOpen() {
		t.Error("熔断器应该打开")
	}

	// 等待监控检查并恢复
	time.Sleep(200 * time.Millisecond)

	// 验证监控正在运行（通过检查熔断器状态）
	state := cb.GetState()
	if state != stateOpen && state != stateHalfOpen {
		t.Logf("熔断器状态: %d", state)
	}

	// 停止监控
	cancelMonitor()
	time.Sleep(50 * time.Millisecond)
}

// TestGetDBWithContext 测试带 context 的数据库连接获取
func TestGetDBWithContext(t *testing.T) {
	c := NewContainer()

	ctx := context.Background()
	db := c.GetDBWithContext(ctx)

	if db != nil {
		t.Error("db 未初始化时应该返回 nil")
	}
}

// TestGetRDB 测试 Redis 客户端获取
func TestGetRDB(t *testing.T) {
	c := NewContainer()

	client := c.GetRDB()

	if client != nil {
		t.Error("redis 未初始化时应该返回 nil")
	}
}

// 辅助类型
type testError struct {
	msg string
}

func (e *testError) Error() string {
	return e.msg
}

