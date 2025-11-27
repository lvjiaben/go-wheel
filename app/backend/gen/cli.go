package gen

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"

	"gorm.io/gorm"
)

// CLI 命令行交互工具
type CLI struct {
	db      *gorm.DB
	scanner *bufio.Scanner
}

// NewCLI 创建命令行工具
func NewCLI(db *gorm.DB) *CLI {
	return &CLI{
		db:      db,
		scanner: bufio.NewScanner(os.Stdin),
	}
}

// Run 运行命令行工具
func (cli *CLI) Run() error {
	fmt.Println("=== Go-Wheel CRUD 代码生成器 ===")
	fmt.Println()

	// 1. 选择表
	tableName, err := cli.selectTable()
	if err != nil {
		return err
	}

	// 2. 获取表配置（自动应用智能默认值）
	tableService := NewTableService(cli.db)
	config, err := tableService.GetTableConfig(tableName)
	if err != nil {
		return err
	}

	fmt.Printf("\n表名：%s\n", config.TableName)
	fmt.Printf("表注释：%s\n", config.TableComment)

	// 3. 配置字段（可选）
	if cli.confirm("\n是否自定义字段配置？（默认使用智能识别）") {
		if err := cli.configureFields(config); err != nil {
			return err
		}
	}

	// 4. 配置方法
	if err := cli.configureMethods(config); err != nil {
		return err
	}

	// 5. 配置排序
	if err := cli.configureSorting(config); err != nil {
		return err
	}

	// 6. 配置菜单
	if err := cli.configureMenu(config); err != nil {
		return err
	}

	// 7. 预览代码
	if cli.confirm("是否预览生成的代码？") {
		cli.previewCode(config)
	}

	// 8. 生成代码
	if cli.confirm("确认生成代码？") {
		return cli.generateCode(config)
	}

	fmt.Println("已取消生成")
	return nil
}

// selectTable 选择表
func (cli *CLI) selectTable() (string, error) {
	tableService := NewTableService(cli.db)
	tables, err := tableService.GetAllTables()
	if err != nil {
		return "", err
	}

	fmt.Println("可用的数据库表：")
	for i, table := range tables {
		fmt.Printf("%d. %s\n", i+1, table)
	}

	fmt.Print("\n请选择表（输入序号或表名）：")
	input := cli.readLine()

	// 尝试解析为序号
	if idx, err := strconv.Atoi(input); err == nil && idx > 0 && idx <= len(tables) {
		return tables[idx-1], nil
	}

	// 直接使用表名
	return input, nil
}

// configureFields 配置字段
func (cli *CLI) configureFields(config *GenConfig) error {
	fmt.Println("\n=== 字段配置 ===")
	fmt.Println("字段已根据类型智能识别，当前配置如下：")

	// 显示当前字段配置
	cli.showFieldsTable(config)

	// 是否修改单个字段配置
	for {
		fmt.Print("\n输入字段序号进行配置（留空结束）：")
		input := cli.readLine()
		if input == "" {
			break
		}

		idx, err := strconv.Atoi(input)
		if err != nil || idx < 1 || idx > len(config.Fields) {
			fmt.Println("无效的序号")
			continue
		}

		cli.configureField(&config.Fields[idx-1])
	}

	return nil
}

// showFieldsTable 显示字段配置表格
func (cli *CLI) showFieldsTable(config *GenConfig) {
	fmt.Println("\n序号 | 字段名           | 类型       | 注释         | 表格 | 搜索 | 表单 | 组件")
	fmt.Println("-----|------------------|------------|--------------|------|------|------|------------")
	for i, field := range config.Fields {
		showTable := "✗"
		if field.ShowInTable {
			showTable = "✓"
		}
		searchable := "✗"
		if field.TableSearchable {
			searchable = "✓"
		}
		showForm := "✗"
		if field.ShowInForm {
			showForm = "✓"
		}
		fmt.Printf("%-4d | %-16s | %-10s | %-12s | %-4s | %-4s | %-4s | %s\n",
			i+1, truncate(field.ColumnName, 16), truncate(field.ColumnType, 10),
			truncate(field.ColumnComment, 12), showTable, searchable, showForm, field.FormComponent)
	}
}

