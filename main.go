package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"admin/pkg/container"
	"admin/pkg/initialize"
	"admin/pkg/middleware"
	"admin/routes"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var (
	configFile = flag.String("c", "configs/config.yaml", "配置文件路径")
)

func main() {
	flag.Parse()

	// 检查配置文件是否存在
	if _, err := os.Stat(*configFile); os.IsNotExist(err) {
		log.Fatalf("配置文件 %s 不存在", *configFile)
	}

	// 创建容器
	container := container.NewContainer()
	defer container.Shutdown()

	// 初始化配置
	config := initialize.InitConfig(*configFile)
	container.SetConfig(config)

	// 初始化日志
	logger := initialize.ZapLoad()
	defer logger.Sync()
	logger.Debug("Logger init success")
	container.SetLogger(logger)

	// 初始化数据库
	db := initialize.InitDB()
	if db != nil {
		sqlDB, _ := db.DB()
		defer sqlDB.Close()
		container.SetDB(db)
	}

	// 初始化Redis
	if config.Redis.State {
		redis := initialize.RedisLoad()
		defer redis.Close()
		container.SetRedis(redis)
	}

	// 初始化验证器
	initialize.ValidateLoad()

	// 初始化多语言
	i18n := initialize.NewI18n()
	if err := i18n.LoadTranslations("translations"); err != nil {
		logger.Error("加载翻译文件失败", zap.Error(err))
	}
	container.SetI18n(i18n)

	// 初始化消息队列
	mq := initialize.NewMessageQueue()
	container.SetMessageQueue(mq)

	// 初始化延迟队列
	dq := initialize.NewDelayQueue()
	container.SetDelayQueue(dq)

	// 初始化定时任务
	cron := initialize.NewCronManager()
	cron.Start()
	defer cron.Stop()
	container.SetCron(cron)

	// 启动GIN
	gin.SetMode(config.App.Mode)
	r := gin.New()
	r.Use(middleware.GinLogger(), middleware.GinRecovery(true))

	// 注册路由
	routes.RegisterRoutes(r, container)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", config.App.Port),
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 优雅退出
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	logger.Info("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server Shutdown", zap.Error(err))
	}
	logger.Info("Server exiting")
}
