package types

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisWrapper Redis客户端包装器
type RedisWrapper struct {
	Client *redis.Client
}

// Close 关闭连接
func (r *RedisWrapper) Close() error {
	return r.Client.Close()
}

// Ping 检查连接
func (r *RedisWrapper) Ping(ctx context.Context) error {
	return r.Client.Ping(ctx).Err()
}

// Set 设置键值对
func (r *RedisWrapper) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return r.Client.Set(ctx, key, value, expiration).Err()
}

// SetNX 设置键值对（仅当键不存在时）用于分布式锁
func (r *RedisWrapper) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return r.Client.SetNX(ctx, key, value, expiration).Result()
}

// Get 获取值
func (r *RedisWrapper) Get(ctx context.Context, key string) (string, error) {
	return r.Client.Get(ctx, key).Result()
}

// Del 删除键
func (r *RedisWrapper) Del(ctx context.Context, keys ...string) error {
	return r.Client.Del(ctx, keys...).Err()
}

// Exists 检查键是否存在
func (r *RedisWrapper) Exists(ctx context.Context, keys ...string) (int64, error) {
	return r.Client.Exists(ctx, keys...).Result()
}

// Expire 设置过期时间
func (r *RedisWrapper) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return r.Client.Expire(ctx, key, expiration).Err()
}

// TTL 获取过期时间
func (r *RedisWrapper) TTL(ctx context.Context, key string) (time.Duration, error) {
	return r.Client.TTL(ctx, key).Result()
}

// Incr 原子递增
func (r *RedisWrapper) Incr(ctx context.Context, key string) (int64, error) {
	return r.Client.Incr(ctx, key).Result()
}

// Decr 原子递减
func (r *RedisWrapper) Decr(ctx context.Context, key string) (int64, error) {
	return r.Client.Decr(ctx, key).Result()
}

// HSet 设置哈希表字段值
func (r *RedisWrapper) HSet(ctx context.Context, key string, field string, value interface{}) error {
	return r.Client.HSet(ctx, key, field, value).Err()
}

// HGet 获取哈希表字段值
func (r *RedisWrapper) HGet(ctx context.Context, key, field string) (string, error) {
	return r.Client.HGet(ctx, key, field).Result()
}

// HGetAll 获取哈希表所有字段和值
func (r *RedisWrapper) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return r.Client.HGetAll(ctx, key).Result()
}

// HDel 删除哈希表字段
func (r *RedisWrapper) HDel(ctx context.Context, key string, fields ...string) error {
	return r.Client.HDel(ctx, key, fields...).Err()
}

// LPush 将一个或多个值插入到列表头部
func (r *RedisWrapper) LPush(ctx context.Context, key string, values ...interface{}) error {
	return r.Client.LPush(ctx, key, values...).Err()
}

// RPush 将一个或多个值插入到列表尾部
func (r *RedisWrapper) RPush(ctx context.Context, key string, values ...interface{}) error {
	return r.Client.RPush(ctx, key, values...).Err()
}

// LPop 移出并获取列表的第一个元素
func (r *RedisWrapper) LPop(ctx context.Context, key string) (string, error) {
	return r.Client.LPop(ctx, key).Result()
}

// RPop 移出并获取列表的最后一个元素
func (r *RedisWrapper) RPop(ctx context.Context, key string) (string, error) {
	return r.Client.RPop(ctx, key).Result()
}

// LRange 获取列表指定范围内的元素
func (r *RedisClient) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return r.Client.LRange(ctx, key, start, stop).Result()
}

// LLen 获取列表长度
func (r *RedisClient) LLen(ctx context.Context, key string) (int64, error) {
	return r.Client.LLen(ctx, key).Result()
}
