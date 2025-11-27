# Go-Wheel CRUD 代码生成器

类似于 PHP FastAdmin 的一键 CRUD 代码生成功能，支持命令行和可视化界面两种模式。

## 功能特性

### 1. 双模式支持
- **命令行模式**：交互式命令行工具，逐步引导配置
- **可视化界面**：Web 界面配置，支持代码预览和下载

### 2. 智能字段识别
自动识别字段类型和用途：
- **Operate 字段**：`status`、`is_xxx` 等状态字段
- **排序字段**：`sort`、`weigh` 等排序字段
- **时间字段**：`xxx_at`、`xxx_time` 等时间字段
- **关联字段**：`xxx_id` 等外键字段
- **枚举字段**：`enum`、`set` 类型字段
- **图片字段**：`image`、`avatar`、`img` 等图片字段
- **文本字段**：`text`、`longtext` 等长文本字段
- **布尔字段**：`tinyint(1)` 等布尔字段

### 3. 生成内容

#### 后端（Go）
- **Controller**：`app/backend/controller/{package}/{package}.go`
- **Service**：`app/backend/service/{package}/{package}.go`
- **Model**：`app/backend/model/{package}.go`
- **Validate**：`app/backend/validate/{package}/{package}.go`
- **Route**：自动注册到 `routes/routes.go` 的 `region:backend-routes` 区域

#### 前端（Vue3 + TypeScript）
- **API**：`api/{module}.ts`
- **ListView**：`views/{module}/list.vue`
- **DataTS**：`views/{module}/data.ts`
- **FormVue**：`views/{module}/modules/form.vue`

#### 菜单
- 生成菜单 SQL，包含父级菜单检查和权限按钮

### 4. 方法配置
- **List**：列表查询（必选）
- **Create**：创建记录（可选）
- **Update**：更新记录（可选）
- **Delete**：删除记录（可选）
- **Operate**：批量操作字段（可选，需要配置 Operate 字段）

### 5. 字段配置
- **模糊搜索字段**：配置哪些字段可以进行模糊搜索
- **Operate 字段**：配置哪些字段可以通过 Operate 方法批量修改
- **表格显示字段**：配置哪些字段在列表中显示
- **表格排序字段**：配置哪些字段可以排序
- **表格搜索字段**：配置哪些字段可以在表格中搜索
- **表单字段**：配置哪些字段在创建/编辑表单中显示

## 使用方法

### 命令行模式

```bash
# 编译命令行工具
go build -o gen cmd/gen/main.go

# 运行
./gen
```

交互式步骤：
1. 选择数据库表
2. 配置模块名、结构体名、包名
3. 配置搜索字段
4. 配置 Operate 字段
5. 选择要生成的方法
6. 配置菜单信息
7. 预览代码
8. 确认生成

### 可视化界面模式

#### 后端接口

已在 `routes/routes.go` 中注册以下接口：

```
GET  /system/gen/table-list      # 获取数据库表列表
GET  /system/gen/table-info      # 获取表详细信息
GET  /system/gen/table-config    # 获取表的默认配置
POST /system/gen/preview         # 预览生成的代码
POST /system/gen/generate        # 生成代码并写入文件
GET  /system/gen/history         # 获取生成历史
POST /system/gen/delete          # 删除生成的代码
POST /system/gen/download        # 下载生成的代码（ZIP）
```

#### 前端页面（待实现）

位置：`vben-admin/apps/web-antd/src/views/system/gen/`

功能：
- 表选择和字段配置
- 方法选择
- 代码预览
- 生成和下载

## 配置结构

### GenConfig

```go
type GenConfig struct {
    TableName        string       // 表名
    TableComment     string       // 表注释
    ModuleName       string       // 模块名（用于路由和前端目录）
    StructName       string       // 结构体名（PascalCase）
    PackageName      string       // 包名（snake_case）
    FrontendSrcPath  string       // 前端 src 目录绝对路径
    Methods          MethodConfig // 方法配置
    Fields           []FieldConfig // 字段配置
    SearchFields     []string     // 模糊搜索字段
    OperateFields    []string     // Operate 允许操作的字段
    MenuConfig       MenuConfig   // 菜单配置
}
```

### FieldConfig

