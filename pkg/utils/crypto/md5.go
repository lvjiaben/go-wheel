package crypto

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
