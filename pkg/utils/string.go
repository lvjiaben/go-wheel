package utils

import (
	"regexp"
	"strings"
	"unicode"
)

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

// IsStringInSlice 检查字符串是否在切片中
func IsStringInSlice(item string, slice []string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// IsIntInSlice 检查整数是否在切片中
func IsIntInSlice(item int, slice []int) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

// IntToBool 将整数转换为布尔值
func IntToBool(i int) bool {
	return i != 0
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

// ContainsInt 检查整数是否在切片中
func ContainsInt(slice []int, num int) bool {
	for _, v := range slice {
		if v == num {
			return true
		}
	}
	return false
}
