package gen

import (
	"fmt"
	"strings"
)

// BackendGenerator 后端代码生成器
type BackendGenerator struct {
	config *GenConfig
}

// NewBackendGenerator 创建后端代码生成器
func NewBackendGenerator(config *GenConfig) *BackendGenerator {
	return &BackendGenerator{config: config}
}

// Generate 生成后端代码
func (g *BackendGenerator) Generate() BackendCode {
	return BackendCode{
		Controller: g.GenerateController(),
		Service:    g.GenerateService(),
		Model:      g.GenerateModel(),
		Validate:   g.GenerateValidate(),
		Route:      g.GenerateRoute(),
	}
}

// GenerateController 生成 Controller 代码
func (g *BackendGenerator) GenerateController() string {
	var sb strings.Builder

	// 包声明
	sb.WriteString("package controller\n\n")

	// 导入
	sb.WriteString("import (\n")
	sb.WriteString("\t\"github.com/gin-gonic/gin\"\n")
	sb.WriteString("\t\"github.com/lvjiaben/go-wheel/app/backend/service\"\n")
	sb.WriteString("\t\"github.com/lvjiaben/go-wheel/app/backend/validate\"\n")
	sb.WriteString("\t\"github.com/lvjiaben/go-wheel/pkg/container\"\n")
	sb.WriteString("\t\"github.com/lvjiaben/go-wheel/pkg/utils/http\"\n")
	sb.WriteString(")\n\n")

	// 结构体
	sb.WriteString(fmt.Sprintf("type %sController struct {\n", g.config.StructName))
	sb.WriteString("\tcontainer *container.Container\n")
	sb.WriteString(fmt.Sprintf("\t%sService *service.%sService\n",
		ToCamelCase(g.config.StructName), g.config.StructName))
	sb.WriteString("}\n\n")

	// 构造函数
	sb.WriteString(fmt.Sprintf("func New%sController(c *container.Container) *%sController {\n",
		g.config.StructName, g.config.StructName))
	sb.WriteString(fmt.Sprintf("\treturn &%sController{\n", g.config.StructName))
	sb.WriteString("\t\tcontainer: c,\n")
	sb.WriteString(fmt.Sprintf("\t\t%sService: service.New%sService(c),\n",
		ToCamelCase(g.config.StructName), g.config.StructName))
	sb.WriteString("\t}\n")
	sb.WriteString("}\n\n")

	// List 方法（必选）
	sb.WriteString(fmt.Sprintf("func (ctrl *%sController) List(ctx *gin.Context) {\n", g.config.StructName))
	sb.WriteString(fmt.Sprintf("\thttp.SuccessWithI18n(ctx, \"common.success\", ctrl.%sService.List(ctx))\n",
		ToCamelCase(g.config.StructName)))
	sb.WriteString("}\n\n")

	// Create 方法（可选）
	if g.config.Methods.Create {
		sb.WriteString(fmt.Sprintf("func (ctrl *%sController) Create(ctx *gin.Context) {\n", g.config.StructName))
		sb.WriteString(fmt.Sprintf("\tform, valid := validate.Validate%sCreate(ctx)\n", g.config.StructName))
		sb.WriteString("\tif !valid {\n")
		sb.WriteString("\t\treturn\n")
		sb.WriteString("\t}\n")
		sb.WriteString(fmt.Sprintf("\tres, err := ctrl.%sService.Create(ctx, form)\n",
			ToCamelCase(g.config.StructName)))
		sb.WriteString("\tif err != nil {\n")
		sb.WriteString("\t\thttp.ErrorWithI18n(ctx, err.Error(), nil)\n")
		sb.WriteString("\t\treturn\n")
		sb.WriteString("\t}\n")
		sb.WriteString("\thttp.SuccessWithI18n(ctx, \"common.success\", res)\n")
		sb.WriteString("}\n\n")
	}

	// Update 方法（可选）
	if g.config.Methods.Update {
		sb.WriteString(fmt.Sprintf("func (ctrl *%sController) Update(ctx *gin.Context) {\n", g.config.StructName))
		sb.WriteString(fmt.Sprintf("\tform, valid := validate.Validate%sUpdate(ctx)\n", g.config.StructName))
		sb.WriteString("\tif !valid {\n")
		sb.WriteString("\t\treturn\n")
		sb.WriteString("\t}\n")
		sb.WriteString(fmt.Sprintf("\tres, err := ctrl.%sService.Update(ctx, form)\n",
			ToCamelCase(g.config.StructName)))
		sb.WriteString("\tif err != nil {\n")
		sb.WriteString("\t\thttp.ErrorWithI18n(ctx, err.Error(), nil)\n")
		sb.WriteString("\t\treturn\n")
		sb.WriteString("\t}\n")
		sb.WriteString("\thttp.SuccessWithI18n(ctx, \"common.success\", res)\n")
		sb.WriteString("}\n\n")
	}

	// Delete 方法（可选）
	if g.config.Methods.Delete {
		sb.WriteString(fmt.Sprintf("func (ctrl *%sController) Delete(ctx *gin.Context) {\n", g.config.StructName))
		sb.WriteString(fmt.Sprintf("\tform, valid := validate.Validate%sDelete(ctx)\n", g.config.StructName))
		sb.WriteString("\tif !valid {\n")
		sb.WriteString("\t\treturn\n")
		sb.WriteString("\t}\n")
		sb.WriteString(fmt.Sprintf("\tif err := ctrl.%sService.Delete(ctx, form.Ids); err != nil {\n",
			ToCamelCase(g.config.StructName)))
		sb.WriteString("\t\thttp.ErrorWithI18n(ctx, err.Error(), nil)\n")
		sb.WriteString("\t\treturn\n")
		sb.WriteString("\t}\n")
		sb.WriteString("\thttp.SuccessWithI18n(ctx, \"common.success\", nil)\n")
		sb.WriteString("}\n\n")
	}

	// Operate 方法（可选）
	if g.config.Methods.Operate {
		sb.WriteString(fmt.Sprintf("func (ctrl *%sController) Operate(ctx *gin.Context) {\n", g.config.StructName))
		sb.WriteString(fmt.Sprintf("\tform, valid := validate.Validate%sOperate(ctx)\n", g.config.StructName))
		sb.WriteString("\tif !valid {\n")
		sb.WriteString("\t\treturn\n")
		sb.WriteString("\t}\n")
		sb.WriteString(fmt.Sprintf("\tif err := ctrl.%sService.Operate(form.Ids, form.Field, form.Value); err != nil {\n",
			ToCamelCase(g.config.StructName)))
		sb.WriteString("\t\thttp.ErrorWithI18n(ctx, err.Error(), nil)\n")
		sb.WriteString("\t\treturn\n")
		sb.WriteString("\t}\n")
		sb.WriteString("\thttp.SuccessWithI18n(ctx, \"common.success\", nil)\n")
		sb.WriteString("}\n")
	}

	return sb.String()
}

