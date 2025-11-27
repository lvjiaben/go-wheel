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
	"time"

	_ "github.com/lvjiaben/go-wheel/app/cron"  // 导入任务包以触发 init()
	_ "github.com/lvjiaben/go-wheel/app/queue" // 导入队列消费者包以触发 init()
	commonService "github.com/lvjiaben/go-wheel/app/common/service"
	cronPkg "github.com/lvjiaben/go-wheel/pkg/cron"
	queuePkg "github.com/lvjiaben/go-wheel/pkg/queue"
	"github.com/lvjiaben/go-wheel/pkg/constants"
	"github.com/lvjiaben/go-wheel/pkg/middleware"
	"github.com/lvjiaben/go-wheel/pkg/monitor"
	"github.com/lvjiaben/go-wheel/routes"

	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/pkg/container"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)



func main() {
	// 初始化容器
	c := container.NewContainer()
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

	// 初始化 Prometheus 监控
	prometheusMetrics := monitor.NewPrometheusMetrics(c, "goweb")
	c.Set("prometheus_metrics", prometheusMetrics) // 存储到 container 供其他组件使用
	c.GetLogger().Info("Prometheus 监控已初始化")

	// 启动资源监控
	resourceMonitor := monitor.NewResourceMonitor(c, constants.DefaultResourceMonitorInterval)
	c.SetResourceMonitor(resourceMonitor)
	c.TrackGoroutine()
	go func() {
		defer c.UntrackGoroutine()
		resourceMonitor.Start()
	}()
	defer resourceMonitor.Stop()
	c.GetLogger().Info("资源监控已启动", zap.Duration("收集间隔", constants.DefaultResourceMonitorInterval))

	// 启动 Prometheus 指标收集
	c.TrackGoroutine()
	go func() {
		defer c.UntrackGoroutine()
		ticker := time.NewTicker(constants.DefaultPrometheusCollectInterval)
		defer ticker.Stop()
		for {
			select {
			case <-c.GetContext().Done():
				return
			case <-ticker.C:
				prometheusMetrics.Collect()
			}
		}
	}()
	c.GetLogger().Info("Prometheus 指标收集已启动", zap.Duration("收集间隔", constants.DefaultPrometheusCollectInterval))

	// 注册gin
	gin.SetMode(c.GetConfig().App.Mode)
	r := gin.New()
	// 使用自定义中间件 - 将日志输出到Zap
	r.Use(middleware.GinLogger(c))
	r.Use(middleware.GinRecovery(c))
	// 使用请求体大小限制中间件
	r.Use(middleware.RequestBodyLimitMiddleware(c))
	// 使用 Prometheus 监控中间件
	r.Use(middleware.PrometheusMiddleware(prometheusMetrics))

	// 注册 Prometheus metrics 端点
	r.GET("/metrics", gin.WrapH(promhttp.Handler()))

	// 注册路由
	routes.RegisterRoutes(r, c)
	// 检查配置的端口是否可用
	addr := fmt.Sprintf(":%d", c.GetConfig().App.Port)
	listener, err := net.Listen("tcp", addr)
	listener.Close()
	if err != nil {

		log.Fatalf("程序端口被占用: %v", err)
	}

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
