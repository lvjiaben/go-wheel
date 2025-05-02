package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/lvjiaben/go-wheel/pkg/container"
	"github.com/redis/go-redis/v9"
)

// Redis客户端配置
type RedisConfig struct {
	Host     string
	Port     int
	Password string
	DB       int
}

// Redis客户端包装器
type RedisClient struct {
	client *redis.Client
	ctx    context.Context
}

// 创建新的Redis客户端
func NewRedisClient(config RedisConfig) *RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.Host, config.Port),
		Password: config.Password,
		DB:       config.DB,
	})

	return &RedisClient{
		client: client,
		ctx:    context.Background(),
	}
}

// 关闭Redis连接
func (r *RedisClient) Close() error {
	return r.client.Close()
}

// 设置字符串值
func (r *RedisClient) Set(key string, value interface{}, expiration time.Duration) error {
	return r.client.Set(r.ctx, key, value, expiration).Err()
}

// 获取字符串值
func (r *RedisClient) Get(key string) (string, error) {
	return r.client.Get(r.ctx, key).Result()
}

// 检查键是否存在
func (r *RedisClient) Exists(key string) (bool, error) {
	result, err := r.client.Exists(r.ctx, key).Result()
	return result > 0, err
}

// 删除键
func (r *RedisClient) Delete(key string) error {
	return r.client.Del(r.ctx, key).Err()
}

// 设置过期时间
func (r *RedisClient) Expire(key string, expiration time.Duration) error {
	return r.client.Expire(r.ctx, key, expiration).Err()
}

// 哈希表操作 - 设置字段值
func (r *RedisClient) HSet(key, field string, value interface{}) error {
	return r.client.HSet(r.ctx, key, field, value).Err()
}

// 哈希表操作 - 获取字段值
func (r *RedisClient) HGet(key, field string) (string, error) {
	return r.client.HGet(r.ctx, key, field).Result()
}

// 哈希表操作 - 获取所有字段和值
func (r *RedisClient) HGetAll(key string) (map[string]string, error) {
	return r.client.HGetAll(r.ctx, key).Result()
}

// 列表操作 - 左侧推入
func (r *RedisClient) LPush(key string, values ...interface{}) error {
	return r.client.LPush(r.ctx, key, values...).Err()
}

// 列表操作 - 右侧推入
func (r *RedisClient) RPush(key string, values ...interface{}) error {
	return r.client.RPush(r.ctx, key, values...).Err()
}

// 列表操作 - 左侧弹出
func (r *RedisClient) LPop(key string) (string, error) {
	return r.client.LPop(r.ctx, key).Result()
}

// 列表操作 - 获取列表范围
func (r *RedisClient) LRange(key string, start, stop int64) ([]string, error) {
	return r.client.LRange(r.ctx, key, start, stop).Result()
}

// 集合操作 - 添加成员
func (r *RedisClient) SAdd(key string, members ...interface{}) error {
	return r.client.SAdd(r.ctx, key, members...).Err()
}

// 集合操作 - 获取所有成员
func (r *RedisClient) SMembers(key string) ([]string, error) {
	return r.client.SMembers(r.ctx, key).Result()
}

// 集合操作 - 判断成员是否存在
func (r *RedisClient) SIsMember(key string, member interface{}) (bool, error) {
	return r.client.SIsMember(r.ctx, key, member).Result()
}

// 有序集合操作 - 添加成员
func (r *RedisClient) ZAdd(key string, score float64, member interface{}) error {
	return r.client.ZAdd(r.ctx, key, redis.Z{Score: score, Member: member}).Err()
}

// 有序集合操作 - 获取范围
func (r *RedisClient) ZRange(key string, start, stop int64) ([]string, error) {
	return r.client.ZRange(r.ctx, key, start, stop).Result()
}

// 测试Redis缓存扩展类，用于测试
type TestRedisCache struct {
	redisClient *RedisClient // 使用我们自己的客户端
}

// 创建测试用Redis缓存
func NewTestRedisCache(config RedisConfig) (*TestRedisCache, error) {
	client := NewRedisClient(config)

	return &TestRedisCache{
		redisClient: client,
	}, nil
}

// Set 实现缓存设置
func (t *TestRedisCache) Set(ctx context.Context, key string, value interface{}, ttl time.Duration) error {
	// 序列化值为JSON
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return t.redisClient.Set(key, string(data), ttl)
}

