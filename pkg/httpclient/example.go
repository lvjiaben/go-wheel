package httpclient

import (
	"context"
	"fmt"
	"time"
)

// Example 使用示例（仅供参考，不会被编译）
func Example() {
	// ==================== 1. 创建客户端 ====================

	// 基础客户端
	client := NewClient()

	// 带配置的客户端
	client = NewClient(
		WithBaseURL("https://api.example.com"),
		WithTimeout(30*time.Second),
		WithHeaders(map[string]string{
			"User-Agent": "MyApp/1.0",
		}),
		WithRetry(3, time.Second),
	)

	// 带代理的客户端
	client = NewClient(
		WithProxy("http://proxy.example.com:8080"),
		// 或 SOCKS5 代理
		// WithSocks5Proxy("127.0.0.1:1080"),
	)

	// 跳过 SSL 验证（不推荐生产环境）
	client = NewClient(
		WithInsecureSkipVerify(),
	)

	// 使用客户端证书
	client = NewClient(
		WithCertificates("/path/to/cert.pem", "/path/to/key.pem"),
	)

	// ==================== 2. GET 请求 ====================

	// 简单 GET 请求
	resp, err := client.Get("https://api.example.com/users").Send()
	if err != nil {
		panic(err)
	}
	defer resp.Close()

	// 获取响应字符串
	body, _ := resp.String()
	fmt.Println(body)

	// 解析 JSON 响应
	var users []User
	err = resp.JSON(&users)

	// 带查询参数的 GET 请求
	resp, err = client.Get("/users").
		SetQuery("page", "1").
		SetQuery("limit", "10").
		SetQueryParams(map[string]string{
			"sort": "created_at",
		}).
		Send()

	// ==================== 3. POST 请求 ====================

	// POST JSON 数据
	loginData := map[string]string{
		"username": "admin",
		"password": "123456",
	}
	resp, err = client.Post("/login").
		SetJSON(loginData).
		Send()

	// POST 表单数据
	resp, err = client.Post("/form").
		SetForm("name", "John").
		SetForm("email", "john@example.com").
		Send()

	// POST 表单数据（批量）
	resp, err = client.Post("/form").
		SetFormData(map[string]string{
			"name":  "John",
			"email": "john@example.com",
		}).
		Send()

	// ==================== 4. 设置 Headers ====================

	resp, err = client.Get("/api/data").
		SetHeader("Authorization", "Bearer token123").
		SetHeader("Content-Type", "application/json").
		SetHeaders(map[string]string{
			"X-Custom-Header": "value",
			"User-Agent":      "MyApp/1.0",
		}).
		Send()

	// ==================== 5. 认证 ====================

	// Bearer Token
	resp, err = client.Get("/api/protected").
		SetBearerToken("your-token-here").
		Send()

	// Basic Auth
	resp, err = client.Get("/api/protected").
		SetBasicAuth("username", "password").
		Send()

	// ==================== 6. 文件上传 ====================

	// 上传单个文件
	resp, err = client.UploadFile("/upload", "file", "/path/to/file.jpg")

	// 上传多个文件
	resp, err = client.UploadFiles("/upload", map[string]string{
		"file1": "/path/to/file1.jpg",
		"file2": "/path/to/file2.jpg",
	})

	// 上传文件并附带表单数据
	resp, err = client.UploadFileWithData(
		"/upload",
		"file",
		"/path/to/file.jpg",
		map[string]string{
			"title":       "My Photo",
			"description": "A beautiful photo",
		},
	)

	// 使用 Request 方式上传
	resp, err = client.Post("/upload").
		SetFile("avatar", "/path/to/avatar.jpg").
		SetForm("user_id", "123").
		Send()

	// ==================== 7. 文件下载 ====================

	// 下载文件
	err = client.DownloadFile("https://example.com/file.zip", "/path/to/save.zip")

	// 下载文件（带进度）
	err = client.DownloadFileWithProgress(
		"https://example.com/large-file.zip",
		"/path/to/save.zip",
		func(downloaded, total int64) {
			percent := float64(downloaded) / float64(total) * 100
			fmt.Printf("下载进度: %.2f%%\n", percent)
		},
	)

	// 断点续传下载
	err = client.DownloadFileResume("https://example.com/large-file.zip", "/path/to/save.zip")

	// 获取文件信息
	fileInfo, err := client.GetFileInfo("https://example.com/file.zip")
	if err == nil {
		fmt.Printf("文件大小: %d bytes\n", fileInfo.Size)
		fmt.Printf("文件类型: %s\n", fileInfo.ContentType)
	}

	// ==================== 8. 超时和重试 ====================

	resp, err = client.Get("/api/slow").
		SetTimeout(5 * time.Second).
		SetRetry(3, time.Second).
		Send()

	// ==================== 9. Context 支持 ====================

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	resp, err = client.Get("/api/data").
		WithContext(ctx).
		Send()

	// ==================== 10. 中间件 ====================

	// 添加日志中间件
	client = NewClient(
		WithMiddleware(NewLoggingMiddleware(nil)),
	)

	// 添加认证中间件
	client = NewClient(
		WithMiddleware(NewAuthMiddleware("your-token")),
	)

	// 添加自定义 Header 中间件
	client = NewClient(
		WithMiddleware(NewCustomHeaderMiddleware(map[string]string{
			"X-API-Key": "your-api-key",
		})),
	)

	// ==================== 11. 响应处理 ====================

	resp, err = client.Get("/api/data").Send()
	if err != nil {
		panic(err)
	}
	defer resp.Close()

	// 检查状态码
	if resp.IsSuccess() {
		fmt.Println("请求成功")
	}

	// 获取响应头
	contentType := resp.GetHeader("Content-Type")
	fmt.Println(contentType)

	// 获取 Cookie
	cookie := resp.GetCookie("session_id")
	if cookie != nil {
		fmt.Println(cookie.Value)
	}

	// 保存响应到文件
	err = resp.SaveToFile("/path/to/response.json")

	// ==================== 12. PUT/PATCH/DELETE 请求 ====================

	// PUT 请求
	updateData := map[string]interface{}{
		"name":  "New Name",
		"email": "new@example.com",
	}
	resp, err = client.Put("/users/123").
		SetJSON(updateData).
		Send()

	// PATCH 请求
	patchData := map[string]interface{}{
		"status": "active",
	}
	resp, err = client.Patch("/users/123").
		SetJSON(patchData).
		Send()

	// DELETE 请求
	resp, err = client.Delete("/users/123").
		SetBearerToken("token").
		Send()

	// ==================== 13. XML 请求和响应 ====================

	// 发送 XML 数据
	xmlData := struct {
		Name  string `xml:"name"`
		Email string `xml:"email"`
	}{
		Name:  "John",
		Email: "john@example.com",
	}
	resp, err = client.Post("/api/xml").
		SetXML(xmlData).
		Send()

	// 解析 XML 响应
	var result struct {
		Status  string `xml:"status"`
		Message string `xml:"message"`
	}
	err = resp.XML(&result)

	// ==================== 14. Cookie 管理 ====================

	resp, err = client.Get("/api/data").
		SetCookie("session_id", "abc123").
		SetCookie("user_id", "456").
		Send()

	// ==================== 15. 链式调用完整示例 ====================

	resp, err = client.Post("/api/users").
		SetHeader("Authorization", "Bearer token123").
		SetHeader("Content-Type", "application/json").
		SetQuery("notify", "true").
		SetJSON(map[string]interface{}{
			"name":  "John Doe",
			"email": "john@example.com",
			"age":   30,
		}).
		SetTimeout(10 * time.Second).
		SetRetry(3, time.Second).
		Send()

	if err != nil {
		panic(err)
	}
	defer resp.Close()

	if resp.IsSuccess() {
		var user User
		resp.JSON(&user)
		fmt.Printf("创建用户成功: %+v\n", user)
	}
}

// User 示例用户结构体
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