// GenerateService 生成 Service 代码
func (g *BackendGenerator) GenerateService() string {
	var sb strings.Builder

	// 包声明
	sb.WriteString("package service\n\n")

	// 导入
	sb.WriteString("import (\n")
	if g.config.Methods.Operate {
		sb.WriteString("\t\"fmt\"\n\n")
	}
	sb.WriteString("\t\"github.com/gin-gonic/gin\"\n")
	if g.config.Methods.Create || g.config.Methods.Update {
		sb.WriteString("\t\"github.com/lvjiaben/go-wheel/app/backend/validate\"\n")
	}
	sb.WriteString("\t\"github.com/lvjiaben/go-wheel/app/common/builder\"\n")

	// 判断 Model 的位置
	modelPackage := "commonModel"
	modelImport := "github.com/lvjiaben/go-wheel/app/common/model"
	// 如果表名不在 common/model 中，使用 backend/model
	if !g.isCommonModel() {
		modelPackage = "backendModel"
		modelImport = "github.com/lvjiaben/go-wheel/app/backend/model"
	}
	sb.WriteString(fmt.Sprintf("\t%s \"%s\"\n", modelPackage, modelImport))

	sb.WriteString("\t\"github.com/lvjiaben/go-wheel/pkg/container\"\n")
	if g.config.Methods.Operate {
		sb.WriteString("\t\"github.com/lvjiaben/go-wheel/pkg/utils/datatype\"\n")
	}
	sb.WriteString("\t\"gorm.io/gorm\"\n")
	sb.WriteString(")\n\n")

	// 结构体
	sb.WriteString(fmt.Sprintf("type %sService struct {\n", g.config.StructName))
	sb.WriteString("\tcontainer   *container.Container\n")
	sb.WriteString(fmt.Sprintf("\tcrudService *builder.CRUDBuilder[%s.%s]\n", modelPackage, g.config.StructName))
	sb.WriteString("}\n\n")

	// 构造函数
	sb.WriteString(fmt.Sprintf("func New%sService(c *container.Container) *%sService {\n",
		g.config.StructName, g.config.StructName))
	sb.WriteString(fmt.Sprintf("\treturn &%sService{\n", g.config.StructName))
	sb.WriteString("\t\tcontainer: c,\n")
	sb.WriteString(fmt.Sprintf("\t\tcrudService: builder.NewCRUDBuilderWithProvider[%s.%s](\n",
		modelPackage, g.config.StructName))
	sb.WriteString("\t\t\tfunc(ctx *gin.Context) *gorm.DB {\n")
	sb.WriteString("\t\t\t\treturn c.GetDB().WithContext(ctx.Request.Context())\n")
	sb.WriteString("\t\t\t},\n")

	// 添加搜索字段
	if len(g.config.SearchFields) > 0 {
		sb.WriteString(fmt.Sprintf("\t\t).WithSearchFields(\"%s\"),\n",
			strings.Join(g.config.SearchFields, "\", \"")))
	} else {
		sb.WriteString("\t\t),\n")
	}

	sb.WriteString("\t}\n")
	sb.WriteString("}\n\n")

	// List 方法
	sb.WriteString(fmt.Sprintf("func (s *%sService) List(ctx *gin.Context) map[string]interface{} {\n",
		g.config.StructName))
	sb.WriteString("\treturn s.crudService.WithContext(ctx).List()\n")
	sb.WriteString("}\n\n")

	// Create 方法
	if g.config.Methods.Create {
		sb.WriteString(fmt.Sprintf("func (s *%sService) Create(ctx *gin.Context, form *validate.%sCreate) (*%s.%s, error) {\n",
			g.config.StructName, g.config.StructName, modelPackage, g.config.StructName))
		sb.WriteString("\treturn s.crudService.WithContext(ctx).Create(form)\n")
		sb.WriteString("}\n\n")
	}

	// Update 方法
	if g.config.Methods.Update {
		sb.WriteString(fmt.Sprintf("func (s *%sService) Update(ctx *gin.Context, form *validate.%sUpdate) (*%s.%s, error) {\n",
			g.config.StructName, g.config.StructName, modelPackage, g.config.StructName))
		sb.WriteString("\treturn s.crudService.WithContext(ctx).Update(form.Id, form)\n")
		sb.WriteString("}\n\n")
	}

	// Delete 方法
	if g.config.Methods.Delete {
		sb.WriteString(fmt.Sprintf("func (s *%sService) Delete(ctx *gin.Context, ids []int) error {\n",
			g.config.StructName))
		sb.WriteString("\treturn s.crudService.WithContext(ctx).Delete(ids)\n")
		sb.WriteString("}\n\n")
	}

	// Operate 方法
	if g.config.Methods.Operate {
		sb.WriteString(fmt.Sprintf("func (s *%sService) Operate(ids []int, field string, value int) error {\n",
			g.config.StructName))
		sb.WriteString("\t// 检查字段是否允许操作\n")
		sb.WriteString(fmt.Sprintf("\tallowedFields := []string{\"%s\"}\n",
			strings.Join(g.config.OperateFields, "\", \"")))
		sb.WriteString("\tif !datatype.Contains(allowedFields, field) {\n")
		sb.WriteString("\t\treturn fmt.Errorf(\"common.server_error\")\n")
		sb.WriteString("\t}\n")
		sb.WriteString(fmt.Sprintf("\ts.container.GetDB().Model(&%s.%s{}).Where(\"id IN ?\", ids).Update(field, value)\n",
			modelPackage, g.config.StructName))
		sb.WriteString("\treturn nil\n")
		sb.WriteString("}\n")
	}

	return sb.String()
}

