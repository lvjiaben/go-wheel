package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	commonService "github.com/lvjiaben/go-wheel/app/common/service"
	_ "github.com/lvjiaben/go-wheel/app/cron"  // 导入任务包以触发 init()
	_ "github.com/lvjiaben/go-wheel/app/queue" // 导入队列消费者包以触发 init()
	"github.com/lvjiaben/go-wheel/pkg/constants"
	cronPkg "github.com/lvjiaben/go-wheel/pkg/cron"
	"github.com/lvjiaben/go-wheel/pkg/middleware"
	queuePkg "github.com/lvjiaben/go-wheel/pkg/queue"
	"github.com/lvjiaben/go-wheel/routes"

	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/pkg/container"
	"go.uber.org/zap"
)

func main() {
	// 初始化容器（传入嵌入的文件系统）
	c := container.NewContainer(&container.EmbedFS{
		ConfigFS: ConfigFS,
		I18nFS:   I18nFS,
		ViewsFS:  ViewsFS,
	})
	defer c.Shutdown()
	if err := c.Initialize(); err != nil {
		log.Fatalf("初始化失败: %v", err)
	}

	// 启动配置缓存服务
	configCache := commonService.NewConfigCacheService(c)
	configCache.Start(constants.DefaultConfigCacheInterval)
	defer configCache.Stop()

	c.GetLogger().Info("配置缓存服务已启动", zap.Duration("轮询间隔", constants.DefaultConfigCacheInterval))

	// 启动定时任务
	if err := cronPkg.RegisterAllTasks(c.AsCronContainer()); err != nil {
		log.Fatalf("注册定时任务失败: %v", err)
	}
	c.GetCron().Start()
	defer c.GetCron().Stop()

	// 启动消息队列消费者
	if c.GetRabbitMQ() != nil {
		if err := queuePkg.StartAllConsumers(c.AsQueueContainer(), c.GetRabbitMQ()); err != nil {
			log.Fatalf("启动消息队列消费者失败: %v", err)
		}
	}

	// 注册gin
	gin.SetMode(c.GetConfig().App.Mode)
	r := gin.New()
	// 使用自定义中间件 - 将日志输出到Zap
	r.Use(middleware.GinLogger(c))
	r.Use(middleware.GinRecovery(c))
	// 使用请求体大小限制中间件
	r.Use(middleware.RequestBodyLimitMiddleware(c))

	// 注册路由
	routes.RegisterRoutes(r, c)
	// 检查配置的端口是否可用
	addr := fmt.Sprintf("0.0.0.0:%d", c.GetConfig().App.Port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("程序端口被占用: %v", err)
	}
	listener.Close()

	// 启动服务器
	srv := &http.Server{
		Addr:    addr,
		Handler: r,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			c.GetLogger().Fatal("服务器启动失败", zap.Error(err))
		}
	}()
	// 优雅关机，等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	c.GetLogger().Info("Shutdown Server ...")

	// 设置超时时间
	ctx, cancel := context.WithTimeout(context.Background(), constants.GracefulShutdownTimeout)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		c.GetLogger().Fatal("Server Shutdown", zap.Error(err))
	}
	c.GetLogger().Info("Server exiting")
}
