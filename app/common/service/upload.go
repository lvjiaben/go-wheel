package service

import (
	"context"
	"crypto/sha256"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/lvjiaben/go-wheel/pkg/utils/datatype"

	"github.com/tencentyun/cos-go-sdk-v5"

	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/lvjiaben/go-wheel/pkg/container"
)

// UploadResult 上传结果
type UploadResult struct {
	Path      string `json:"path"`      // 存储路径
	Parent    string `json:"parent"`    // 父级文件夹
	URL       string `json:"url"`       // 在线HTTP链接
	Filename  string `json:"filename"`  // 文件名称
	Size      int64  `json:"size"`      // 文件大小
	MediaType string `json:"mediatype"` // 文件类型
	Extension string `json:"extension"` // 文件后缀
}

// UploadService 上传服务
type UploadService struct {
	container *container.Container
}

// NewUploadService 创建上传服务
func NewUploadService(c *container.Container) *UploadService {
	return &UploadService{
		container: c,
	}
}

// Delete 删除文件
func (s *UploadService) Delete(path, typeName string) {
	path = strings.TrimLeft(strings.ReplaceAll(path, "\\", "/"), "/")
	if typeName == "oss" {
		// 创建OSS客户端
		client, err := oss.New(s.container.GetConfig().Upload.Oss.Endpoint, s.container.GetConfig().Upload.Oss.AccessKeyID, s.container.GetConfig().Upload.Oss.AccessKeySecret)
		if err != nil {
			s.container.GetLogger().Error(fmt.Sprintf("创建OSS客户端失败: %v", err))
			return
		}
		bucket, err := client.Bucket(s.container.GetConfig().Upload.Oss.BucketName)
		if err != nil {
			s.container.GetLogger().Error(fmt.Sprintf("获取OSS bucket失败: %v", err))
			return
		}
		err = bucket.DeleteObject(path)
		if err != nil {
			s.container.GetLogger().Error(fmt.Sprintf("删除OSS文件失败: %v", err))
			return
		}
	} else if typeName == "qiniu" {
		// 创建认证
		mac := qbox.NewMac(s.container.GetConfig().Upload.Qiniu.AccessKey, s.container.GetConfig().Upload.Qiniu.SecretKey)
		bm := storage.NewBucketManager(mac, &storage.Config{})
		err := bm.Delete(s.container.GetConfig().Upload.Qiniu.Bucket, path)
		if err != nil {
			s.container.GetLogger().Error(fmt.Sprintf("删除QINIU文件失败: %v", err))
			return
		}
	} else if typeName == "cos" {
		bucketURL := s.container.GetConfig().Upload.Cos.Bucket
		// 解析bucket URL
		u, err := url.Parse("https://" + bucketURL + ".cos." + s.container.GetConfig().Upload.Cos.Region + ".myqcloud.com")
		if err != nil {
			s.container.GetLogger().Error(fmt.Sprintf("URL解析失败: %v", err))
			return
		}
		// 创建COS客户端
		b := &cos.BaseURL{BucketURL: u}
		client := cos.NewClient(b, &http.Client{
			Transport: &cos.AuthorizationTransport{
				SecretID:  s.container.GetConfig().Upload.Cos.SecretId,
				SecretKey: s.container.GetConfig().Upload.Cos.SecretKey,
			},
		})
		// 上传文件
		_, err = client.Object.Delete(context.Background(), path)
		if err != nil {
			s.container.GetLogger().Error(fmt.Sprintf("删除COS文件失败: %v", err))
			return
		}
	} else {
		pwd, _ := os.Getwd()
		filePath := filepath.Join(pwd, path)
		err := os.Remove(filePath)
		if err != nil {
			s.container.GetLogger().Error("删除文件失败（" + path + "）：" + err.Error())
		}
	}
}