// isCommonModel 判断是否为 common/model
func (g *BackendGenerator) isCommonModel() bool {
	// 常见的 common model 表名
	commonTables := []string{"user", "user_money_log", "user_score_log"}
	for _, table := range commonTables {
		if g.config.TableName == table {
			return true
		}
	}
	return false
}

// GenerateModel 生成 Model 代码
func (g *BackendGenerator) GenerateModel() string {
	var sb strings.Builder

	// 包声明
	sb.WriteString("package model\n\n")

	// 导入
	sb.WriteString("import (\n")
	sb.WriteString("\t\"gorm.io/gorm\"\n")
	sb.WriteString(")\n\n")

	// 结构体注释
	sb.WriteString(fmt.Sprintf("// %s %s\n", g.config.StructName, g.config.TableComment))
	sb.WriteString(fmt.Sprintf("type %s struct {\n", g.config.StructName))

	// 字段
	for _, field := range g.config.Fields {
		// 跳过不需要的字段
		if field.ColumnName == "" {
			continue
		}

		// 字段注释
		comment := field.ColumnComment
		if comment == "" {
			comment = field.FieldName
		}

		// 字段定义
		sb.WriteString(fmt.Sprintf("\t%s %s `json:\"%s\" gorm:\"%s\"` // %s\n",
			field.FieldName,
			field.FieldType,
			field.JsonTag,
			field.GormTag,
			comment,
		))
	}

	sb.WriteString("}\n\n")

	// TableName 方法
	sb.WriteString("// TableName 表名\n")
	sb.WriteString(fmt.Sprintf("func (%s) TableName() string {\n", g.config.StructName))
	sb.WriteString(fmt.Sprintf("\treturn \"%s\"\n", g.config.TableName))
	sb.WriteString("}\n\n")

	// BeforeCreate 钩子
	sb.WriteString("// BeforeCreate 创建前的钩子\n")
	sb.WriteString(fmt.Sprintf("func (m *%s) BeforeCreate(tx *gorm.DB) error {\n", g.config.StructName))
	sb.WriteString("\treturn nil\n")
	sb.WriteString("}\n\n")

	// BeforeUpdate 钩子
	sb.WriteString("// BeforeUpdate 更新前的钩子\n")
	sb.WriteString(fmt.Sprintf("func (m *%s) BeforeUpdate(tx *gorm.DB) error {\n", g.config.StructName))
	sb.WriteString("\treturn nil\n")
	sb.WriteString("}\n")

	return sb.String()
}

