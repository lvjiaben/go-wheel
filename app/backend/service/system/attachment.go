package system

import (
	"fmt"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/app/backend/model/system"
	"github.com/lvjiaben/go-wheel/app/common/builder"
	commonService "github.com/lvjiaben/go-wheel/app/common/service"
	"github.com/lvjiaben/go-wheel/pkg/container"
)

// AttachmentService 附件服务
type AttachmentService struct {
	container     *container.Container
	uploadService *commonService.UploadService
}

// NewAttachmentService 创建附件服务
func NewAttachmentService(c *container.Container) *AttachmentService {
	return &AttachmentService{
		container:     c,
		uploadService: commonService.NewUploadService(c),
	}
}

func (s *AttachmentService) List(ctx *gin.Context) map[string]interface{} {
	db := s.container.GetDB().WithContext(ctx.Request.Context())
	return builder.NewCRUDBuilder[system.Attachment](db).WithFilter(false).WithContext(ctx).WithSearchFields("filename", "url").List()
}

// GetDirectories 获取目录列表（用于左侧栏）
func (s *AttachmentService) GetDirectories() ([]map[string]interface{}, error) {
	var directories []map[string]interface{}

	// 查询所有不重复的parent值
	err := s.container.GetDB().Model(&system.Attachment{}).
		Select("DISTINCT parent").
		Where("parent != ''").
		Scan(&directories).Error

	if err != nil {
		return nil, fmt.Errorf("backend.attachment.query_failed")
	}

	// 构建目录树结构
	result := []map[string]interface{}{
		{
			"name": "全部",
			"path": "",
		},
	}

	// 添加其他目录
	for _, dir := range directories {
		if parent, ok := dir["parent"].(string); ok && parent != "" {
			result = append(result, map[string]interface{}{
				"name": parent,
				"path": parent,
			})
		}
	}

	return result, nil
}

// Upload 上传附件
func (s *AttachmentService) Upload(file *multipart.FileHeader, adminId int) (*system.Attachment, error) {

	// 使用公共上传服务上传文件（固定使用空字符串作为parent）
	result, err := s.uploadService.Upload(file)
	if err != nil {
		return nil, fmt.Errorf("backend.attachment.upload_failed: %v", err)
	}
	s.container.GetDB().Where("filename = ?", result.Filename).Delete(&system.Attachment{})
	// 保存到数据库
	attachment := &system.Attachment{
		Type:      s.container.GetConfig().Upload.Type,
		AdminId:   adminId,
		UserId:    0,
		Path:      result.Path,
		Parent:    result.Parent,
		URL:       result.URL,
		Filename:  result.Filename,
		Size:      result.Size,
		MediaType: result.MediaType,
		Extension: result.Extension,
		CreatedAt: int(time.Now().Unix()),
		UpdatedAt: int(time.Now().Unix()),
	}

	if err := s.container.GetDB().Create(attachment).Error; err != nil {
		return nil, fmt.Errorf("backend.attachment.save_failed")
	}

	return attachment, nil
}

// Delete 删除附件
func (s *AttachmentService) Delete(id int) {
	var attachment system.Attachment

	// 查找附件
	if err := s.container.GetDB().Where("id = ?", id).First(&attachment).Error; err != nil {
		return
	}

	// 数据库直接删除
	s.container.GetDB().Where("id = ?", id).Delete(&attachment)

	// 实体文件删除
	filePath := filepath.Join(attachment.Path, attachment.Filename)
	s.uploadService.Delete(filePath, attachment.Type)

}
