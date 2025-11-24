package monitor

import (
	"testing"
	"time"

	"github.com/lvjiaben/go-wheel/pkg/container"
)

// TestNewResourceMonitor 测试创建资源监控器
func TestNewResourceMonitor(t *testing.T) {
	c := container.NewContainer()

	monitor := NewResourceMonitor(c, 1*time.Second)

	if monitor == nil {
		t.Error("NewResourceMonitor 应该返回非 nil")
	}

	if monitor.interval != 1*time.Second {
		t.Errorf("interval 应该是 1s，实际: %v", monitor.interval)
	}

	if monitor.maxHistorySize != 100 {
		t.Error("maxHistorySize 应该是 100")
	}
}

// TestResourceMonitorCollect 测试资源收集
func TestResourceMonitorCollect(t *testing.T) {
	c := container.NewContainer()
	monitor := NewResourceMonitor(c, 100*time.Millisecond)

	// 执行收集（不会panic就说明成功）
	monitor.collect()

	// 验证统计信息
	stats := monitor.GetCurrentStats()
	if stats == nil {
		t.Error("collect 后应该有统计信息")
	}

	if stats.Timestamp.IsZero() {
		t.Error("Timestamp 不应该为零")
	}

	// 验证 Goroutine 统计
	if stats.GoroutineStats.Count <= 0 {
		t.Error("Goroutine 数量应该大于 0")
	}

	if stats.GoroutineStats.CPUCount <= 0 {
		t.Error("CPU 数量应该大于 0")
	}

	// 验证内存统计
	if stats.MemoryStats.Alloc == 0 {
		t.Error("内存分配应该大于 0")
	}
}

// TestResourceMonitorHistory 测试历史记录
func TestResourceMonitorHistory(t *testing.T) {
	c := container.NewContainer()
	monitor := NewResourceMonitor(c, 10*time.Millisecond)

	// 收集多次
	for i := 0; i < 15; i++ {
		monitor.collect()
		time.Sleep(5 * time.Millisecond)
	}

	history := monitor.GetStatsHistory()

	// 验证历史记录不超过最大值
	if len(history) > 100 {
		t.Errorf("历史记录应该不超过 100 条，实际: %d", len(history))
	}

	// 验证历史记录是按时间排序的
	for i := 1; i < len(history); i++ {
		if history[i].Timestamp.Before(history[i-1].Timestamp) {
			t.Error("历史记录应该按时间排序")
		}
	}
}

// TestResourceMonitorStartStop 测试启动和停止
func TestResourceMonitorStartStop(t *testing.T) {
	c := container.NewContainer()
	monitor := NewResourceMonitor(c, 50*time.Millisecond)

	// 在 goroutine 中启动
	go monitor.Start()

	// 等待几次收集
	time.Sleep(150 * time.Millisecond)

	// 停止监控
	monitor.Stop()

	// 验证有收集到数据
	history := monitor.GetStatsHistory()
	if len(history) < 2 {
		t.Errorf("应该至少收集 2 次数据，实际: %d", len(history))
	}
}

// TestResourceMonitorConcurrency 测试并发安全
func TestResourceMonitorConcurrency(t *testing.T) {
	c := container.NewContainer()
	monitor := NewResourceMonitor(c, 10*time.Millisecond)

	// 启动监控
	go monitor.Start()

	// 并发读取统计信息
	done := make(chan bool)
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				_ = monitor.GetCurrentStats()
				_ = monitor.GetStatsHistory()
				time.Sleep(1 * time.Millisecond)
			}
			done <- true
		}()
	}

	// 等待所有 goroutine 完成
	for i := 0; i < 10; i++ {
		<-done
	}

	monitor.Stop()
}

// TestMemoryStatsCollection 测试内存统计收集
func TestMemoryStatsCollection(t *testing.T) {
	c := container.NewContainer()
	monitor := NewResourceMonitor(c, 100*time.Millisecond)

	stats := &ResourceStats{}
	monitor.collectMemoryStats(stats)

	// 验证内存统计
	if stats.MemoryStats.Alloc == 0 {
		t.Error("Alloc 应该大于 0")
	}

	if stats.MemoryStats.Sys == 0 {
		t.Error("Sys 应该大于 0")
	}

	if stats.MemoryStats.HeapAlloc == 0 {
		t.Error("HeapAlloc 应该大于 0")
	}

	t.Logf("内存统计: Alloc=%dMB, Sys=%dMB, HeapAlloc=%dMB",
		stats.MemoryStats.Alloc,
		stats.MemoryStats.Sys,
		stats.MemoryStats.HeapAlloc)
}

// TestGoroutineStatsCollection 测试 Goroutine 统计收集
func TestGoroutineStatsCollection(t *testing.T) {
	c := container.NewContainer()
	monitor := NewResourceMonitor(c, 100*time.Millisecond)

	stats := &ResourceStats{}
	monitor.collectGoroutineStats(stats)

	// 验证 Goroutine 统计
	if stats.GoroutineStats.Count <= 0 {
		t.Error("Goroutine 数量应该大于 0")
	}

	if stats.GoroutineStats.CPUCount <= 0 {
		t.Error("CPU 数量应该大于 0")
	}

	if stats.GoroutineStats.MaxCount < stats.GoroutineStats.Count {
		t.Error("MaxCount 应该大于等于 Count")
	}

	t.Logf("Goroutine 统计: Count=%d, CPUCount=%d, MaxCount=%d",
		stats.GoroutineStats.Count,
		stats.GoroutineStats.CPUCount,
		stats.GoroutineStats.MaxCount)
}

