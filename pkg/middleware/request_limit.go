package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/pkg/container"
)

// RequestBodyLimitMiddleware 请求体大小限制中间件
func RequestBodyLimitMiddleware(c *container.Container) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 从配置获取最大请求体大小（MB）
		maxSizeMB := c.GetConfig().App.MaxRequestBody
		if maxSizeMB <= 0 {
			maxSizeMB = 10 // 默认 10MB
		}
		
		maxSize := maxSizeMB << 20 // 转换为字节（MB * 1024 * 1024）
		
		// 限制请求体大小
		ctx.Request.Body = http.MaxBytesReader(ctx.Writer, ctx.Request.Body, maxSize)
		
		ctx.Next()
	}
}

