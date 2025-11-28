# CRUDBuilder 构建器

`CRUDBuilder` 是一个泛型 CRUD 构建器，提供标准化的增删改查操作。

## 基本用法

### 创建构建器

**推荐方式 - 使用 DB Provider：**

```go
import "github.com/lvjiaben/go-wheel/app/common/builder"

type UserService struct {
    container   *container.Container
    crudService *builder.CRUDBuilder[model.User]
}

func NewUserService(c *container.Container) *UserService {
    return &UserService{
        container: c,
        crudService: builder.NewCRUDBuilderWithProvider[model.User](
            func(ctx *gin.Context) *gorm.DB {
                return c.GetDB().WithContext(ctx.Request.Context())
            },
        ).WithSearchFields("username", "email", "mobile"),
    }
}
```

### List 列表查询

```go
func (s *UserService) List(ctx *gin.Context) map[string]interface{} {
    return s.crudService.WithContext(ctx).List()
}
```

返回格式：

```json
{
  "list": [...],
  "total": 100,
  "page": 1,
  "page_size": 20
}
```

支持的 URL 参数：

| 参数 | 说明 | 示例 |
|------|------|------|
| `page` | 页码 | `?page=1` |
| `page_size` | 每页数量 | `?page_size=20` |
| `search` | 关键词搜索 | `?search=张三` |
| `filter` | 筛选条件 (JSON) | `?filter={"status":1}` |
| `sort_by` | 排序字段 | `?sort_by=created_at` |
| `sort_order` | 排序方向 | `?sort_order=desc` |
| `pid` | 父级 ID | `?pid=0` |

### Create 创建

```go
func (s *UserService) Create(ctx *gin.Context, form *validate.UserCreate) (*model.User, error) {
    return s.crudService.WithContext(ctx).Create(form)
}
```

### Update 更新

```go
func (s *UserService) Update(ctx *gin.Context, form *validate.UserUpdate) (*model.User, error) {
    return s.crudService.WithContext(ctx).Update(form.Id, form)
}
```

### Delete 删除

```go
func (s *UserService) Delete(ctx *gin.Context, ids []int) error {
    return s.crudService.WithContext(ctx).Delete(ids)
}
```

## 配置选项

### WithSearchFields 搜索字段

```go
builder.NewCRUDBuilderWithProvider[model.User](...).
    WithSearchFields("username", "email", "mobile")
```

### WithTransaction 事务控制

```go
// 禁用事务（默认开启）
crudService.WithTransaction(false)
```

### WithPagination 分页控制

```go
// 禁用分页（默认开启）
crudService.WithPagination(false)
```

### WithFilter 筛选器控制

```go
// 禁用筛选器（默认开启）
crudService.WithFilter(false)
```

## 回调钩子

### Before 操作前回调

```go
crudService.Before(func(query interface{}, db *gorm.DB) error {
    // 在操作前执行
    return nil
})
```

### After 操作后回调

```go
crudService.After(func(query interface{}, db *gorm.DB) error {
    // 在操作后执行
    return nil
})
```

## Filter 筛选器

Filter 参数支持多种格式：

```json
// 精确匹配
{"status": 1}

// 范围查询（两个值）
{"created_at": [1700000000, 1700100000]}

// IN 查询（多个值）
{"status": [1, 2, 3]}
```

## 字段验证

构建器会自动验证字段是否在模型中定义，防止非法字段注入：

```go
// 如果传入模型中不存在的字段，会返回错误
err := crudService.ValidateFields(data)
```

## 完整示例

```go
type ArticleService struct {
    container   *container.Container
    crudService *builder.CRUDBuilder[model.Article]
}

func NewArticleService(c *container.Container) *ArticleService {
    return &ArticleService{
        container: c,
        crudService: builder.NewCRUDBuilderWithProvider[model.Article](
            func(ctx *gin.Context) *gorm.DB {
                return c.GetDB().WithContext(ctx.Request.Context())
            },
        ).WithSearchFields("title", "content"),
    }
}

func (s *ArticleService) List(ctx *gin.Context) map[string]interface{} {
    return s.crudService.WithContext(ctx).List()
}

func (s *ArticleService) Create(ctx *gin.Context, form *validate.ArticleCreate) (*model.Article, error) {
    return s.crudService.WithContext(ctx).Create(form)
}
```

