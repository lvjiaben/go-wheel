package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/pkg/container"
	"github.com/lvjiaben/go-wheel/pkg/initialize"
	"go.uber.org/zap"
)

var (
	configFile = flag.String("c", "configs/config.yaml", "配置文件路径")
)

// 检查端口是否可用
func isPortAvailable(port int) bool {
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return false
	}
	listener.Close()
	return true
}

// 查找下一个可用端口
func findAvailablePort(startPort int) int {
	port := startPort
	for i := 0; i < 10; i++ { // 尝试10个端口
		if isPortAvailable(port) {
			return port
		}
		port++
	}
	return 0 // 没有找到可用端口
}

func main() {
	flag.Parse()

	// 检查配置文件是否存在
	if _, err := os.Stat(*configFile); os.IsNotExist(err) {
		log.Fatalf("配置文件 %s 不存在", *configFile)
	}

	// 初始化容器
	c := container.NewContainer()
	defer c.Shutdown()

	// 初始化配置
	if err := initialize.ViperLoad(c); err != nil {
		log.Fatalf("配置初始化失败: %v", err)
	}

	// 初始化日志
	if err := initialize.ZapLoad(c); err != nil {
		log.Fatalf("日志初始化失败: %v", err)
	}

	c.GetLogger().Debug("Logger init success")

	// 初始化数据库
	if err := initialize.MysqlLoad(c); err != nil {
		c.GetLogger().Error("数据库初始化失败", zap.Error(err))
	} else {
		db, _ := c.GetDB().DB()
		defer db.Close()
	}

	// 初始化Redis
	if c.GetConfig().Redis.State {
		if err := initialize.RedisLoad(c); err != nil {
			c.GetLogger().Error("Redis初始化失败", zap.Error(err))
		}
	}

	// 初始化验证器
	if err := initialize.ValidateLoad(c); err != nil {
		c.GetLogger().Error("验证器初始化失败", zap.Error(err))
	}

	// 初始化多语言
	i18n := initialize.NewI18n()
	if err := i18n.LoadTranslations("configs/i18n"); err != nil {
		c.GetLogger().Error("加载语言文件失败", zap.Error(err))
	}
	c.SetI18n(i18n)

	// 初始化消息队列
	mq := initialize.NewMessageQueue()
	c.SetMessageQueue(mq)

	// 初始化延迟队列
	dq := initialize.NewDelayQueue()
	c.SetDelayQueue(dq)

	// 初始化定时任务
	cron := initialize.NewCronManager()
	cron.Start()
	defer cron.Stop()
	c.SetCron(cron)

	// 设置gin模式
	gin.SetMode(c.GetConfig().App.Mode)

	// 初始化路由
	router := initialize.RoutersLoad(c)

	// 检查配置的端口是否可用，如果不可用则自动选择一个可用端口
	port := c.GetConfig().App.Port
	if !isPortAvailable(port) {
		newPort := findAvailablePort(port + 1)
		if newPort == 0 {
			c.GetLogger().Fatal("无法找到可用端口")
		}
		c.GetLogger().Info("原端口被占用，切换到新端口", zap.Int("original_port", port), zap.Int("new_port", newPort))
		port = newPort
	}

	// 启动服务器
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			c.GetLogger().Fatal("服务器启动失败", zap.Error(err))
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	c.GetLogger().Info("Shutdown Server ...")

	// 设置超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		c.GetLogger().Fatal("Server Shutdown", zap.Error(err))
	}
	c.GetLogger().Info("Server exiting")
}
