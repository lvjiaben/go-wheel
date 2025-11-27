package cron

import (
	"context"
	"time"
)

// Task 定时任务接口
type Task interface {
	// GetName 获取任务名称
	GetName() string

	// GetSpec 获取 Cron 表达式
	GetSpec() string

	// Run 执行任务
	Run(ctx context.Context) error

	// GetDescription 获取任务描述（可选）
	GetDescription() string

	// UseDistributedLock 是否使用分布式锁（防止多实例重复执行）
	UseDistributedLock() bool

	// GetLockTimeout 获取分布式锁超时时间（默认30秒）
	GetLockTimeout() time.Duration
}

// BaseTask 基础任务结构（可选继承）
type BaseTask struct {
	Name               string
	Spec               string
	Description        string
	DistributedLock    bool          // 是否启用分布式锁
	LockTimeout        time.Duration // 锁超时时间
}

func (t *BaseTask) GetName() string {
	return t.Name
}

func (t *BaseTask) GetSpec() string {
	return t.Spec
}

func (t *BaseTask) GetDescription() string {
	return t.Description
}

func (t *BaseTask) UseDistributedLock() bool {
	return t.DistributedLock
}

func (t *BaseTask) GetLockTimeout() time.Duration {
	if t.LockTimeout > 0 {
		return t.LockTimeout
	}
	return 30 * time.Second // 默认30秒
}