// GenerateValidate 生成 Validate 代码
func (g *BackendGenerator) GenerateValidate() string {
	var sb strings.Builder

	// 包声明
	sb.WriteString("package validate\n\n")

	// 导入
	sb.WriteString("import (\n")
	sb.WriteString("\t\"github.com/gin-gonic/gin\"\n")
	sb.WriteString("\t\"github.com/lvjiaben/go-wheel/app/common/validator\"\n")
	sb.WriteString(")\n\n")

	// Create 结构体
	if g.config.Methods.Create {
		sb.WriteString(fmt.Sprintf("// %sCreate 创建%s\n", g.config.StructName, g.config.TableComment))
		sb.WriteString(fmt.Sprintf("type %sCreate struct {\n", g.config.StructName))
		for _, field := range g.config.Fields {
			if field.InCreate && !field.IsPrimaryKey && !field.IsAutoIncrement {
				label := field.ColumnComment
				if label == "" {
					label = field.FieldName
				}
				sb.WriteString(fmt.Sprintf("\t%s %s `json:\"%s\" binding:\"%s\" label:\"%s\"`\n",
					field.FieldName,
					field.FieldType,
					field.JsonTag,
					field.ValidateRules,
					label,
				))
			}
		}
		sb.WriteString("}\n\n")
	}

	// Update 结构体
	if g.config.Methods.Update {
		sb.WriteString(fmt.Sprintf("// %sUpdate 更新%s\n", g.config.StructName, g.config.TableComment))
		sb.WriteString(fmt.Sprintf("type %sUpdate struct {\n", g.config.StructName))
		// ID 字段（必须）
		sb.WriteString("\tId int `json:\"id\" binding:\"required,min=1\" label:\"ID\"`\n")
		for _, field := range g.config.Fields {
			if field.InUpdate && !field.IsPrimaryKey && !field.IsAutoIncrement {
				label := field.ColumnComment
				if label == "" {
					label = field.FieldName
				}
				sb.WriteString(fmt.Sprintf("\t%s %s `json:\"%s\" binding:\"%s\" label:\"%s\"`\n",
					field.FieldName,
					field.FieldType,
					field.JsonTag,
					field.ValidateRules,
					label,
				))
			}
		}
		sb.WriteString("}\n\n")
	}

	// Operate 结构体
	if g.config.Methods.Operate {
		sb.WriteString(fmt.Sprintf("// %sOperate 操作%s字段\n", g.config.StructName, g.config.TableComment))
		sb.WriteString(fmt.Sprintf("type %sOperate struct {\n", g.config.StructName))
		sb.WriteString("\tIds   []int  `json:\"ids\" binding:\"required,dive,min=1\" label:\"ID列表\"`\n")
		sb.WriteString("\tField string `json:\"field\" binding:\"required\" label:\"字段名\"`\n")
		sb.WriteString("\tValue int    `json:\"value\" binding:\"min=0\" label:\"字段值\"`\n")
		sb.WriteString("}\n\n")
	}

	// Delete 结构体
	if g.config.Methods.Delete {
		sb.WriteString(fmt.Sprintf("// %sDelete 删除%s\n", g.config.StructName, g.config.TableComment))
		sb.WriteString(fmt.Sprintf("type %sDelete struct {\n", g.config.StructName))
		sb.WriteString("\tIds []int `json:\"ids\" binding:\"required,min=1,dive,min=1\" label:\"ID列表\"`\n")
		sb.WriteString("}\n\n")
	}

	// 验证函数
	if g.config.Methods.Create {
		sb.WriteString(fmt.Sprintf("// Validate%sCreate 验证创建%s\n", g.config.StructName, g.config.TableComment))
		sb.WriteString(fmt.Sprintf("func Validate%sCreate(c *gin.Context) (*%sCreate, bool) {\n",
			g.config.StructName, g.config.StructName))
		sb.WriteString(fmt.Sprintf("\treturn validator.ValidateStructWithConvert[%sCreate](c)\n", g.config.StructName))
		sb.WriteString("}\n\n")
	}

	if g.config.Methods.Update {
		sb.WriteString(fmt.Sprintf("// Validate%sUpdate 验证更新%s\n", g.config.StructName, g.config.TableComment))
		sb.WriteString(fmt.Sprintf("func Validate%sUpdate(c *gin.Context) (*%sUpdate, bool) {\n",
			g.config.StructName, g.config.StructName))
		sb.WriteString(fmt.Sprintf("\treturn validator.ValidateStructWithConvert[%sUpdate](c)\n", g.config.StructName))
		sb.WriteString("}\n\n")
	}

	if g.config.Methods.Operate {
		sb.WriteString(fmt.Sprintf("func Validate%sOperate(c *gin.Context) (*%sOperate, bool) {\n",
			g.config.StructName, g.config.StructName))
		sb.WriteString(fmt.Sprintf("\treturn validator.ValidateStructWithConvert[%sOperate](c)\n", g.config.StructName))
		sb.WriteString("}\n\n")
	}

	if g.config.Methods.Delete {
		sb.WriteString(fmt.Sprintf("// Validate%sDelete 验证删除%s\n", g.config.StructName, g.config.TableComment))
		sb.WriteString(fmt.Sprintf("func Validate%sDelete(c *gin.Context) (*%sDelete, bool) {\n",
			g.config.StructName, g.config.StructName))
		sb.WriteString(fmt.Sprintf("\treturn validator.ValidateStruct[%sDelete](c)\n", g.config.StructName))
		sb.WriteString("}\n")
	}

	return sb.String()
}

