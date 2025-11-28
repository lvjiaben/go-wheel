package controller

import (
	"strings"

	"github.com/gin-gonic/gin"
	apiService "github.com/lvjiaben/go-wheel/app/api/v1/service"
	"github.com/lvjiaben/go-wheel/app/api/v1/validate"
	"github.com/lvjiaben/go-wheel/pkg/captcha"
	"github.com/lvjiaben/go-wheel/pkg/container"
	"github.com/lvjiaben/go-wheel/pkg/utils/http"
)

// UserController 用户控制器
type UserController struct {
	container      *container.Container
	service        *apiService.UserService
	captchaService *captcha.CaptchaService
}

// NewUserController 创建用户控制器
func NewUserController(c *container.Container) *UserController {
	return &UserController{
		container:      c,
		service:        apiService.NewUserService(c),
		captchaService: captcha.NewCaptchaService(c.GetRDB()),
	}
}

// Login 用户登录（账号密码）
func (c *UserController) Login(ctx *gin.Context) {
	// 验证请求参数
	req, valid := validate.ValidateLogin(ctx)
	if !valid {
		return
	}

	// 图形验证码验证
	if !c.captchaService.Verify(req.Captcha.ID, req.Captcha.Code) {
		http.ErrorWithI18n(ctx, "api.user.captcha_invalid", nil)
		return
	}

	// 调用服务层登录
	resp, err := c.service.Login(req)
	if err != nil {
		http.ErrorWithI18n(ctx, "api.user."+err.Error(), nil)
		return
	}

	http.SuccessWithI18n(ctx, "api.user.login_success", resp)
}

// MobileLogin 手机号登录
func (c *UserController) MobileLogin(ctx *gin.Context) {
	// 验证请求参数
	req, valid := validate.ValidateMobileLogin(ctx)
	if !valid {
		return
	}

	// 调用服务层登录
	resp, err := c.service.MobileLogin(req)
	if err != nil {
		http.ErrorWithI18n(ctx, "api.user."+err.Error(), nil)
		return
	}

	http.SuccessWithI18n(ctx, "api.user.login_success", resp)
}

// Reg 用户注册
func (c *UserController) Reg(ctx *gin.Context) {
	// 验证请求参数
	req, valid := validate.ValidateReg(ctx)
	if !valid {
		return
	}

	// 调用服务层注册
	resp, err := c.service.Reg(req)
	if err != nil {
		http.ErrorWithI18n(ctx, "api.user."+err.Error(), nil)
		return
	}

	http.SuccessWithI18n(ctx, "api.user.reg_success", resp)
}

// Logout 退出登录
func (c *UserController) Logout(ctx *gin.Context) {
	userId := ctx.GetInt("user_id")
	// 从请求头获取token
	token := ctx.GetHeader("Authorization")
	if token != "" {
		// 移除 Bearer 前缀
		token = strings.TrimPrefix(token, "Bearer ")
	}

	if err := c.service.Logout(userId, token); err != nil {
		http.ErrorWithI18n(ctx, "api.user."+err.Error(), nil)
		return
	}

	http.SuccessWithI18n(ctx, "api.user.logout_success", nil)
}

// ResetPwd 重置密码（忘记密码）
func (c *UserController) ResetPwd(ctx *gin.Context) {
	// 验证请求参数
	req, valid := validate.ValidateResetPwd(ctx)
	if !valid {
		return
	}

	// 调用服务层重置密码
	if err := c.service.ResetPwd(req); err != nil {
		http.ErrorWithI18n(ctx, "api.user."+err.Error(), nil)
		return
	}

	http.SuccessWithI18n(ctx, "api.user.reset_pwd_success", nil)
}

// ChangeMobile 修改手机号
func (c *UserController) ChangeMobile(ctx *gin.Context) {
	// 验证请求参数
	req, valid := validate.ValidateChangeMobile(ctx)
	if !valid {
		return
	}

	userId := ctx.GetInt("user_id")

	// 调用服务层修改手机号
	if err := c.service.ChangeMobile(userId, req); err != nil {
		http.ErrorWithI18n(ctx, "api.user."+err.Error(), nil)
		return
	}

	http.SuccessWithI18n(ctx, "api.user.change_mobile_success", nil)
}

// ChangePwd 修改密码
func (c *UserController) ChangePwd(ctx *gin.Context) {
	// 验证请求参数
	req, valid := validate.ValidateChangePwd(ctx)
	if !valid {
		return
	}

	userId := ctx.GetInt("user_id")

	// 调用服务层修改密码
	if err := c.service.ChangePwd(userId, req); err != nil {
		http.ErrorWithI18n(ctx, "api.user."+err.Error(), nil)
		return
	}

	http.SuccessWithI18n(ctx, "api.user.change_pwd_success", nil)
}

// Info 获取用户信息
func (c *UserController) Info(ctx *gin.Context) {
	userId := ctx.GetInt("user_id")
	user, err := c.service.GetById(userId)
	if err != nil {
		http.ErrorWithI18n(ctx, "api.user.user_not_found", nil)
		return
	}
	http.SuccessWithI18n(ctx, "common.success", gin.H{
		"id":       user.Id,
		"username": user.Username,
		"avatar":   user.Avatar,
		"mobile":   user.Mobile,
		"realName": user.Username,
		"score":    user.Score,
		"money":    user.Money,
	})
}