// Upload 上传文件
func (s *UploadService) Upload(file *multipart.FileHeader) (*UploadResult, error) {
	// 验证文件
	if err := s.validateFile(file); err != nil {
		return nil, err
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("Open File Fail: %v", err)
	}
	defer src.Close()

	// 计算文件SHA-256哈希
	hash := sha256.New()
	if _, err := io.Copy(hash, src); err != nil {
		return nil, fmt.Errorf("计算文件哈希失败: %v", err)
	}
	src.Seek(0, 0) // 重置文件指针

	// 生成文件名和路径
	fileExt := strings.ToLower(filepath.Ext(file.Filename))
	fileName := fmt.Sprintf("%x%s", hash.Sum(nil), fileExt)

	dateDir := filepath.Join(s.container.GetConfig().Upload.UploadPath, time.Now().Format("2006-01-02"))
	relative := fmt.Sprintf("%s/%s", strings.TrimLeft(strings.ReplaceAll(dateDir, "\\", "/"), "/"), fileName)
	// 根据配置类型上传
	typeName := s.container.GetConfig().Upload.Type
	if typeName == "oss" {
		// 创建OSS客户端
		client, err := oss.New(s.container.GetConfig().Upload.Oss.Endpoint, s.container.GetConfig().Upload.Oss.AccessKeyID, s.container.GetConfig().Upload.Oss.AccessKeySecret)
		if err != nil {
			return nil, fmt.Errorf("创建OSS客户端失败: %v", err)
		}
		bucket, err := client.Bucket(s.container.GetConfig().Upload.Oss.BucketName)
		if err != nil {
			return nil, fmt.Errorf("获取OSS bucket失败: %v", err)
		}
		// 上传文件
		err = bucket.PutObject(strings.TrimLeft(relative, "/"), src)
		if err != nil {
			return nil, fmt.Errorf("上传到OSS失败: %v", err)
		}
	} else if typeName == "qiniu" {
		// 创建认证
		mac := qbox.NewMac(s.container.GetConfig().Upload.Qiniu.AccessKey, s.container.GetConfig().Upload.Qiniu.SecretKey)
		putPolicy := storage.PutPolicy{
			Scope: s.container.GetConfig().Upload.Qiniu.Bucket,
		}
		upToken := putPolicy.UploadToken(mac)
		formUploader := storage.NewFormUploader(&storage.Config{})
		// 上传文件
		err := formUploader.Put(nil, &storage.PutRet{}, upToken, relative, src, file.Size, nil)
		if err != nil {
			return nil, fmt.Errorf("上传到七牛云失败: %v", err)
		}
	} else if typeName == "cos" {
		bucketURL := s.container.GetConfig().Upload.Cos.Bucket
		// 解析bucket URL
		u, err := url.Parse("https://" + bucketURL + ".cos." + s.container.GetConfig().Upload.Cos.Region + ".myqcloud.com")
		if err != nil {
			return nil, fmt.Errorf("解析COS URL失败: %v", err)
		}
		// 创建COS客户端
		b := &cos.BaseURL{BucketURL: u}
		client := cos.NewClient(b, &http.Client{
			Transport: &cos.AuthorizationTransport{
				SecretID:  s.container.GetConfig().Upload.Cos.SecretId,
				SecretKey: s.container.GetConfig().Upload.Cos.SecretKey,
			},
		})
		// 上传文件
		_, err = client.Object.Put(nil, relative, src, nil)
		if err != nil {
			return nil, fmt.Errorf("上传到腾讯云COS失败: %v", err)
		}
	} else {
		// 创建目录
		pwd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("Error getting current directory: %v", err)
		}
		fullDir := filepath.Join(pwd, dateDir)
		if err := os.MkdirAll(fullDir, 0755); err != nil {
			return nil, fmt.Errorf("创建目录失败: %v", err)
		}
		// 保存文件
		fullPath := filepath.Join(fullDir, fileName)
		dst, err := os.Create(fullPath)
		if err != nil {
			return nil, fmt.Errorf("创建文件失败: %v", err)
		}
		defer dst.Close()

		if _, err := io.Copy(dst, src); err != nil {
			return nil, fmt.Errorf("保存文件失败: %v", err)
		}
	}
	// 文件网络地址
	fileURL := strings.TrimRight(s.container.GetConfig().Upload.BaseUrl, "/") + "/" + relative
	return &UploadResult{
		Path:      dateDir,
		Parent:    filepath.Base(dateDir),
		URL:       fileURL,
		Filename:  fileName,
		Size:      file.Size,
		MediaType: file.Header.Get("Content-Type"),
		Extension: strings.ToLower(fileExt[1:]),
	}, nil
}

