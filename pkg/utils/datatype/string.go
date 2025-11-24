package datatype

import (
	"math/rand"
	"regexp"
	"strings"
	"time"
	"unicode"
)

// init 初始化随机数种子
func init() {
	rand.Seed(time.Now().UnixNano())
}

// GenerateRandomString 生成指定长度的随机字符串
func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = charset[rand.Intn(len(charset))]
	}
	return string(result)
}

// TruncateString 截断字符串到指定长度，并添加省略号
func TruncateString(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}

// IsEmptyString 检查字符串是否为空
func IsEmptyString(s string) bool {
	return s == ""
}

// Marshal 将下划线命名转换为驼峰命名
func Marshal(name string) string {
	if name == "" {
		return ""
	}

	temp := strings.Split(name, "_")
	var s string
	for i, v := range temp {
		if i == 0 {
			s = v
			continue
		}
		s += strings.Title(v)
	}
	return s
}

// CamelToSnake 将驼峰命名转换为下划线命名
func CamelToSnake(s string) string {
	var result string
	for i, r := range s {
		if i > 0 && unicode.IsUpper(r) {
			result += "_"
		}
		result += strings.ToLower(string(r))
	}
	return result
}

// ToSnakeCase 驼峰转下划线
func ToSnakeCase(str string) string {
	var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
	var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

// ToCamelCase 下划线转驼峰
func ToCamelCase(str string) string {
	str = strings.Replace(str, "_", " ", -1)
	str = strings.Title(str)
	return strings.Replace(str, " ", "", -1)
}

// Contains 检查字符串是否在切片中
func Contains(slice []string, str string) bool {
	for _, v := range slice {
		if v == str {
			return true
		}
	}
	return false
}
