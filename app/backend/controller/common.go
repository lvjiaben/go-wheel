package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/pkg/captcha"
	"github.com/lvjiaben/go-wheel/pkg/container"
	"github.com/lvjiaben/go-wheel/pkg/utils/http"
)

// CommonController 认证控制器
type CommonController struct {
	container      *container.Container
	captchaService *captcha.CaptchaService
}

// NewCommonController 创建认证控制器
func NewCommonController(c *container.Container) *CommonController {
	return &CommonController{
		container:      c,
		captchaService: captcha.NewCaptchaService(c.GetRDB()),
	}
}

// Captcha 生成验证码
func (c *CommonController) Captcha(ctx *gin.Context) {
	// 创建验证码
	result, err := c.captchaService.Generate()
	if err != nil {
		http.ErrorWithI18n(ctx, "common.server_error", nil)
		return
	}
	http.SuccessWithI18n(ctx, "common.success", gin.H{
		"captcha_id":   result.ID,
		"captcha_data": result.Base64,
	})
}
