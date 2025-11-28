package routes

import (
	"html/template"
	"os"

	"github.com/gin-gonic/gin"
	v1 "github.com/lvjiaben/go-wheel/app/api/v1/controller"
	apiMiddleware "github.com/lvjiaben/go-wheel/app/api/v1/middleware"
	"github.com/lvjiaben/go-wheel/app/backend/controller"
	"github.com/lvjiaben/go-wheel/app/backend/controller/admin"
	"github.com/lvjiaben/go-wheel/app/backend/controller/system"
	"github.com/lvjiaben/go-wheel/app/backend/controller/user"
	backendMiddleware "github.com/lvjiaben/go-wheel/app/backend/middleware"
	websocketController "github.com/lvjiaben/go-wheel/app/websocket/controller"
	"github.com/lvjiaben/go-wheel/pkg/container"
	"github.com/lvjiaben/go-wheel/pkg/middleware"
)

// RegisterRoutes 注册所有路由
func RegisterRoutes(r *gin.Engine, c *container.Container) {
	// CORS 跨域中间件（最先执行）
	r.Use(middleware.CorsMiddleware())

	// 注册全局中间件
	r.Use(middleware.I18nMiddleware(c))

	// 添加容器中间件，确保所有路由都能访问container
	r.Use(middleware.ContainerMiddleware(c))

	// 加载HTML模板（优先本地文件，支持热更新）
	if _, err := os.Stat("app/views"); err == nil {
		// 本地文件存在，从文件系统加载（开发环境）
		r.LoadHTMLGlob("app/views/*.html")
	} else if viewsFS := c.GetViewsFS(); viewsFS != nil {
		// 从嵌入的文件系统加载（生产环境）
		tmpl := template.Must(template.ParseFS(viewsFS, "app/views/*.html"))
		r.SetHTMLTemplate(tmpl)
	}

	// 注册前端页面路由
	registerFrontendRoutes(r, c)

	// 注册后台路由
	registerBackendRoutes(r, c)

	// 注册API路由
	registerApiRoutes(r, c)

	// 注册 WebSocket 路由
	registerWebSocketRoutes(r, c)

}

// registerFrontendRoutes 注册前端页面路由
func registerFrontendRoutes(r *gin.Engine, c *container.Container) {
	// 静态资源
	r.Static("/public", "./public")

	// favicon.ico 返回空响应，避免404日志
	r.GET("/favicon.ico", func(ctx *gin.Context) {
		ctx.Status(204)
	})

	// 创建控制器实例
	indexController := v1.NewIndexController(c)

	// 首页路由
	r.GET("/", indexController.Index)
}

// registerBackendRoutes 注册后台路由
func registerBackendRoutes(r *gin.Engine, c *container.Container) {
	// 创建控制器实例
	adminController := admin.NewAdminController(c)
	roleController := admin.NewRoleController(c)
	menuController := admin.NewMenuController(c)
	authController := controller.NewAuthController(c)

	// 获取中间件
	authMiddleware := backendMiddleware.NewAuthMiddleware(c)

	// API分组
	api := r.Group("/backend")
	{
		// 公共控制器
		commonController := controller.NewCommonController(c)
		common := api.Group("/common")
		{
			common.POST("/captcha", commonController.Captcha)
		}
		// 控制台控制器
		homeController := controller.NewHomeController(c)
		home := api.Group("/home", authMiddleware.JWTAuthCheck(), authMiddleware.PermissionCheck())
		{
			home.GET("/index", homeController.Index)
		}
		// 认证相关路由 - 无需JWT认证
		auth := api.Group("/auth")
		{
			auth.POST("/login", authController.Login)
			auth.POST("/logout", authMiddleware.JWTAuthCheck(), authController.Logout)
			auth.GET("/menus", authMiddleware.JWTAuthCheck(), authController.Menus)
			auth.GET("/permission", authMiddleware.JWTAuthCheck(), authController.Permission)
			auth.POST("/profile", authMiddleware.JWTAuthCheck(), authController.Profile)
			auth.POST("/password", authMiddleware.JWTAuthCheck(), authController.Password)
			auth.GET("/log", authMiddleware.JWTAuthCheck(), authController.Log)
			auth.GET("/info", authMiddleware.JWTAuthCheck(), authController.Info)
		}

		// 需要JWT认证的路由
		adminGroup := api.Group("/admin", authMiddleware.JWTAuthCheck(), authMiddleware.PermissionCheck())
		{
			// 管理员管理
			adminRoutes := adminGroup.Group("/admin")
			{
				adminRoutes.GET("/list", adminController.List)
				adminRoutes.POST("/save", adminController.Save)
				adminRoutes.DELETE("/delete/:id", adminController.Delete)
			}

			// 角色管理
			role := adminGroup.Group("/role")
			{
				role.GET("/list", roleController.List)
				role.POST("/save", roleController.Save)
				role.DELETE("/delete/:id", roleController.Delete)
				role.GET("/my-menus", roleController.GetMyMenus)
			}

			// 菜单管理
			menu := adminGroup.Group("/menu")
			{
				menu.GET("/list", menuController.List)
				menu.POST("/save", menuController.Save)
				menu.POST("/delete", menuController.Delete)
			}
		}

		// 系统管理
		systemGroup := api.Group("/system").Use(authMiddleware.JWTAuthCheck()).Use(authMiddleware.PermissionCheck())
		{
			// 附件管理
			attachmentController := system.NewAttachmentController(c)
			systemGroup.GET("/attachment/directories", attachmentController.Directories)
			systemGroup.GET("/attachment/list", attachmentController.List)
			systemGroup.POST("/attachment/upload", attachmentController.Upload)
			systemGroup.POST("/attachment/delete", attachmentController.Delete)

			// 配置管理
			configController := system.NewConfigController(c)
			systemGroup.GET("/config/list", configController.List)
			systemGroup.POST("/config/create", configController.Create)
			systemGroup.POST("/config/update", configController.Update)
			systemGroup.DELETE("/config/delete/:id", configController.Delete)

			// 代码生成
			genController := system.NewGenController(c)
			systemGroup.GET("/gen/table-list", genController.TableList)
			systemGroup.GET("/gen/table-info", genController.TableInfo)
			systemGroup.GET("/gen/table-config", genController.TableConfig)
			systemGroup.POST("/gen/preview", genController.Preview)
			systemGroup.POST("/gen/generate", genController.Generate)
			systemGroup.GET("/gen/history", genController.History)
			systemGroup.POST("/gen/delete", genController.Delete)
			systemGroup.POST("/gen/download", genController.Download)
		}

		// 用户管理
		userGroup := api.Group("/user").Use(authMiddleware.JWTAuthCheck()).Use(authMiddleware.PermissionCheck())
		{
			userController := user.NewUserController(c)
			userGroup.GET("/list", userController.List)
			userGroup.POST("/create", userController.Create)
			userGroup.POST("/update", userController.Update)
			userGroup.POST("/delete", userController.Delete)
			userGroup.POST("/update-money", userController.UpdateMoney)
			userGroup.POST("/update-score", userController.UpdateScore)
			userGroup.POST("/operate", userController.Operate)
		}

		// 下方注释不要删除 为代码生成的预留位置
		// region:backend-routes

		// endregion:backend-routes
	}

}

