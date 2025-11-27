// package cron

// import (
// 	"context"
// 	"time"

// 	cronPkg "github.com/lvjiaben/go-wheel/pkg/cron"
// 	"go.uber.org/zap"
// )

// // 自动注册任务
// func init() {
// 	cronPkg.Register("分布式任务示例", NewExampleDistributedTask)
// }

// type ExampleDistributedTask struct {
// 	cronPkg.BaseTask
// 	container cronPkg.Container
// }

// // NewExampleDistributedTask 创建分布式任务示例
// func NewExampleDistributedTask(c cronPkg.Container) cronPkg.Task {
// 	return &ExampleDistributedTask{
// 		BaseTask: cronPkg.BaseTask{
// 			Name:            "分布式任务示例",
// 			Spec:            "0 */5 * * * *",    // 每5分钟执行一次
// 			Description:     "演示分布式锁的使用",
// 			DistributedLock: true,               // 启用分布式锁
// 			LockTimeout:     60 * time.Second,   // 锁超时时间60秒
// 		},
// 		container: c,
// 	}
// }

// // Run 执行任务
// func (t *ExampleDistributedTask) Run(ctx context.Context) error {
// 	logger := t.container.GetLogger()
	
// 	logger.Info("开始执行分布式任务",
// 		zap.String("task", t.GetName()),
// 	)
	
// 	// 模拟耗时操作
// 	select {
// 	case <-ctx.Done():
// 		logger.Warn("任务被取消")
// 		return ctx.Err()
// 	case <-time.After(5 * time.Second):
// 		// 业务逻辑
// 		logger.Info("分布式任务执行完成")
// 	}
	
// 	return nil
// }

