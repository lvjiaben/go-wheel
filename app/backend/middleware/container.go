package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/pkg/container"
)

// ContainerMiddleware 设置容器到上下文的中间件
func ContainerMiddleware(c *container.Container) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 设置容器到gin上下文
		ctx.Set("container", c)
		ctx.Next()
	}
}
