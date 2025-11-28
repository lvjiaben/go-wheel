package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/pkg/captcha"
	"github.com/lvjiaben/go-wheel/pkg/container"
	"github.com/lvjiaben/go-wheel/pkg/utils/http"
)

// CaptchaController 验证码控制器
type CaptchaController struct {
	container      *container.Container
	captchaService *captcha.CaptchaService
}

// NewCaptchaController 创建验证码控制器
func NewCaptchaController(c *container.Container) *CaptchaController {
	return &CaptchaController{
		container:      c,
		captchaService: captcha.NewCaptchaService(c.GetRDB()),
	}
}

// Generate 生成图形验证码
func (c *CaptchaController) Generate(ctx *gin.Context) {
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

