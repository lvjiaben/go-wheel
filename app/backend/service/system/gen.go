package system

import (
	"archive/zip"
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/lvjiaben/go-wheel/app/backend/gen"
	"github.com/lvjiaben/go-wheel/app/backend/model/system"
	"github.com/lvjiaben/go-wheel/pkg/container"
)

// GenService 代码生成服务
type GenService struct {
	container *container.Container
}

// NewGenService 创建代码生成服务
func NewGenService(c *container.Container) *GenService {
	return &GenService{container: c}
}

// GetTableList 获取数据库表列表
func (s *GenService) GetTableList(search string) ([]map[string]interface{}, error) {
	tableService := gen.NewTableService(s.container.GetDB())

	// 获取所有表
	tables, err := tableService.GetAllTables()
	if err != nil {
		return nil, err
	}

	// 获取每个表的注释
	var result []map[string]interface{}
	for _, tableName := range tables {
		// 搜索过滤
		if search != "" && !strings.Contains(tableName, search) {
			continue
		}

		tableInfo, err := tableService.GetTableInfo(tableName)
		if err != nil {
			continue
		}

		result = append(result, map[string]interface{}{
			"table_name":    tableInfo.TableName,
			"table_comment": tableInfo.TableComment,
		})
	}

	return result, nil
}

// GetTableInfo 获取表详细信息
func (s *GenService) GetTableInfo(tableName string) (*gen.TableInfo, error) {
	tableService := gen.NewTableService(s.container.GetDB())
	return tableService.GetTableInfo(tableName)
}

// GetTableConfig 获取表的默认配置
func (s *GenService) GetTableConfig(tableName string) (*gen.GenConfig, error) {
	// 获取表信息
	tableInfo, err := s.GetTableInfo(tableName)
	if err != nil {
		return nil, err
	}

	// 生成默认配置
	config := &gen.GenConfig{
		TableName:       tableInfo.TableName,
		TableComment:    tableInfo.TableComment,
		ModuleName:      tableInfo.TableName,
		PackageName:     tableInfo.TableName,
		StructName:      gen.ToPascalCase(tableInfo.TableName),
		FrontendSrcPath: "vben-admin/apps/web-antd/src", // 默认前端路径

		// 默认生成所有方法
		Methods: gen.MethodConfig{
			List:    true,
			Create:  true,
			Update:  true,
			Delete:  true,
			Operate: false, // 默认不生成 Operate
		},

		// 字段配置
		Fields: make([]gen.FieldConfig, 0),

		// 搜索字段（默认为空，用户可选）
		SearchFields: []string{},

		// Operate 字段（默认为空）
		OperateFields: []string{},

		// 菜单配置
		MenuConfig: gen.MenuConfig{
			ParentMenuName: "AutoPlay",
			MenuName:       tableInfo.TableComment,
			MenuIcon:       "",
			MenuSort:       50,
		},
	}

	// 转换字段配置
	var operateFields []string
	hasSortField := false
	hasWeighField := false
	for _, column := range tableInfo.Columns {
		field := gen.ConvertToFieldConfig(column)
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
	config.DefaultSortOrder = "desc"

	return config, nil
}

// fillConfigDefaults 自动填充配置的默认值
func (s *GenService) fillConfigDefaults(config *gen.GenConfig) {
	// 自动生成 ModuleName（使用表名）
	if config.ModuleName == "" {
		config.ModuleName = config.TableName
	}

	// 自动生成 PackageName（使用表名）
	if config.PackageName == "" {
		config.PackageName = config.TableName
	}

	// 自动生成 StructName（表名转驼峰）
	if config.StructName == "" {
		config.StructName = gen.ToPascalCase(config.TableName)
	}

	// 自动设置前端路径
	if config.FrontendSrcPath == "" {
		config.FrontendSrcPath = "vben-admin/apps/web-antd/src"
	}
}

// PreviewCode 预览生成的代码
func (s *GenService) PreviewCode(config *gen.GenConfig) (*gen.GeneratedCode, error) {
	// 自动填充必要字段
	s.fillConfigDefaults(config)

	generator := gen.NewGenerator(s.container.GetDB(), config)
	return generator.Generate()
}

// GenerateCode 生成代码并写入文件
func (s *GenService) GenerateCode(config *gen.GenConfig) error {
	// 自动填充必要字段
	s.fillConfigDefaults(config)

	// 生成代码
	generator := gen.NewGenerator(s.container.GetDB(), config)
	code, err := generator.Generate()
	if err != nil {
		return err
	}

	// 写入文件（获取当前工作目录）
	workDir := "."
	if err := generator.WriteFiles(code, workDir); err != nil {
		return err
	}

	// 执行菜单 SQL
	if err := s.executeMenuSQL(code.Menu.SQL); err != nil {
		return fmt.Errorf("执行菜单 SQL 失败: %w", err)
	}

	// 保存生成历史
	if err := s.saveHistory(config); err != nil {
		return err
	}

	return nil
}

// executeMenuSQL 执行菜单 SQL
func (s *GenService) executeMenuSQL(sql string) error {
	// 移除注释行
	lines := strings.Split(sql, "\n")
	var sqlStatements []string
	var currentStatement strings.Builder

	for _, line := range lines {
		line = strings.TrimSpace(line)
		// 跳过注释和空行
		if line == "" || strings.HasPrefix(line, "--") {
			continue
		}

		currentStatement.WriteString(line)
		currentStatement.WriteString(" ")

		// 如果遇到分号，执行当前语句
		if strings.HasSuffix(line, ";") {
			sqlStatements = append(sqlStatements, currentStatement.String())
			currentStatement.Reset()
		}
	}

	// 执行所有 SQL 语句
	for _, stmt := range sqlStatements {
		stmt = strings.TrimSpace(stmt)
		if stmt != "" {
			if err := s.container.GetDB().Exec(stmt).Error; err != nil {
				return fmt.Errorf("执行 SQL 失败: %s, 错误: %w", stmt, err)
			}
		}
	}

	return nil
}

// saveHistory 保存生成历史
func (s *GenService) saveHistory(config *gen.GenConfig) error {
	// 将配置转换为 JSON
	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}

	history := &system.GenHistory{
		GenTableName:    config.TableName,
		TableComment:    config.TableComment,
		ModuleName:      config.ModuleName,
		StructName:      config.StructName,
		PackageName:     config.PackageName,
		FrontendSrcPath: config.FrontendSrcPath,
		Config:          string(configJSON),
	}

	return s.container.GetDB().Create(history).Error
}

