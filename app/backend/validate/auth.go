package validate

import (
	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/app/common/validator"
)

type CaptchaData struct {
	ID   string `json:"id" binding:"required"`
	Code string `json:"code" binding:"required"`
}

// LoginRequest 登录请求验证器
type LoginRequest struct {
	Username string      `json:"username" binding:"required" label:"backend.auth.username" msg:"backend.auth.username_required"`
	Password string      `json:"password" binding:"required" label:"backend.auth.password" msg:"backend.auth.password_required"`
	Captcha  CaptchaData `json:"captcha" binding:"required" label:"backend.auth.captcha" msg:"backend.auth.captcha_required"`
}

// ProfileRequest 修改资料请求验证器
type ProfileRequest struct {
	Avatar string `json:"avatar" binding:"required" label:"backend.auth.avatar" msg:"backend.auth.avatar_required"`
	Email  string `json:"email" binding:"required" label:"backend.auth.email" msg:"backend.auth.email_required"`
}

// ValidateLogin 验证登录请求
func ValidateLogin(c *gin.Context) (*LoginRequest, bool) {
	return validator.ValidateStruct[LoginRequest](c)
}

// ValidateProfile 验证资料请求
func ValidateProfile(c *gin.Context) (*ProfileRequest, bool) {
	return validator.ValidateStruct[ProfileRequest](c)
}
