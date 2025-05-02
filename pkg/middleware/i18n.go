package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/pkg/container"
	"go.uber.org/zap"
)

// 用于在Container中保存Gin上下文的键
const GinContextKey = "gin_context"

// WithGinContext 将gin.Context保存到context中
func WithGinContext(parent context.Context, c *gin.Context) context.Context {
	return context.WithValue(parent, GinContextKey, c)
}

// GinContextFromContext 从context中获取gin.Context
func GinContextFromContext(ctx context.Context) (*gin.Context, bool) {
	ginCtx, ok := ctx.Value(GinContextKey).(*gin.Context)
	return ginCtx, ok
}

// I18nMiddleware 语言中间件
func I18nMiddleware(c *container.Container) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从请求头获取语言
		lang := ctx.GetHeader("Accept-Language")
		if lang == "" {
			lang = "zh-CN" // 默认中文
		}

		// 支持不同格式的语言代码
		// 如果使用简写格式"zh"，转换为完整格式"zh-CN"
		if lang == "zh" {
			lang = "zh-CN"
		} else if lang == "en" {
			lang = "en-US"
		}

		// 记录检测到的语言
		c.GetLogger().Debug("检测到的语言",
			zap.String("accept-language", ctx.GetHeader("Accept-Language")),
			zap.String("language", lang))

		// 设置语言到上下文
		ctx.Set("lang", lang)

		// 将gin上下文包装到Container的上下文中
		c.SetContext(WithGinContext(c.GetContext(), ctx))

		ctx.Next()
	}
}

// Translate 翻译函数
func Translate(c *container.Container, key string) string {
	lang := "zh-CN" // 默认中文

	// 尝试从Container的上下文中获取Gin上下文
	if ginCtx, ok := GinContextFromContext(c.GetContext()); ok && ginCtx != nil {
		if l, exists := ginCtx.Get("lang"); exists {
			if l, ok := l.(string); ok {
				lang = l
			}
		}
	}

	// 获取翻译结果
	result := c.GetI18n().Get(key, lang)

	// 记录翻译过程
	c.GetLogger().Debug("翻译过程",
		zap.String("key", key),
		zap.String("lang", lang),
		zap.String("result", result))

	return result
}
