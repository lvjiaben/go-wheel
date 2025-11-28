# Captcha 验证码

项目使用 [base64Captcha](https://github.com/mojocn/base64Captcha) 生成图形验证码，支持 Redis 存储。

## 功能特点

- 生成 Base64 格式图片验证码
- 支持内存存储和 Redis 存储
- 可配置验证码样式
- 自动过期清理

## 基本用法

### 创建服务

```go
import "github.com/lvjiaben/go-wheel/pkg/captcha"

// 使用 Redis 存储（推荐）
captchaService := captcha.NewCaptchaService(redisClient)

// 使用内存存储
captchaService := captcha.NewCaptchaService(nil)
```

### 生成验证码

```go
result, err := captchaService.Generate()
if err != nil {
    // 处理错误
}

// result.ID     - 验证码ID（用于验证）
// result.Base64 - Base64编码的图片
```

### 验证验证码

```go
isValid := captchaService.Verify(captchaId, userInput)
if !isValid {
    // 验证失败
}
```

## 在控制器中使用

```go
// app/backend/controller/common.go
type CommonController struct {
    captchaService *captcha.CaptchaService
}

func NewCommonController(c *container.Container) *CommonController {
    return &CommonController{
        captchaService: captcha.NewCaptchaService(c.GetRedis().Client),
    }
}

// Captcha 获取验证码
func (c *CommonController) Captcha(ctx *gin.Context) {
    result, err := c.captchaService.Generate()
    if err != nil {
        http.Error(ctx, "生成验证码失败")
        return
    }
    
    http.Success(ctx, gin.H{
        "id":     result.ID,
        "base64": result.Base64,
    })
}

// Login 登录验证
func (c *CommonController) Login(ctx *gin.Context) {
    var req LoginRequest
    ctx.ShouldBindJSON(&req)
    
    // 验证验证码
    if !c.captchaService.Verify(req.CaptchaId, req.CaptchaCode) {
        http.Error(ctx, "验证码错误")
        return
    }
    
    // 继续登录逻辑...
}
```

## 前端使用

```vue
<template>
  <div class="captcha">
    <img :src="captchaBase64" @click="refreshCaptcha" />
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { getCaptcha } from '@/api/common'

const captchaId = ref('')
const captchaBase64 = ref('')

const refreshCaptcha = async () => {
  const { data } = await getCaptcha()
  captchaId.value = data.id
  captchaBase64.value = data.base64
}

onMounted(() => {
  refreshCaptcha()
})
</script>
```

## 配置说明

验证码默认配置：

```go
// pkg/captcha/captcha.go
captchaConfig = struct {
    Height          int     // 高度: 80
    Width           int     // 宽度: 240
    Length          int     // 字符数: 4
    NoiseCount      float64 // 干扰点: 0
    ShowLineOptions int     // 干扰线: 0
    BgColor         struct {
        R, G, B, A int      // 背景色: 白色
    }
    Fonts []string          // 字体
}
```

## Redis 存储

使用 Redis 存储验证码，支持分布式部署：

```go
// pkg/captcha/captcha.go
type RedisStore struct {
    client *redis.Client
    ttl    time.Duration  // 默认 5 分钟
}

// 存储验证码
func (r *RedisStore) Set(id string, value string) error

// 获取验证码
func (r *RedisStore) Get(id string, clear bool) string

// 验证验证码
func (r *RedisStore) Verify(id, answer string, clear bool) bool
```

## 返回格式

```json
{
    "code": 200,
    "message": "success",
    "data": {
        "id": "abc123",
        "base64": "data:image/png;base64,iVBORw0KGgo..."
    }
}
```

## 安全建议

1. **限制请求频率** - 防止恶意刷验证码
2. **设置过期时间** - 验证码 5 分钟内有效
3. **一次性使用** - 验证后立即删除
4. **大小写不敏感** - 提升用户体验

