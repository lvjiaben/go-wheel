package gorm

import (
	"strings"

	"github.com/lvjiaben/go-wheel/pkg/actions"
	"gorm.io/gorm"
)

// 获取具体表单
func GetTables(db *gorm.DB, tableNames []string, dbName string) []Table {

	// 字符串拼接生成表名范围
	tableNamesStr := "'" + strings.Join(tableNames, "','") + "'"
	// 获取指定表信息
	var tables []Table
	if tableNamesStr == "''" {
		db.Raw("SELECT TABLE_NAME as Name,TABLE_COMMENT as Comment FROM information_schema.TABLES " +
			"WHERE table_schema='" + dbName + "';").Find(&tables)
	} else {
		db.Raw("SELECT TABLE_NAME as Name,TABLE_COMMENT as Comment FROM information_schema.TABLES " +
			"WHERE TABLE_NAME IN (" + tableNamesStr + ") AND " +
			"table_schema='" + dbName + "';").Find(&tables)
	}
	return tables
}

// 获取字段的详情信息
func GetFields(db *gorm.DB, tableName string) []Field {
	var fields []Field
	db.Raw("show FULL COLUMNS from " + tableName + ";").Find(&fields)
	return fields
}

// 获取字段json描述
func GetFieldJson(field Field) string {
	return `json:"` + field.Field + `" ` + `gorm:"column:` + field.Field + `"`
}

// 获取字段说明
func GetFieldComment(field Field) string {
	if len(field.Comment) > 0 {
		return "// " + field.Comment
	}
	return ""
}

func GetFieldZh(field Field) string {
	if len(field.Comment) > 0 {
		return field.Comment
	}
	return actions.Marshal(field.Field)
}

// 获取字段类型
func GetFiledType(field Field) string {
	typeArr := strings.Split(field.Type, "(")
	switch typeArr[0] {
	case "int":
		return "int"
	case "integer":
		return "int"
	case "mediumint":
		return "int"
	case "bit":
		return "int"
	case "year":
		return "int"
	case "smallint":
		return "int"
	case "tinyint":
		return "int"
	case "bigint":
		return "int64"
	case "decimal":
		return "float32"
	case "double":
		return "float32"
	case "float":
		return "float32"
	case "real":
		return "float32"
	case "numeric":
		return "float32"
	case "timestamp":
		return "time.Time"
	case "datetime":
		return "time.Time"
	case "time":
		return "time.Time"
	default:
		return "string"
	}
}
