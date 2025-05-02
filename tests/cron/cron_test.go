package cron

import (
	"fmt"
	"log"
	"sync"
	"testing"
	"time"

	"github.com/lvjiaben/go-wheel/pkg/container"
	pkg_cron "github.com/lvjiaben/go-wheel/pkg/cron"
	"github.com/robfig/cron/v3"
)

// 任务执行器
type TaskExecutor struct {
	cron           *cron.Cron
	executionCount map[string]int
	mu             sync.Mutex
}

// 创建新的任务执行器
func NewTaskExecutor() *TaskExecutor {
	return &TaskExecutor{
		cron:           cron.New(cron.WithSeconds()),
		executionCount: make(map[string]int),
	}
}

// 添加定时任务
func (e *TaskExecutor) AddJob(name, spec string, fn func()) (cron.EntryID, error) {
	wrappedFn := func() {
		e.mu.Lock()
		e.executionCount[name]++
		e.mu.Unlock()

		fn()
	}

	return e.cron.AddFunc(spec, wrappedFn)
}

// 开始执行任务
func (e *TaskExecutor) Start() {
	e.cron.Start()
}

// 停止执行任务
func (e *TaskExecutor) Stop() {
	e.cron.Stop()
}

// 获取任务执行次数
func (e *TaskExecutor) GetExecutionCount(name string) int {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.executionCount[name]
}

// 清除执行计数
func (e *TaskExecutor) ClearExecutionCounts() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.executionCount = make(map[string]int)
}

// 使用依赖注入的任务服务
type TaskService struct {
	container   *container.Container
	cronManager *pkg_cron.CronManager
	tasks       map[string]string // 任务ID -> 任务名称映射
}

// 创建新的任务服务
func NewTaskService(c *container.Container) *TaskService {
	return &TaskService{
		container:   c,
		cronManager: pkg_cron.New(),
		tasks:       make(map[string]string),
	}
}

// 开始所有任务
func (s *TaskService) Start() {
	s.cronManager.Start()
}

// 停止所有任务
func (s *TaskService) Stop() {
	s.cronManager.Stop()
}

// 添加定时任务
func (s *TaskService) AddTask(name, spec string, handler func() error) error {
	// 创建任务
	task := &pkg_cron.Task{
		ID:      name,
		Name:    name,
		Spec:    spec,
		Handler: handler,
	}

	// 添加到管理器
	if err := s.cronManager.AddTask(task); err != nil {
		return fmt.Errorf("添加任务失败: %v", err)
	}

	// 记录任务
	s.tasks[name] = name

	return nil
}

// 移除任务
func (s *TaskService) RemoveTask(name string) error {
	// 从管理器中移除
	return s.cronManager.RemoveTask(name)
}

// 手动触发任务执行
func (s *TaskService) RunTask(name string) error {
	// 查找任务
	task, err := s.cronManager.GetTask(name)
	if err != nil {
		return fmt.Errorf("查找任务失败: %v", err)
	}

	// 执行任务处理函数
	if err := task.Handler(); err != nil {
		return fmt.Errorf("任务执行失败: %v", err)
	}

	return nil
}

// 统计服务示例
type StatsService struct {
	taskService *TaskService
	stats       map[string]int
	mu          sync.Mutex
}

// 创建新的统计服务
func NewStatsService(taskService *TaskService) *StatsService {
	return &StatsService{
		taskService: taskService,
		stats:       make(map[string]int),
	}
}

// 设置定时统计任务
func (s *StatsService) SetupDailyStats() error {
	// 每天凌晨1点执行统计
	return s.taskService.AddTask("daily-stats", "0 0 1 * * *", func() error {
		// 在实际应用中，会在此处理统计逻辑
		s.mu.Lock()
		defer s.mu.Unlock()

		s.stats["daily_runs"]++
		log.Printf("执行每日统计，已累计执行 %d 次", s.stats["daily_runs"])
		return nil
	})
}

// 设置每小时统计任务
func (s *StatsService) SetupHourlyStats() error {
	// 每小时执行一次统计
	return s.taskService.AddTask("hourly-stats", "0 0 * * * *", func() error {
		// 在实际应用中，会在此处理统计逻辑
		s.mu.Lock()
		defer s.mu.Unlock()

		s.stats["hourly_runs"]++
		log.Printf("执行每小时统计，已累计执行 %d 次", s.stats["hourly_runs"])
		return nil
	})
}

// 获取统计数据
func (s *StatsService) GetStats() map[string]int {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 复制一份数据返回
	result := make(map[string]int)
	for k, v := range s.stats {
		result[k] = v
	}
	return result
}

// 测试基本定时任务
func TestCronTaskWithDI(t *testing.T) {
	// 创建容器
	c := container.NewContainer()

	// 创建任务服务
	taskService := NewTaskService(c)

	// 启动定时任务
	taskService.Start()
	defer taskService.Stop()

	// 变量用于验证任务是否执行
	var jobExecuted bool
	var mu sync.Mutex

	// 添加每秒执行一次的任务
	err := taskService.AddTask("test-job", "* * * * * *", func() error {
		mu.Lock()
		jobExecuted = true
		mu.Unlock()
		return nil
	})

	if err != nil {
		t.Fatalf("添加任务失败: %v", err)
	}

	// 等待任务执行
	time.Sleep(1200 * time.Millisecond)

	mu.Lock()
	executed := jobExecuted
	mu.Unlock()

	if !executed {
		t.Error("任务未按预期执行")
	} else {
		t.Log("任务成功执行")
	}
}

// 测试统计服务
func TestStatsService(t *testing.T) {
	// 创建容器
	c := container.NewContainer()

	// 创建任务服务
	taskService := NewTaskService(c)

	// 创建统计服务
	statsService := NewStatsService(taskService)

	// 设置统计任务
	if err := statsService.SetupHourlyStats(); err != nil {
		t.Fatalf("设置每小时统计任务失败: %v", err)
	}

	// 启动定时任务
	taskService.Start()
	defer taskService.Stop()

	// 手动触发任务以验证功能
	err := taskService.RunTask("hourly-stats")
	if err != nil {
		t.Fatalf("手动触发任务失败: %v", err)
	}

	// 验证统计数据
	stats := statsService.GetStats()
	if hourlyRuns := stats["hourly_runs"]; hourlyRuns != 1 {
		t.Errorf("每小时统计任务执行次数错误，期望: 1, 实际: %d", hourlyRuns)
	} else {
		t.Log("统计服务正常工作")
	}
}

// 示例：实际使用场景的定时任务
func ExampleCronUsage() {
	// 在实际应用中，容器会在应用启动时初始化
	c := container.NewContainer()

	// 创建任务服务
	taskService := NewTaskService(c)

	// 创建统计服务
	statsService := NewStatsService(taskService)

	// 设置定时任务
	statsService.SetupDailyStats()
	statsService.SetupHourlyStats()

	// 启动定时任务
	taskService.Start()
	defer taskService.Stop()

	fmt.Println("定时统计任务已启动")

	// 运行一段时间...
	time.Sleep(1 * time.Second)
}
