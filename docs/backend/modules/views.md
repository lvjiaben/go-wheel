# Views 视图模板

项目支持 HTML 模板渲染，用于服务端渲染页面。

## 目录结构

```
app/views/
├── index.html          # 首页模板
├── error.html          # 错误页面
└── email/              # 邮件模板
    └── welcome.html
```

## 配置

在 `routes/routes.go` 中加载模板：

```go
func RegisterRoutes(r *gin.Engine, c *container.Container) {
    // 加载 HTML 模板
    r.LoadHTMLGlob("app/views/*")
    
    // 或加载多级目录
    r.LoadHTMLGlob("app/views/**/*")
}
```

## 基本用法

### 渲染模板

```go
// app/api/v1/controller/index.go
func (c *IndexController) Index(ctx *gin.Context) {
    ctx.HTML(http.StatusOK, "index.html", gin.H{
        "title":   "首页",
        "message": "欢迎访问",
        "user":    user,
    })
}
```

### 模板语法

```html
<!-- app/views/index.html -->
<!DOCTYPE html>
<html>
<head>
    <title>{{ .title }}</title>
</head>
<body>
    <h1>{{ .message }}</h1>
    
    <!-- 条件判断 -->
    {{ if .user }}
        <p>欢迎, {{ .user.Name }}</p>
    {{ else }}
        <p>请登录</p>
    {{ end }}
    
    <!-- 循环 -->
    <ul>
    {{ range .items }}
        <li>{{ .Name }} - {{ .Price }}</li>
    {{ end }}
    </ul>
    
    <!-- 模板函数 -->
    <p>{{ .date | formatDate }}</p>
</body>
</html>
```

## 自定义函数

```go
// 注册自定义模板函数
r.SetFuncMap(template.FuncMap{
    "formatDate": func(t time.Time) string {
        return t.Format("2006-01-02 15:04:05")
    },
    "safeHTML": func(s string) template.HTML {
        return template.HTML(s)
    },
})

// 必须在 LoadHTMLGlob 之前调用
r.LoadHTMLGlob("app/views/*")
```

## 模板继承

### 基础模板

```html
<!-- app/views/layouts/base.html -->
{{ define "base" }}
<!DOCTYPE html>
<html>
<head>
    <title>{{ .title }}</title>
    {{ template "head" . }}
</head>
<body>
    {{ template "header" . }}
    
    <main>
        {{ template "content" . }}
    </main>
    
    {{ template "footer" . }}
</body>
</html>
{{ end }}
```

### 子模板

```html
<!-- app/views/home.html -->
{{ template "base" . }}

{{ define "head" }}
<link rel="stylesheet" href="/css/home.css">
{{ end }}

{{ define "content" }}
<h1>首页内容</h1>
{{ end }}
```

## 静态资源

```go
// 配置静态资源目录
r.Static("/public", "./public")
r.Static("/uploads", "./storage/uploads")

// 单个文件
r.StaticFile("/favicon.ico", "./public/favicon.ico")
```

## 错误页面

```go
// 自定义错误处理
r.NoRoute(func(c *gin.Context) {
    c.HTML(http.StatusNotFound, "error.html", gin.H{
        "code":    404,
        "message": "页面不存在",
    })
})

r.NoMethod(func(c *gin.Context) {
    c.HTML(http.StatusMethodNotAllowed, "error.html", gin.H{
        "code":    405,
        "message": "方法不允许",
    })
})
```

## 邮件模板

```go
// 渲染邮件模板
func RenderEmailTemplate(templateName string, data interface{}) (string, error) {
    tmpl, err := template.ParseFiles("app/views/email/" + templateName)
    if err != nil {
        return "", err
    }
    
    var buf bytes.Buffer
    if err := tmpl.Execute(&buf, data); err != nil {
        return "", err
    }
    
    return buf.String(), nil
}

// 使用
html, _ := RenderEmailTemplate("welcome.html", gin.H{
    "username": "张三",
    "link":     "https://example.com/verify",
})
```

## 最佳实践

1. **模板复用** - 使用模板继承减少重复代码
2. **XSS 防护** - 默认自动转义，需要原始 HTML 时使用 `safeHTML`
3. **性能优化** - 生产环境开启模板缓存
4. **分离关注点** - 模板只负责展示，逻辑在控制器中处理

