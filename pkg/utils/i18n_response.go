package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/pkg/container"
)

// GetLang 从Context获取语言
func GetLang(c *gin.Context) string {
	lang := "zh-CN" // 默认中文
	if l, exists := c.Get("lang"); exists {
		if lStr, ok := l.(string); ok {
			lang = lStr
		}
	}
	return lang
}

// GetContainer 从Context获取Container
func GetContainer(c *gin.Context) *container.Container {
	if cont, exists := c.Get("container"); exists {
		if container, ok := cont.(*container.Container); ok {
			return container
		}
	}
	return nil
}

// I18nResponse 国际化响应结构体
type I18nResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// ResponseSuccess 成功响应，默认code为200
func ResponseSuccess(c *gin.Context, data interface{}, code ...int) {
	statusCode := 200
	if len(code) > 0 {
		statusCode = code[0]
	}

	c.JSON(http.StatusOK, I18nResponse{
		Code:    statusCode,
		Message: "success",
		Data:    data,
	})
}

// ResponseError 错误响应，默认code为500
func ResponseError(c *gin.Context, message string, data interface{}, code ...int) {
	statusCode := 500
	if len(code) > 0 {
		statusCode = code[0]
	}

	c.JSON(http.StatusOK, I18nResponse{
		Code:    statusCode,
		Message: message,
		Data:    data,
	})
}

// SuccessWithI18n 带国际化的成功响应，默认code为200
func SuccessWithI18n(c *gin.Context, messageKey string, data interface{}, code ...int) {
	statusCode := 200
	if len(code) > 0 {
		statusCode = code[0]
	}

	container := GetContainer(c)
	if container == nil {
		ResponseSuccess(c, data, statusCode)
		return
	}

	lang := GetLang(c)
	message := container.GetI18n().Get(messageKey, lang)

	c.JSON(http.StatusOK, I18nResponse{
		Code:    statusCode,
		Message: message,
		Data:    data,
	})
}

// ErrorWithI18n 带国际化的错误响应，默认code为500
func ErrorWithI18n(c *gin.Context, messageKey string, data interface{}, code ...int) {
	statusCode := 500
	if len(code) > 0 {
		statusCode = code[0]
	}

	container := GetContainer(c)
	if container == nil {
		ResponseError(c, messageKey, data, statusCode) // 如果获取不到container，就使用key作为错误信息
		return
	}

	lang := GetLang(c)
	message := container.GetI18n().Get(messageKey, lang)

	c.JSON(http.StatusOK, I18nResponse{
		Code:    statusCode,
		Message: message,
		Data:    data,
	})
}

// ValidateErrorI18n 带国际化的验证错误响应，code为400
func ValidateErrorI18n(c *gin.Context, messageKey string, data ...interface{}) {
	var dataValue interface{} = nil
	if len(data) > 0 {
		dataValue = data[0]
	}
	ErrorWithI18n(c, messageKey, dataValue, 400)
}

// AuthErrorI18n 带国际化的认证错误响应，code为401
func AuthErrorI18n(c *gin.Context, messageKey string, data ...interface{}) {
	var dataValue interface{} = nil
	if len(data) > 0 {
		dataValue = data[0]
	}
	ErrorWithI18n(c, messageKey, dataValue, 401)
}

// ForbiddenErrorI18n 带国际化的权限错误响应，code为403
func ForbiddenErrorI18n(c *gin.Context, messageKey string, data ...interface{}) {
	var dataValue interface{} = nil
	if len(data) > 0 {
		dataValue = data[0]
	}
	ErrorWithI18n(c, messageKey, dataValue, 403)
}

// NotFoundErrorI18n 带国际化的未找到错误响应，code为404
func NotFoundErrorI18n(c *gin.Context, messageKey string, data ...interface{}) {
	var dataValue interface{} = nil
	if len(data) > 0 {
		dataValue = data[0]
	}
	ErrorWithI18n(c, messageKey, dataValue, 404)
}
