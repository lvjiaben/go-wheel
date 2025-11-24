package system

import (
	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/app/common/validator"
	"github.com/lvjiaben/go-wheel/pkg/utils/http"
)

// AttachmentUpload 附件上传请求
type AttachmentUpload struct {
	// 不需要parent字段
}

type AttachmentIds struct {
	Ids []int `json:"ids" binding:"required,min=1" label:"Ids"`
}

func ValidateAttachmentIds(c *gin.Context) (*AttachmentIds, bool) {
	return validator.ValidateStruct[AttachmentIds](c)
}

// ValidateAttachmentUpload 验证附件上传请求
func ValidateAttachmentUpload(ctx *gin.Context) (*AttachmentUpload, bool) {
	var req AttachmentUpload

	// 检查是否有文件
	_, header, err := ctx.Request.FormFile("file")
	if err != nil {
		http.ErrorWithI18n(ctx, "backend.attachment.no_file", nil)
		return nil, false
	}

	if header == nil {
		http.ErrorWithI18n(ctx, "backend.attachment.no_file", nil)
		return nil, false
	}

	return &req, true
}
