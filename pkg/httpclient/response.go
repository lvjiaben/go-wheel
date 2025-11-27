package httpclient

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"os"
)

// Response HTTP 响应
type Response struct {
	StatusCode int
	Status     string
	Headers    http.Header
	Cookies    []*http.Cookie
	body       io.ReadCloser
	bodyBytes  []byte
	request    *http.Request
}

// Body 获取响应体（字节数组）
func (r *Response) Body() ([]byte, error) {
	if r.bodyBytes != nil {
		return r.bodyBytes, nil
	}

	defer r.body.Close()
	data, err := io.ReadAll(r.body)
	if err != nil {
		return nil, err
	}

	r.bodyBytes = data
	return data, nil
}

// String 获取响应体（字符串）
func (r *Response) String() (string, error) {
	data, err := r.Body()
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// JSON 解析 JSON 响应
func (r *Response) JSON(v interface{}) error {
	data, err := r.Body()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// XML 解析 XML 响应
func (r *Response) XML(v interface{}) error {
	data, err := r.Body()
	if err != nil {
		return err
	}
	return xml.Unmarshal(data, v)
}

// SaveToFile 保存响应到文件
func (r *Response) SaveToFile(filePath string) error {
	data, err := r.Body()
	if err != nil {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	return nil
}

// SaveToFileStream 流式保存到文件（适合大文件）
func (r *Response) SaveToFileStream(filePath string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer file.Close()

	defer r.body.Close()
	_, err = io.Copy(file, r.body)
	if err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	return nil
}

// GetHeader 获取响应头
func (r *Response) GetHeader(key string) string {
	return r.Headers.Get(key)
}

// GetCookie 获取 Cookie
func (r *Response) GetCookie(name string) *http.Cookie {
	for _, cookie := range r.Cookies {
		if cookie.Name == name {
			return cookie
		}
	}
	return nil
}

// IsSuccess 判断是否成功（2xx）
func (r *Response) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

// IsRedirect 判断是否重定向（3xx）
func (r *Response) IsRedirect() bool {
	return r.StatusCode >= 300 && r.StatusCode < 400
}

// IsClientError 判断是否客户端错误（4xx）
func (r *Response) IsClientError() bool {
	return r.StatusCode >= 400 && r.StatusCode < 500
}

// IsServerError 判断是否服务器错误（5xx）
func (r *Response) IsServerError() bool {
	return r.StatusCode >= 500 && r.StatusCode < 600
}

// Close 关闭响应体
func (r *Response) Close() error {
	if r.body != nil {
		return r.body.Close()
	}
	return nil
}

// Request 获取原始请求
func (r *Response) Request() *http.Request {
	return r.request
}

