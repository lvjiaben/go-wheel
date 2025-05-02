package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Response 响应结构体
type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Success 成功响应
func Success(c *gin.Context, data interface{}) {
	c.JSON(http.StatusOK, Response{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// Error 错误响应
func Error(c *gin.Context, code int, message string) {
	c.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

// ValidateError 验证错误响应
func ValidateError(c *gin.Context, message string) {
	Error(c, 400, message)
}

// AuthError 认证错误响应
func AuthError(c *gin.Context, message string) {
	Error(c, 401, message)
}

// ForbiddenError 权限错误响应
func ForbiddenError(c *gin.Context, message string) {
	Error(c, 403, message)
}

// NotFoundError 未找到错误响应
func NotFoundError(c *gin.Context, message string) {
	Error(c, 404, message)
}

// ServerError 服务器错误响应
func ServerError(c *gin.Context, message string) {
	Error(c, 500, message)
}
