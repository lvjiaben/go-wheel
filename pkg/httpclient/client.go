package httpclient

import (
	"context"
	"crypto/tls"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"golang.org/x/net/proxy"
)

// Client HTTP 客户端（类似 Guzzle）
type Client struct {
	httpClient  *http.Client
	baseURL     string
	headers     map[string]string
	cookies     map[string]string
	timeout     time.Duration
	retryCount  int
	retryDelay  time.Duration
	middlewares []Middleware
	debug       bool
	logger      Logger
}

// Logger 日志接口
type Logger interface {
	Debug(msg string, fields ...interface{})
	Info(msg string, fields ...interface{})
	Error(msg string, fields ...interface{})
}

// ClientOption 客户端配置选项
type ClientOption func(*Client)

// NewClient 创建 HTTP 客户端
func NewClient(options ...ClientOption) *Client {
	jar, _ := cookiejar.New(nil)

	client := &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Jar:     jar,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
				DialContext: (&net.Dialer{
					Timeout:   30 * time.Second,
					KeepAlive: 30 * time.Second,
				}).DialContext,
			},
		},
		headers:     make(map[string]string),
		cookies:     make(map[string]string),
		timeout:     30 * time.Second,
		retryCount:  0,
		retryDelay:  time.Second,
		middlewares: []Middleware{},
		debug:       false,
	}

	// 应用配置选项
	for _, option := range options {
		option(client)
	}

	return client
}

// WithBaseURL 设置基础 URL
func WithBaseURL(baseURL string) ClientOption {
	return func(c *Client) {
		c.baseURL = baseURL
	}
}

// WithTimeout 设置超时时间
func WithTimeout(timeout time.Duration) ClientOption {
	return func(c *Client) {
		c.timeout = timeout
		c.httpClient.Timeout = timeout
	}
}

// WithHeaders 设置默认 Headers
func WithHeaders(headers map[string]string) ClientOption {
	return func(c *Client) {
		for k, v := range headers {
			c.headers[k] = v
		}
	}
}

// WithRetry 设置重试配置
func WithRetry(count int, delay time.Duration) ClientOption {
	return func(c *Client) {
		c.retryCount = count
		c.retryDelay = delay
	}
}

// WithProxy 设置 HTTP 代理
func WithProxy(proxyURL string) ClientOption {
	return func(c *Client) {
		if proxyURL == "" {
			return
		}
		pURL, err := url.Parse(proxyURL)
		if err != nil {
			return
		}
		transport := c.httpClient.Transport.(*http.Transport)
		transport.Proxy = http.ProxyURL(pURL)
	}
}

// WithSocks5Proxy 设置 SOCKS5 代理
func WithSocks5Proxy(proxyAddr string) ClientOption {
	return func(c *Client) {
		if proxyAddr == "" {
			return
		}
		dialer, err := proxy.SOCKS5("tcp", proxyAddr, nil, proxy.Direct)
		if err != nil {
			return
		}
		transport := c.httpClient.Transport.(*http.Transport)
		transport.DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.Dial(network, addr)
		}
	}
}

// WithTLSConfig 设置 TLS 配置
func WithTLSConfig(config *tls.Config) ClientOption {
	return func(c *Client) {
		transport := c.httpClient.Transport.(*http.Transport)
		transport.TLSClientConfig = config
	}
}

// WithInsecureSkipVerify 跳过 SSL 证书验证（不推荐生产环境）
func WithInsecureSkipVerify() ClientOption {
	return func(c *Client) {
		transport := c.httpClient.Transport.(*http.Transport)
		if transport.TLSClientConfig == nil {
			transport.TLSClientConfig = &tls.Config{}
		}
		transport.TLSClientConfig.InsecureSkipVerify = true
	}
}

