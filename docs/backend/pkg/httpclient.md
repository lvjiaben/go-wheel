# HTTPClient HTTP 客户端

项目内置功能丰富的 HTTP 客户端，类似 PHP 的 Guzzle。

## 功能特点

- 支持 GET/POST/PUT/DELETE 等方法
- 自动重试机制
- 请求/响应中间件
- Cookie 管理
- 代理支持
- 调试模式

## 基本用法

### 创建客户端

```go
import "github.com/lvjiaben/go-wheel/pkg/httpclient"

// 基础创建
client := httpclient.NewClient()

// 带配置创建
client := httpclient.NewClient(
    httpclient.WithBaseURL("https://api.example.com"),
    httpclient.WithTimeout(30 * time.Second),
    httpclient.WithHeaders(map[string]string{
        "Authorization": "Bearer xxx",
    }),
    httpclient.WithRetry(3, time.Second),
)
```

### GET 请求

```go
// 简单 GET
resp, err := client.Get("/users")

// 带参数
resp, err := client.Get("/users", httpclient.WithQuery(map[string]string{
    "page": "1",
    "size": "10",
}))

// 获取响应
body := resp.Body()
statusCode := resp.StatusCode()
```

### POST 请求

```go
// JSON 请求
resp, err := client.PostJSON("/users", map[string]interface{}{
    "name":  "张三",
    "email": "zhangsan@example.com",
})

// 表单请求
resp, err := client.PostForm("/login", map[string]string{
    "username": "admin",
    "password": "123456",
})

// 原始数据
resp, err := client.Post("/data", []byte("raw data"))
```

### PUT/DELETE 请求

```go
// PUT
resp, err := client.PutJSON("/users/1", map[string]interface{}{
    "name": "李四",
})

// DELETE
resp, err := client.Delete("/users/1")
```

## 配置选项

```go
// 设置基础 URL
httpclient.WithBaseURL("https://api.example.com")

// 设置超时
httpclient.WithTimeout(30 * time.Second)

// 设置默认 Headers
httpclient.WithHeaders(map[string]string{
    "Content-Type": "application/json",
})

// 设置重试
httpclient.WithRetry(3, time.Second)

// 设置代理
httpclient.WithProxy("http://127.0.0.1:7890")

// 设置 SOCKS5 代理
httpclient.WithSocks5Proxy("127.0.0.1:1080", "user", "pass")

// 跳过 SSL 验证
httpclient.WithInsecureSkipVerify()

// 开启调试
httpclient.WithDebug(true)
```

## 中间件

```go
// 添加请求中间件
client.Use(func(req *http.Request, next httpclient.Handler) (*httpclient.Response, error) {
    // 请求前处理
    req.Header.Set("X-Request-ID", uuid.New().String())
    
    // 执行请求
    resp, err := next(req)
    
    // 响应后处理
    log.Printf("请求耗时: %v", resp.Duration())
    
    return resp, err
})
```

## 响应处理

```go
resp, err := client.Get("/users/1")

// 获取状态码
statusCode := resp.StatusCode()

// 获取响应体
body := resp.Body()

// 解析 JSON
var user User
resp.JSON(&user)

// 获取响应头
contentType := resp.Header("Content-Type")

// 获取请求耗时
duration := resp.Duration()
```

## 在控制器中使用

```go
type PaymentController struct {
    container *container.Container
}

func NewPaymentController(c *container.Container) *PaymentController {
    return &PaymentController{container: c}
}

func (ctrl *PaymentController) CreatePayment(ctx *gin.Context) {
    // 通过容器获取 HTTP 客户端
    httpClient := ctrl.container.GetHTTPClient()

    // 发送 GET 请求
    resp, err := httpClient.Get("https://api.payment.com/status").Send()
    if err != nil {
        // 处理错误
        return
    }

    // 发送 POST 请求
    resp, err = httpClient.Post("https://api.payment.com/create").
        SetJSON(map[string]interface{}{
            "amount": 100,
            "order_id": "123",
        }).Send()

    // 解析响应
    var result map[string]interface{}
    resp.JSON(&result)
}
```

## 在服务层中使用

```go
type PaymentService struct {
    container *container.Container
}

func NewPaymentService(c *container.Container) *PaymentService {
    return &PaymentService{container: c}
}

func (s *PaymentService) NotifyThirdParty(ctx context.Context, data interface{}) error {
    httpClient := s.container.GetHTTPClient()

    resp, err := httpClient.Post("https://api.third-party.com/notify").
        SetHeader("Authorization", "Bearer xxx").
        SetJSON(data).
        Send()

    if err != nil {
        return err
    }

    if resp.StatusCode() != 200 {
        return fmt.Errorf("请求失败: %d", resp.StatusCode())
    }

    return nil
}
```

## 文件上传

```go
// 上传文件
resp, err := client.PostMultipart("/upload", map[string]interface{}{
    "file": httpclient.File{
        Name:     "avatar.png",
        Filename: "avatar.png",
        Content:  fileBytes,
    },
    "name": "张三",
})
```

## 下载文件

```go
// 下载文件
resp, err := client.Get("/files/document.pdf")
if err == nil {
    ioutil.WriteFile("document.pdf", resp.Body(), 0644)
}
```

## 错误处理

```go
resp, err := client.Get("/users")
if err != nil {
    // 网络错误、超时等
    log.Printf("请求失败: %v", err)
    return
}

if resp.StatusCode() >= 400 {
    // HTTP 错误
    log.Printf("HTTP 错误: %d", resp.StatusCode())
}
```

## 最佳实践

1. **复用客户端** - 不要每次请求都创建新客户端
2. **设置超时** - 避免请求无限等待
3. **错误重试** - 对临时错误进行重试
4. **日志记录** - 开启调试模式排查问题

