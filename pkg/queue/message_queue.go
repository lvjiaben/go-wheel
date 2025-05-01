package queue

import (
	"context"
	"encoding/json"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// Message 消息结构体
type Message struct {
	ID        string                 `json:"id"`
	Topic     string                 `json:"topic"`
	Data      map[string]interface{} `json:"data"`
	CreatedAt time.Time              `json:"created_at"`
}

// MessageQueue 消息队列接口
type MessageQueue interface {
	Publish(ctx context.Context, topic string, data map[string]interface{}) error
	Subscribe(ctx context.Context, topic string, handler func(*Message) error) error
}

// RedisMessageQueue Redis消息队列实现
type RedisMessageQueue struct {
	client *redis.Client
}

// NewRedisMessageQueue 创建Redis消息队列实例
func NewRedisMessageQueue(client *redis.Client) *RedisMessageQueue {
	return &RedisMessageQueue{client: client}
}

// Publish 发布消息
func (q *RedisMessageQueue) Publish(ctx context.Context, topic string, data map[string]interface{}) error {
	message := Message{
		ID:        time.Now().Format("20060102150405.000"),
		Topic:     topic,
		Data:      data,
		CreatedAt: time.Now(),
	}

	bytes, err := json.Marshal(message)
	if err != nil {
		return err
	}

	return q.client.LPush(ctx, "queue:"+topic, bytes).Err()
}

// Subscribe 订阅消息
func (q *RedisMessageQueue) Subscribe(ctx context.Context, topic string, handler func(*Message) error) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			result, err := q.client.BRPop(ctx, 0, "queue:"+topic).Result()
			if err != nil {
				if err == redis.Nil {
					continue
				}
				return err
			}

			var message Message
			if err := json.Unmarshal([]byte(result[1]), &message); err != nil {
				continue
			}

			if err := handler(&message); err != nil {
				// 处理失败，将消息重新放入队列
				bytes, _ := json.Marshal(message)
				_ = q.client.LPush(ctx, "queue:"+topic, bytes).Err()
			}
		}
	}
}

// MemoryMessageQueue 内存消息队列实现
type MemoryMessageQueue struct {
	queues map[string]chan *Message
	mu     sync.RWMutex
}

// NewMemoryMessageQueue 创建内存消息队列实例
func NewMemoryMessageQueue() *MemoryMessageQueue {
	return &MemoryMessageQueue{
		queues: make(map[string]chan *Message),
	}
}

// Publish 发布消息
func (q *MemoryMessageQueue) Publish(ctx context.Context, topic string, data map[string]interface{}) error {
	message := &Message{
		ID:        time.Now().Format("20060102150405.000"),
		Topic:     topic,
		Data:      data,
		CreatedAt: time.Now(),
	}

	q.mu.Lock()
	queue, ok := q.queues[topic]
	if !ok {
		queue = make(chan *Message, 1000)
		q.queues[topic] = queue
	}
	q.mu.Unlock()

	select {
	case queue <- message:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// Subscribe 订阅消息
func (q *MemoryMessageQueue) Subscribe(ctx context.Context, topic string, handler func(*Message) error) error {
	q.mu.Lock()
	queue, ok := q.queues[topic]
	if !ok {
		queue = make(chan *Message, 1000)
		q.queues[topic] = queue
	}
	q.mu.Unlock()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case message := <-queue:
			if err := handler(message); err != nil {
				// 处理失败，将消息重新放入队列
				select {
				case queue <- message:
				default:
					// 队列已满，丢弃消息
				}
			}
		}
	}
}
