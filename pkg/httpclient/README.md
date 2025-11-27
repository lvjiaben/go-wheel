# HTTP Client - 类似 PHP Guzzle 的 Go HTTP 客户端

一个功能强大、易用的 Go HTTP 客户端库，类似于 PHP 的 Guzzle，支持链式调用、中间件、文件上传下载、代理、证书等完整功能。

---

## ✨ 特性

- ✅ **链式调用** - 优雅的 API 设计
- ✅ **GET/POST/PUT/PATCH/DELETE** - 支持所有 HTTP 方法
- ✅ **JSON/XML/Form** - 多种数据格式支持
- ✅ **文件上传/下载** - 支持单文件、多文件、断点续传
- ✅ **中间件系统** - 请求/响应拦截器
- ✅ **代理支持** - HTTP/SOCKS5 代理
- ✅ **SSL/TLS** - 客户端证书、跳过验证
- ✅ **超时和重试** - 可配置的超时和重试机制
- ✅ **Cookie 管理** - 自动 Cookie 管理
- ✅ **Context 支持** - 支持 context 取消
- ✅ **日志记录** - 可选的请求/响应日志

---

## 📦 安装

```bash
go get golang.org/x/net/proxy
```

---

## 🚀 快速开始

### 1. 从容器获取客户端（推荐）

```go
// 从容器获取已配置的 HTTP 客户端
client := container.GetHTTPClient()

// 发送 GET 请求
resp, err := client.Get("https://api.example.com/users").Send()
if err != nil {
    panic(err)
}
defer resp.Close()

// 解析 JSON 响应
var users []User
resp.JSON(&users)
```

### 2. 创建独立客户端

```go
import "github.com/lvjiaben/go-wheel/pkg/httpclient"

// 创建基础客户端
client := httpclient.NewClient()

// 创建带配置的客户端
client := httpclient.NewClient(
    httpclient.WithBaseURL("https://api.example.com"),
    httpclient.WithTimeout(30*time.Second),
    httpclient.WithHeaders(map[string]string{
        "User-Agent": "MyApp/1.0",
    }),
)
```

---

## 📖 使用示例

### GET 请求

```go
// 简单 GET 请求
resp, err := client.Get("/users").Send()

// 带查询参数
resp, err := client.Get("/users").
    SetQuery("page", "1").
    SetQuery("limit", "10").
    Send()

// 批量设置查询参数
resp, err := client.Get("/users").
    SetQueryParams(map[string]string{
        "page":  "1",
        "limit": "10",
        "sort":  "created_at",
    }).
    Send()
```

### POST 请求

```go
// POST JSON 数据
data := map[string]interface{}{
    "username": "admin",
    "password": "123456",
}
resp, err := client.Post("/login").
    SetJSON(data).
    Send()

// POST 表单数据
resp, err := client.Post("/form").
    SetForm("name", "John").
    SetForm("email", "john@example.com").
    Send()

// 批量设置表单数据
resp, err := client.Post("/form").
    SetFormData(map[string]string{
        "name":  "John",
        "email": "john@example.com",
    }).
    Send()
```

### 设置 Headers

```go
resp, err := client.Get("/api/data").
    SetHeader("Authorization", "Bearer token123").
    SetHeader("Content-Type", "application/json").
    Send()

// 批量设置
resp, err := client.Get("/api/data").
    SetHeaders(map[string]string{
        "Authorization": "Bearer token123",
        "X-Custom-Header": "value",
    }).
    Send()
```

### 认证

```go
// Bearer Token
resp, err := client.Get("/api/protected").
    SetBearerToken("your-token-here").
    Send()

// Basic Auth
resp, err := client.Get("/api/protected").
    SetBasicAuth("username", "password").
    Send()
```

### 文件上传

```go
// 上传单个文件
resp, err := client.UploadFile("/upload", "file", "/path/to/file.jpg")

// 上传多个文件
resp, err := client.UploadFiles("/upload", map[string]string{
    "file1": "/path/to/file1.jpg",
    "file2": "/path/to/file2.jpg",
})

// 上传文件并附带表单数据
resp, err := client.UploadFileWithData(
    "/upload",
    "file",
    "/path/to/file.jpg",
    map[string]string{
        "title": "My Photo",
        "description": "A beautiful photo",
    },
)

// 使用 Request 方式
resp, err := client.Post("/upload").
    SetFile("avatar", "/path/to/avatar.jpg").
    SetForm("user_id", "123").
    Send()
```

### 文件下载

```go
// 下载文件
err := client.DownloadFile("https://example.com/file.zip", "/path/to/save.zip")

// 下载文件（带进度）
err := client.DownloadFileWithProgress(
    "https://example.com/large-file.zip",
    "/path/to/save.zip",
    func(downloaded, total int64) {
        percent := float64(downloaded) / float64(total) * 100
        fmt.Printf("下载进度: %.2f%%\n", percent)
    },
)

// 断点续传下载
err := client.DownloadFileResume("https://example.com/large-file.zip", "/path/to/save.zip")

// 获取文件信息
fileInfo, err := client.GetFileInfo("https://example.com/file.zip")
fmt.Printf("文件大小: %d bytes\n", fileInfo.Size)
```

