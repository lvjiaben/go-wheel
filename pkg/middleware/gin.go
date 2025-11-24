package middleware

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/pkg/container"
	"go.uber.org/zap"
)

// GinLogger 自定义Gin日志中间件，使用Zap记录请求日志
func GinLogger(c *container.Container) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		start := time.Now()
		path := ctx.Request.URL.Path
		query := ctx.Request.URL.RawQuery
		method := ctx.Request.Method
		ip := ctx.ClientIP()

		ctx.Next()

		// 请求结束后记录日志
		cost := time.Since(start)
		status := ctx.Writer.Status()

		// 根据状态码决定日志级别
		if status >= 400 {
			// 错误请求使用Error级别
			c.GetLogger().Error("[GIN]",
				zap.String("method", method),
				zap.Int("status", status),
				zap.Duration("cost", cost),
				zap.String("ip", ip),
				zap.String("path", path),
				zap.String("query", query),
				zap.String("error", ctx.Errors.String()),
			)
		} else {
			// 正常请求使用Info级别
			c.GetLogger().Info("[GIN]",
				zap.String("method", method),
				zap.Int("status", status),
				zap.Duration("cost", cost),
				zap.String("ip", ip),
				zap.String("path", path),
				zap.String("query", query),
			)
		}
	}
}

// GinRecovery 自定义Gin恢复中间件，使用Zap记录panic信息
func GinRecovery(c *container.Container) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				// 获取请求信息（不包含查询参数，避免泄露敏感信息）
				httpRequest := ctx.Request.Method + " " + ctx.Request.URL.Path

				// 记录错误日志（不记录请求体，避免泄露密码等敏感信息）
				c.GetLogger().Error("[Recovery] panic recovered",
					zap.Any("error", err),
					zap.String("request", httpRequest),
					zap.String("ip", ctx.ClientIP()),
				)
				// 返回500错误
				ctx.AbortWithStatus(500)
			}
		}()
		ctx.Next()
	}
}