// validateFile 验证文件
func (s *UploadService) validateFile(file *multipart.FileHeader) error {
	// 检查文件大小
	maxSize := s.container.GetConfig().Upload.MaxSize
	if file.Size > maxSize {
		return fmt.Errorf("文件大小超过限制，最大允许 %d 字节", maxSize)
	}

	// 清理文件名，防止路径遍历攻击
	cleanFilename := filepath.Base(filepath.Clean(file.Filename))
	if cleanFilename != file.Filename {
		return fmt.Errorf("文件名包含非法字符或路径")
	}

	// 检查文件名是否包含危险字符
	if strings.Contains(cleanFilename, "..") || strings.Contains(cleanFilename, "/") || strings.Contains(cleanFilename, "\\") {
		return fmt.Errorf("文件名包含非法路径字符")
	}

	// 检查文件扩展名
	fileExt := strings.ToLower(filepath.Ext(cleanFilename))
	if fileExt != "" {
		fileExt = fileExt[1:] // 去掉点号
	}

	// 严格的扩展名白名单验证
	allowedExts := s.container.GetConfig().Upload.AllowedExtensions
	if len(allowedExts) == 0 {
		// 如果没有配置白名单，使用默认安全的扩展名列表
		allowedExts = []string{"jpg", "jpeg", "png", "gif", "pdf", "doc", "docx", "xls", "xlsx"}
	}

	if !datatype.Contains(allowedExts, fileExt) {
		return fmt.Errorf("不支持的文件类型: %s，仅允许: %v", fileExt, allowedExts)
	}

	// 检查文件扩展名是否为可执行文件
	dangerousExts := []string{"exe", "bat", "cmd", "sh", "php", "jsp", "asp", "aspx", "js", "vbs", "ps1"}
	if datatype.Contains(dangerousExts, fileExt) {
		return fmt.Errorf("禁止上传可执行文件: %s", fileExt)
	}

	// 检查MIME类型
	src, err := file.Open()
	if err != nil {
		return fmt.Errorf("打开文件失败: %v", err)
	}
	defer src.Close()

	buffer := make([]byte, 512)
	_, err = src.Read(buffer)
	if err != nil && err != io.EOF {
		return fmt.Errorf("读取文件失败: %v", err)
	}

	mimeType := file.Header.Get("Content-Type")
	allowedTypes := s.container.GetConfig().Upload.AllowedTypes
	if len(allowedTypes) > 0 && !datatype.Contains(allowedTypes, mimeType) {
		return fmt.Errorf("不支持的文件MIME类型: %s", mimeType)
	}

	// 简单的病毒扫描：检查文件内容中是否包含可疑特征
	if err := s.scanFileContent(buffer); err != nil {
		return err
	}

	return nil
}

// scanFileContent 扫描文件内容中的可疑特征（简单的启发式检测）
func (s *UploadService) scanFileContent(content []byte) error {
	// 检查是否包含可执行文件的魔数（文件头特征）
	suspiciousSignatures := map[string][]byte{
		"PE可执行文件":    {0x4D, 0x5A},                   // MZ (Windows PE)
		"ELF可执行文件":   {0x7F, 0x45, 0x4C, 0x46},       // ELF (Linux)
		"Mach-O可执行文件": {0xFE, 0xED, 0xFA, 0xCE},       // Mach-O (macOS)
		"脚本文件":       {0x23, 0x21},                   // #! (Shebang)
		"批处理文件":      {0x40, 0x65, 0x63, 0x68, 0x6F}, // @echo
	}

	for name, signature := range suspiciousSignatures {
		if len(content) >= len(signature) {
			if string(content[:len(signature)]) == string(signature) {
				return fmt.Errorf("检测到可疑文件类型: %s", name)
			}
		}
	}

	// 检查是否包含可疑的脚本代码特征
	suspiciousPatterns := []string{
		"<script",
		"javascript:",
		"eval(",
		"exec(",
		"system(",
		"shell_exec(",
		"passthru(",
		"<?php",
		"<%",
	}

	contentStr := string(content)
	contentLower := strings.ToLower(contentStr)

	for _, pattern := range suspiciousPatterns {
		if strings.Contains(contentLower, strings.ToLower(pattern)) {
			return fmt.Errorf("检测到可疑代码特征: %s", pattern)
		}
	}

	return nil
}
