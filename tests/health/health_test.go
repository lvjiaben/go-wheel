package health

import (
	"context"
	"fmt"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/lvjiaben/go-wheel/pkg/container"
	"github.com/lvjiaben/go-wheel/pkg/health"
)

// 健康检查配置
type HealthCheckConfig struct {
	URL              string        // 要检查的端点URL
	Interval         time.Duration // 检查间隔
	Timeout          time.Duration // 请求超时时间
	ExpectedStatus   int           // 预期的HTTP状态码
	FailureThreshold int           // 连续失败次数阈值，超过此阈值则认为服务不健康
}

// 健康检查器
type HealthChecker struct {
	config    HealthCheckConfig
	isHealthy bool
	failures  int
	client    *http.Client
	mu        sync.Mutex
	stopChan  chan struct{}
}

// 创建一个新的健康检查器
func NewHealthChecker(config HealthCheckConfig) *HealthChecker {
	if config.Timeout == 0 {
		config.Timeout = 5 * time.Second
	}
	if config.FailureThreshold == 0 {
		config.FailureThreshold = 3
	}

	return &HealthChecker{
		config:    config,
		isHealthy: true, // 初始状态设置为健康
		failures:  0,
		client: &http.Client{
			Timeout: config.Timeout,
		},
		stopChan: make(chan struct{}),
	}
}

// 开始定期健康检查
func (hc *HealthChecker) Start() {
	ticker := time.NewTicker(hc.config.Interval)
	defer ticker.Stop()

	// 立即进行第一次检查
	hc.checkHealth()

	for {
		select {
		case <-ticker.C:
			hc.checkHealth()
		case <-hc.stopChan:
			return
		}
	}
}

// 停止健康检查
func (hc *HealthChecker) Stop() {
	close(hc.stopChan)
}

// 获取当前健康状态
func (hc *HealthChecker) IsHealthy() bool {
	hc.mu.Lock()
	defer hc.mu.Unlock()
	return hc.isHealthy
}

// 执行健康检查
func (hc *HealthChecker) checkHealth() {
	resp, err := hc.client.Get(hc.config.URL)

	hc.mu.Lock()
	defer hc.mu.Unlock()

	if err != nil || resp.StatusCode != hc.config.ExpectedStatus {
		hc.failures++
		fmt.Printf("健康检查失败: %v, 连续失败次数: %d\n", err, hc.failures)

		if hc.failures >= hc.config.FailureThreshold {
			hc.isHealthy = false
			fmt.Println("服务被标记为不健康")
		}
		return
	}

	// 如果请求成功，重置失败计数
	if hc.failures > 0 {
		fmt.Println("服务恢复正常")
	}
	hc.failures = 0
	hc.isHealthy = true

	if resp != nil && resp.Body != nil {
		resp.Body.Close()
	}
}

// 模拟服务层
type MonitorService struct {
	container   *container.Container
	healthCheck *health.Health
}

// 创建新的监控服务
func NewMonitorService(c *container.Container) *MonitorService {
	return &MonitorService{
		container:   c,
		healthCheck: health.New(),
	}
}

// 注册健康检查项
func (s *MonitorService) RegisterChecks() {
	// 数据库健康检查
	s.healthCheck.Register(&DatabaseChecker{
		container: s.container,
	})

	// Redis健康检查
	s.healthCheck.Register(&RedisChecker{
		container: s.container,
	})

	// 外部API健康检查
	s.healthCheck.Register(&APIChecker{
		url: "http://localhost:8089/health",
	})
}

// 获取健康状态
func (s *MonitorService) GetHealthStatus(ctx context.Context) map[string]error {
	return s.healthCheck.Check(ctx)
}

// 检查是否所有组件都健康
func (s *MonitorService) IsHealthy(ctx context.Context) bool {
	return s.healthCheck.IsHealthy(ctx)
}

// 数据库健康检查器
type DatabaseChecker struct {
	container *container.Container
}

func (c *DatabaseChecker) Name() string {
	return "database"
}

func (c *DatabaseChecker) Check(ctx context.Context) error {
	// 在实际应用中，我们会使用container.GetDB()
	// 为了测试，这里模拟数据库检查
	return nil
}

// Redis健康检查器
type RedisChecker struct {
	container *container.Container
}

func (c *RedisChecker) Name() string {
	return "redis"
}

func (c *RedisChecker) Check(ctx context.Context) error {
	// 为了测试，这里仅仅模拟Redis检查
	// 在实际应用中，我们会使用container.GetRedis().Ping()
	// 由于这是测试环境，我们假设Redis总是可用的
	return nil
}

