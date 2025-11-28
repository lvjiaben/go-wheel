# 代码生成器

项目内置 CRUD 代码生成器，可以根据数据库表自动生成前后端代码。

## 功能特点

- 自动读取数据库表结构
- 生成后端 Controller、Service、Model、Validate
- 生成前端 API、列表页、表单页、国际化文件
- 自动生成菜单 SQL
- 支持命令行和 Web 界面两种方式

## 命令行使用

### 运行生成器

```bash
go run app/backend/gen/cmd/main.go
```

### 交互流程

```
=== Go-Wheel CRUD 代码生成器 ===

可用的数据库表：
1. user
2. admin
3. role
4. menu

请选择表（输入序号或表名）：1

表名：user
表注释：用户表

是否自定义字段配置？（默认使用智能识别）[y/N]：n

选择要生成的方法：
[x] List   - 列表查询
[x] Create - 创建
[x] Update - 更新
[x] Delete - 删除
[ ] Operate - 批量操作

是否预览生成的代码？[y/N]：y

确认生成代码？[y/N]：y

✓ 代码生成成功！
✓ 菜单 SQL 已执行！
```

## Web 界面使用

### API 接口

| 接口 | 方法 | 说明 |
|------|------|------|
| `/backend/gen/tables` | GET | 获取所有表 |
| `/backend/gen/table-info` | GET | 获取表详情 |
| `/backend/gen/preview` | POST | 预览代码 |
| `/backend/gen/generate` | POST | 生成代码 |
| `/backend/gen/download` | POST | 下载代码 |
| `/backend/gen/history` | GET | 生成历史 |

### 请求示例

```json
POST /backend/gen/generate
{
    "config": {
        "table_name": "user",
        "table_comment": "用户管理",
        "module_name": "user",
        "struct_name": "User",
        "frontend_src_path": "vben-admin/apps/web-antd/src",
        "methods": {
            "list": true,
            "create": true,
            "update": true,
            "delete": true,
            "operate": false
        },
        "search_fields": ["username", "mobile"],
        "menu_config": {
            "parent_menu_name": "系统管理",
            "menu_name": "用户管理",
            "menu_icon": "user",
            "menu_sort": 10
        }
    }
}
```

## 生成的文件

### 后端文件

```
app/backend/
├── controller/user/
│   └── user.go          # 控制器
├── service/user/
│   └── user.go          # 服务层
├── model/user/
│   └── user.go          # 数据模型
└── validator/user/
    └── user.go          # 验证器
```

### 前端文件

```
vben-admin/apps/web-antd/src/
├── api/user/
│   └── index.ts         # API 接口
├── views/user/
│   ├── index.vue        # 列表页
│   ├── data.ts          # 表格配置
│   └── form.vue         # 表单弹窗
└── locales/langs/
    ├── zh-CN/user.json  # 中文
    └── en-US/user.json  # 英文
```

## 配置说明

### GenConfig 配置

```go
type GenConfig struct {
    TableName       string       // 表名
    TableComment    string       // 表注释
    ModuleName      string       // 模块名
    StructName      string       // 结构体名
    FrontendSrcPath string       // 前端路径
    Methods         MethodConfig // 方法配置
    Fields          []FieldConfig // 字段配置
    SearchFields    []string     // 搜索字段
    MenuConfig      MenuConfig   // 菜单配置
}
```

### 方法配置

```go
type MethodConfig struct {
    List    bool // 列表查询
    Create  bool // 创建
    Update  bool // 更新
    Delete  bool // 删除
    Operate bool // 批量操作
}
```

### 字段配置

```go
type FieldConfig struct {
    ColumnName   string // 列名
    FieldName    string // 字段名
    FieldType    string // 字段类型
    Comment      string // 注释
    ShowInList   bool   // 列表显示
    ShowInForm   bool   // 表单显示
    Searchable   bool   // 可搜索
    Required     bool   // 必填
    FormType     string // 表单类型
}
```

## 智能识别

生成器会自动识别：

- **主键字段** - 自动排除 `id`
- **时间字段** - `created_at`, `updated_at` 等
- **状态字段** - 自动使用 Switch 组件
- **图片字段** - 自动使用上传组件
- **长文本** - 自动使用 Textarea

## 自定义模板

模板文件位于 `app/backend/gen/` 目录：

- `backend_generator.go` - 后端代码模板
- `frontend_generator.go` - 前端代码模板

可以根据项目需求修改模板。

