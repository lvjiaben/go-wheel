package queue

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/lvjiaben/go-wheel/pkg/global"
)

type DelayQueue struct {
	client *redis.Client
}

type DelayTask struct {
	ID        string      `json:"id"`
	Topic     string      `json:"topic"`
	Data      interface{} `json:"data"`
	ExecuteAt int64       `json:"execute_at"`
}

func NewDelayQueue() *DelayQueue {
	return &DelayQueue{
		client: global.RDB,
	}
}

// AddTask 添加延迟任务
func (q *DelayQueue) AddTask(ctx context.Context, task *DelayTask) error {
	data, err := json.Marshal(task)
	if err != nil {
		return err
	}

	// 使用ZADD命令将任务添加到有序集合
	// 使用执行时间作为分数
	_, err = q.client.ZAdd(ctx, fmt.Sprintf("delay_queue:%s", task.Topic), redis.Z{
		Score:  float64(task.ExecuteAt),
		Member: data,
	}).Result()
	return err
}

// GetTask 获取到期的任务
func (q *DelayQueue) GetTask(ctx context.Context, topic string) (*DelayTask, error) {
	// 使用ZRANGEBYSCORE获取当前时间之前的任务
	now := time.Now().Unix()
	results, err := q.client.ZRangeByScore(ctx, fmt.Sprintf("delay_queue:%s", topic), &redis.ZRangeBy{
		Min: "0",
		Max: fmt.Sprintf("%d", now),
	}).Result()
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, nil
	}

	// 获取并删除第一个任务
	taskData := results[0]
	var task DelayTask
	if err := json.Unmarshal([]byte(taskData), &task); err != nil {
		return nil, err
	}

	// 从队列中删除任务
	_, err = q.client.ZRem(ctx, fmt.Sprintf("delay_queue:%s", topic), taskData).Result()
	if err != nil {
		return nil, err
	}

	return &task, nil
}

// StartConsumer 启动消费者
func (q *DelayQueue) StartConsumer(ctx context.Context, topic string, handler func(task *DelayTask) error) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			default:
				task, err := q.GetTask(ctx, topic)
				if err != nil {
					global.LOG.Error("获取延迟任务失败", err)
					time.Sleep(time.Second)
					continue
				}

				if task == nil {
					time.Sleep(time.Second)
					continue
				}

				if err := handler(task); err != nil {
					global.LOG.Error("处理延迟任务失败", err)
					// 处理失败的任务可以重新入队或记录到失败队列
				}
			}
		}
	}()
}