// Get 实现缓存获取
func (t *TestRedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	// 获取缓存值
	data, err := t.redisClient.Get(key)
	if err != nil {
		return err
	}

	// 反序列化JSON到目标对象
	return json.Unmarshal([]byte(data), dest)
}

// Del 实现缓存删除
func (t *TestRedisCache) Del(ctx context.Context, key string) error {
	return t.redisClient.Delete(key)
}

// Exists 检查键是否存在
func (t *TestRedisCache) Exists(ctx context.Context, key string) (bool, error) {
	return t.redisClient.Exists(key)
}

// 清理资源
func (t *TestRedisCache) Close() error {
	if t.redisClient != nil {
		t.redisClient.Close()
	}
	return nil
}

// 用户模型，用于测试JSON序列化
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Age      int    `json:"age"`
}

// 用户缓存服务示例
type UserCacheService struct {
	container   *container.Container
	cacheClient *TestRedisCache // 使用测试缓存客户端
}

// 创建新的用户缓存服务
func NewUserCacheService(c *container.Container, cacheClient *TestRedisCache) *UserCacheService {
	return &UserCacheService{
		container:   c,
		cacheClient: cacheClient,
	}
}

// 缓存用户信息
func (s *UserCacheService) CacheUser(ctx context.Context, user *User) error {
	// 用户缓存键
	cacheKey := fmt.Sprintf("user:%s", user.ID)

	// 缓存用户数据，有效期30分钟
	return s.cacheClient.Set(ctx, cacheKey, user, 30*time.Minute)
}

