package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/app/backend/controller"
	backendMiddleware "github.com/lvjiaben/go-wheel/app/backend/middleware"
	"github.com/lvjiaben/go-wheel/pkg/container"
	globalMiddleware "github.com/lvjiaben/go-wheel/pkg/middleware"
)

func RegisterRoutes(r *gin.Engine, c *container.Container) {
	// 注册全局中间件
	r.Use(globalMiddleware.I18nMiddleware(c))

	// 创建控制器实例
	authController := controller.NewAuthController(c)
	authMiddleware := backendMiddleware.NewAuthMiddleware(c)

	// 后端路由组
	backend := r.Group("/backend")
	{
		// 认证路由组
		auth := backend.Group("/auth")
		{
			auth.POST("/login", authController.Login)
			auth.POST("/logout", authMiddleware.JWTAuthCheck(), authController.Logout)
			auth.GET("/codes", authMiddleware.JWTAuthCheck(), authController.Codes)
			auth.GET("/menus", authMiddleware.JWTAuthCheck(), authController.Menus)
			auth.GET("/user", authMiddleware.JWTAuthCheck(), authController.User)
		}
	}
}
