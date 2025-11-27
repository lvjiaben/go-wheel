package gen

import (
	"fmt"
	"strings"

	"gorm.io/gorm"
)

// TableService 表信息服务
type TableService struct {
	db *gorm.DB
}

// NewTableService 创建表信息服务
func NewTableService(db *gorm.DB) *TableService {
	return &TableService{db: db}
}

// GetTableInfo 获取表信息
func (s *TableService) GetTableInfo(tableName string) (*TableInfo, error) {
	// 获取数据库驱动类型
	driver := s.db.Dialector.Name()

	var tableInfo TableInfo
	tableInfo.TableName = tableName

	// 获取表注释
	comment, err := s.getTableComment(tableName, driver)
	if err != nil {
		return nil, err
	}
	tableInfo.TableComment = comment

	// 获取列信息
	columns, err := s.getColumns(tableName, driver)
	if err != nil {
		return nil, err
	}
	tableInfo.Columns = columns

	return &tableInfo, nil
}

// GetAllTables 获取所有表名
func (s *TableService) GetAllTables() ([]string, error) {
	driver := s.db.Dialector.Name()

	var tables []string
	var err error

	switch driver {
	case "mysql":
		// MySQL 查询
		var result []struct {
			TableName string `gorm:"column:TABLE_NAME"`
		}
		err = s.db.Raw("SELECT TABLE_NAME FROM information_schema.TABLES WHERE TABLE_SCHEMA = DATABASE()").Scan(&result).Error
		if err != nil {
			return nil, fmt.Errorf("查询表列表失败: %v", err)
		}
		for _, r := range result {
			tables = append(tables, r.TableName)
		}

	case "postgres", "postgresql":
		// PostgreSQL 查询
		var result []struct {
			TableName string `gorm:"column:tablename"`
		}
		err = s.db.Raw("SELECT tablename FROM pg_tables WHERE schemaname = 'public'").Scan(&result).Error
		if err != nil {
			return nil, fmt.Errorf("查询表列表失败: %v", err)
		}
		for _, r := range result {
			tables = append(tables, r.TableName)
		}

	default:
		return nil, fmt.Errorf("不支持的数据库类型: %s", driver)
	}

	return tables, nil
}

// GetTableConfig 获取表配置（包含智能识别的默认值）
func (s *TableService) GetTableConfig(tableName string) (*GenConfig, error) {
	tableInfo, err := s.GetTableInfo(tableName)
	if err != nil {
		return nil, err
	}

	config := &GenConfig{
		TableName:        tableInfo.TableName,
		TableComment:     tableInfo.TableComment,
		StructName:       ToPascalCase(tableInfo.TableName),
		ModuleName:       tableInfo.TableName,
		PackageName:      tableInfo.TableName,
		FrontendSrcPath:  "vben-admin/apps/web-antd/src",
		DefaultSortField: "id",
		DefaultSortOrder: "desc",
		Methods: MethodConfig{
			List:    true,
			Create:  true,
			Update:  true,
			Delete:  true,
			Operate: false,
		},
		Fields:        []FieldConfig{},
		SearchFields:  []string{},
		OperateFields: []string{},
		MenuConfig: MenuConfig{
			ParentMenuName: "AutoPlay",
			MenuName:       tableInfo.TableComment,
			MenuIcon:       "",
			MenuSort:       50,
		},
	}

	// 转换字段配置（自动应用智能识别）
	var operateFields []string
	hasSortField := false
	hasWeighField := false
	for _, column := range tableInfo.Columns {
		field := ConvertToFieldConfig(column)
		config.Fields = append(config.Fields, field)

		// 自动识别 Operate 字段
		if field.IsOperateField {
			operateFields = append(operateFields, field.ColumnName)
		}

		// 检查是否有 sort 或 weigh 字段
		if column.ColumnName == "sort" {
			hasSortField = true
		}
		if column.ColumnName == "weigh" {
			hasWeighField = true
		}
	}

	// 如果有 Operate 字段，默认生成 Operate 方法
	if len(operateFields) > 0 {
		config.Methods.Operate = true
		config.OperateFields = operateFields
	}

	// 设置默认排序字段（优先使用 sort 或 weigh，否则使用 id）
	if hasSortField {
		config.DefaultSortField = "sort"
	} else if hasWeighField {
		config.DefaultSortField = "weigh"
	} else {
		config.DefaultSortField = "id"
	}

	return config, nil
}

// getTableComment 获取表注释
func (s *TableService) getTableComment(tableName, driver string) (string, error) {
	var comment string

	switch driver {
	case "mysql":
		// MySQL 查询表注释
		var result struct {
			TableComment string `gorm:"column:TABLE_COMMENT"`
		}
		err := s.db.Raw(`
			SELECT TABLE_COMMENT 
			FROM information_schema.TABLES 
			WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?
		`, tableName).Scan(&result).Error
		if err != nil {
			return "", fmt.Errorf("查询表注释失败: %v", err)
		}
		comment = result.TableComment

	case "postgres", "postgresql":
		// PostgreSQL 查询表注释
		var result struct {
			Description string `gorm:"column:description"`
		}
		err := s.db.Raw(`
			SELECT obj_description(oid) as description
			FROM pg_class
			WHERE relname = ? AND relkind = 'r'
		`, tableName).Scan(&result).Error
		if err != nil {
			return "", fmt.Errorf("查询表注释失败: %v", err)
		}
		comment = result.Description

	default:
		return "", fmt.Errorf("不支持的数据库类型: %s", driver)
	}

	return comment, nil
}

