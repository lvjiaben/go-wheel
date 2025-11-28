package validate

import (
	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/app/common/validator"
)

// SendSmsRequest 发送短信验证码请求
type SendSmsRequest struct {
	Mobile string `json:"mobile" msg:"api.sms.mobile_required"`
	Event  string `json:"event" msg:"api.sms.event_required"`
}

// ValidateSendSms 验证发送短信请求
func ValidateSendSms(c *gin.Context) (*SendSmsRequest, bool) {
	return validator.ValidateStructWithConvert[SendSmsRequest](c)
}

