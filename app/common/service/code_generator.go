package service

import (
	"fmt"

	"github.com/lvjiaben/go-wheel/pkg/constants"
	"github.com/lvjiaben/go-wheel/pkg/container"
	"github.com/lvjiaben/go-wheel/pkg/utils/crypto"
)

// CodeGeneratorService 代码生成服务（邀请码、验证码等）
type CodeGeneratorService struct {
	container *container.Container
}

// NewCodeGeneratorService 创建代码生成服务
func NewCodeGeneratorService(c *container.Container) *CodeGeneratorService {
	return &CodeGeneratorService{
		container: c,
	}
}

// GenerateInviteCode 生成邀请码（使用密码学安全的随机数生成器）
// 返回格式：10位大写字母+数字组合
func (s *CodeGeneratorService) GenerateInviteCode() string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	length := constants.InviteCodeLength

	code := make([]byte, length)
	randomBytes, err := crypto.GenerateRandomBytes(length)
	if err != nil {
		// 如果生成失败，记录错误并返回时间戳作为后备方案
		s.container.GetLogger().Error("生成随机邀请码失败: " + err.Error())
		return fmt.Sprintf("INV%d", crypto.GenerateTimestampCode())
	}

	for i := range code {
		code[i] = charset[randomBytes[i]%byte(len(charset))]
	}

	return string(code)
}

// GenerateUniqueInviteCode 生成唯一邀请码（检查数据库去重）
// tableName: 表名（如 "user"）
// columnName: 列名（如 "code"）
// maxRetries: 最大重试次数，0 表示使用默认值
func (s *CodeGeneratorService) GenerateUniqueInviteCode(tableName, columnName string, maxRetries int) string {
	if maxRetries == 0 {
		maxRetries = constants.MaxInviteCodeRetries
	}

	for i := 0; i < maxRetries; i++ {
		code := s.GenerateInviteCode()
		
		// 检查是否重复
		var count int64
		s.container.GetDB().Table(tableName).Where(columnName+" = ?", code).Count(&count)
		
		if count == 0 {
			return code
		}
	}

	// 如果重试多次仍然失败，使用时间戳
	s.container.GetLogger().Warn(fmt.Sprintf("邀请码生成重试 %d 次后仍重复，使用时间戳后备方案", maxRetries))
	return fmt.Sprintf("INV%d", crypto.GenerateTimestampCode())
}

// GenerateRandomPassword 生成随机密码（使用密码学安全的随机数生成器）
// 返回格式：12位字母+数字+特殊字符组合
func (s *CodeGeneratorService) GenerateRandomPassword() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*"
	length := constants.RandomPasswordLength

	password := make([]byte, length)
	randomBytes, err := crypto.GenerateRandomBytes(length)
	if err != nil {
		// 如果生成失败，记录错误并返回时间戳作为后备方案
		s.container.GetLogger().Error("生成随机密码失败: " + err.Error())
		return fmt.Sprintf("Pass%d!", crypto.GenerateTimestampCode())
	}

	for i := range password {
		password[i] = charset[randomBytes[i]%byte(len(charset))]
	}

	return string(password)
}

// GenerateNumericCode 生成纯数字验证码
// length: 验证码长度
func (s *CodeGeneratorService) GenerateNumericCode(length int) string {
	const charset = "0123456789"
	
	if length <= 0 {
		length = 6 // 默认6位
	}

	code := make([]byte, length)
	randomBytes, err := crypto.GenerateRandomBytes(length)
	if err != nil {
		// 如果生成失败，记录错误并返回时间戳的后几位
		s.container.GetLogger().Error("生成数字验证码失败: " + err.Error())
		timestamp := crypto.GenerateTimestampCode()
		return fmt.Sprintf("%0*d", length, timestamp%int64(pow(10, length)))
	}

	for i := range code {
		code[i] = charset[randomBytes[i]%byte(len(charset))]
	}

	return string(code)
}

// GenerateAlphanumericCode 生成字母数字混合验证码
// length: 验证码长度
// uppercase: 是否全部大写
func (s *CodeGeneratorService) GenerateAlphanumericCode(length int, uppercase bool) string {
	var charset string
	if uppercase {
		charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	} else {
		charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	}
	
	if length <= 0 {
		length = 6 // 默认6位
	}

	code := make([]byte, length)
	randomBytes, err := crypto.GenerateRandomBytes(length)
	if err != nil {
		// 如果生成失败，记录错误并返回时间戳
		s.container.GetLogger().Error("生成字母数字验证码失败: " + err.Error())
		return fmt.Sprintf("CODE%d", crypto.GenerateTimestampCode())
	}

	for i := range code {
		code[i] = charset[randomBytes[i]%byte(len(charset))]
	}

	return string(code)
}

// pow 计算整数幂（辅助函数）
func pow(base, exp int) int64 {
	result := int64(1)
	for i := 0; i < exp; i++ {
		result *= int64(base)
	}
	return result
}

