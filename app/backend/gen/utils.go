package gen

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// ToPascalCase 转换为 PascalCase（大驼峰）
func ToPascalCase(s string) string {
	// 处理下划线分隔的字符串
	words := strings.Split(s, "_")
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + strings.ToLower(word[1:])
		}
	}
	return strings.Join(words, "")
}

// ToCamelCase 转换为 camelCase（小驼峰）
func ToCamelCase(s string) string {
	pascal := ToPascalCase(s)
	if len(pascal) > 0 {
		return strings.ToLower(pascal[:1]) + pascal[1:]
	}
	return pascal
}

// ToSnakeCase 转换为 snake_case（下划线）
func ToSnakeCase(s string) string {
	var result []rune
	for i, r := range s {
		if unicode.IsUpper(r) {
			if i > 0 {
				result = append(result, '_')
			}
			result = append(result, unicode.ToLower(r))
		} else {
			result = append(result, r)
		}
	}
	return string(result)
}

// ConvertDBTypeToGoType 将数据库类型转换为 Go 类型
func ConvertDBTypeToGoType(dbType string, nullable bool) string {
	dbType = strings.ToLower(dbType)

	// 整数类型
	if strings.Contains(dbType, "tinyint(1)") {
		return "int8" // 布尔类型使用 int8
	}
	if strings.Contains(dbType, "tinyint") {
		return "int8"
	}
	if strings.Contains(dbType, "smallint") {
		return "int16"
	}
	if strings.Contains(dbType, "mediumint") || strings.Contains(dbType, "int") {
		return "int"
	}
	if strings.Contains(dbType, "bigint") {
		return "int64"
	}

	// 浮点类型
	if strings.Contains(dbType, "float") {
		return "float32"
	}
	if strings.Contains(dbType, "double") || strings.Contains(dbType, "decimal") || strings.Contains(dbType, "numeric") {
		return "float64"
	}

	// 字符串类型
	if strings.Contains(dbType, "char") || strings.Contains(dbType, "text") || strings.Contains(dbType, "enum") || strings.Contains(dbType, "set") {
		return "string"
	}

	// 时间类型
	if strings.Contains(dbType, "date") || strings.Contains(dbType, "time") {
		return "int" // 使用时间戳
	}

	// 二进制类型
	if strings.Contains(dbType, "blob") || strings.Contains(dbType, "binary") {
		return "[]byte"
	}

	// 默认返回 string
	return "string"
}

// BuildGormTag 构建 GORM 标签
func BuildGormTag(column ColumnInfo) string {
	var tags []string

	// 主键
	if column.ColumnKey == "PRI" {
		tags = append(tags, "primaryKey")
	}

	// 自增
	if strings.Contains(column.Extra, "auto_increment") {
		tags = append(tags, "autoIncrement")
	}

	// 类型
	if column.DataType != "" {
		// 特殊类型处理
		if strings.HasPrefix(column.DataType, "varchar") || strings.HasPrefix(column.DataType, "char") {
			if column.CharacterMaxLength > 0 {
				tags = append(tags, fmt.Sprintf("size:%d", column.CharacterMaxLength))
			}
		} else if strings.HasPrefix(column.DataType, "decimal") {
			tags = append(tags, fmt.Sprintf("type:%s", column.ColumnType))
		} else if column.DataType == "text" || column.DataType == "longtext" {
			tags = append(tags, fmt.Sprintf("type:%s", column.DataType))
		}
	}

	// 非空
	if column.IsNullable == "NO" && column.ColumnKey != "PRI" && !strings.Contains(column.Extra, "auto_increment") {
		tags = append(tags, "not null")
	}

	// 默认值
	if column.ColumnDefault != "" && column.ColumnDefault != "NULL" {
		// 清理默认值
		defaultValue := strings.Trim(column.ColumnDefault, "'\"")
		tags = append(tags, fmt.Sprintf("default:%s", defaultValue))
	}

	// 索引
	if column.ColumnKey == "UNI" {
		tags = append(tags, "uniqueIndex")
	} else if column.ColumnKey == "MUL" {
		tags = append(tags, "index")
	}

	// 自动时间戳
	if column.ColumnName == "created_at" {
		tags = append(tags, "autoCreateTime")
	} else if column.ColumnName == "updated_at" {
		tags = append(tags, "autoUpdateTime")
	}

	return strings.Join(tags, ";")
}