### 响应处理

```go
resp, err := client.Get("/api/data").Send()
if err != nil {
    panic(err)
}
defer resp.Close()

// 获取响应字符串
body, _ := resp.String()

// 解析 JSON
var result map[string]interface{}
resp.JSON(&result)

// 解析 XML
var xmlResult struct {
    Status string `xml:"status"`
}
resp.XML(&xmlResult)

// 检查状态码
if resp.IsSuccess() {
    fmt.Println("请求成功")
}

// 获取响应头
contentType := resp.GetHeader("Content-Type")

// 获取 Cookie
cookie := resp.GetCookie("session_id")

// 保存到文件
resp.SaveToFile("/path/to/response.json")
```

### 超时和重试

```go
resp, err := client.Get("/api/slow").
    SetTimeout(5 * time.Second).
    SetRetry(3, time.Second).
    Send()
```

### Context 支持

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

resp, err := client.Get("/api/data").
    WithContext(ctx).
    Send()
```

---

## 🔧 高级配置

### 代理

```go
// HTTP 代理
client := httpclient.NewClient(
    httpclient.WithProxy("http://proxy.example.com:8080"),
)

// SOCKS5 代理
client := httpclient.NewClient(
    httpclient.WithSocks5Proxy("127.0.0.1:1080"),
)
```

### SSL/TLS

```go
// 跳过 SSL 验证（不推荐生产环境）
client := httpclient.NewClient(
    httpclient.WithInsecureSkipVerify(),
)

// 使用客户端证书
client := httpclient.NewClient(
    httpclient.WithCertificates("/path/to/cert.pem", "/path/to/key.pem"),
)

// 自定义 TLS 配置
tlsConfig := &tls.Config{
    MinVersion: tls.VersionTLS12,
}
client := httpclient.NewClient(
    httpclient.WithTLSConfig(tlsConfig),
)
```

### 中间件

```go
// 日志中间件
client := httpclient.NewClient(
    httpclient.WithMiddleware(httpclient.NewLoggingMiddleware(logger)),
)

// 认证中间件
client := httpclient.NewClient(
    httpclient.WithMiddleware(httpclient.NewAuthMiddleware("your-token")),
)

// 自定义 Header 中间件
client := httpclient.NewClient(
    httpclient.WithMiddleware(httpclient.NewCustomHeaderMiddleware(map[string]string{
        "X-API-Key": "your-api-key",
    })),
)

// 调试中间件
client := httpclient.NewClient(
    httpclient.WithMiddleware(httpclient.NewDebugMiddleware(logger)),
)
```

---

## 📝 完整示例

```go
package main

import (
    "fmt"
    "time"
    "github.com/lvjiaben/go-wheel/pkg/httpclient"
)

func main() {
    // 创建客户端
    client := httpclient.NewClient(
        httpclient.WithBaseURL("https://api.example.com"),
        httpclient.WithTimeout(30*time.Second),
        httpclient.WithRetry(3, time.Second),
        httpclient.WithHeaders(map[string]string{
            "User-Agent": "MyApp/1.0",
        }),
    )

    // 发送请求
    resp, err := client.Post("/api/users").
        SetHeader("Authorization", "Bearer token123").
        SetQuery("notify", "true").
        SetJSON(map[string]interface{}{
            "name":  "John Doe",
            "email": "john@example.com",
            "age":   30,
        }).
        Send()

    if err != nil {
        panic(err)
    }
    defer resp.Close()

    // 处理响应
    if resp.IsSuccess() {
        var user User
        resp.JSON(&user)
        fmt.Printf("创建用户成功: %+v\n", user)
    } else {
        body, _ := resp.String()
        fmt.Printf("请求失败: %s\n", body)
    }
}

type User struct {
    ID    int    `json:"id"`
    Name  string `json:"name"`
    Email string `json:"email"`
}
```

---

## 🎯 与 PHP Guzzle 对比

| 功能 | Guzzle (PHP) | httpclient (Go) |
|------|--------------|-----------------|
| 链式调用 | ✅ | ✅ |
| 中间件 | ✅ | ✅ |
| 文件上传 | ✅ | ✅ |
| 文件下载 | ✅ | ✅ |
| 代理支持 | ✅ | ✅ |
| SSL/TLS | ✅ | ✅ |
| 超时重试 | ✅ | ✅ |
| Cookie 管理 | ✅ | ✅ |
| 异步请求 | ✅ | ✅ (通过 goroutine) |

---

## 📚 更多示例

查看 `example.go` 文件获取更多使用示例。

---

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

---

## 📄 许可证

MIT License

