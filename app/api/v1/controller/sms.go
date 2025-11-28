package controller

import (
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/app/api/v1/validate"
	commonService "github.com/lvjiaben/go-wheel/app/common/service"
	"github.com/lvjiaben/go-wheel/pkg/container"
	"github.com/lvjiaben/go-wheel/pkg/utils/http"
)

// SmsController 短信控制器
type SmsController struct {
	container  *container.Container
	smsService *commonService.SmsService
}

// NewSmsController 创建短信控制器
func NewSmsController(c *container.Container) *SmsController {
	return &SmsController{
		container:  c,
		smsService: commonService.NewSmsService(c),
	}
}

// Send 发送短信验证码
func (c *SmsController) Send(ctx *gin.Context) {
	req, ok := validate.ValidateSendSms(ctx)
	if !ok {
		return
	}

	event := req.Event
	mobile := req.Mobile

	if event == "" {
		event = "default"
	}

	// 定义不需要登录的事件类型
	noLoginEvents := []string{"login", "register", "resetpwd"}

	// 检查是否是不需要登录的事件
	isNoLoginEvent := false
	for _, e := range noLoginEvents {
		if event == e {
			isNoLoginEvent = true
			break
		}
	}

	if isNoLoginEvent {
		// 不需要登录的事件，验证手机号
		if mobile == "" {
			http.ErrorWithI18n(ctx, "api.sms.mobile_required", nil)
			return
		}

		// 验证手机号格式
		matched, _ := regexp.MatchString(`^1\d{10}$`, mobile)
		if !matched {
			http.ErrorWithI18n(ctx, "api.sms.mobile_invalid", nil)
			return
		}
	} else {
		// 需要登录的事件，从数据库查询用户手机号
		userId := ctx.GetInt("user_id")
		if userId == 0 {
			http.ErrorWithI18n(ctx, "api.sms.not_logged_in", nil)
			return
		}

		// 从数据库获取用户手机号
		var userMobile string
		if err := c.container.GetDB().Table("user").Where("id = ?", userId).Pluck("mobile", &userMobile).Error; err != nil || userMobile == "" {
			http.ErrorWithI18n(ctx, "api.sms.user_mobile_not_found", nil)
			return
		}
		mobile = userMobile
	}

	// 重置密码事件，检查手机号是否注册
	if event == "resetpwd" {
		var count int64
		c.container.GetDB().Table("user").Where("mobile = ?", mobile).Count(&count)
		if count == 0 {
			http.ErrorWithI18n(ctx, "api.sms.mobile_not_registered", nil)
			return
		}
	}

	// 注册事件，检查手机号是否已注册
	if event == "register" {
		var count int64
		c.container.GetDB().Table("user").Where("mobile = ?", mobile).Count(&count)
		if count > 0 {
			http.ErrorWithI18n(ctx, "api.sms.mobile_already_registered", nil)
			return
		}
	}

	// 发送短信验证码
	if err := c.smsService.Send(mobile, "", event); err != nil {
		// 处理特定错误
		errMsg := err.Error()
		if errMsg == "send_too_frequent" {
			http.ErrorWithI18n(ctx, "api.sms.send_too_frequent", nil)
			return
		}
		http.ErrorWithI18n(ctx, "api.sms.send_failed", nil)
		return
	}

	http.SuccessWithI18n(ctx, "api.sms.send_success", map[string]any{
		"expire": c.smsService.GetExpire(),
		"mobile": mobile,
		"event":  event,
	})
}
