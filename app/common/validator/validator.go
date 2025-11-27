package validator

import (
	"bytes"
	"encoding/json"
	"io"
	"reflect"
	"strconv"
	"strings"

	"github.com/lvjiaben/go-wheel/pkg/utils/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// 自动转换前端的数据结构类型
func ValidateStructWithConvert[T any](c *gin.Context) (*T, bool) {
	// 读取原始 Body（只读一次）
	bodyBytes, err := io.ReadAll(c.Request.Body)
	if err != nil {
		http.ErrorWithI18n(c, "common.paramError", nil)
		return nil, false
	}

	// 解析为 map
	var rawData map[string]interface{}
	if err := json.Unmarshal(bodyBytes, &rawData); err != nil {
		http.ErrorWithI18n(c, "common.paramError", nil)
		return nil, false
	}

	// 获取结构体类型信息
	var t T
	structType := reflect.TypeOf(t)

	// 自动转换结构
	convertedData := convertByStructType(rawData, structType)

	// 重新序列化
	jsonBytes, _ := json.Marshal(convertedData)

	// 设置新的 Body
	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))

	return ValidateStruct[T](c)
}

// ValidateStruct 通用验证函数，接受任何结构体类型
// T 为要验证的结构体类型
// 自动根据请求方法选择绑定方式：GET/DELETE 使用 Query，POST/PUT/PATCH 使用 JSON
func ValidateStruct[T any](c *gin.Context) (*T, bool) {
	var data T
	var err error

	// 根据请求方法选择绑定方式
	if c.Request.Method == "GET" || c.Request.Method == "DELETE" {
		err = c.ShouldBindQuery(&data)
	} else {
		err = c.ShouldBindJSON(&data)
	}

	if err != nil {
		validationErrors, ok := err.(validator.ValidationErrors)
		if !ok {
			// 非验证错误，如JSON解析错误等
			http.ErrorWithI18n(c, "common.invalid_params", nil)
			return nil, false
		}

		// 处理验证错误
		t := reflect.TypeOf(data)
		for _, e := range validationErrors {
			// 获取字段的msg标签
			field := e.Field()
			structField, exists := t.FieldByName(field)
			if !exists {
				http.ErrorWithI18n(c, "common.invalid_params", nil)
				return nil, false
			}

			// 获取msg标签值，如果不存在则使用通用错误消息
			msgKey := structField.Tag.Get("msg")
			if msgKey == "" {
				msgKey = "common.invalid_params"
			}

			// 确保container已设置到上下文中
			if _, exists := c.Get("container"); !exists {
				// 如果controller没有设置container，直接返回错误码和消息键（不翻译）
				c.JSON(200, gin.H{
					"code":    400,
					"message": "请求参数无效，需要填写" + field + "字段", // 提供一个通用错误信息
					"data":    nil,
				})
				return nil, false
			}

			// 返回特定的错误信息
			http.ErrorWithI18n(c, msgKey, nil)
			return nil, false
		}
	}

	return &data, true
}

// GetFieldLabel 根据字段名获取标签值
func GetFieldLabel(s interface{}, fieldName string) string {
	t := reflect.TypeOf(s)
	if field, exists := t.FieldByName(fieldName); exists {
		return field.Tag.Get("label")
	}
	return ""
}

// GetFieldMsg 根据字段名获取错误消息
func GetFieldMsg(s interface{}, fieldName string) string {
	t := reflect.TypeOf(s)
	if field, exists := t.FieldByName(fieldName); exists {
		return field.Tag.Get("msg")
	}
	return ""
}

// convertByStructType 根据结构体定义的类型智能转换
func convertByStructType(data map[string]interface{}, structType reflect.Type) map[string]interface{} {
	result := make(map[string]interface{})

	// 构建 json tag 到字段类型的映射
	fieldTypeMap := make(map[string]reflect.Type)
	for i := 0; i < structType.NumField(); i++ {
		field := structType.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" || jsonTag == "-" {
			continue
		}

		// 处理 json tag 中的选项（如 `json:"name,omitempty"`）
		jsonName := jsonTag
		if idx := strings.Index(jsonTag, ","); idx != -1 {
			jsonName = jsonTag[:idx]
		}

		fieldTypeMap[jsonName] = field.Type
	}

	// 遍历数据，根据结构体定义的类型进行转换
	for key, value := range data {
		// 获取该字段在结构体中定义的类型
		targetType, exists := fieldTypeMap[key]
		if !exists {
			// 结构体中没有定义这个字段，保持原样
			result[key] = value
			continue
		}

		// 根据目标类型进行转换
		result[key] = convertToTargetType(value, targetType)
	}

	return result
}

