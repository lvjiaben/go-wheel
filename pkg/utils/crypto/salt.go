package crypto

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/lvjiaben/go-wheel/pkg/constants"
)

// GenerateSalt 生成随机盐值
func GenerateSalt() (string, error) {
	bytes := make([]byte, constants.SaltLength)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// GenerateRandomBytes 生成指定长度的随机字节
func GenerateRandomBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// GenerateTimestampCode 生成基于时间戳的唯一码
func GenerateTimestampCode() int64 {
	return time.Now().UnixNano()
}
