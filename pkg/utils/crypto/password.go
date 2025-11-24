package crypto

import (
	"github.com/lvjiaben/go-wheel/pkg/constants"
	"golang.org/x/crypto/bcrypt"
)

// PasswordHashWithSalt 使用盐值加密密码
func PasswordHashWithSalt(password, salt string) (string, error) {
	// 将密码和盐值组合
	saltedPassword := password + salt
	hashedBytes, err := bcrypt.GenerateFromPassword([]byte(saltedPassword), constants.BcryptCost)
	if err != nil {
		return "", err
	}
	return string(hashedBytes), nil
}

// HashPassword 使用bcrypt和盐值加密密码（别名，保持兼容性）
func HashPassword(password string, salt string) (string, error) {
	return PasswordHashWithSalt(password, salt)
}

// PasswordVerifyWithSalt 验证带盐值的密码
func PasswordVerifyWithSalt(password, salt, hashedPassword string) bool {
	// 将密码和盐值组合
	saltedPassword := password + salt
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(saltedPassword))
	return err == nil
}
