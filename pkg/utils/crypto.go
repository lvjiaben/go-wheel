package utils

import (
	"crypto/md5"
	"encoding/hex"
)

// MD5 计算MD5值
func MD5(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}

// PasswordHash 密码加密
func PasswordHash(password string) string {
	return MD5(password)
}

// PasswordVerify 密码验证
func PasswordVerify(password, hash string) bool {
	return MD5(password) == hash
}
