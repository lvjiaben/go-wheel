package httpclient

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// DownloadFile 下载文件
func (c *Client) DownloadFile(url, savePath string) error {
	resp, err := c.Get(url).Send()
	if err != nil {
		return fmt.Errorf("下载文件失败: %v", err)
	}
	defer resp.Close()

	if !resp.IsSuccess() {
		return fmt.Errorf("下载文件失败，状态码: %d", resp.StatusCode)
	}

	return resp.SaveToFileStream(savePath)
}

// DownloadFileWithProgress 下载文件（带进度回调）
func (c *Client) DownloadFileWithProgress(url, savePath string, callback func(downloaded, total int64)) error {
	resp, err := c.Get(url).Send()
	if err != nil {
		return fmt.Errorf("下载文件失败: %v", err)
	}
	defer resp.Close()

	if !resp.IsSuccess() {
		return fmt.Errorf("下载文件失败，状态码: %d", resp.StatusCode)
	}

	// 获取文件大小
	total := resp.request.ContentLength

	// 创建文件
	file, err := os.Create(savePath)
	if err != nil {
		return fmt.Errorf("创建文件失败: %v", err)
	}
	defer file.Close()

	// 创建进度读取器
	reader := &progressReader{
		reader:   resp.body,
		total:    total,
		callback: callback,
	}

	// 复制数据
	_, err = io.Copy(file, reader)
	if err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	return nil
}

// UploadFile 上传单个文件
func (c *Client) UploadFile(url, fieldName, filePath string) (*Response, error) {
	return c.Post(url).SetFile(fieldName, filePath).Send()
}

// UploadFiles 上传多个文件
func (c *Client) UploadFiles(url string, files map[string]string) (*Response, error) {
	return c.Post(url).SetFiles(files).Send()
}

// UploadFileWithData 上传文件并附带表单数据
func (c *Client) UploadFileWithData(url, fieldName, filePath string, formData map[string]string) (*Response, error) {
	req := c.Post(url).SetFile(fieldName, filePath)
	if formData != nil {
		req.SetFormData(formData)
	}
	return req.Send()
}

// progressReader 进度读取器
type progressReader struct {
	reader     io.Reader
	total      int64
	downloaded int64
	callback   func(downloaded, total int64)
}

// Read 实现 io.Reader 接口
func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.reader.Read(p)
	pr.downloaded += int64(n)

	if pr.callback != nil {
		pr.callback(pr.downloaded, pr.total)
	}

	return n, err
}

// DownloadToWriter 下载到 Writer
func (c *Client) DownloadToWriter(url string, writer io.Writer) error {
	resp, err := c.Get(url).Send()
	if err != nil {
		return fmt.Errorf("下载失败: %v", err)
	}
	defer resp.Close()

	if !resp.IsSuccess() {
		return fmt.Errorf("下载失败，状态码: %d", resp.StatusCode)
	}

	_, err = io.Copy(writer, resp.body)
	if err != nil {
		return fmt.Errorf("写入数据失败: %v", err)
	}

	return nil
}

// UploadFromReader 从 Reader 上传数据
func (c *Client) UploadFromReader(url string, reader io.Reader, contentType string) (*Response, error) {
	req := c.Post(url).SetBodyReader(reader)
	if contentType != "" {
		req.SetHeader("Content-Type", contentType)
	}
	return req.Send()
}

// DownloadFileResume 断点续传下载
func (c *Client) DownloadFileResume(url, savePath string) error {
	// 检查文件是否存在
	var downloaded int64 = 0
	if fileInfo, err := os.Stat(savePath); err == nil {
		downloaded = fileInfo.Size()
	}

	// 创建请求
	req := c.Get(url)
	if downloaded > 0 {
		req.SetHeader("Range", fmt.Sprintf("bytes=%d-", downloaded))
	}

	// 发送请求
	resp, err := req.Send()
	if err != nil {
		return fmt.Errorf("下载文件失败: %v", err)
	}
	defer resp.Close()

	// 检查状态码（206 表示部分内容，200 表示完整内容）
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusPartialContent {
		return fmt.Errorf("下载文件失败，状态码: %d", resp.StatusCode)
	}

	// 打开文件（追加模式）
	var file *os.File
	if downloaded > 0 && resp.StatusCode == http.StatusPartialContent {
		file, err = os.OpenFile(savePath, os.O_APPEND|os.O_WRONLY, 0644)
	} else {
		file, err = os.Create(savePath)
	}
	if err != nil {
		return fmt.Errorf("打开文件失败: %v", err)
	}
	defer file.Close()

	// 复制数据
	_, err = io.Copy(file, resp.body)
	if err != nil {
		return fmt.Errorf("写入文件失败: %v", err)
	}

	return nil
}

// GetFileInfo 获取远程文件信息
func (c *Client) GetFileInfo(url string) (*FileInfo, error) {
	resp, err := c.Head(url).Send()
	if err != nil {
		return nil, fmt.Errorf("获取文件信息失败: %v", err)
	}
	defer resp.Close()

	if !resp.IsSuccess() {
		return nil, fmt.Errorf("获取文件信息失败，状态码: %d", resp.StatusCode)
	}

	info := &FileInfo{
		Size:        resp.request.ContentLength,
		ContentType: resp.GetHeader("Content-Type"),
		FileName:    getFileNameFromResponse(resp),
		AcceptRange: resp.GetHeader("Accept-Ranges") == "bytes",
	}

	return info, nil
}

// FileInfo 文件信息
type FileInfo struct {
	Size        int64
	ContentType string
	FileName    string
	AcceptRange bool
}

// getFileNameFromResponse 从响应中获取文件名
func getFileNameFromResponse(resp *Response) string {
	// 从 Content-Disposition 获取
	disposition := resp.GetHeader("Content-Disposition")
	if disposition != "" {
		// 简单解析，实际应该更严格
		if len(disposition) > 0 {
			// 这里简化处理，实际应该解析 Content-Disposition
			return "downloaded_file"
		}
	}

	// 从 URL 获取
	if resp.request != nil && resp.request.URL != nil {
		return filepath.Base(resp.request.URL.Path)
	}

	return "downloaded_file"
}

