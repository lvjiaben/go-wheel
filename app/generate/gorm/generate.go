package gorm

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lvjiaben/go-wheel/pkg/container"
	"github.com/lvjiaben/go-wheel/pkg/utils"
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
	return utils.Marshal(field.Field)
}

// Generator 生成器
type Generator struct {
	container *container.Container
}

// NewGenerator 创建生成器
func NewGenerator(c *container.Container) *Generator {
	return &Generator{
		container: c,
	}
}

// Generate 生成模型
func (g *Generator) Generate() error {
	// 获取数据库连接
	db := g.container.GetDB()
	if db == nil {
		return fmt.Errorf("数据库连接为空")
	}

	// 获取所有表名
	var tables []string
	if err := db.Raw("SHOW TABLES").Scan(&tables).Error; err != nil {
		return fmt.Errorf("获取表名失败: %v", err)
	}

	// 创建 models 目录
	modelsDir := filepath.Join("app", "backend", "model")
	if err := os.MkdirAll(modelsDir, 0755); err != nil {
		return fmt.Errorf("创建 models 目录失败: %v", err)
	}

	// 为每个表生成模型
	for _, table := range tables {
		// 获取表结构
		var columns []struct {
			Field   string
			Type    string
			Null    string
			Key     string
			Default string
			Extra   string
		}
		if err := db.Raw(fmt.Sprintf("DESCRIBE %s", table)).Scan(&columns).Error; err != nil {
			return fmt.Errorf("获取表 %s 结构失败: %v", table, err)
		}

		// 生成模型代码
		modelCode := g.generateModelCode(table, columns)

		// 写入文件
		modelFile := filepath.Join(modelsDir, fmt.Sprintf("%s.go", strings.ToLower(table)))
		if err := os.WriteFile(modelFile, []byte(modelCode), 0644); err != nil {
			return fmt.Errorf("写入模型文件失败: %v", err)
		}
	}

	return nil
}

func (g *Generator) generateModelCode(table string, columns []struct {
	Field   string
	Type    string
	Null    string
	Key     string
	Default string
	Extra   string
}) string {
	// 生成模型代码
	code := fmt.Sprintf("package model\n\n")
	code += fmt.Sprintf("type %s struct {\n", strings.Title(table))

	for _, col := range columns {
		fieldType := g.getGoType(col.Type)
		fieldName := strings.Title(col.Field)
		code += fmt.Sprintf("\t%s %s `gorm:\"column:%s\"`\n", fieldName, fieldType, col.Field)
	}

	code += "}\n\n"
	code += fmt.Sprintf("func (%s) TableName() string {\n", strings.Title(table))
	code += fmt.Sprintf("\treturn \"%s\"\n", table)
	code += "}\n"

	return code
}

func (g *Generator) getGoType(mysqlType string) string {
	switch {
	case strings.Contains(mysqlType, "int"):
		return "int"
	case strings.Contains(mysqlType, "varchar"), strings.Contains(mysqlType, "text"):
		return "string"
	case strings.Contains(mysqlType, "datetime"), strings.Contains(mysqlType, "timestamp"):
		return "time.Time"
	case strings.Contains(mysqlType, "decimal"), strings.Contains(mysqlType, "float"):
		return "float64"
	default:
		return "string"
	}
}

// Column 列信息
type Column struct {
	Name          string `gorm:"column:COLUMN_NAME"`
	DataType      string `gorm:"column:DATA_TYPE"`
	Comment       string `gorm:"column:COLUMN_COMMENT"`
	IsNullable    string `gorm:"column:IS_NULLABLE"`
	ColumnKey     string `gorm:"column:COLUMN_KEY"`
	ColumnDefault string `gorm:"column:COLUMN_DEFAULT"`
}