// GenerateRoute 生成路由代码
func (g *BackendGenerator) GenerateRoute() string {
	var sb strings.Builder

	// 注释
	sb.WriteString(fmt.Sprintf("\t\t// %s管理\n", g.config.TableComment))
	sb.WriteString(fmt.Sprintf("\t\t%sController := controller.New%sController(c)\n",
		ToCamelCase(g.config.StructName), g.config.StructName))
	sb.WriteString(fmt.Sprintf("\t\t%sGroup := api.Group(\"/%s\").Use(authMiddleware.JWTAuthCheck()).Use(authMiddleware.PermissionCheck())\n",
		ToCamelCase(g.config.StructName), g.config.ModuleName))
	sb.WriteString("\t\t{\n")

	// List 路由
	sb.WriteString(fmt.Sprintf("\t\t\t%sGroup.GET(\"/list\", %sController.List)\n",
		ToCamelCase(g.config.StructName), ToCamelCase(g.config.StructName)))

	// Create 路由
	if g.config.Methods.Create {
		sb.WriteString(fmt.Sprintf("\t\t\t%sGroup.POST(\"/create\", %sController.Create)\n",
			ToCamelCase(g.config.StructName), ToCamelCase(g.config.StructName)))
	}

	// Update 路由
	if g.config.Methods.Update {
		sb.WriteString(fmt.Sprintf("\t\t\t%sGroup.POST(\"/update\", %sController.Update)\n",
			ToCamelCase(g.config.StructName), ToCamelCase(g.config.StructName)))
	}

	// Delete 路由
	if g.config.Methods.Delete {
		sb.WriteString(fmt.Sprintf("\t\t\t%sGroup.POST(\"/delete\", %sController.Delete)\n",
			ToCamelCase(g.config.StructName), ToCamelCase(g.config.StructName)))
	}

	// Operate 路由
	if g.config.Methods.Operate {
		sb.WriteString(fmt.Sprintf("\t\t\t%sGroup.POST(\"/operate\", %sController.Operate)\n",
			ToCamelCase(g.config.StructName), ToCamelCase(g.config.StructName)))
	}

	sb.WriteString("\t\t}\n")

	return sb.String()
}
