package gen

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gorm.io/gorm"
)

// Generator 代码生成器
type Generator struct {
	db     *gorm.DB
	config *GenConfig
}

// NewGenerator 创建代码生成器
func NewGenerator(db *gorm.DB, config *GenConfig) *Generator {
	return &Generator{
		db:     db,
		config: config,
	}
}

// Generate 生成代码
func (g *Generator) Generate() (*GeneratedCode, error) {
	// 验证配置
	if err := g.validateConfig(); err != nil {
		return nil, err
	}

	// 生成后端代码
	backendGen := NewBackendGenerator(g.config)
	backendCode := backendGen.Generate()

	// 生成前端代码
	frontendGen := NewFrontendGenerator(g.config)
	frontendCode := frontendGen.Generate()

	// 生成菜单 SQL
	menuCode := g.generateMenuSQL()

	return &GeneratedCode{
		Backend:  backendCode,
		Frontend: frontendCode,
		Menu:     menuCode,
	}, nil
}

// validateConfig 验证配置
func (g *Generator) validateConfig() error {
	if g.config.TableName == "" {
		return fmt.Errorf("表名不能为空")
	}
	if g.config.ModuleName == "" {
		return fmt.Errorf("模块名不能为空")
	}
	if g.config.StructName == "" {
		return fmt.Errorf("结构体名不能为空")
	}
	if g.config.PackageName == "" {
		return fmt.Errorf("包名不能为空")
	}
	return nil
}

// generateMenuSQL 生成菜单 SQL
func (g *Generator) generateMenuSQL() MenuCode {
	var sb strings.Builder

	// 查询父级菜单
	parentMenuName := g.config.MenuConfig.ParentMenuName
	if parentMenuName == "" {
		parentMenuName = "代码生成"
	}

	menuName := g.config.MenuConfig.MenuName
	if menuName == "" {
		if g.config.TableComment != "" {
			menuName = g.config.TableComment
		} else {
			menuName = g.config.TableName
		}
	}

	// 生成 SQL
	sb.WriteString("-- 菜单 SQL\n")
	sb.WriteString("-- 请在数据库中执行以下 SQL\n\n")

	// 生成 enname（驼峰格式）
	enname := ToCamelCase(g.config.TableName)

	// 1. 检查并创建父级菜单
	sb.WriteString(fmt.Sprintf("-- 1. 检查父级菜单 '%s' 是否存在，不存在则创建\n", parentMenuName))
	sb.WriteString("SET @parent_id = (SELECT id FROM admin_menu WHERE name = '" + parentMenuName + "' LIMIT 1);\n")
	sb.WriteString("INSERT INTO admin_menu (pid, type, name, enname, icon, sort, show_tag, fixed_tag, created_at, updated_at)\n")
	sb.WriteString("SELECT 0, 'menu', '" + parentMenuName + "', 'codeGen', '', 100, 0, 0, UNIX_TIMESTAMP(), UNIX_TIMESTAMP()\n")
	sb.WriteString("WHERE @parent_id IS NULL;\n")
	sb.WriteString("SET @parent_id = (SELECT id FROM admin_menu WHERE name = '" + parentMenuName + "' LIMIT 1);\n\n")

	// 2. 创建菜单
	menuIcon := g.config.MenuConfig.MenuIcon
	if menuIcon == "" {
		menuIcon = "mdi:file-document"
	}
	component := fmt.Sprintf("/%s/list", g.config.ModuleName)

	sb.WriteString(fmt.Sprintf("-- 2. 创建 '%s' 菜单\n", menuName))
	sb.WriteString("INSERT INTO admin_menu (pid, type, name, enname, path, component, icon, sort, show_tag, fixed_tag, created_at, updated_at) VALUES\n")
	sb.WriteString(fmt.Sprintf("(@parent_id, 'menu', '%s', '%s', '/%s', '%s', '%s', %d, 0, 0, UNIX_TIMESTAMP(), UNIX_TIMESTAMP());\n",
		menuName, enname, g.config.ModuleName, component, menuIcon, g.config.MenuConfig.MenuSort))
	sb.WriteString("SET @menu_id = LAST_INSERT_ID();\n\n")

	// 3. 创建权限按钮
	sb.WriteString("-- 3. 创建权限按钮\n")
	sb.WriteString("INSERT INTO admin_menu (pid, type, name, enname, permission, route, sort, created_at, updated_at) VALUES\n")

	var permissions []string

	// List 权限（必选）
	permissions = append(permissions, fmt.Sprintf("(@menu_id, 'button', '查看', '%s:list', '%s:list', '/backend/%s/list', 100, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())",
		g.config.ModuleName, g.config.ModuleName, g.config.ModuleName))

	// Create 权限
	if g.config.Methods.Create {
		permissions = append(permissions, fmt.Sprintf("(@menu_id, 'button', '添加', '%s:create', '%s:create', '/backend/%s/create', 90, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())",
			g.config.ModuleName, g.config.ModuleName, g.config.ModuleName))
	}

	// Update 权限
	if g.config.Methods.Update {
		permissions = append(permissions, fmt.Sprintf("(@menu_id, 'button', '编辑', '%s:update', '%s:update', '/backend/%s/update', 80, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())",
			g.config.ModuleName, g.config.ModuleName, g.config.ModuleName))
	}

	// Delete 权限
	if g.config.Methods.Delete {
		permissions = append(permissions, fmt.Sprintf("(@menu_id, 'button', '删除', '%s:delete', '%s:delete', '/backend/%s/delete', 70, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())",
			g.config.ModuleName, g.config.ModuleName, g.config.ModuleName))
	}

	// Operate 权限
	if g.config.Methods.Operate {
		permissions = append(permissions, fmt.Sprintf("(@menu_id, 'button', '操作', '%s:operate', '%s:operate', '/backend/%s/operate', 60, UNIX_TIMESTAMP(), UNIX_TIMESTAMP())",
			g.config.ModuleName, g.config.ModuleName, g.config.ModuleName))
	}

	sb.WriteString(strings.Join(permissions, ",\n"))
	sb.WriteString(";\n")

	return MenuCode{SQL: sb.String()}
}

