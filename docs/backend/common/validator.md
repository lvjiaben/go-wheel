# Validator 验证器

项目封装了通用验证器，支持泛型和自动类型转换。

## 基本用法

### ValidateStruct 通用验证

```go
import "github.com/lvjiaben/go-wheel/app/common/validator"

// 定义验证结构体
type LoginForm struct {
    Username string `json:"username" binding:"required" msg:"auth.username_required"`
    Password string `json:"password" binding:"required,min=6" msg:"auth.password_required"`
}

// 在控制器中使用
func (ctrl *AuthController) Login(ctx *gin.Context) {
    form, valid := validator.ValidateStruct[LoginForm](ctx)
    if !valid {
        return  // 验证失败会自动返回错误响应
    }
    // 使用 form.Username, form.Password
}
```

### ValidateStructWithConvert 带类型转换

自动将前端传来的字符串转换为结构体定义的类型：

```go
type UserForm struct {
    Status int    `json:"status"`  // 前端传 "1"，自动转换为 int
    Score  float64 `json:"score"` // 前端传 "99.5"，自动转换为 float64
    Active bool   `json:"active"` // 前端传 "true"，自动转换为 bool
}

form, valid := validator.ValidateStructWithConvert[UserForm](ctx)
```

## 结构体标签

### binding 验证规则

```go
type Form struct {
    Name     string `binding:"required"`           // 必填
    Email    string `binding:"required,email"`     // 必填，邮箱格式
    Age      int    `binding:"min=1,max=150"`     // 最小1，最大150
    Password string `binding:"required,min=6"`    // 必填，最少6位
    Phone    string `binding:"len=11"`             // 固定11位
    URL      string `binding:"url"`                // URL 格式
}
```

### msg 错误消息

`msg` 标签指定验证失败时的国际化消息键：

```go
type Form struct {
    Username string `binding:"required" msg:"user.username_required"`
    Email    string `binding:"email" msg:"user.email_invalid"`
}
```

### json 字段映射

```go
type Form struct {
    UserId   int    `json:"user_id"`    // JSON 字段名为 user_id
    UserName string `json:"username"`   // JSON 字段名为 username
}
```

## 自动请求方法识别

验证器会根据 HTTP 方法自动选择绑定方式：

| 方法 | 绑定方式 |
|------|----------|
| GET | Query 参数 |
| DELETE | Query 参数 |
| POST | JSON Body |
| PUT | JSON Body |
| PATCH | JSON Body |

## 类型转换支持

`ValidateStructWithConvert` 支持以下类型转换：

| 源类型 | 目标类型 |
|--------|----------|
| string | int/int8/int16/int32/int64 |
| string | uint/uint8/uint16/uint32/uint64 |
| string | float32/float64 |
| string | bool |
| float64 | int/uint |
| bool | int（true=1, false=0）|
| int | string |

## 验证函数示例

为每个验证器创建独立的验证函数：

```go
// app/backend/validate/user.go
package validate

import (
    "github.com/gin-gonic/gin"
    "github.com/lvjiaben/go-wheel/app/common/validator"
)

// UserCreate 创建用户验证
type UserCreate struct {
    Username string `json:"username" binding:"required,min=3,max=50" msg:"user.username_required"`
    Password string `json:"password" binding:"required,min=6" msg:"user.password_required"`
    Email    string `json:"email" binding:"email" msg:"user.email_invalid"`
    Mobile   string `json:"mobile" binding:"len=11" msg:"user.mobile_invalid"`
    Status   int    `json:"status"`
}

func ValidateUserCreate(ctx *gin.Context) (*UserCreate, bool) {
    return validator.ValidateStructWithConvert[UserCreate](ctx)
}

// UserUpdate 更新用户验证
type UserUpdate struct {
    Id       int    `json:"id" binding:"required" msg:"common.id_required"`
    Username string `json:"username" binding:"min=3,max=50" msg:"user.username_required"`
    Email    string `json:"email" binding:"omitempty,email" msg:"user.email_invalid"`
    Mobile   string `json:"mobile" binding:"omitempty,len=11" msg:"user.mobile_invalid"`
    Status   int    `json:"status"`
}

func ValidateUserUpdate(ctx *gin.Context) (*UserUpdate, bool) {
    return validator.ValidateStructWithConvert[UserUpdate](ctx)
}

// UserDelete 删除用户验证
type UserDelete struct {
    Ids []int `json:"ids" binding:"required,min=1" msg:"common.ids_required"`
}

func ValidateUserDelete(ctx *gin.Context) (*UserDelete, bool) {
    return validator.ValidateStruct[UserDelete](ctx)
}
```

## 在控制器中使用

```go
func (ctrl *UserController) Create(ctx *gin.Context) {
    form, valid := validate.ValidateUserCreate(ctx)
    if !valid {
        return
    }
    
    user, err := ctrl.service.Create(ctx, form)
    if err != nil {
        http.ErrorWithI18n(ctx, err.Error(), nil)
        return
    }
    
    http.SuccessWithI18n(ctx, "common.success", user)
}
```