// WithCertificates 设置客户端证书
func WithCertificates(certFile, keyFile string) ClientOption {
	return func(c *Client) {
		cert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return
		}
		transport := c.httpClient.Transport.(*http.Transport)
		if transport.TLSClientConfig == nil {
			transport.TLSClientConfig = &tls.Config{}
		}
		transport.TLSClientConfig.Certificates = []tls.Certificate{cert}
	}
}

// WithTransport 自定义 Transport
func WithTransport(transport *http.Transport) ClientOption {
	return func(c *Client) {
		c.httpClient.Transport = transport
	}
}

// WithCookieJar 设置 Cookie Jar
func WithCookieJar(jar http.CookieJar) ClientOption {
	return func(c *Client) {
		c.httpClient.Jar = jar
	}
}

// WithDebug 启用调试模式
func WithDebug(debug bool) ClientOption {
	return func(c *Client) {
		c.debug = debug
	}
}

// WithLogger 设置日志记录器
func WithLogger(logger Logger) ClientOption {
	return func(c *Client) {
		c.logger = logger
	}
}

// WithMiddleware 添加中间件
func WithMiddleware(middleware Middleware) ClientOption {
	return func(c *Client) {
		c.middlewares = append(c.middlewares, middleware)
	}
}

// SetHeader 设置单个 Header
func (c *Client) SetHeader(key, value string) *Client {
	c.headers[key] = value
	return c
}

// SetHeaders 批量设置 Headers
func (c *Client) SetHeaders(headers map[string]string) *Client {
	for k, v := range headers {
		c.headers[k] = v
	}
	return c
}

// SetCookie 设置 Cookie
func (c *Client) SetCookie(key, value string) *Client {
	c.cookies[key] = value
	return c
}

// SetTimeout 设置超时时间
func (c *Client) SetTimeout(timeout time.Duration) *Client {
	c.timeout = timeout
	c.httpClient.Timeout = timeout
	return c
}

// SetRetry 设置重试配置
func (c *Client) SetRetry(count int, delay time.Duration) *Client {
	c.retryCount = count
	c.retryDelay = delay
	return c
}

// Get 发送 GET 请求
func (c *Client) Get(url string) *Request {
	return c.NewRequest("GET", url)
}

// Post 发送 POST 请求
func (c *Client) Post(url string) *Request {
	return c.NewRequest("POST", url)
}

// Put 发送 PUT 请求
func (c *Client) Put(url string) *Request {
	return c.NewRequest("PUT", url)
}

// Patch 发送 PATCH 请求
func (c *Client) Patch(url string) *Request {
	return c.NewRequest("PATCH", url)
}

// Delete 发送 DELETE 请求
func (c *Client) Delete(url string) *Request {
	return c.NewRequest("DELETE", url)
}

// Head 发送 HEAD 请求
func (c *Client) Head(url string) *Request {
	return c.NewRequest("HEAD", url)
}

// Options 发送 OPTIONS 请求
func (c *Client) Options(url string) *Request {
	return c.NewRequest("OPTIONS", url)
}

// NewRequest 创建新请求
func (c *Client) NewRequest(method, url string) *Request {
	// 如果是相对路径，拼接 baseURL
	if c.baseURL != "" && !isAbsoluteURL(url) {
		url = c.baseURL + url
	}

	return &Request{
		client:      c,
		method:      method,
		url:         url,
		headers:     copyMap(c.headers),
		cookies:     copyMap(c.cookies),
		queryParams: make(map[string]string),
		formData:    make(map[string]string),
		timeout:     c.timeout,
		retryCount:  c.retryCount,
		retryDelay:  c.retryDelay,
		ctx:         context.Background(),
	}
}

// isAbsoluteURL 判断是否为绝对 URL
func isAbsoluteURL(urlStr string) bool {
	u, err := url.Parse(urlStr)
	if err != nil {
		return false
	}
	return u.Scheme != "" && u.Host != ""
}

// copyMap 复制 map
func copyMap(m map[string]string) map[string]string {
	result := make(map[string]string, len(m))
	for k, v := range m {
		result[k] = v
	}
	return result
}

