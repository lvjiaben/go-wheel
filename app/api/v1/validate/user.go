package validate

import (
	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/app/common/validator"
)

// CaptchaInfo 图形验证码信息
type CaptchaInfo struct {
	ID   string `json:"id" binding:"required" msg:"api.user.captcha_id_required"`
	Code string `json:"code" binding:"required" msg:"api.user.captcha_code_required"`
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username string      `json:"username" binding:"required" msg:"api.user.username_required"`
	Password string      `json:"password" binding:"required" msg:"api.user.password_required"`
	Captcha  CaptchaInfo `json:"captcha" binding:"required" msg:"api.user.captcha_required"`
}

// ValidateLogin 验证登录请求
func ValidateLogin(c *gin.Context) (*LoginRequest, bool) {
	return validator.ValidateStructWithConvert[LoginRequest](c)
}

// MobileLoginRequest 手机号登录请求
type MobileLoginRequest struct {
	Mobile string `json:"mobile" binding:"required" msg:"api.user.mobile_required"`
	Code   string `json:"code" binding:"required" msg:"api.user.sms_code_required"`
}

// ValidateMobileLogin 验证手机号登录请求
func ValidateMobileLogin(c *gin.Context) (*MobileLoginRequest, bool) {
	return validator.ValidateStructWithConvert[MobileLoginRequest](c)
}

// RegRequest 注册请求
type RegRequest struct {
	Mobile     string `json:"mobile" binding:"required" msg:"api.user.mobile_required"`
	Password   string `json:"password" binding:"required,min=6" msg:"api.user.password_required"`
	Code       string `json:"code" binding:"required" msg:"api.user.sms_code_required"`
	InviteCode string `json:"invite_code" msg:"api.user.invite_code_invalid"`
}

// ValidateReg 验证注册请求
func ValidateReg(c *gin.Context) (*RegRequest, bool) {
	return validator.ValidateStructWithConvert[RegRequest](c)
}

// ResetPwdRequest 重置密码请求
type ResetPwdRequest struct {
	Mobile      string `json:"mobile" binding:"required" msg:"api.user.mobile_required"`
	Code        string `json:"code" binding:"required" msg:"api.user.sms_code_required"`
	NewPassword string `json:"new_password" binding:"required,min=6" msg:"api.user.new_password_required"`
}

// ValidateResetPwd 验证重置密码请求
func ValidateResetPwd(c *gin.Context) (*ResetPwdRequest, bool) {
	return validator.ValidateStructWithConvert[ResetPwdRequest](c)
}

// ChangeMobileRequest 修改手机号请求
type ChangeMobileRequest struct {
	Mobile    string `json:"mobile" binding:"required" msg:"api.user.mobile_required"`
	Code      string `json:"code" binding:"required" msg:"api.user.sms_code_required"`
	NewMobile string `json:"new_mobile" binding:"required" msg:"api.user.new_mobile_required"`
}

// ValidateChangeMobile 验证修改手机号请求
func ValidateChangeMobile(c *gin.Context) (*ChangeMobileRequest, bool) {
	return validator.ValidateStructWithConvert[ChangeMobileRequest](c)
}

// ChangePwdRequest 修改密码请求
type ChangePwdRequest struct {
	OldPassword string `json:"old_password" binding:"required" msg:"api.user.old_password_required"`
	NewPassword string `json:"new_password" binding:"required,min=6" msg:"api.user.new_password_required"`
}

// ValidateChangePwd 验证修改密码请求
func ValidateChangePwd(c *gin.Context) (*ChangePwdRequest, bool) {
	return validator.ValidateStructWithConvert[ChangePwdRequest](c)
}
