package gorm

import (
	"fmt"
	"strings"

	"github.com/lvjiaben/go-wheel/pkg/actions"
	"github.com/lvjiaben/go-wheel/pkg/global"
	"github.com/lvjiaben/go-wheel/pkg/initialize"
	"github.com/lvjiaben/go-wheel/pkg/utils/file"
	"gorm.io/gorm"
)

// Table 表结构
type Table struct {
	Name    string
	Comment string
}

// Field 字段结构
type Field struct {
	Field      string // 字段名
	Type       string // 字段类型
	Null       string // 是否可为空
	Key        string // 键类型
	Default    string // 默认值
	Extra      string // 额外信息
	Privileges string // 权限
	Comment    string // 注释
}

// GetTables 获取表结构
func GetTables(db *gorm.DB, tableNames []string, dbName string) []Table {
	var tables []Table
	query := "SELECT TABLE_NAME, TABLE_COMMENT FROM information_schema.TABLES WHERE TABLE_SCHEMA = ?"
	if len(tableNames) > 0 {
		query += " AND TABLE_NAME IN (?)"
		db = db.Raw(query, dbName, tableNames)
	} else {
		db = db.Raw(query, dbName)
	}
	rows, err := db.Rows()
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var table Table
		rows.Scan(&table.Name, &table.Comment)
		tables = append(tables, table)
	}
	return tables
}

// GetFields 获取字段结构
func GetFields(db *gorm.DB, tableName string) []Field {
	var fields []Field
	query := "SHOW FULL COLUMNS FROM " + tableName
	rows, err := db.Raw(query).Rows()
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		var field Field
		rows.Scan(&field.Field, &field.Type, &field.Null, &field.Key, &field.Default, &field.Extra, &field.Privileges, &field.Comment)
		fields = append(fields, field)
	}
	return fields
}

// GetFiledType 获取字段类型
func GetFiledType(field Field) string {
	switch field.Type {
	case "varchar", "char", "text", "longtext", "mediumtext", "tinytext":
		return "string"
	case "int", "tinyint", "smallint", "mediumint":
		return "int"
	case "bigint":
		return "int64"
	case "float", "double", "decimal":
		return "float64"
	case "datetime", "timestamp":
		return "time.Time"
	case "date":
		return "time.Time"
	case "time":
		return "time.Time"
	default:
		return "string"
	}
}

// GetFieldJson 获取字段json标签
func GetFieldJson(field Field) string {
	return fmt.Sprintf(`json:"%s" gorm:"column:%s"`, field.Field, field.Field)
}

// GetFieldComment 获取字段注释
func GetFieldComment(field Field) string {
	if field.Comment != "" {
		return "// " + field.Comment
	}
	return ""
}

// GetFieldZh 获取字段中文名
func GetFieldZh(field Field) string {
	if field.Comment != "" {
		return field.Comment
	}
	return actions.Marshal(field.Field)
}

// 初始化需要传入一个model
func Genertate(TableName string, PackageName string, Path string, Cover bool) {
	tableNames := strings.Split(TableName, ",")
	initialize.ViperLoad()
	if global.DB != nil {
		db, _ := global.DB.DB()
		defer db.Close()
	}
	db := global.DB
	tables := GetTables(db, tableNames, global.CONFIG.Mysql.Dbname)
	for _, table := range tables {
		fields := GetFields(db, table.Name)
		generateModel(PackageName, Path, Cover, table, fields)
	}
}

// 生成Model
func generateModel(PackageName string, Path string, Cover bool, table Table, fields []Field) {
	var builder strings.Builder
	builder.WriteString("package " + PackageName + "\n\n")

	// 表注释
	if len(table.Comment) > 0 {
		builder.WriteString("// " + table.Comment + "\n")
	}

	// 生成结构体
	builder.WriteString("type " + actions.Marshal(table.Name) + " struct {\n")

	// 文件内容填充
	for _, field := range fields {
		fieldName := field.Field
		/**
		字段名 字段类型 `json:"字段名" gorm:"column:字段名"` //注释
		*/
		builder.WriteString("\t" + actions.Marshal(fieldName) + "\t" + GetFiledType(field) + "\t" +
			"`" + GetFieldJson(field) + "`\t" + GetFieldComment(field) + "\n")
	}
	builder.WriteString("}\n")

	// 函数名称返回自身
	/**
	func (e *结构体名) TableName() string {
	    return 结构体名
	}
	*/
	builder.WriteString("func (e *" + actions.Marshal(table.Name) +
		") TableName() string { \n    return \"" + table.Name + "\"\n}")

	// 文件生成
	fileName := Path + table.Name + ".go"
	file.MakeFile(Path, fileName, builder.String(), Cover)
}