// BuildValidateRules 构建验证规则
func BuildValidateRules(field FieldConfig) string {
	var rules []string

	// 必填
	if field.IsRequired {
		rules = append(rules, "required")
	}

	// 字符串长度
	if field.FieldType == "string" && field.ColumnType != "" {
		// 提取长度
		re := regexp.MustCompile(`\((\d+)\)`)
		matches := re.FindStringSubmatch(field.ColumnType)
		if len(matches) > 1 {
			maxLen := matches[1]
			if field.IsRequired {
				rules = append(rules, fmt.Sprintf("min=1,max=%s", maxLen))
			} else {
				rules = append(rules, fmt.Sprintf("omitempty,max=%s", maxLen))
			}
		}
	}

	// 数字范围
	if strings.Contains(field.FieldType, "int") {
		if field.IsRequired {
			rules = append(rules, "min=0")
		}
	}

	// 邮箱
	if strings.Contains(field.ColumnName, "email") {
		rules = append(rules, "omitempty,email")
	}

	// 手机号
	if strings.Contains(field.ColumnName, "mobile") || strings.Contains(field.ColumnName, "phone") {
		rules = append(rules, "omitempty,len=11")
	}

	// 枚举
	if field.IsEnumField || field.IsOperateField {
		// 从 column_type 中提取枚举值
		if strings.Contains(field.ColumnType, "enum") {
			re := regexp.MustCompile(`enum\((.*?)\)`)
			matches := re.FindStringSubmatch(field.ColumnType)
			if len(matches) > 1 {
				values := strings.Split(strings.ReplaceAll(matches[1], "'", ""), ",")
				rules = append(rules, fmt.Sprintf("oneof=%s", strings.Join(values, " ")))
			}
		} else if field.ColumnName == "status" {
			rules = append(rules, "oneof=0 1")
		}
	}

	return strings.Join(rules, ",")
}

// IsOperateField 判断是否为 Operate 字段
func IsOperateField(columnName string) bool {
	operateFields := []string{"status", "is_", "enable", "disable", "visible", "hidden"}
	for _, field := range operateFields {
		if columnName == field || strings.HasPrefix(columnName, field) {
			return true
		}
	}
	return false
}

// IsSortField 判断是否为排序字段
func IsSortField(columnName string) bool {
	sortFields := []string{"sort", "weigh", "order", "sequence"}
	for _, field := range sortFields {
		if columnName == field || strings.Contains(columnName, field) {
			return true
		}
	}
	return false
}

// IsTimeField 判断是否为时间字段
func IsTimeField(columnName string) bool {
	timeFields := []string{"_at", "_time", "date"}
	for _, field := range timeFields {
		if strings.HasSuffix(columnName, field) || strings.Contains(columnName, field) {
			return true
		}
	}
	return false
}

// IsRelationField 判断是否为关联字段
func IsRelationField(columnName string) bool {
	return strings.HasSuffix(columnName, "_id") && columnName != "id"
}

// IsTextField 判断是否为长文本字段
func IsTextField(dataType string) bool {
	textTypes := []string{"text", "longtext", "mediumtext"}
	for _, t := range textTypes {
		if dataType == t {
			return true
		}
	}
	return false
}

// IsBoolField 判断是否为布尔字段
func IsBoolField(dataType string) bool {
	return strings.Contains(dataType, "tinyint(1)")
}

// IsImageField 判断是否为单图字段
func IsImageField(columnName string) bool {
	imageFields := []string{"image", "avatar", "img", "icon", "logo", "cover", "thumb"}
	for _, field := range imageFields {
		if columnName == field || strings.HasSuffix(columnName, field) {
			// 排除复数形式
			if !strings.HasSuffix(columnName, "s") {
				return true
			}
		}
	}
	return false
}

// IsImagesField 判断是否为多图字段
func IsImagesField(columnName string) bool {
	imageFields := []string{"images", "avatars", "imgs", "icons", "logos", "covers", "thumbs", "gallery", "photos"}
	for _, field := range imageFields {
		if columnName == field || strings.HasSuffix(columnName, field) {
			return true
		}
	}
	return false
}

// GetDefaultTableDisplayType 获取默认的表格显示类型
func GetDefaultTableDisplayType(field FieldConfig) string {
	if field.IsOperateField {
		return "tag" // Tag 标签
	}
	if field.IsSortField {
		return "editable" // 双击编辑
	}
	if field.IsTimeField {
		return "datetime" // 时间格式化
	}
	if field.IsImageField {
		return "image" // 图片
	}
	if field.IsImagesField {
		return "images" // 多图
	}
	return "text" // 普通文本
}

// GetDefaultFormComponent 获取默认的表单组件
func GetDefaultFormComponent(field FieldConfig) string {
	if field.IsImageField {
		return "AttachmentInput" // 单图上传
	}
	if field.IsImagesField {
		return "AttachmentInput" // 多图上传
	}
	if field.IsTextField {
		return "RichEditor" // 富文本编辑器
	}
	if field.IsBoolField || (field.IsOperateField && field.FieldType == "int8") {
		return "Switch" // 开关
	}
	if field.IsEnumField {
		return "Radio" // 单选
	}
	if field.IsSetField {
		return "Select" // 下拉多选
	}
	if field.IsRelationField {
		return "TableSelect" // 表格选择
	}
	if field.IsTimeField {
		return "RangePicker" // 时间范围选择
	}
	if strings.Contains(field.FieldType, "int") || strings.Contains(field.FieldType, "float") {
		return "InputNumber" // 数字输入
	}
	return "Input" // 默认文本输入
}

// GetDefaultSearchFormComponent 获取默认的搜索表单组件
func GetDefaultSearchFormComponent(field FieldConfig) string {
	// status, is_xxx 等布尔类字段在搜索时使用 Select
	if field.IsBoolField || field.IsOperateField {
		return "Select"
	}
	// 其他情况与表单组件一致
	return GetDefaultFormComponent(field)
}
