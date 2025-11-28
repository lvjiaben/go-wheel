# HTTP 响应

项目封装了统一的 HTTP 响应方法，支持国际化。

## 响应结构

```json
{
  "code": 200,
  "message": "操作成功",
  "data": {}
}
```

## 响应方法

### SuccessWithI18n

成功响应，默认 code 为 200：

```go
import "github.com/lvjiaben/go-wheel/pkg/utils/http"

func (ctrl *Controller) Get(ctx *gin.Context) {
    data := map[string]interface{}{
        "name": "张三",
        "age":  18,
    }
    http.SuccessWithI18n(ctx, "common.success", data)
}
```

自定义状态码：

```go
http.SuccessWithI18n(ctx, "common.created", data, 201)
```

### ErrorWithI18n

错误响应，默认 code 为 500：

```go
http.ErrorWithI18n(ctx, "common.error", nil)

// 自定义状态码
http.ErrorWithI18n(ctx, "user.not_found", nil, 404)
```

### 常用快捷方法

```go
// 验证错误（400）
http.ValidateErrorI18n(ctx, "common.invalid_params")

// 认证错误（401）
http.AuthErrorI18n(ctx, "auth.unauthorized")

// 权限错误（403）
http.ForbiddenErrorI18n(ctx, "auth.forbidden")

// 未找到错误（404）
http.NotFoundErrorI18n(ctx, "common.not_found")
```

## 非国际化响应

如果不需要国际化：

```go
// 成功响应
http.ResponseSuccess(ctx, data)

// 错误响应
http.ResponseError(ctx, "操作失败", nil)
```

## 工具方法

### 获取语言

```go
lang := http.GetLang(ctx)  // 返回 "zh-CN" 或 "en-US"
```

### 获取容器

```go
container := http.GetContainer(ctx)
```

## 消息键配置

消息键定义在 `configs/i18n/` 目录下：

**zh-CN.yaml:**
```yaml
common:
  success: "操作成功"
  error: "操作失败"
  invalid_params: "参数无效"
  not_found: "数据不存在"

auth:
  unauthorized: "请先登录"
  forbidden: "权限不足"

user:
  not_found: "用户不存在"
  create_success: "用户创建成功"
```

**en-US.yaml:**
```yaml
common:
  success: "Success"
  error: "Error"
  invalid_params: "Invalid parameters"
  not_found: "Not found"

auth:
  unauthorized: "Please login first"
  forbidden: "Access denied"

user:
  not_found: "User not found"
  create_success: "User created successfully"
```

## 最佳实践

1. **统一使用国际化方法** - 方便后续多语言支持
2. **合理的错误码** - 400 参数错误，401 未登录，403 无权限，404 不存在，500 服务器错误
3. **消息键命名** - 使用 `模块.操作` 格式，如 `user.create_success`
4. **data 字段** - 成功时返回数据，错误时可返回详细信息或 `nil`

