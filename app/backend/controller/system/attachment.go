package system

import (
	"github.com/gin-gonic/gin"
	serviceSystem "github.com/lvjiaben/go-wheel/app/backend/service/system"
	validateSystem "github.com/lvjiaben/go-wheel/app/backend/validate/system"
	"github.com/lvjiaben/go-wheel/pkg/container"
	"github.com/lvjiaben/go-wheel/pkg/utils/http"
)

// AttachmentController 附件控制器
type AttachmentController struct {
	container         *container.Container
	attachmentService *serviceSystem.AttachmentService
}

// NewAttachmentController 创建附件控制器
func NewAttachmentController(c *container.Container) *AttachmentController {
	return &AttachmentController{
		container:         c,
		attachmentService: serviceSystem.NewAttachmentService(c),
	}
}

// List 获取附件列表
func (c *AttachmentController) List(ctx *gin.Context) {
	http.SuccessWithI18n(ctx, "common.success", c.attachmentService.List(ctx))
}

// Directories 获取目录列表
func (c *AttachmentController) Directories(ctx *gin.Context) {
	// 获取目录列表
	result, err := c.attachmentService.GetDirectories()
	if err != nil {
		http.ErrorWithI18n(ctx, err.Error(), nil)
		return
	}

	http.SuccessWithI18n(ctx, "common.success", result)
}

// Upload 上传附件
func (c *AttachmentController) Upload(ctx *gin.Context) {
	// 验证请求参数
	_, valid := validateSystem.ValidateAttachmentUpload(ctx)
	if !valid {
		return
	}

	// 获取上传的文件
	file, err := ctx.FormFile("file")
	if err != nil {
		http.ErrorWithI18n(ctx, "common.error", nil)
		return
	}

	// 上传文件
	attachment, err := c.attachmentService.Upload(file, ctx.GetInt("admin_id"))
	if err != nil {
		http.ErrorWithI18n(ctx, err.Error(), nil)
		return
	}

	http.SuccessWithI18n(ctx, "common.success", attachment)
}

// Delete 删除附件
func (c *AttachmentController) Delete(ctx *gin.Context) {
	data, valid := validateSystem.ValidateAttachmentIds(ctx)
	if !valid {
		return
	}
	for _, id := range data.Ids {
		c.attachmentService.Delete(id)
	}
	http.SuccessWithI18n(ctx, "common.success", nil)
}