// 获取缓存的用户信息
func (s *UserCacheService) GetCachedUser(ctx context.Context, userID string) (*User, error) {
	// 用户缓存键
	cacheKey := fmt.Sprintf("user:%s", userID)

	// 从缓存获取用户
	var user User
	err := s.cacheClient.Get(ctx, cacheKey, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// 移除用户缓存
func (s *UserCacheService) InvalidateUserCache(ctx context.Context, userID string) error {
	// 用户缓存键
	cacheKey := fmt.Sprintf("user:%s", userID)

	// 删除缓存
	return s.cacheClient.Del(ctx, cacheKey)
}

// 产品缓存服务示例
type ProductCacheService struct {
	container   *container.Container
	cacheClient *TestRedisCache
}

// 产品模型
type Product struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
	Stock       int     `json:"stock"`
}

// 创建产品缓存服务
func NewProductCacheService(c *container.Container, cacheClient *TestRedisCache) *ProductCacheService {
	return &ProductCacheService{
		container:   c,
		cacheClient: cacheClient,
	}
}

// 缓存产品
func (s *ProductCacheService) CacheProduct(ctx context.Context, product *Product) error {
	cacheKey := fmt.Sprintf("product:%s", product.ID)
	return s.cacheClient.Set(ctx, cacheKey, product, 1*time.Hour)
}

// 获取缓存产品
func (s *ProductCacheService) GetCachedProduct(ctx context.Context, productID string) (*Product, error) {
	cacheKey := fmt.Sprintf("product:%s", productID)

	var product Product
	err := s.cacheClient.Get(ctx, cacheKey, &product)
	if err != nil {
		return nil, err
	}

	return &product, nil
}

// 更新产品库存
func (s *ProductCacheService) UpdateProductStock(ctx context.Context, productID string, newStock int) error {
	product, err := s.GetCachedProduct(ctx, productID)
	if err != nil {
		return err
	}

	product.Stock = newStock
	return s.CacheProduct(ctx, product)
}

// 使用依赖注入测试Redis缓存
func TestRedisCacheWithDI(t *testing.T) {
	// 如果没有Redis环境，则跳过测试
	t.Skip("需要Redis环境才能运行此测试")

	// 创建容器
	c := container.NewContainer()

	// 创建Redis配置
	redisConfig := RedisConfig{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		DB:       0,
	}

	// 创建测试Redis缓存
	testRedisCache, err := NewTestRedisCache(redisConfig)
	if err != nil {
		t.Fatalf("创建测试Redis缓存失败: %v", err)
	}

	// 创建用户缓存服务
	userCacheService := NewUserCacheService(c, testRedisCache)

	// 测试用户数据
	testUser := &User{
		ID:       "user123",
		Username: "testuser",
		Email:    "test@example.com",
		Age:      30,
	}

	// 测试缓存写入
	ctx := context.Background()
	err = userCacheService.CacheUser(ctx, testUser)
	if err != nil {
		t.Fatalf("缓存用户数据失败: %v", err)
	}

	// 检查键是否存在
	exists, err := testRedisCache.Exists(ctx, "user:user123")
	if err != nil {
		t.Fatalf("检查缓存键存在失败: %v", err)
	}

	if !exists {
		t.Error("用户缓存键应该存在，但返回不存在")
	}

	// 测试缓存读取
	cachedUser, err := userCacheService.GetCachedUser(ctx, testUser.ID)
	if err != nil {
		t.Fatalf("获取缓存用户失败: %v", err)
	}

	if cachedUser.Username != testUser.Username || cachedUser.Email != testUser.Email {
		t.Errorf("缓存的用户数据不匹配，期望: %+v, 实际: %+v", testUser, cachedUser)
	}

	// 测试缓存删除
	err = userCacheService.InvalidateUserCache(ctx, testUser.ID)
	if err != nil {
		t.Fatalf("移除用户缓存失败: %v", err)
	}

	// 验证缓存已删除
	exists, err = testRedisCache.Exists(ctx, "user:user123")
	if err != nil {
		t.Fatalf("检查缓存键存在失败: %v", err)
	}

	if exists {
		t.Error("用户缓存键应该已被删除，但仍然存在")
	}

	t.Log("Redis用户缓存服务测试通过")
}

// 测试产品缓存服务
func TestProductCacheService(t *testing.T) {
	// 如果没有Redis环境，则跳过测试
	t.Skip("需要Redis环境才能运行此测试")

	// 创建容器
	c := container.NewContainer()

	// 创建Redis配置
	redisConfig := RedisConfig{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		DB:       0,
	}

	// 创建测试Redis缓存
	testRedisCache, err := NewTestRedisCache(redisConfig)
	if err != nil {
		t.Fatalf("创建测试Redis缓存失败: %v", err)
	}

	// 创建产品缓存服务
	productCacheService := NewProductCacheService(c, testRedisCache)

	// 测试产品数据
	testProduct := &Product{
		ID:          "product-123",
		Name:        "测试产品",
		Price:       99.99,
		Description: "这是一个测试产品",
		Stock:       100,
	}

	// 测试缓存写入
	ctx := context.Background()
	err = productCacheService.CacheProduct(ctx, testProduct)
	if err != nil {
		t.Fatalf("缓存产品数据失败: %v", err)
	}

	// 测试缓存读取
	cachedProduct, err := productCacheService.GetCachedProduct(ctx, testProduct.ID)
	if err != nil {
		t.Fatalf("获取缓存产品失败: %v", err)
	}

	if cachedProduct.Name != testProduct.Name || cachedProduct.Price != testProduct.Price {
		t.Errorf("缓存的产品数据不匹配，期望: %+v, 实际: %+v", testProduct, cachedProduct)
	}

	// 测试更新库存
	newStock := 50
	err = productCacheService.UpdateProductStock(ctx, testProduct.ID, newStock)
	if err != nil {
		t.Fatalf("更新产品库存失败: %v", err)
	}

	// 验证库存已更新
	updatedProduct, err := productCacheService.GetCachedProduct(ctx, testProduct.ID)
	if err != nil {
		t.Fatalf("获取更新后的产品失败: %v", err)
	}

	if updatedProduct.Stock != newStock {
		t.Errorf("产品库存未正确更新，期望: %d, 实际: %d", newStock, updatedProduct.Stock)
	}

	t.Log("Redis产品缓存服务测试通过")
}

// 计数器服务示例
type CounterService struct {
	container   *container.Container
	cacheClient *TestRedisCache
}

// 创建计数器服务
func NewCounterService(c *container.Container, cacheClient *TestRedisCache) *CounterService {
	return &CounterService{
		container:   c,
		cacheClient: cacheClient,
	}
}

// 增加计数
func (s *CounterService) Increment(ctx context.Context, key string, value int64) (int64, error) {
	// 这里应该调用 INCRBY 操作，但我们模拟实现
	var currentValue int64
	err := s.cacheClient.Get(ctx, key, &currentValue)
	if err != nil && err.Error() != "redis: nil" {
		return 0, err
	}

	newValue := currentValue + value
	err = s.cacheClient.Set(ctx, key, newValue, 0) // 无过期时间
	if err != nil {
		return 0, err
	}

	return newValue, nil
}

// 获取当前计数
func (s *CounterService) Get(ctx context.Context, key string) (int64, error) {
	var value int64
	err := s.cacheClient.Get(ctx, key, &value)
	if err != nil {
		return 0, err
	}
	return value, nil
}

// 重置计数
func (s *CounterService) Reset(ctx context.Context, key string) error {
	return s.cacheClient.Set(ctx, key, int64(0), 0)
}

// 测试计数器服务
func TestCounterService(t *testing.T) {
	// 如果没有Redis环境，则跳过测试
	t.Skip("需要Redis环境才能运行此测试")

	// 创建容器
	c := container.NewContainer()

	// 创建Redis配置
	redisConfig := RedisConfig{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		DB:       0,
	}

	// 创建测试Redis缓存
	testRedisCache, err := NewTestRedisCache(redisConfig)
	if err != nil {
		t.Fatalf("创建测试Redis缓存失败: %v", err)
	}

	// 创建计数器服务
	counterService := NewCounterService(c, testRedisCache)

	// 测试计数器
	ctx := context.Background()
	counterKey := "test:counter"

	// 重置计数器
	err = counterService.Reset(ctx, counterKey)
	if err != nil {
		t.Fatalf("重置计数器失败: %v", err)
	}

	// 测试增加计数
	newValue, err := counterService.Increment(ctx, counterKey, 5)
	if err != nil {
		t.Fatalf("增加计数失败: %v", err)
	}

	if newValue != 5 {
		t.Errorf("计数器值不正确，期望: 5, 实际: %d", newValue)
	}

	// 再次增加
	newValue, err = counterService.Increment(ctx, counterKey, 3)
	if err != nil {
		t.Fatalf("增加计数失败: %v", err)
	}

	if newValue != 8 {
		t.Errorf("计数器值不正确，期望: 8, 实际: %d", newValue)
	}

	// 测试获取当前计数
	value, err := counterService.Get(ctx, counterKey)
	if err != nil {
		t.Fatalf("获取计数器值失败: %v", err)
	}

	if value != 8 {
		t.Errorf("计数器当前值不正确，期望: 8, 实际: %d", value)
	}

	t.Log("Redis计数器服务测试通过")
}

// 实际使用Redis缓存的示例
func ExampleRedisCacheUsage() {
	// 在实际应用中，容器会在应用启动时初始化
	c := container.NewContainer()

	// 创建Redis配置
	redisConfig := RedisConfig{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		DB:       0,
	}

	// 创建测试Redis缓存
	testRedisCache, _ := NewTestRedisCache(redisConfig)

	// 创建用户缓存服务
	userCacheService := NewUserCacheService(c, testRedisCache)

	// 模拟用户数据
	user := &User{
		ID:       "example123",
		Username: "exampleuser",
		Email:    "example@example.com",
		Age:      25,
	}

	// 缓存用户
	ctx := context.Background()
	userCacheService.CacheUser(ctx, user)

	// 从缓存获取用户
	cachedUser, _ := userCacheService.GetCachedUser(ctx, user.ID)
	fmt.Printf("从缓存获取用户: %s (%s)\n", cachedUser.Username, cachedUser.Email)
}

// 计数器使用示例
func ExampleRedisCounter() {
	// 在实际应用中，容器会在应用启动时初始化
	c := container.NewContainer()

	// 创建Redis配置
	redisConfig := RedisConfig{
		Host:     "localhost",
		Port:     6379,
		Password: "",
		DB:       0,
	}

	// 创建测试Redis缓存
	testRedisCache, _ := NewTestRedisCache(redisConfig)

	// 创建计数器服务
	counterService := NewCounterService(c, testRedisCache)

	// 使用计数器记录页面访问
	ctx := context.Background()
	pageKey := "visits:homepage"

	// 增加访问计数
	newCount, _ := counterService.Increment(ctx, pageKey, 1)

	fmt.Printf("页面访问次数: %d\n", newCount)
}
