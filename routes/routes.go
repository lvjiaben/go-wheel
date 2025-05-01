package routes

import (
	"admin/app/backend/controller"
	"admin/app/backend/middleware"
	"admin/pkg/container"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine, container *container.Container) {
	r.GET("/", func(c *gin.Context) {
		c.String(200, "Hello, Gin!")
	})
	//后端接口
	backend := r.Group("/backend", middleware.SetLang())
	{
		// Auth控制器
		auth := backend.Group("/auth")
		{
			authController := &controller.AuthController{}
			auth.POST("/login", authController.Login)
			auth.GET("/codes", middleware.JWTAuthCheck(), authController.Codes)
			auth.GET("/menus", middleware.JWTAuthCheck(), authController.Menus)
			auth.GET("/user", middleware.JWTAuthCheck(), authController.User)
		}
	}
}