// GetFilePaths 获取生成的文件路径
func (g *Generator) GetFilePaths() FilePath {
	return FilePath{
		Backend: BackendFilePath{
			Controller: fmt.Sprintf("app/backend/controller/%s.go", g.config.PackageName),
			Service:    fmt.Sprintf("app/backend/service/%s.go", g.config.PackageName),
			Model:      fmt.Sprintf("app/backend/model/%s.go", g.config.PackageName),
			Validate:   fmt.Sprintf("app/backend/validate/%s.go", g.config.PackageName),
		},
		Frontend: FrontendFilePath{
			API:        fmt.Sprintf("api/%s.ts", g.config.ModuleName),
			ListView:   fmt.Sprintf("views/%s/list.vue", g.config.ModuleName),
			DataTS:     fmt.Sprintf("views/%s/data.ts", g.config.ModuleName),
			FormVue:    fmt.Sprintf("views/%s/modules/form.vue", g.config.ModuleName),
			LocaleZhCN: fmt.Sprintf("locales/langs/zh-cn/%s.json", g.config.ModuleName),
			LocaleEnUS: fmt.Sprintf("locales/langs/en-us/%s.json", g.config.ModuleName),
		},
	}
}

// WriteFiles 写入文件
func (g *Generator) WriteFiles(code *GeneratedCode, workDir string) error {
	paths := g.GetFilePaths()

	// 写入后端文件
	if err := g.writeBackendFiles(code.Backend, paths.Backend, workDir); err != nil {
		return err
	}

	// 写入前端文件
	if err := g.writeFrontendFiles(code.Frontend, paths.Frontend, g.config.FrontendSrcPath); err != nil {
		return err
	}

	// 更新路由文件
	if err := g.updateRoutes(code.Backend.Route, workDir); err != nil {
		return err
	}

	return nil
}

// writeBackendFiles 写入后端文件
func (g *Generator) writeBackendFiles(code BackendCode, paths BackendFilePath, workDir string) error {
	// Controller
	if err := g.writeFile(filepath.Join(workDir, paths.Controller), code.Controller); err != nil {
		return fmt.Errorf("写入 Controller 失败: %v", err)
	}

	// Service
	if err := g.writeFile(filepath.Join(workDir, paths.Service), code.Service); err != nil {
		return fmt.Errorf("写入 Service 失败: %v", err)
	}

	// Model
	if err := g.writeFile(filepath.Join(workDir, paths.Model), code.Model); err != nil {
		return fmt.Errorf("写入 Model 失败: %v", err)
	}

	// Validate
	if err := g.writeFile(filepath.Join(workDir, paths.Validate), code.Validate); err != nil {
		return fmt.Errorf("写入 Validate 失败: %v", err)
	}

	return nil
}

// writeFrontendFiles 写入前端文件
func (g *Generator) writeFrontendFiles(code FrontendCode, paths FrontendFilePath, frontendSrcPath string) error {
	if frontendSrcPath == "" {
		return fmt.Errorf("前端 src 路径未设置")
	}

	// API
	if err := g.writeFile(filepath.Join(frontendSrcPath, paths.API), code.API); err != nil {
		return fmt.Errorf("写入 API 失败: %v", err)
	}

	// ListView
	if err := g.writeFile(filepath.Join(frontendSrcPath, paths.ListView), code.ListView); err != nil {
		return fmt.Errorf("写入 ListView 失败: %v", err)
	}

	// DataTS
	if err := g.writeFile(filepath.Join(frontendSrcPath, paths.DataTS), code.DataTS); err != nil {
		return fmt.Errorf("写入 DataTS 失败: %v", err)
	}

	// FormVue
	if err := g.writeFile(filepath.Join(frontendSrcPath, paths.FormVue), code.FormVue); err != nil {
		return fmt.Errorf("写入 FormVue 失败: %v", err)
	}

	// LocaleZhCN (中文语言包)
	if err := g.writeFile(filepath.Join(frontendSrcPath, paths.LocaleZhCN), code.LocaleZhCN); err != nil {
		return fmt.Errorf("写入 LocaleZhCN 失败: %v", err)
	}

	// LocaleEnUS (英文语言包)
	if err := g.writeFile(filepath.Join(frontendSrcPath, paths.LocaleEnUS), code.LocaleEnUS); err != nil {
		return fmt.Errorf("写入 LocaleEnUS 失败: %v", err)
	}

	return nil
}

