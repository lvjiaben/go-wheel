package health

import (
	"context"
	"sync"
	"time"
)

// Checker 健康检查接口
type Checker interface {
	Name() string
	Check(ctx context.Context) error
}

// Health 健康检查管理器
type Health struct {
	checkers []Checker
	mu       sync.RWMutex
}

// New 创建健康检查管理器
func New() *Health {
	return &Health{
		checkers: make([]Checker, 0),
	}
}

// Register 注册健康检查器
func (h *Health) Register(checker Checker) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.checkers = append(h.checkers, checker)
}

// Check 执行健康检查
func (h *Health) Check(ctx context.Context) map[string]error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	results := make(map[string]error)
	for _, checker := range h.checkers {
		results[checker.Name()] = checker.Check(ctx)
	}
	return results
}

// CheckWithTimeout 带超时的健康检查
func (h *Health) CheckWithTimeout(ctx context.Context, timeout time.Duration) map[string]error {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	return h.Check(ctx)
}

// IsHealthy 判断是否健康
func (h *Health) IsHealthy(ctx context.Context) bool {
	results := h.Check(ctx)
	for _, err := range results {
		if err != nil {
			return false
		}
	}
	return true
}

// IsHealthyWithTimeout 带超时的健康判断
func (h *Health) IsHealthyWithTimeout(ctx context.Context, timeout time.Duration) bool {
	results := h.CheckWithTimeout(ctx, timeout)
	for _, err := range results {
		if err != nil {
			return false
		}
	}
	return true
}