// configureField 配置单个字段
func (cli *CLI) configureField(field *FieldConfig) {
	fmt.Printf("\n配置字段：%s (%s)\n", field.ColumnName, field.ColumnComment)

	// 是否显示在表格
	fmt.Printf("是否在表格显示（当前：%v）(y/n/留空跳过)：", field.ShowInTable)
	if input := cli.readLine(); input != "" {
		field.ShowInTable = strings.ToLower(input) == "y"
	}

	// 表格显示类型
	if field.ShowInTable {
		fmt.Printf("表格显示类型（当前：%s）\n", field.TableDisplayType)
		fmt.Println("  1.text 2.editable 3.tag 4.datetime 5.image 6.link 7.links")
		fmt.Print("选择（留空跳过）：")
		if input := cli.readLine(); input != "" {
			types := map[string]string{"1": "text", "2": "editable", "3": "tag", "4": "datetime", "5": "image", "6": "link", "7": "links"}
			if t, ok := types[input]; ok {
				field.TableDisplayType = t
			}
		}
	}

	// 是否可搜索
	fmt.Printf("是否可搜索（当前：%v）(y/n/留空跳过)：", field.TableSearchable)
	if input := cli.readLine(); input != "" {
		field.TableSearchable = strings.ToLower(input) == "y"
	}

	// 搜索表单类型
	if field.TableSearchable {
		fmt.Printf("搜索表单类型（当前：%s）\n", field.SearchFormType)
		fmt.Println("  1.Input 2.InputNumber 3.DatePicker 4.RangePicker 5.Switch 6.Select")
		fmt.Print("选择（留空跳过）：")
		if input := cli.readLine(); input != "" {
			types := map[string]string{"1": "Input", "2": "InputNumber", "3": "DatePicker", "4": "RangePicker", "5": "Switch", "6": "Select"}
			if t, ok := types[input]; ok {
				field.SearchFormType = t
			}
		}
	}

	// 是否显示在表单
	fmt.Printf("是否在表单显示（当前：%v）(y/n/留空跳过)：", field.ShowInForm)
	if input := cli.readLine(); input != "" {
		field.ShowInForm = strings.ToLower(input) == "y"
	}

	// 表单组件类型
	if field.ShowInForm {
		fmt.Printf("表单组件（当前：%s）\n", field.FormComponent)
		fmt.Println("  1.Input 2.Textarea 3.InputNumber 4.DatePicker 5.Switch")
		fmt.Println("  6.Select 7.Upload 8.RichTextEditor 9.ImageUpload")
		fmt.Print("选择（留空跳过）：")
		if input := cli.readLine(); input != "" {
			types := map[string]string{
				"1": "Input", "2": "Textarea", "3": "InputNumber", "4": "DatePicker",
				"5": "Switch", "6": "Select", "7": "Upload", "8": "RichTextEditor", "9": "ImageUpload",
			}
			if t, ok := types[input]; ok {
				field.FormComponent = t
			}
		}

		// 是否必填
		fmt.Printf("是否必填（当前：%v）(y/n/留空跳过)：", field.IsRequired)
		if input := cli.readLine(); input != "" {
			field.IsRequired = strings.ToLower(input) == "y"
		}

		// 组件配置
		fmt.Printf("组件配置 JSON（当前：%s）（留空跳过）：", field.ComponentProps)
		if input := cli.readLine(); input != "" {
			field.ComponentProps = input
		}
	}
}

// truncate 截断字符串
func truncate(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) > maxLen {
		return string(runes[:maxLen-2]) + ".."
	}
	return s
}

// configureMethods 配置方法
func (cli *CLI) configureMethods(config *GenConfig) error {
	fmt.Println("\n=== 方法配置 ===")
	fmt.Println("List 方法（必选）：已启用")

	config.Methods.Create = cli.confirm("是否生成 Create 方法？")
	config.Methods.Update = cli.confirm("是否生成 Update 方法？")
	config.Methods.Delete = cli.confirm("是否生成 Delete 方法？")

	// 自动识别 Operate 字段
	var operateFields []string
	for _, field := range config.Fields {
		if field.IsOperateField {
			operateFields = append(operateFields, field.ColumnName)
		}
	}

	if len(operateFields) > 0 {
		fmt.Printf("检测到 Operate 字段：%s\n", strings.Join(operateFields, ", "))
		config.Methods.Operate = cli.confirm("是否生成 Operate 方法？")
		if config.Methods.Operate {
			config.OperateFields = operateFields
		}
	}

	return nil
}

// configureSorting 配置排序
func (cli *CLI) configureSorting(config *GenConfig) error {
	fmt.Println("\n=== 排序配置 ===")

	// 默认排序字段
	fmt.Printf("默认排序字段（默认：%s，留空跳过）：", config.DefaultSortField)
	if input := cli.readLine(); input != "" {
		config.DefaultSortField = input
	}

	// 默认排序方向
	fmt.Printf("默认排序方向（1.desc 2.asc，当前：%s，留空跳过）：", config.DefaultSortOrder)
	if input := cli.readLine(); input != "" {
		if input == "1" {
			config.DefaultSortOrder = "desc"
		} else if input == "2" {
			config.DefaultSortOrder = "asc"
		}
	}

	return nil
}

