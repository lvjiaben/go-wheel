# Redis

项目使用 [go-redis](https://github.com/redis/go-redis) 作为 Redis 客户端，并提供了封装的 `RedisWrapper`。

## 配置

在 `configs/config.yaml` 中配置：

```yaml
redis:
  state: true             # 是否启用
  host: "127.0.0.1"
  port: 6379
  pass: ""                # 密码
  db: 0                   # 数据库编号
  pool_size: 100          # 连接池大小
  min_idle_conns: 10      # 最小空闲连接
  max_conn_age: 3600      # 连接最大生命周期（秒）
  idle_timeout: 300       # 空闲超时（秒）
  pool_timeout: 5         # 获取连接超时（秒）
  read_timeout: 3         # 读取超时（秒）
  write_timeout: 3        # 写入超时（秒）
```

## 基本用法

```go
redis := container.GetRedis()
ctx := context.Background()
```

### 字符串操作

```go
// 设置值
redis.Set(ctx, "key", "value", time.Hour)

// 获取值
val, err := redis.Get(ctx, "key")

// 设置值（仅当键不存在时）- 用于分布式锁
ok, err := redis.SetNX(ctx, "lock:key", "1", 10*time.Second)

// 删除键
redis.Del(ctx, "key1", "key2")

// 检查键是否存在
count, err := redis.Exists(ctx, "key")

// 设置过期时间
redis.Expire(ctx, "key", time.Hour)

// 获取过期时间
ttl, err := redis.TTL(ctx, "key")
```

### 计数器

```go
// 递增
val, err := redis.Incr(ctx, "counter")

// 递减
val, err := redis.Decr(ctx, "counter")
```

### 哈希表

```go
// 设置字段
redis.HSet(ctx, "user:1", "name", "张三")

// 获取字段
name, err := redis.HGet(ctx, "user:1", "name")

// 获取所有字段
data, err := redis.HGetAll(ctx, "user:1")

// 删除字段
redis.HDel(ctx, "user:1", "name")
```

### 列表

```go
// 左侧插入
redis.LPush(ctx, "queue", "item1", "item2")

// 右侧插入
redis.RPush(ctx, "queue", "item1", "item2")

// 左侧弹出
val, err := redis.LPop(ctx, "queue")

// 右侧弹出
val, err := redis.RPop(ctx, "queue")

// 获取范围
items, err := redis.LRange(ctx, "queue", 0, -1)

// 获取长度
length, err := redis.LLen(ctx, "queue")
```

## 分布式锁

```go
// 获取锁
locked, err := redis.SetNX(ctx, "lock:order:123", "1", 30*time.Second)
if !locked {
    return errors.New("获取锁失败")
}

// 释放锁
defer redis.Del(ctx, "lock:order:123")

// 执行业务逻辑...
```

## 缓存示例

```go
// 获取用户信息（带缓存）
func (s *UserService) GetUserWithCache(userId int) (*User, error) {
    ctx := context.Background()
    cacheKey := fmt.Sprintf("user:%d", userId)
    
    // 尝试从缓存获取
    data, err := s.redis.Get(ctx, cacheKey)
    if err == nil {
        var user User
        json.Unmarshal([]byte(data), &user)
        return &user, nil
    }
    
    // 从数据库获取
    var user User
    if err := s.db.First(&user, userId).Error; err != nil {
        return nil, err
    }
    
    // 写入缓存
    jsonData, _ := json.Marshal(user)
    s.redis.Set(ctx, cacheKey, jsonData, time.Hour)
    
    return &user, nil
}
```

## 使用原生客户端

如需使用更多 Redis 命令，可以直接使用原生客户端：

```go
// 获取原生客户端
client := container.GetRedis().Client

// 使用原生命令
client.SAdd(ctx, "set:key", "member1", "member2")
client.SMembers(ctx, "set:key")
client.ZAdd(ctx, "zset:key", redis.Z{Score: 1, Member: "a"})
```

## 连接池监控

```go
// 获取连接池状态
stats := container.GetRedis().Client.PoolStats()

fmt.Printf("总连接数: %d\n", stats.TotalConns)
fmt.Printf("空闲连接: %d\n", stats.IdleConns)
fmt.Printf("过期连接: %d\n", stats.StaleConns)
```

## 最佳实践

1. **合理设置过期时间** - 避免缓存永不过期
2. **使用前缀** - 如 `user:`, `order:` 区分不同业务
3. **序列化** - 复杂对象使用 JSON 序列化
4. **错误处理** - 缓存失败不应影响主流程
5. **连接池** - 根据并发量调整连接池大小

