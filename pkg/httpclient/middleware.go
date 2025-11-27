package httpclient

import (
	"fmt"
	"net/http"
	"time"
)

// Middleware 中间件接口
type Middleware interface {
	BeforeRequest(req *http.Request) error
	AfterResponse(resp *Response) error
}

// LoggingMiddleware 日志中间件
type LoggingMiddleware struct {
	logger Logger
}

// NewLoggingMiddleware 创建日志中间件
func NewLoggingMiddleware(logger Logger) *LoggingMiddleware {
	return &LoggingMiddleware{logger: logger}
}

// BeforeRequest 请求前记录日志
func (m *LoggingMiddleware) BeforeRequest(req *http.Request) error {
	if m.logger != nil {
		m.logger.Info("发送HTTP请求",
			"method", req.Method,
			"url", req.URL.String(),
			"headers", req.Header,
		)
	}
	return nil
}

// AfterResponse 响应后记录日志
func (m *LoggingMiddleware) AfterResponse(resp *Response) error {
	if m.logger != nil {
		m.logger.Info("收到HTTP响应",
			"status", resp.StatusCode,
			"url", resp.request.URL.String(),
		)
	}
	return nil
}

// RetryMiddleware 重试中间件
type RetryMiddleware struct {
	maxRetries int
	retryDelay time.Duration
}

// NewRetryMiddleware 创建重试中间件
func NewRetryMiddleware(maxRetries int, retryDelay time.Duration) *RetryMiddleware {
	return &RetryMiddleware{
		maxRetries: maxRetries,
		retryDelay: retryDelay,
	}
}

// BeforeRequest 请求前处理
func (m *RetryMiddleware) BeforeRequest(req *http.Request) error {
	return nil
}

// AfterResponse 响应后处理
func (m *RetryMiddleware) AfterResponse(resp *Response) error {
	return nil
}

// AuthMiddleware 认证中间件
type AuthMiddleware struct {
	token string
}

// NewAuthMiddleware 创建认证中间件
func NewAuthMiddleware(token string) *AuthMiddleware {
	return &AuthMiddleware{token: token}
}

// BeforeRequest 请求前添加认证信息
func (m *AuthMiddleware) BeforeRequest(req *http.Request) error {
	if m.token != "" {
		req.Header.Set("Authorization", "Bearer "+m.token)
	}
	return nil
}

// AfterResponse 响应后处理
func (m *AuthMiddleware) AfterResponse(resp *Response) error {
	return nil
}

// RateLimitMiddleware 限流中间件
type RateLimitMiddleware struct {
	limiter chan struct{}
}

// NewRateLimitMiddleware 创建限流中间件
func NewRateLimitMiddleware(requestsPerSecond int) *RateLimitMiddleware {
	limiter := make(chan struct{}, requestsPerSecond)
	
	// 定时释放令牌
	go func() {
		ticker := time.NewTicker(time.Second / time.Duration(requestsPerSecond))
		defer ticker.Stop()
		
		for range ticker.C {
			select {
			case limiter <- struct{}{}:
			default:
			}
		}
	}()
	
	return &RateLimitMiddleware{limiter: limiter}
}

// BeforeRequest 请求前限流
func (m *RateLimitMiddleware) BeforeRequest(req *http.Request) error {
	<-m.limiter // 等待令牌
	return nil
}

// AfterResponse 响应后处理
func (m *RateLimitMiddleware) AfterResponse(resp *Response) error {
	return nil
}

// ErrorHandlerMiddleware 错误处理中间件
type ErrorHandlerMiddleware struct {
	handler func(*Response) error
}

// NewErrorHandlerMiddleware 创建错误处理中间件
func NewErrorHandlerMiddleware(handler func(*Response) error) *ErrorHandlerMiddleware {
	return &ErrorHandlerMiddleware{handler: handler}
}

// BeforeRequest 请求前处理
func (m *ErrorHandlerMiddleware) BeforeRequest(req *http.Request) error {
	return nil
}

// AfterResponse 响应后处理错误
func (m *ErrorHandlerMiddleware) AfterResponse(resp *Response) error {
	if !resp.IsSuccess() && m.handler != nil {
		return m.handler(resp)
	}
	return nil
}

// TimeoutMiddleware 超时中间件
type TimeoutMiddleware struct {
	timeout time.Duration
}

// NewTimeoutMiddleware 创建超时中间件
func NewTimeoutMiddleware(timeout time.Duration) *TimeoutMiddleware {
	return &TimeoutMiddleware{timeout: timeout}
}

// BeforeRequest 请求前设置超时
func (m *TimeoutMiddleware) BeforeRequest(req *http.Request) error {
	// 超时在 Client 层面处理
	return nil
}

// AfterResponse 响应后处理
func (m *TimeoutMiddleware) AfterResponse(resp *Response) error {
	return nil
}

// CustomHeaderMiddleware 自定义 Header 中间件
type CustomHeaderMiddleware struct {
	headers map[string]string
}

// NewCustomHeaderMiddleware 创建自定义 Header 中间件
func NewCustomHeaderMiddleware(headers map[string]string) *CustomHeaderMiddleware {
	return &CustomHeaderMiddleware{headers: headers}
}

// BeforeRequest 请求前添加自定义 Header
func (m *CustomHeaderMiddleware) BeforeRequest(req *http.Request) error {
	for k, v := range m.headers {
		req.Header.Set(k, v)
	}
	return nil
}

// AfterResponse 响应后处理
func (m *CustomHeaderMiddleware) AfterResponse(resp *Response) error {
	return nil
}

// DebugMiddleware 调试中间件
type DebugMiddleware struct {
	logger Logger
}

// NewDebugMiddleware 创建调试中间件
func NewDebugMiddleware(logger Logger) *DebugMiddleware {
	return &DebugMiddleware{logger: logger}
}

// BeforeRequest 请求前打印调试信息
func (m *DebugMiddleware) BeforeRequest(req *http.Request) error {
	if m.logger != nil {
		m.logger.Debug("=== HTTP Request ===",
			"method", req.Method,
			"url", req.URL.String(),
			"headers", req.Header,
		)
	}
	return nil
}

// AfterResponse 响应后打印调试信息
func (m *DebugMiddleware) AfterResponse(resp *Response) error {
	if m.logger != nil {
		body, _ := resp.String()
		m.logger.Debug("=== HTTP Response ===",
			"status", resp.StatusCode,
			"headers", resp.Headers,
			"body", body,
		)
	}
	return nil
}

// StatusCodeMiddleware 状态码检查中间件
type StatusCodeMiddleware struct {
	allowedCodes []int
}

// NewStatusCodeMiddleware 创建状态码检查中间件
func NewStatusCodeMiddleware(allowedCodes []int) *StatusCodeMiddleware {
	return &StatusCodeMiddleware{allowedCodes: allowedCodes}
}

// BeforeRequest 请求前处理
func (m *StatusCodeMiddleware) BeforeRequest(req *http.Request) error {
	return nil
}

// AfterResponse 响应后检查状态码
func (m *StatusCodeMiddleware) AfterResponse(resp *Response) error {
	if len(m.allowedCodes) == 0 {
		return nil
	}

	for _, code := range m.allowedCodes {
		if resp.StatusCode == code {
			return nil
		}
	}

	return fmt.Errorf("不允许的状态码: %d", resp.StatusCode)
}