// configureMenu 配置菜单
func (cli *CLI) configureMenu(config *GenConfig) error {
	fmt.Println("\n=== 菜单配置 ===")

	fmt.Printf("父级菜单名（默认：%s，留空跳过）：", config.MenuConfig.ParentMenuName)
	if input := cli.readLine(); input != "" {
		config.MenuConfig.ParentMenuName = input
	}

	fmt.Printf("菜单名（默认：%s，留空跳过）：", config.MenuConfig.MenuName)
	if input := cli.readLine(); input != "" {
		config.MenuConfig.MenuName = input
	}

	fmt.Printf("菜单图标（默认：%s，留空跳过）：", config.MenuConfig.MenuIcon)
	if input := cli.readLine(); input != "" {
		config.MenuConfig.MenuIcon = input
	}

	fmt.Printf("菜单排序（默认：%d，留空跳过）：", config.MenuConfig.MenuSort)
	if input := cli.readLine(); input != "" {
		if sort, err := strconv.Atoi(input); err == nil {
			config.MenuConfig.MenuSort = sort
		}
	}

	return nil
}

// previewCode 预览代码
func (cli *CLI) previewCode(config *GenConfig) {
	generator := NewGenerator(cli.db, config)
	code, err := generator.Generate()
	if err != nil {
		fmt.Printf("生成代码失败：%v\n", err)
		return
	}

	fmt.Println("\n=== 生成的文件 ===")
	paths := generator.GetFilePaths()

	fmt.Println("\n后端文件：")
	fmt.Printf("  - %s\n", paths.Backend.Controller)
	fmt.Printf("  - %s\n", paths.Backend.Service)
	fmt.Printf("  - %s\n", paths.Backend.Model)
	fmt.Printf("  - %s\n", paths.Backend.Validate)

	if config.FrontendSrcPath != "" {
		fmt.Println("\n前端文件：")
		fmt.Printf("  - %s\n", paths.Frontend.API)
		fmt.Printf("  - %s\n", paths.Frontend.ListView)
		fmt.Printf("  - %s\n", paths.Frontend.DataTS)
		fmt.Printf("  - %s\n", paths.Frontend.FormVue)
		fmt.Printf("  - %s\n", paths.Frontend.LocaleZhCN)
		fmt.Printf("  - %s\n", paths.Frontend.LocaleEnUS)
	}

	fmt.Println("\n菜单 SQL：")
	fmt.Println(code.Menu.SQL)
}

// generateCode 生成代码
func (cli *CLI) generateCode(config *GenConfig) error {
	generator := NewGenerator(cli.db, config)
	code, err := generator.Generate()
	if err != nil {
		return err
	}

	workDir := "."
	if err := generator.WriteFiles(code, workDir); err != nil {
		return err
	}

	// 执行菜单 SQL
	if err := cli.executeMenuSQL(code.Menu.SQL); err != nil {
		return fmt.Errorf("执行菜单 SQL 失败: %w", err)
	}

	// 保存生成历史
	if err := cli.saveHistory(config); err != nil {
		return fmt.Errorf("保存生成历史失败: %w", err)
	}

	fmt.Println("\n✓ 代码生成成功！")
	fmt.Println("✓ 菜单 SQL 已执行！")
	fmt.Println("✓ 生成历史已保存！")

	return nil
}

// saveHistory 保存生成历史
func (cli *CLI) saveHistory(config *GenConfig) error {
	configJSON, err := json.Marshal(config)
	if err != nil {
		return err
	}

	history := map[string]interface{}{
		"table_name":        config.TableName,
		"table_comment":     config.TableComment,
		"module_name":       config.ModuleName,
		"struct_name":       config.StructName,
		"package_name":      config.PackageName,
		"frontend_src_path": config.FrontendSrcPath,
		"config":            string(configJSON),
	}

	return cli.db.Table("gen_history").Create(history).Error
}

// executeMenuSQL 执行菜单 SQL
func (cli *CLI) executeMenuSQL(sql string) error {
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
			if err := cli.db.Exec(stmt).Error; err != nil {
				return fmt.Errorf("执行 SQL 失败: %s, 错误: %w", stmt, err)
			}
		}
	}

	return nil
}

// readLine 读取一行输入
func (cli *CLI) readLine() string {
	cli.scanner.Scan()
	return strings.TrimSpace(cli.scanner.Text())
}

// confirm 确认操作
func (cli *CLI) confirm(message string) bool {
	fmt.Printf("%s (Y/n)：", message)
	input := cli.readLine()
	return input == "" || strings.ToLower(input) == "y"
}
