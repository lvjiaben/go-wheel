package errors

import (
	"fmt"
	"net/http"
)

// 错误码定义
const (
	CodeSuccess = 0

	// 通用错误码 1000-1999
	CodeBadRequest   = 1001
	CodeUnauthorized = 1002
	CodeForbidden    = 1003
	CodeNotFound     = 1004
	CodeServerError  = 1005

	// 业务错误码 2000-2999
	CodeUserNotFound      = 2001
	CodeUsernameExists    = 2002
	CodeEmailExists       = 2003
	CodePasswordIncorrect = 2004
	CodeTokenExpired      = 2005
	CodeTokenInvalid      = 2006
)

// Error 自定义错误类型
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Detail  string `json:"detail,omitempty"`
}

// New 创建新的错误
func New(code int, message string) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}

// Error 实现error接口
func (e *Error) Error() string {
	return fmt.Sprintf("code: %d, message: %s", e.Code, e.Message)
}

// WithDetail 添加错误详情
func (e *Error) WithDetail(detail string) *Error {
	e.Detail = detail
	return e
}

// HTTPStatus 获取HTTP状态码
func (e *Error) HTTPStatus() int {
	switch e.Code {
	case CodeSuccess:
		return http.StatusOK
	case CodeBadRequest:
		return http.StatusBadRequest
	case CodeUnauthorized:
		return http.StatusUnauthorized
	case CodeForbidden:
		return http.StatusForbidden
	case CodeNotFound:
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}

// 预定义错误
var (
	ErrBadRequest   = New(CodeBadRequest, "请求参数错误")
	ErrUnauthorized = New(CodeUnauthorized, "未授权")
	ErrForbidden    = New(CodeForbidden, "禁止访问")
	ErrNotFound     = New(CodeNotFound, "资源不存在")
	ErrServerError  = New(CodeServerError, "服务器内部错误")

	ErrUserNotFound      = New(CodeUserNotFound, "用户不存在")
	ErrUsernameExists    = New(CodeUsernameExists, "用户名已存在")
	ErrEmailExists       = New(CodeEmailExists, "邮箱已存在")
	ErrPasswordIncorrect = New(CodePasswordIncorrect, "密码错误")
	ErrTokenExpired      = New(CodeTokenExpired, "令牌已过期")
	ErrTokenInvalid      = New(CodeTokenInvalid, "令牌无效")
)
