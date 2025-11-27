package container

import (
	cronPkg "github.com/lvjiaben/go-wheel/pkg/cron"
	queuePkg "github.com/lvjiaben/go-wheel/pkg/queue"
)

// 确保 Container 实现 cronPkg.ContainerWithCron 接口
var _ cronPkg.ContainerWithCron = (*CronAdapter)(nil)

// 确保 Container 实现 queuePkg.ContainerWithQueue 接口
var _ queuePkg.ContainerWithQueue = (*QueueAdapter)(nil)

// CronAdapter Cron适配器（嵌入Container，重写返回类型）
type CronAdapter struct {
	*Container
}

// AsCronContainer 将Container转换为CronAdapter
func (c *Container) AsCronContainer() cronPkg.ContainerWithCron {
	return &CronAdapter{Container: c}
}

// GetLogger 实现 cronPkg.Logger 接口
func (ca *CronAdapter) GetLogger() cronPkg.Logger {
	return ca.Container.GetLogger()
}

// GetDB 返回 interface{} 类型
func (ca *CronAdapter) GetDB() interface{} {
	return ca.Container.GetDB()
}

// GetRedis 返回 interface{} 类型
func (ca *CronAdapter) GetRedis() interface{} {
	return ca.Container.GetRedis()
}

// GetCron 返回 interface{} 类型
func (ca *CronAdapter) GetCron() interface{} {
	return ca.Container.GetCron()
}

// QueueAdapter 队列适配器（嵌入Container，重写返回类型）
type QueueAdapter struct {
	*Container
}

// AsQueueContainer 将Container转换为QueueAdapter
func (c *Container) AsQueueContainer() queuePkg.ContainerWithQueue {
	return &QueueAdapter{Container: c}
}

// GetLogger 实现 queuePkg.Logger 接口
func (qa *QueueAdapter) GetLogger() queuePkg.Logger {
	return qa.Container.GetLogger()
}

// GetDB 返回 interface{} 类型
func (qa *QueueAdapter) GetDB() interface{} {
	return qa.Container.GetDB()
}

// GetRedis 返回 interface{} 类型
func (qa *QueueAdapter) GetRedis() interface{} {
	return qa.Container.GetRedis()
}

// GetMessageQueue 返回 interface{} 类型
func (qa *QueueAdapter) GetMessageQueue() interface{} {
	return qa.Container.GetMessageQueue()
}

// GetDelayQueue 返回 interface{} 类型
func (qa *QueueAdapter) GetDelayQueue() interface{} {
	return qa.Container.GetDelayQueue()
}

