package httpclient

import "encoding/base64"

// base64Encode Base64 编码
func base64Encode(data []byte) string {
	return base64.StdEncoding.EncodeToString(data)
}

// base64Decode Base64 解码
func base64Decode(data string) ([]byte, error) {
	return base64.StdEncoding.DecodeString(data)
}

