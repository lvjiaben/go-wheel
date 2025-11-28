# 后端 CRUD 模板

标准的后端 CRUD 模块包含：Model、Validate、Service、Controller、Route。

## 目录结构

```
app/backend/
├── controller/
│   └── article/
│       └── article.go      # 控制器
├── model/
│   └── article.go          # 模型
├── service/
│   └── article/
│       └── article.go      # 服务
├── validate/
│   └── article.go          # 验证器
└── route/
    └── article.go          # 路由
```

## 1. Model 模型

```go
// app/backend/model/article.go
package model

import "time"

type Article struct {
    Id        int       `json:"id" gorm:"primaryKey;autoIncrement"`
    Title     string    `json:"title" gorm:"size:200;not null;comment:标题"`
    Content   string    `json:"content" gorm:"type:text;comment:内容"`
    Status    int       `json:"status" gorm:"default:1;comment:状态 0禁用 1启用"`
    Sort      int       `json:"sort" gorm:"default:0;comment:排序"`
    CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
    UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
}

func (Article) TableName() string {
    return "article"
}
```

## 2. Validate 验证器

```go
// app/backend/validate/article.go
package validate

import (
    "github.com/gin-gonic/gin"
    "github.com/lvjiaben/go-wheel/app/common/validator"
)

type ArticleCreate struct {
    Title   string `json:"title" binding:"required,max=200" msg:"article.title_required"`
    Content string `json:"content" binding:"required" msg:"article.content_required"`
    Status  int    `json:"status"`
    Sort    int    `json:"sort"`
}

func ValidateArticleCreate(ctx *gin.Context) (*ArticleCreate, bool) {
    return validator.ValidateStructWithConvert[ArticleCreate](ctx)
}

type ArticleUpdate struct {
    Id      int    `json:"id" binding:"required" msg:"common.id_required"`
    Title   string `json:"title" binding:"max=200"`
    Content string `json:"content"`
    Status  int    `json:"status"`
    Sort    int    `json:"sort"`
}

func ValidateArticleUpdate(ctx *gin.Context) (*ArticleUpdate, bool) {
    return validator.ValidateStructWithConvert[ArticleUpdate](ctx)
}

type ArticleDelete struct {
    Ids []int `json:"ids" binding:"required,min=1" msg:"common.ids_required"`
}

func ValidateArticleDelete(ctx *gin.Context) (*ArticleDelete, bool) {
    return validator.ValidateStruct[ArticleDelete](ctx)
}
```

## 3. Service 服务

```go
// app/backend/service/article/article.go
package article

import (
    "github.com/gin-gonic/gin"
    "github.com/lvjiaben/go-wheel/app/backend/model"
    "github.com/lvjiaben/go-wheel/app/backend/validate"
    "github.com/lvjiaben/go-wheel/app/common/builder"
    "github.com/lvjiaben/go-wheel/pkg/container"
    "gorm.io/gorm"
)

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

func (s *ArticleService) Update(ctx *gin.Context, form *validate.ArticleUpdate) (*model.Article, error) {
    return s.crudService.WithContext(ctx).Update(form.Id, form)
}

func (s *ArticleService) Delete(ctx *gin.Context, ids []int) error {
    return s.crudService.WithContext(ctx).Delete(ids)
}
```

## 4. Controller 控制器

```go
// app/backend/controller/article/article.go
package article

import (
    "github.com/gin-gonic/gin"
    serviceArticle "github.com/lvjiaben/go-wheel/app/backend/service/article"
    "github.com/lvjiaben/go-wheel/app/backend/validate"
    "github.com/lvjiaben/go-wheel/pkg/container"
    "github.com/lvjiaben/go-wheel/pkg/utils/http"
)

type ArticleController struct {
    container *container.Container
    service   *serviceArticle.ArticleService
}

func NewArticleController(c *container.Container) *ArticleController {
    return &ArticleController{
        container: c,
        service:   serviceArticle.NewArticleService(c),
    }
}

func (ctrl *ArticleController) List(ctx *gin.Context) {
    http.SuccessWithI18n(ctx, "common.success", ctrl.service.List(ctx))
}

func (ctrl *ArticleController) Create(ctx *gin.Context) {
    form, valid := validate.ValidateArticleCreate(ctx)
    if !valid {
        return
    }
    res, err := ctrl.service.Create(ctx, form)
    if err != nil {
        http.ErrorWithI18n(ctx, err.Error(), nil)
        return
    }
    http.SuccessWithI18n(ctx, "common.success", res)
}

func (ctrl *ArticleController) Update(ctx *gin.Context) {
    form, valid := validate.ValidateArticleUpdate(ctx)
    if !valid {
        return
    }
    res, err := ctrl.service.Update(ctx, form)
    if err != nil {
        http.ErrorWithI18n(ctx, err.Error(), nil)
        return
    }
    http.SuccessWithI18n(ctx, "common.success", res)
}

func (ctrl *ArticleController) Delete(ctx *gin.Context) {
    form, valid := validate.ValidateArticleDelete(ctx)
    if !valid {
        return
    }
    if err := ctrl.service.Delete(ctx, form.Ids); err != nil {
        http.ErrorWithI18n(ctx, err.Error(), nil)
        return
    }
    http.SuccessWithI18n(ctx, "common.success", nil)
}
```

## 5. Route 路由

```go
// app/backend/route/article.go
package route

import (
    "github.com/gin-gonic/gin"
    "github.com/lvjiaben/go-wheel/app/backend/controller/article"
    "github.com/lvjiaben/go-wheel/pkg/container"
)

func RegisterArticleRoutes(r *gin.RouterGroup, c *container.Container) {
    ctrl := article.NewArticleController(c)

    group := r.Group("/article")
    {
        group.GET("", ctrl.List)
        group.POST("", ctrl.Create)
        group.PUT("", ctrl.Update)
        group.DELETE("", ctrl.Delete)
    }
}
```

## 6. 注册路由

在 `app/backend/route/route.go` 中注册：

```go
func RegisterRoutes(r *gin.Engine, c *container.Container) {
    api := r.Group("/backend")
    api.Use(middleware.Auth(c))

    // 注册文章路由
    RegisterArticleRoutes(api, c)
}
```

## 7. 国际化消息

在 `configs/i18n/zh-CN.yaml` 添加：

```yaml
article:
  title_required: "标题不能为空"
  content_required: "内容不能为空"
```

在 `configs/i18n/en-US.yaml` 添加：

```yaml
article:
  title_required: "Title is required"
  content_required: "Content is required"
```

