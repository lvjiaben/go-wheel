package cron

import (
	"errors"
	"sync"
	"time"

	"github.com/robfig/cron/v3"
)

// Task 定时任务结构体
type Task struct {
	ID        string
	Name      string
	Spec      string
	Handler   func() error
	CreatedAt time.Time
	UpdatedAt time.Time
	EntryID   cron.EntryID
}

// CronManager 定时任务管理器
type CronManager struct {
	cron  *cron.Cron
	tasks map[string]*Task
	mu    sync.RWMutex
}

// New 创建定时任务管理器
func New() *CronManager {
	return &CronManager{
		cron:  cron.New(),
		tasks: make(map[string]*Task),
	}
}

// Start 启动定时任务管理器
func (m *CronManager) Start() {
	m.cron.Start()
}

// Stop 停止定时任务管理器
func (m *CronManager) Stop() {
	m.cron.Stop()
}

// AddTask 添加定时任务
func (m *CronManager) AddTask(task *Task) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 检查任务是否已存在
	if _, ok := m.tasks[task.ID]; ok {
		return ErrTaskExists
	}

	// 添加任务到cron
	entryID, err := m.cron.AddFunc(task.Spec, func() {
		if err := task.Handler(); err != nil {
			// 记录错误日志
		}
	})
	if err != nil {
		return err
	}

	// 保存任务信息
	task.CreatedAt = time.Now()
	task.UpdatedAt = time.Now()
	task.EntryID = entryID
	m.tasks[task.ID] = task

	return nil
}

// RemoveTask 移除定时任务
func (m *CronManager) RemoveTask(id string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 检查任务是否存在
	task, ok := m.tasks[id]
	if !ok {
		return ErrTaskNotFound
	}

	// 从cron中移除任务
	m.cron.Remove(task.EntryID)

	// 删除任务信息
	delete(m.tasks, id)

	return nil
}

// UpdateTask 更新定时任务
func (m *CronManager) UpdateTask(task *Task) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 检查任务是否存在
	oldTask, ok := m.tasks[task.ID]
	if !ok {
		return ErrTaskNotFound
	}

	// 从cron中移除旧任务
	m.cron.Remove(oldTask.EntryID)

	// 添加新任务到cron
	entryID, err := m.cron.AddFunc(task.Spec, func() {
		if err := task.Handler(); err != nil {
			// 记录错误日志
		}
	})
	if err != nil {
		return err
	}

	// 更新任务信息
	task.UpdatedAt = time.Now()
	task.EntryID = entryID
	m.tasks[task.ID] = task

	return nil
}

// GetTask 获取定时任务
func (m *CronManager) GetTask(id string) (*Task, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	task, ok := m.tasks[id]
	if !ok {
		return nil, ErrTaskNotFound
	}

	return task, nil
}

// ListTasks 列出所有定时任务
func (m *CronManager) ListTasks() []*Task {
	m.mu.RLock()
	defer m.mu.RUnlock()

	tasks := make([]*Task, 0, len(m.tasks))
	for _, task := range m.tasks {
		tasks = append(tasks, task)
	}

	return tasks
}

// 错误定义
var (
	ErrTaskExists   = errors.New("task already exists")
	ErrTaskNotFound = errors.New("task not found")
)
