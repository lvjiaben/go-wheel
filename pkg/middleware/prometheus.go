package middleware

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/pkg/monitor"
)

// PrometheusMiddleware Prometheus 监控中间件
func PrometheusMiddleware(metrics *monitor.PrometheusMetrics) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.FullPath()
		if path == "" {
			path = c.Request.URL.Path
		}

		// 获取请求大小
		reqSize := computeApproximateRequestSize(c.Request)

		// 处理请求
		c.Next()

		// 计算响应时间
		duration := time.Since(start).Seconds()

		// 获取响应状态码
		status := strconv.Itoa(c.Writer.Status())

		// 获取响应大小
		respSize := c.Writer.Size()
		if respSize < 0 {
			respSize = 0
		}

		// 记录指标
		metrics.RecordHTTPRequest(
			c.Request.Method,
			path,
			status,
			duration,
			reqSize,
			respSize,
		)
	}
}

// computeApproximateRequestSize 计算请求大小（近似值）
func computeApproximateRequestSize(r *http.Request) int {
	s := 0
	if r.URL != nil {
		s += len(r.URL.String())
	}

	s += len(r.Method)
	s += len(r.Proto)
	for name, values := range r.Header {
		s += len(name)
		for _, value := range values {
			s += len(value)
		}
	}
	s += len(r.Host)

	// 如果有 Content-Length，使用它
	if r.ContentLength > 0 {
		s += int(r.ContentLength)
	}

	return s
}