// registerApiRoutes 注册API路由（前端用户接口）
func registerApiRoutes(r *gin.Engine, c *container.Container) {
	// 创建控制器实例
	userController := v1.NewUserController(c)
	captchaController := v1.NewCaptchaController(c)
	smsController := v1.NewSmsController(c)
	menuController := v1.NewMenuController(c)

	// 获取中间件
	authMiddleware := apiMiddleware.NewAuthMiddleware(c)

	// API分组
	api := r.Group("/api")
	{
		common := api.Group("/common")
		{
			common.POST("/captcha", captchaController.Generate)
		}

		// 短信 - 可选登录（登录和未登录都可访问）
		sms := api.Group("/sms", authMiddleware.OptionalJWTAuth())
		{
			sms.POST("/send", smsController.Send)
		}

		// 用户认证 - 无需登录
		userAuth := api.Group("/user")
		{
			userAuth.POST("/login", userController.Login)
			userAuth.POST("/mobile-login", userController.MobileLogin)
			userAuth.POST("/reg", userController.Reg)
			userAuth.POST("/reset-pwd", userController.ResetPwd)
		}

		// 用户操作 - 需要登录
		userGroup := api.Group("/user", authMiddleware.JWTAuthCheck())
		{
			userGroup.POST("/logout", userController.Logout)
			userGroup.POST("/change-mobile", userController.ChangeMobile)
			userGroup.POST("/change-pwd", userController.ChangePwd)
			userGroup.GET("/info", userController.Info)
		}

		// 菜单 - 需要登录
		menu := api.Group("/menu", authMiddleware.JWTAuthCheck())
		{
			menu.GET("/all", menuController.All)
		}
	}
}

// registerWebSocketRoutes 注册 WebSocket 路由
func registerWebSocketRoutes(r *gin.Engine, c *container.Container) {
	// 获取 WebSocket Hub
	wsHub := c.GetWebSocketHub()

	// WebSocket 路由组（需要认证）
	// 中间件执行顺序：
	// 1. ContainerMiddleware - 注入容器
	// 2. WebSocketAuthMiddleware - JWT 认证，设置 user_id 和 username 到上下文
	// 3. WebSocketUpgradeMiddleware - 升级为 WebSocket，从上下文获取用户信息
	ws := r.Group("/ws").Use(
		middleware.ContainerMiddleware(c),
		middleware.WebSocketAuthMiddleware(c),
		middleware.WebSocketUpgradeMiddleware(wsHub),
	)
	{
		pingController := websocketController.NewPingController(c, wsHub)
		ws.GET("/ping",
			pingController.Connect, // 控制器处理
		)
	}
}