// convertToTargetType 将值转换为目标类型
func convertToTargetType(value interface{}, targetType reflect.Type) interface{} {
	if value == nil {
		return getZeroValue(targetType.Kind())
	}

	// 获取目标类型的 Kind
	targetKind := targetType.Kind()

	switch v := value.(type) {
	case string:
		return convertStringToType(v, targetKind)

	case float64:
		// JSON 数字默认是 float64
		return convertFloat64ToType(v, targetKind)

	case bool:
		return convertBoolToType(v, targetKind)

	case int:
		// 添加 int 类型的处理
		return convertIntToType(int64(v), targetKind)

	case int64:
		// 添加 int64 类型的处理
		return convertIntToType(v, targetKind)

	default:
		return v
	}
}

// getZeroValue 获取类型的零值
func getZeroValue(kind reflect.Kind) interface{} {
	switch kind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return uint(0)
	case reflect.Float32, reflect.Float64:
		return 0.0
	case reflect.Bool:
		return false
	case reflect.String:
		return ""
	default:
		return nil
	}
}

// 添加新函数：convertIntToType
func convertIntToType(num int64, targetKind reflect.Kind) interface{} {
	switch targetKind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return convertIntToKind(num, targetKind)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return convertUintToKind(uint64(num), targetKind)

	case reflect.Float32:
		return float32(num)

	case reflect.Float64:
		return float64(num)

	case reflect.Bool:
		return num != 0

	case reflect.String:
		return strconv.FormatInt(num, 10) // ← int 转 string
	}

	return num
}

// convertStringToType 将字符串转换为目标类型
func convertStringToType(str string, targetKind reflect.Kind) interface{} {
	if str == "" {
		return getZeroValue(targetKind)
	}

	switch targetKind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if num, err := strconv.ParseInt(str, 10, 64); err == nil {
			return convertIntToKind(num, targetKind)
		}

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if num, err := strconv.ParseUint(str, 10, 64); err == nil {
			return convertUintToKind(num, targetKind)
		}

	case reflect.Float32, reflect.Float64:
		if num, err := strconv.ParseFloat(str, 64); err == nil {
			if targetKind == reflect.Float32 {
				return float32(num)
			}
			return num
		}

	case reflect.Bool:
		switch str {
		case "true", "True", "TRUE", "1", "yes", "Yes", "YES":
			return true
		case "false", "False", "FALSE", "0", "no", "No", "NO":
			return false
		}

	case reflect.String:
		return str
	}

	return str
}

// convertFloat64ToType 将 float64 转换为目标类型
func convertFloat64ToType(num float64, targetKind reflect.Kind) interface{} {
	switch targetKind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return convertIntToKind(int64(num), targetKind)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return convertUintToKind(uint64(num), targetKind)

	case reflect.Float32:
		return float32(num)

	case reflect.Float64:
		return num

	case reflect.Bool:
		return num != 0

	case reflect.String:
		return strconv.FormatFloat(num, 'f', -1, 64)
	}

	return num
}

// convertBoolToType 将 bool 转换为目标类型
func convertBoolToType(b bool, targetKind reflect.Kind) interface{} {
	switch targetKind {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		if b {
			return convertIntToKind(1, targetKind)
		}
		return convertIntToKind(0, targetKind)

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		if b {
			return convertUintToKind(1, targetKind)
		}
		return convertUintToKind(0, targetKind)

	case reflect.Float32, reflect.Float64:
		if b {
			if targetKind == reflect.Float32 {
				return float32(1)
			}
			return float64(1)
		}
		if targetKind == reflect.Float32 {
			return float32(0)
		}
		return float64(0)

	case reflect.Bool:
		return b

	case reflect.String:
		if b {
			return "true"
		}
		return "false"
	}

	return b
}

// convertIntToKind 将 int64 转换为具体的 int 类型
func convertIntToKind(num int64, kind reflect.Kind) interface{} {
	switch kind {
	case reflect.Int:
		return int(num)
	case reflect.Int8:
		return int8(num)
	case reflect.Int16:
		return int16(num)
	case reflect.Int32:
		return int32(num)
	case reflect.Int64:
		return num
	default:
		return int(num)
	}
}

// convertUintToKind 将 uint64 转换为具体的 uint 类型
func convertUintToKind(num uint64, kind reflect.Kind) interface{} {
	switch kind {
	case reflect.Uint:
		return uint(num)
	case reflect.Uint8:
		return uint8(num)
	case reflect.Uint16:
		return uint16(num)
	case reflect.Uint32:
		return uint32(num)
	case reflect.Uint64:
		return num
	default:
		return uint(num)
	}
}