// writeFile 写入文件
func (g *Generator) writeFile(filePath, content string) error {
	// 创建目录
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// 写入文件
	return os.WriteFile(filePath, []byte(content), 0644)
}

// updateRoutes 更新路由文件
func (g *Generator) updateRoutes(routeCode, workDir string) error {
	routeFile := filepath.Join(workDir, "routes/routes.go")

	// 读取文件
	content, err := os.ReadFile(routeFile)
	if err != nil {
		return fmt.Errorf("读取路由文件失败: %v", err)
	}

	contentStr := string(content)

	// 查找插入位置
	regionStart := "// region:backend-routes"
	regionEnd := "// endregion:backend-routes"

	startIdx := strings.Index(contentStr, regionStart)
	endIdx := strings.Index(contentStr, regionEnd)

	if startIdx == -1 || endIdx == -1 {
		return fmt.Errorf("未找到路由插入区域标记")
	}

	// 插入路由代码
	newContent := contentStr[:startIdx+len(regionStart)] + "\n" + routeCode + "\n\t\t" + contentStr[endIdx:]

	// 写回文件
	return os.WriteFile(routeFile, []byte(newContent), 0644)
}

// DeleteFiles 删除生成的文件
func (g *Generator) DeleteFiles(workDir string) error {
	paths := g.GetFilePaths()

	// 删除后端文件
	os.Remove(filepath.Join(workDir, paths.Backend.Controller))
	os.Remove(filepath.Join(workDir, paths.Backend.Service))
	os.Remove(filepath.Join(workDir, paths.Backend.Model))
	os.Remove(filepath.Join(workDir, paths.Backend.Validate))

	// 删除前端文件
	if g.config.FrontendSrcPath != "" {
		os.Remove(filepath.Join(g.config.FrontendSrcPath, paths.Frontend.API))
		os.Remove(filepath.Join(g.config.FrontendSrcPath, paths.Frontend.ListView))
		os.Remove(filepath.Join(g.config.FrontendSrcPath, paths.Frontend.DataTS))
		os.Remove(filepath.Join(g.config.FrontendSrcPath, paths.Frontend.FormVue))
		os.Remove(filepath.Join(g.config.FrontendSrcPath, paths.Frontend.LocaleZhCN))
		os.Remove(filepath.Join(g.config.FrontendSrcPath, paths.Frontend.LocaleEnUS))

		// 删除空目录
		os.Remove(filepath.Join(g.config.FrontendSrcPath, fmt.Sprintf("views/%s/modules", g.config.ModuleName)))
		os.Remove(filepath.Join(g.config.FrontendSrcPath, fmt.Sprintf("views/%s", g.config.ModuleName)))
	}

	// 删除路由中的代码
	g.deleteRoutes(workDir)

	return nil
}

// deleteRoutes 删除路由文件中的模块代码
func (g *Generator) deleteRoutes(workDir string) error {
	routeFile := filepath.Join(workDir, "routes/routes.go")

	// 读取文件
	content, err := os.ReadFile(routeFile)
	if err != nil {
		return fmt.Errorf("读取路由文件失败: %v", err)
	}

	contentStr := string(content)

	// 查找 region 范围
	regionStart := "// region:backend-routes"
	regionEnd := "// endregion:backend-routes"

	startIdx := strings.Index(contentStr, regionStart)
	endIdx := strings.Index(contentStr, regionEnd)

	if startIdx == -1 || endIdx == -1 {
		return nil // 没有找到 region，跳过
	}

	// 提取 region 内的内容
	regionContent := contentStr[startIdx+len(regionStart) : endIdx]

	// 使用 moduleName 来匹配路由（api.Group("/moduleName")）
	moduleName := g.config.ModuleName
	camelStructName := ToCamelCase(g.config.StructName)

	// 匹配模式：
	// // xxx管理
	// xxxController := controller.NewXxxController(c)
	// xxxGroup := api.Group("/xxx")...
	// {
	//     ...
	// }
	//
	// 使用更灵活的正则：匹配包含 api.Group("/moduleName") 的整个代码块
	pattern := fmt.Sprintf(`(?s)\n\t*//[^\n]*\n\t*%sController\s*:=\s*controller\.New%sController\(c\)\n\t*%sGroup\s*:=\s*api\.Group\("/%s"\)[^\{]*\{[^\}]*\}\n*`,
		camelStructName, g.config.StructName, camelStructName, moduleName)

	re := regexp.MustCompile(pattern)
	newRegionContent := re.ReplaceAllString(regionContent, "\n")

	// 重新组装内容
	newContent := contentStr[:startIdx+len(regionStart)] + newRegionContent + contentStr[endIdx:]

	// 写回文件
	return os.WriteFile(routeFile, []byte(newContent), 0644)
}
