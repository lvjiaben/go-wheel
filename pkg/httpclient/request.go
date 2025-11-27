package httpclient

import (
	"bytes"
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Request HTTP 请求构建器
type Request struct {
	client      *Client
	method      string
	url         string
	headers     map[string]string
	cookies     map[string]string
	queryParams map[string]string
	formData    map[string]string
	body        interface{}
	bodyReader  io.Reader
	files       map[string]string
	timeout     time.Duration
	retryCount  int
	retryDelay  time.Duration
	ctx         context.Context
}

// WithContext 设置上下文
func (r *Request) WithContext(ctx context.Context) *Request {
	r.ctx = ctx
	return r
}

// SetHeader 设置单个 Header
func (r *Request) SetHeader(key, value string) *Request {
	r.headers[key] = value
	return r
}

// SetHeaders 批量设置 Headers
func (r *Request) SetHeaders(headers map[string]string) *Request {
	for k, v := range headers {
		r.headers[k] = v
	}
	return r
}

// SetCookie 设置 Cookie
func (r *Request) SetCookie(key, value string) *Request {
	r.cookies[key] = value
	return r
}

// SetQuery 设置单个查询参数
func (r *Request) SetQuery(key, value string) *Request {
	r.queryParams[key] = value
	return r
}

// SetQueryParams 批量设置查询参数
func (r *Request) SetQueryParams(params map[string]string) *Request {
	for k, v := range params {
		r.queryParams[k] = v
	}
	return r
}

// SetForm 设置单个表单字段
func (r *Request) SetForm(key, value string) *Request {
	r.formData[key] = value
	return r
}

// SetFormData 批量设置表单数据
func (r *Request) SetFormData(data map[string]string) *Request {
	for k, v := range data {
		r.formData[k] = v
	}
	return r
}

// SetJSON 设置 JSON 请求体
func (r *Request) SetJSON(data interface{}) *Request {
	r.body = data
	r.headers["Content-Type"] = "application/json"
	return r
}

// SetXML 设置 XML 请求体
func (r *Request) SetXML(data interface{}) *Request {
	r.body = data
	r.headers["Content-Type"] = "application/xml"
	return r
}

// SetBody 设置原始请求体
func (r *Request) SetBody(body interface{}) *Request {
	r.body = body
	return r
}

// SetBodyReader 设置请求体 Reader
func (r *Request) SetBodyReader(reader io.Reader) *Request {
	r.bodyReader = reader
	return r
}

// SetFile 设置单个文件上传
func (r *Request) SetFile(fieldName, filePath string) *Request {
	if r.files == nil {
		r.files = make(map[string]string)
	}
	r.files[fieldName] = filePath
	return r
}

// SetFiles 批量设置文件上传
func (r *Request) SetFiles(files map[string]string) *Request {
	if r.files == nil {
		r.files = make(map[string]string)
	}
	for k, v := range files {
		r.files[k] = v
	}
	return r
}

// SetTimeout 设置超时时间
func (r *Request) SetTimeout(timeout time.Duration) *Request {
	r.timeout = timeout
	return r
}

// SetRetry 设置重试配置
func (r *Request) SetRetry(count int, delay time.Duration) *Request {
	r.retryCount = count
	r.retryDelay = delay
	return r
}

// SetBasicAuth 设置 Basic 认证
func (r *Request) SetBasicAuth(username, password string) *Request {
	auth := username + ":" + password
	encoded := "Basic " + base64Encode([]byte(auth))
	r.headers["Authorization"] = encoded
	return r
}

// SetBearerToken 设置 Bearer Token
func (r *Request) SetBearerToken(token string) *Request {
	r.headers["Authorization"] = "Bearer " + token
	return r
}

// Send 发送请求
func (r *Request) Send() (*Response, error) {
	// 构建请求体
	bodyReader, contentType, err := r.buildBody()
	if err != nil {
		return nil, fmt.Errorf("构建请求体失败: %v", err)
	}

	// 构建完整 URL（包含查询参数）
	fullURL := r.buildURL()

	// 创建 HTTP 请求
	req, err := http.NewRequestWithContext(r.ctx, r.method, fullURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %v", err)
	}

	// 设置 Headers
	for k, v := range r.headers {
		req.Header.Set(k, v)
	}

	// 如果有 Content-Type，设置它
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	// 设置 Cookies
	for k, v := range r.cookies {
		req.AddCookie(&http.Cookie{Name: k, Value: v})
	}

	// 执行中间件（请求前）
	for _, middleware := range r.client.middlewares {
		if err := middleware.BeforeRequest(req); err != nil {
			return nil, fmt.Errorf("中间件执行失败: %v", err)
		}
	}

	// 发送请求（带重试）
	var resp *http.Response
	var lastErr error

	for i := 0; i <= r.retryCount; i++ {
		if i > 0 {
			time.Sleep(r.retryDelay)
			if r.client.logger != nil {
				r.client.logger.Info("重试请求", "attempt", i, "url", fullURL)
			}
		}

		// 发送请求
		resp, lastErr = r.client.httpClient.Do(req)
		if lastErr == nil {
			break
		}
	}

	if lastErr != nil {
		return nil, fmt.Errorf("请求失败: %v", lastErr)
	}

	// 创建响应对象
	response := &Response{
		StatusCode: resp.StatusCode,
		Status:     resp.Status,
		Headers:    resp.Header,
		Cookies:    resp.Cookies(),
		body:       resp.Body,
		request:    req,
	}

	// 执行中间件（响应后）
	for _, middleware := range r.client.middlewares {
		if err := middleware.AfterResponse(response); err != nil {
			resp.Body.Close()
			return nil, fmt.Errorf("中间件执行失败: %v", err)
		}
	}

	return response, nil
}

// buildBody 构建请求体
func (r *Request) buildBody() (io.Reader, string, error) {
	// 如果有自定义 bodyReader，直接使用
	if r.bodyReader != nil {
		return r.bodyReader, "", nil
	}

	// 如果有文件上传，构建 multipart
	if len(r.files) > 0 {
		return r.buildMultipartBody()
	}

	// 如果有表单数据
	if len(r.formData) > 0 {
		data := url.Values{}
		for k, v := range r.formData {
			data.Set(k, v)
		}
		return strings.NewReader(data.Encode()), "application/x-www-form-urlencoded", nil
	}

	// 如果有 JSON 数据
	if r.body != nil {
		contentType := r.headers["Content-Type"]
		if contentType == "application/json" {
			jsonData, err := json.Marshal(r.body)
			if err != nil {
				return nil, "", err
			}
			return bytes.NewReader(jsonData), "application/json", nil
		} else if contentType == "application/xml" {
			xmlData, err := xml.Marshal(r.body)
			if err != nil {
				return nil, "", err
			}
			return bytes.NewReader(xmlData), "application/xml", nil
		} else {
			// 默认按字符串处理
			return strings.NewReader(fmt.Sprintf("%v", r.body)), "", nil
		}
	}

	return nil, "", nil
}

// buildMultipartBody 构建 multipart 请求体
func (r *Request) buildMultipartBody() (io.Reader, string, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// 添加表单字段
	for k, v := range r.formData {
		if err := writer.WriteField(k, v); err != nil {
			return nil, "", err
		}
	}

	// 添加文件
	for fieldName, filePath := range r.files {
		file, err := os.Open(filePath)
		if err != nil {
			return nil, "", fmt.Errorf("打开文件失败 %s: %v", filePath, err)
		}
		defer file.Close()

		part, err := writer.CreateFormFile(fieldName, filepath.Base(filePath))
		if err != nil {
			return nil, "", err
		}

		if _, err := io.Copy(part, file); err != nil {
			return nil, "", err
		}
	}

	if err := writer.Close(); err != nil {
		return nil, "", err
	}

	return body, writer.FormDataContentType(), nil
}

// buildURL 构建完整 URL（包含查询参数）
func (r *Request) buildURL() string {
	if len(r.queryParams) == 0 {
		return r.url
	}

	u, err := url.Parse(r.url)
	if err != nil {
		return r.url
	}

	q := u.Query()
	for k, v := range r.queryParams {
		q.Set(k, v)
	}
	u.RawQuery = q.Encode()

	return u.String()
}

