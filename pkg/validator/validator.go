package validator

import (
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Validator 自定义验证器
type Validator struct {
	validate *validator.Validate
}

// New 创建新的验证器
func New() *Validator {
	v := validator.New()
	_ = v.RegisterValidation("username", validateUsername)
	_ = v.RegisterValidation("password", validatePassword)
	_ = v.RegisterValidation("email", validateEmail)
	return &Validator{validate: v}
}

// Validate 验证结构体
func (v *Validator) Validate(s interface{}) error {
	return v.validate.Struct(s)
}

// validateUsername 验证用户名
func validateUsername(fl validator.FieldLevel) bool {
	username := fl.Field().String()
	if len(username) < 3 || len(username) > 20 {
		return false
	}
	// 只允许字母、数字和下划线
	matched, _ := regexp.MatchString("^[a-zA-Z0-9_]+$", username)
	return matched
}

// validatePassword 验证密码
func validatePassword(fl validator.FieldLevel) bool {
	password := fl.Field().String()
	if len(password) < 6 || len(password) > 20 {
		return false
	}
	// 必须包含字母和数字
	hasLetter := false
	hasDigit := false
	for _, c := range password {
		if c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' {
			hasLetter = true
		}
		if c >= '0' && c <= '9' {
			hasDigit = true
		}
	}
	return hasLetter && hasDigit
}

// validateEmail 验证邮箱
func validateEmail(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	if len(email) < 5 || len(email) > 100 {
		return false
	}
	// 简单的邮箱格式验证
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return false
	}
	if len(parts[0]) == 0 || len(parts[1]) == 0 {
		return false
	}
	return true
}
