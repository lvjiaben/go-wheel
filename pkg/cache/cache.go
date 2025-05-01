package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Cache 缓存接口
type Cache interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Del(ctx context.Context, key string) error
	Exists(ctx context.Context, key string) (bool, error)
	Expire(ctx context.Context, key string, expiration time.Duration) error
	TTL(ctx context.Context, key string) (time.Duration, error)
}

// RedisCache Redis缓存实现
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache 创建Redis缓存实例
func NewRedisCache(client *redis.Client) *RedisCache {
	return &RedisCache{client: client}
}

// Get 获取缓存
func (c *RedisCache) Get(ctx context.Context, key string) (string, error) {
	return c.client.Get(ctx, key).Result()
}

// Set 设置缓存
func (c *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	return c.client.Set(ctx, key, value, expiration).Err()
}

// Del 删除缓存
func (c *RedisCache) Del(ctx context.Context, key string) error {
	return c.client.Del(ctx, key).Err()
}

// Exists 判断缓存是否存在
func (c *RedisCache) Exists(ctx context.Context, key string) (bool, error) {
	result, err := c.client.Exists(ctx, key).Result()
	return result > 0, err
}

// Expire 设置过期时间
func (c *RedisCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	return c.client.Expire(ctx, key, expiration).Err()
}

// TTL 获取过期时间
func (c *RedisCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	return c.client.TTL(ctx, key).Result()
}

// MemoryCache 内存缓存实现
type MemoryCache struct {
	data map[string]item
}

type item struct {
	value      interface{}
	expiration time.Time
}

// NewMemoryCache 创建内存缓存实例
func NewMemoryCache() *MemoryCache {
	return &MemoryCache{
		data: make(map[string]item),
	}
}

// Get 获取缓存
func (c *MemoryCache) Get(ctx context.Context, key string) (string, error) {
	if item, ok := c.data[key]; ok {
		if item.expiration.IsZero() || time.Now().Before(item.expiration) {
			if str, ok := item.value.(string); ok {
				return str, nil
			}
		} else {
			delete(c.data, key)
		}
	}
	return "", nil
}

// Set 设置缓存
func (c *MemoryCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	var exp time.Time
	if expiration > 0 {
		exp = time.Now().Add(expiration)
	}
	c.data[key] = item{
		value:      value,
		expiration: exp,
	}
	return nil
}

// Del 删除缓存
func (c *MemoryCache) Del(ctx context.Context, key string) error {
	delete(c.data, key)
	return nil
}

// Exists 判断缓存是否存在
func (c *MemoryCache) Exists(ctx context.Context, key string) (bool, error) {
	if item, ok := c.data[key]; ok {
		if item.expiration.IsZero() || time.Now().Before(item.expiration) {
			return true, nil
		}
		delete(c.data, key)
	}
	return false, nil
}

// Expire 设置过期时间
func (c *MemoryCache) Expire(ctx context.Context, key string, expiration time.Duration) error {
	if item, ok := c.data[key]; ok {
		if expiration > 0 {
			item.expiration = time.Now().Add(expiration)
			c.data[key] = item
		}
	}
	return nil
}

// TTL 获取过期时间
func (c *MemoryCache) TTL(ctx context.Context, key string) (time.Duration, error) {
	if item, ok := c.data[key]; ok {
		if item.expiration.IsZero() {
			return -1, nil
		}
		if time.Now().Before(item.expiration) {
			return item.expiration.Sub(time.Now()), nil
		}
		delete(c.data, key)
	}
	return -2, nil
}