// getColumns 获取列信息
func (s *TableService) getColumns(tableName, driver string) ([]ColumnInfo, error) {
	var columns []ColumnInfo

	switch driver {
	case "mysql":
		// MySQL 查询列信息
		err := s.db.Raw(`
			SELECT 
				COLUMN_NAME as column_name,
				COLUMN_TYPE as column_type,
				DATA_TYPE as data_type,
				COLUMN_COMMENT as column_comment,
				IS_NULLABLE as is_nullable,
				COLUMN_DEFAULT as column_default,
				COLUMN_KEY as column_key,
				EXTRA as extra,
				CHARACTER_MAXIMUM_LENGTH as character_maximum_length,
				NUMERIC_PRECISION as numeric_precision,
				NUMERIC_SCALE as numeric_scale
			FROM information_schema.COLUMNS
			WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME = ?
			ORDER BY ORDINAL_POSITION
		`, tableName).Scan(&columns).Error
		if err != nil {
			return nil, fmt.Errorf("查询列信息失败: %v", err)
		}

	case "postgres", "postgresql":
		// PostgreSQL 查询列信息
		err := s.db.Raw(`
			SELECT 
				a.attname as column_name,
				pg_catalog.format_type(a.atttypid, a.atttypmod) as column_type,
				t.typname as data_type,
				col_description(a.attrelid, a.attnum) as column_comment,
				CASE WHEN a.attnotnull THEN 'NO' ELSE 'YES' END as is_nullable,
				pg_get_expr(d.adbin, d.adrelid) as column_default,
				CASE WHEN pk.contype = 'p' THEN 'PRI' ELSE '' END as column_key,
				CASE WHEN a.attidentity = 'a' THEN 'auto_increment' ELSE '' END as extra,
				a.atttypmod - 4 as character_maximum_length,
				0 as numeric_precision,
				0 as numeric_scale
			FROM pg_attribute a
			LEFT JOIN pg_class c ON a.attrelid = c.oid
			LEFT JOIN pg_type t ON a.atttypid = t.oid
			LEFT JOIN pg_attrdef d ON a.attrelid = d.adrelid AND a.attnum = d.adnum
			LEFT JOIN pg_constraint pk ON a.attrelid = pk.conrelid AND a.attnum = ANY(pk.conkey) AND pk.contype = 'p'
			WHERE c.relname = ? 
				AND a.attnum > 0 
				AND NOT a.attisdropped
			ORDER BY a.attnum
		`, tableName).Scan(&columns).Error
		if err != nil {
			return nil, fmt.Errorf("查询列信息失败: %v", err)
		}

	default:
		return nil, fmt.Errorf("不支持的数据库类型: %s", driver)
	}

	return columns, nil
}

// ConvertToFieldConfig 将列信息转换为字段配置
func ConvertToFieldConfig(column ColumnInfo) FieldConfig {
	field := FieldConfig{
		ColumnName:      column.ColumnName,
		ColumnType:      column.ColumnType,
		ColumnComment:   column.ColumnComment,
		IsNullable:      column.IsNullable == "YES",
		ColumnDefault:   column.ColumnDefault,
		IsPrimaryKey:    column.ColumnKey == "PRI",
		IsAutoIncrement: strings.Contains(column.Extra, "auto_increment"),

		// Go 字段信息（需要转换）
		FieldName: ToPascalCase(column.ColumnName),
		FieldType: ConvertDBTypeToGoType(column.DataType, column.IsNullable == "YES"),
		JsonTag:   column.ColumnName,
		GormTag:   BuildGormTag(column),
	}

	// 智能识别字段类型
	field.IsOperateField = IsOperateField(column.ColumnName)
	field.IsSortField = IsSortField(column.ColumnName)
	field.IsTimeField = IsTimeField(column.ColumnName)
	field.IsRelationField = IsRelationField(column.ColumnName)
	field.IsEnumField = strings.HasPrefix(column.DataType, "enum")
	field.IsSetField = strings.HasPrefix(column.DataType, "set")
	field.IsTextField = IsTextField(column.DataType)
	field.IsBoolField = IsBoolField(column.DataType)
	field.IsImageField = IsImageField(column.ColumnName)
	field.IsImagesField = IsImagesField(column.ColumnName)

	// 设置默认的表格和表单配置
	field.ShowInTable = !field.IsTimeField || column.ColumnName == "created_at" || column.ColumnName == "updated_at"
	field.TableSortable = field.IsSortField || column.ColumnName == "id"
	field.TableDisplayType = GetDefaultTableDisplayType(field)

	// 设置默认可搜索状态
	// 图片、排序、富文本(text类型)、desc/description/content结尾字段默认不可搜索
	field.TableSearchable = !field.IsImageField && !field.IsImagesField && !field.IsSortField &&
		!field.IsTextField &&
		!strings.HasSuffix(strings.ToLower(column.ColumnName), "desc") &&
		!strings.HasSuffix(strings.ToLower(column.ColumnName), "description") &&
		!strings.HasSuffix(strings.ToLower(column.ColumnName), "content")
	field.SearchFormType = GetDefaultSearchFormComponent(field) // 搜索表单类型

	field.ShowInForm = !field.IsPrimaryKey && !field.IsAutoIncrement && !field.IsTimeField
	field.FormComponent = GetDefaultFormComponent(field)

	// 设置验证规则
	field.InCreate = field.ShowInForm
	field.InUpdate = field.ShowInForm
	field.IsRequired = !field.IsNullable && field.ColumnDefault == "" && !field.IsAutoIncrement
	field.ValidateRules = BuildValidateRules(field)

	return field
}
