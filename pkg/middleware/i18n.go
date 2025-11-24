package middleware

import (
	"context"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/pkg/container"
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
		if lang != "" {
			// 处理多种分隔符的情况
			langs := strings.FieldsFunc(lang, func(r rune) bool {
				return r == ',' || r == ' ' || r == ';'
			})
			if len(langs) > 0 {
				// 取第一个语言代码，并移除质量值
				lang = strings.Split(langs[0], ";")[0]
			}
		}
		if lang == "" || lang == "zh" || lang == "zh-CN" || lang == "zh-cn" {
			lang = "zh-CN"
		} else {
			lang = "en-US"
		}

		// 设置语言到 Gin 上下文（推荐方式）
		ctx.Set("lang", lang)
		ctx.Set("isCn", lang == "zh-CN")

		ctx.Next()
	}
}

// Translate 翻译函数（从 Gin Context 获取语言）
func Translate(c *container.Container, ctx *gin.Context, key string) string {
	lang := "zh-CN" // 默认中文

	// 从 Gin Context 获取语言设置
	if ctx != nil {
		if l, exists := ctx.Get("lang"); exists {
			if langStr, ok := l.(string); ok {
				lang = langStr
			}
		}
	}

	// 获取翻译结果
	result := c.GetI18n().Get(key, lang)

	return result
}

// TranslateWithLang 翻译函数（直接指定语言）
func TranslateWithLang(c *container.Container, key string, lang string) string {
	if lang == "" {
		lang = "zh-CN"
	}
	return c.GetI18n().Get(key, lang)
}