// GetHistory 获取生成历史
func (s *GenService) GetHistory() ([]system.GenHistory, error) {
	var histories []system.GenHistory
	err := s.container.GetDB().Order("id DESC").Find(&histories).Error
	return histories, err
}

// DeleteGenerated 删除生成的代码
func (s *GenService) DeleteGenerated(id int) error {
	// 查询历史记录
	var history system.GenHistory
	if err := s.container.GetDB().First(&history, id).Error; err != nil {
		return fmt.Errorf("历史记录不存在")
	}

	// 解析配置
	var config gen.GenConfig
	if err := json.Unmarshal([]byte(history.Config), &config); err != nil {
		return err
	}

	// 删除文件
	generator := gen.NewGenerator(s.container.GetDB(), &config)
	workDir := "."
	if err := generator.DeleteFiles(workDir); err != nil {
		return err
	}

	// 删除历史记录
	return s.container.GetDB().Delete(&history).Error
}

// DownloadCode 下载生成的代码（返回 ZIP 文件）
func (s *GenService) DownloadCode(config *gen.GenConfig) ([]byte, error) {
	// 自动填充必要字段
	s.fillConfigDefaults(config)
	// 生成代码
	generator := gen.NewGenerator(s.container.GetDB(), config)
	code, err := generator.Generate()
	if err != nil {
		return nil, err
	}

	// 创建 ZIP 文件
	var buf bytes.Buffer
	zipWriter := zip.NewWriter(&buf)

	// 获取文件路径
	paths := generator.GetFilePaths()

	// 添加后端文件
	files := map[string]string{
		paths.Backend.Controller:     code.Backend.Controller,
		paths.Backend.Service:        code.Backend.Service,
		paths.Backend.Model:          code.Backend.Model,
		paths.Backend.Validate:       code.Backend.Validate,
		"routes/routes_generated.go": code.Backend.Route,
	}

	// 添加前端文件
	if config.FrontendSrcPath != "" {
		files[paths.Frontend.API] = code.Frontend.API
		files[paths.Frontend.ListView] = code.Frontend.ListView
		files[paths.Frontend.DataTS] = code.Frontend.DataTS
		files[paths.Frontend.FormVue] = code.Frontend.FormVue
		files[paths.Frontend.LocaleZhCN] = code.Frontend.LocaleZhCN
		files[paths.Frontend.LocaleEnUS] = code.Frontend.LocaleEnUS
	}

	// 添加菜单 SQL
	files["menu.sql"] = code.Menu.SQL

	// 写入文件到 ZIP
	for filePath, content := range files {
		f, err := zipWriter.Create(filePath)
		if err != nil {
			zipWriter.Close()
			return nil, err
		}
		_, err = f.Write([]byte(content))
		if err != nil {
			zipWriter.Close()
			return nil, err
		}
	}

	// 关闭 ZIP writer
	if err := zipWriter.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
