package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/pkg/container"
)

// ContainerMiddleware 设置容器到上下文的中间件
// 这是一个公共底层中间件，用于将 Container 注入到 Gin 上下文中
func ContainerMiddleware(c *container.Container) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 设置容器到gin上下文
		ctx.Set("container", c)
		ctx.Next()
	}
}

