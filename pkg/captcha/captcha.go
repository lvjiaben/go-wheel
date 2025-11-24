package captcha

import (
	"time"

	"image/color"

	"github.com/mojocn/base64Captcha"
	"github.com/redis/go-redis/v9"
	"golang.org/x/net/context"
)

// 定义验证码参数
var (
	// 默认存储驱动
	defaultStore base64Captcha.Store
	// 验证码配置
	captchaConfig = struct {
		Height          int     // 高度
		Width           int     // 宽度
		Length          int     // 长度
		NoiseCount      float64 // 干扰数量
		ShowLineOptions int     // 展示线条数量
		BgColor         struct {
			R int
			G int
			B int
			A int
		} // 背景颜色
		Fonts []string // 字体
	}{
		Height:          80,
		Width:           240,
		Length:          4,
		NoiseCount:      0, // 完全禁用噪点
		ShowLineOptions: 0, // 完全禁用线条选项，避免索引越界
		BgColor: struct {
			R int
			G int
			B int
			A int
		}{
			R: 255,
			G: 255,
			B: 255,
			A: 128,
		},
		Fonts: []string{"wqy-microhei.ttc"},
	}
)

// CaptchaService 验证码服务
type CaptchaService struct {
	store base64Captcha.Store
}

// CaptchaResult 验证码结果
type CaptchaResult struct {
	ID     string `json:"id"`     // 验证码ID
	Base64 string `json:"base64"` // Base64编码
}

// NewCaptchaService 创建验证码服务
func NewCaptchaService(redisClient *redis.Client) *CaptchaService {
	var store base64Captcha.Store
	if redisClient != nil {
		store = NewRedisStore(redisClient, time.Minute*5)
	} else {
		store = base64Captcha.DefaultMemStore
	}

	defaultStore = store
	return &CaptchaService{
		store: store,
	}
}

// Generate 生成验证码
func (s *CaptchaService) Generate() (*CaptchaResult, error) {
	// 使用更简单的方法 - 创建字符串验证码配置
	driver := &base64Captcha.DriverString{
		Height:          captchaConfig.Height,
		Width:           captchaConfig.Width,
		Length:          captchaConfig.Length,
		Source:          "234567890abcdefghjkmnpqrstuvwxyz",
		ShowLineOptions: 0, // 不显示干扰线
		NoiseCount:      0, // 不显示噪点
		BgColor: &color.RGBA{
			R: uint8(captchaConfig.BgColor.R),
			G: uint8(captchaConfig.BgColor.G),
			B: uint8(captchaConfig.BgColor.B),
			A: uint8(captchaConfig.BgColor.A),
		},
	}

	// 创建验证码
	c := base64Captcha.NewCaptcha(driver, s.store)
	id, b64s, _, err := c.Generate()
	if err != nil {
		return nil, err
	}

	// 返回结果
	return &CaptchaResult{
		ID:     id,
		Base64: b64s,
	}, nil
}

// Verify 验证验证码
func (s *CaptchaService) Verify(id string, value string) bool {
	return s.store.Verify(id, value, true)
}

// RedisStore Redis存储
type RedisStore struct {
	client *redis.Client
	ttl    time.Duration
}

// NewRedisStore 创建Redis存储
func NewRedisStore(client *redis.Client, ttl time.Duration) *RedisStore {
	return &RedisStore{
		client: client,
		ttl:    ttl,
	}
}

// Set 设置值
func (r *RedisStore) Set(id string, value string) error {
	ctx := context.Background()
	err := r.client.Set(ctx, id, value, r.ttl).Err()
	return err
}

// Get 获取值
func (r *RedisStore) Get(id string, clear bool) string {
	ctx := context.Background()
	val, err := r.client.Get(ctx, id).Result()
	if err != nil {
		return ""
	}
	if clear {
		_, _ = r.client.Del(ctx, id).Result()
	}
	return val
}

// Verify 验证
func (r *RedisStore) Verify(id, answer string, clear bool) bool {
	val := r.Get(id, clear)
	return val == answer
}