```go
type FieldConfig struct {
    // 数据库信息
    ColumnName       string
    ColumnType       string
    ColumnComment    string
    IsNullable       bool
    IsPrimaryKey     bool
    IsAutoIncrement  bool
    
    // Go 信息
    FieldName        string
    FieldType        string
    JsonTag          string
    GormTag          string
    
    // 验证规则
    InCreate         bool   // 是否在创建时使用
    InUpdate         bool   // 是否在更新时使用
    IsRequired       bool   // 是否必填
    ValidateRules    string // 验证规则
    
    // 前端表格配置
    ShowInTable      bool   // 是否在表格中显示
    TableSortable    bool   // 是否可排序
    TableSearchable  bool   // 是否可搜索
    TableSearchType  string // 搜索类型
    TableDisplayType string // 显示类型
    
    // 前端表单配置
    ShowInForm       bool   // 是否在表单中显示
    FormComponent    string // 表单组件
    FormComponentProps map[string]interface{} // 组件属性
    
    // 智能识别标记
    IsOperateField   bool   // 是否为 Operate 字段
    IsSortField      bool   // 是否为排序字段
    IsTimeField      bool   // 是否为时间字段
    IsRelationField  bool   // 是否为关联字段
    IsEnumField      bool   // 是否为枚举字段
    IsSetField       bool   // 是否为集合字段
    IsTextField      bool   // 是否为文本字段
    IsBoolField      bool   // 是否为布尔字段
    IsImageField     bool   // 是否为图片字段
    IsImagesField    bool   // 是否为多图字段
}
```

## 代码示例

### 生成的 Controller

```go
package user_info

import (
    "github.com/gin-gonic/gin"
    serviceUserInfo "github.com/lvjiaben/go-wheel/app/backend/service/user_info"
    validateUserInfo "github.com/lvjiaben/go-wheel/app/backend/validate/user_info"
    "github.com/lvjiaben/go-wheel/pkg/container"
    "github.com/lvjiaben/go-wheel/pkg/utils/http"
)

type UserInfoController struct {
    container *container.Container
    userInfoService *serviceUserInfo.UserInfoService
}

func NewUserInfoController(c *container.Container) *UserInfoController {
    return &UserInfoController{
        container: c,
        userInfoService: serviceUserInfo.NewUserInfoService(c),
    }
}

func (ctrl *UserInfoController) List(ctx *gin.Context) {
    http.SuccessWithI18n(ctx, "common.success", ctrl.userInfoService.List(ctx))
}

func (ctrl *UserInfoController) Create(ctx *gin.Context) {
    form, valid := validateUserInfo.ValidateUserInfoCreate(ctx)
    if !valid {
        return
    }
    http.SuccessWithI18n(ctx, "common.success", ctrl.userInfoService.Create(ctx, form))
}

// ... 其他方法
```

### 生成的 Service

```go
package user_info

import (
    "github.com/gin-gonic/gin"
    "github.com/lvjiaben/go-wheel/app/backend/model"
    "github.com/lvjiaben/go-wheel/pkg/builder"
    "github.com/lvjiaben/go-wheel/pkg/container"
    "gorm.io/gorm"
)

type UserInfoService struct {
    container   *container.Container
    crudService *builder.CRUDBuilder[model.UserInfo]
}

func NewUserInfoService(c *container.Container) *UserInfoService {
    return &UserInfoService{
        container: c,
        crudService: builder.NewCRUDBuilderWithProvider[model.UserInfo](
            func(ctx *gin.Context) *gorm.DB {
                return c.GetDB().WithContext(ctx.Request.Context())
            },
        ).WithSearchFields("username", "email"),
    }
}

func (s *UserInfoService) List(ctx *gin.Context) map[string]interface{} {
    return s.crudService.WithContext(ctx).List()
}

// ... 其他方法
```

## 注意事项

1. **表名命名**：表名使用下划线命名（如 `user_info`），生成的目录和文件名保持下划线，不转换为多级目录
2. **前端路径**：生成前端代码前必须设置正确的前端 `src` 路径
3. **路由区域**：不要删除 `routes/routes.go` 中的 `region:backend-routes` 注释
4. **删除功能**：删除生成的代码会清理所有相关文件，但不会删除路由注册代码（需手动删除）
5. **菜单 SQL**：生成后需要手动执行菜单 SQL 来创建菜单和权限

## 扩展开发

### 添加新的字段类型识别

在 `utils.go` 中添加识别函数：

```go
func IsXxxField(columnName string) bool {
    return strings.Contains(columnName, "xxx")
}
```

### 自定义代码模板

修改 `backend_generator.go` 和 `frontend_generator.go` 中的生成方法。

### 添加新的表单组件

在 `utils.go` 的 `GetDefaultFormComponent()` 中添加映射规则。

## 待完善功能

- [ ] 前端代码生成器完整实现
- [ ] 代码下载功能（ZIP 打包）
- [ ] 路由删除功能（从 routes.go 中移除）
- [ ] 更多字段类型支持
- [ ] 自定义代码模板
- [ ] 批量生成多个表
- [ ] 代码生成历史管理
- [ ] 可视化界面前端页面

## 参考

- [FastAdmin CRUD 文档](https://doc.fastadmin.net/doc/crud.html#toc-4)
- 项目现有代码结构和规范