// 外部API健康检查器
type APIChecker struct {
	url string
}

func (c *APIChecker) Name() string {
	return "api"
}

func (c *APIChecker) Check(ctx context.Context) error {
	req, err := http.NewRequestWithContext(ctx, "GET", c.url, nil)
	if err != nil {
		return err
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API返回非200状态码: %d", resp.StatusCode)
	}

	return nil
}

// 测试基本健康检查功能
func TestBasicHealthCheck(t *testing.T) {
	// 创建容器
	c := container.NewContainer()

	// 创建监控服务
	monitorService := NewMonitorService(c)
	monitorService.RegisterChecks()

	// 创建一个测试服务器，用于模拟API健康检查
	serverHealth := true
	server := http.NewServeMux()
	server.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		if serverHealth {
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(`{"status": "ok"}`))
		} else {
			w.WriteHeader(http.StatusServiceUnavailable)
			w.Write([]byte(`{"status": "error"}`))
		}
	})

	// 启动测试服务器
	go http.ListenAndServe(":8089", server)
	time.Sleep(100 * time.Millisecond) // 给服务器一点启动时间

	// 检查初始健康状态
	ctx := context.Background()
	if !monitorService.IsHealthy(ctx) {
		t.Error("所有服务应该处于健康状态")
	}

	// 模拟API服务不可用
	serverHealth = false

	// 等待健康检查更新
	time.Sleep(100 * time.Millisecond)

	// 再次检查健康状态
	status := monitorService.GetHealthStatus(ctx)
	if status["api"] == nil {
		t.Error("API服务应该被标记为不健康")
	} else {
		t.Logf("API检查正确返回错误: %v", status["api"])
	}

	// 恢复API服务
	serverHealth = true

	// 等待健康检查更新
	time.Sleep(100 * time.Millisecond)

	// 再次检查健康状态
	status = monitorService.GetHealthStatus(ctx)
	if status["api"] != nil {
		t.Errorf("API服务应该已恢复健康，但收到错误: %v", status["api"])
	}
}

// 测试健康检查超时功能
func TestHealthCheckWithTimeout(t *testing.T) {
	// 创建容器
	c := container.NewContainer()

	// 创建一个带延迟响应的测试服务器
	server := http.NewServeMux()
	server.HandleFunc("/slow-health", func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // 延迟2秒响应
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"status": "ok"}`))
	})

	// 启动测试服务器
	go http.ListenAndServe(":8090", server)
	time.Sleep(100 * time.Millisecond) // 给服务器一点启动时间

	// 创建带有慢速API检查的监控服务
	monitorService := NewMonitorService(c)
	monitorService.healthCheck.Register(&APIChecker{
		url: "http://localhost:8090/slow-health",
	})

	// 使用短超时进行检查
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	// 获取健康状态
	status := monitorService.GetHealthStatus(ctx)

	// 验证API检查是否超时
	if status["api"] == nil {
		t.Error("慢速API检查应该超时")
	} else {
		t.Logf("API检查正确返回超时错误: %v", status["api"])
	}
}

// 健康检查控制器示例
type HealthController struct {
	monitorService *MonitorService
}

func NewHealthController(service *MonitorService) *HealthController {
	return &HealthController{
		monitorService: service,
	}
}

func (c *HealthController) GetHealthStatus() map[string]interface{} {
	ctx := context.Background()
	status := c.monitorService.GetHealthStatus(ctx)

	result := make(map[string]interface{})
	allHealthy := true

	for name, err := range status {
		healthStatus := map[string]interface{}{
			"status":  "healthy",
			"message": "服务正常",
		}

		if err != nil {
			healthStatus["status"] = "unhealthy"
			healthStatus["message"] = err.Error()
			allHealthy = false
		}

		result[name] = healthStatus
	}

	result["overall"] = map[string]interface{}{
		"status":    map[bool]string{true: "healthy", false: "unhealthy"}[allHealthy],
		"timestamp": time.Now().Format(time.RFC3339),
	}

	return result
}

// 演示如何在实际应用中使用健康检查
func ExampleHealthCheckUsage() {
	// 创建容器
	c := container.NewContainer()

	// 在实际应用中，容器会被初始化并注入到服务中
	// 这里只是演示依赖注入的使用方式

	// 创建监控服务
	monitorService := NewMonitorService(c)
	monitorService.RegisterChecks()

	// 创建健康检查控制器
	healthController := NewHealthController(monitorService)

	// 获取健康状态
	status := healthController.GetHealthStatus()

	// 在实际应用中，这会返回到HTTP响应
	fmt.Printf("系统健康状态: %v\n", status["overall"].(map[string]interface{})["status"])
}
