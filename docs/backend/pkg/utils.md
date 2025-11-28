# Utils 工具函数

项目提供多种常用工具函数，位于 `pkg/utils` 目录。

## 目录结构

```
pkg/utils/
├── crypto/              # 加密相关
│   ├── md5.go          # MD5 加密
│   ├── password.go     # 密码加密
│   └── salt.go         # 盐值生成
├── datatype/            # 数据类型
│   ├── array.go        # 数组操作
│   ├── int.go          # 整数操作
│   ├── string.go       # 字符串操作
│   └── time.go         # 时间操作
├── file/                # 文件操作
│   └── file.go
├── http/                # HTTP 响应
│   └── response.go
└── security/            # 安全相关
    └── jwt.go
```

## 密码加密

```go
import "github.com/lvjiaben/go-wheel/pkg/utils/crypto"

// 生成盐值
salt := crypto.GenerateSalt()

// 加密密码
hashedPassword, err := crypto.HashPassword(password, salt)

// 验证密码
isValid := crypto.PasswordVerifyWithSalt(password, salt, hashedPassword)
```

## MD5 加密

```go
import "github.com/lvjiaben/go-wheel/pkg/utils/crypto"

// MD5 加密
hash := crypto.MD5(str)
```

## HTTP 响应

```go
import "github.com/lvjiaben/go-wheel/pkg/utils/http"

// 成功响应
http.ResponseSuccess(ctx, data)
http.SuccessWithI18n(ctx, "common.success", data)

// 错误响应
http.ResponseError(ctx, "操作失败", nil)
http.ErrorWithI18n(ctx, "common.error", nil)

// 验证错误 (400)
http.ValidateErrorI18n(ctx, "common.invalid_params")

// 认证错误 (401)
http.AuthErrorI18n(ctx, "common.unauthorized")

// 权限错误 (403)
http.ForbiddenErrorI18n(ctx, "common.forbidden")

// 未找到 (404)
http.NotFoundErrorI18n(ctx, "common.not_found")
```

### 响应格式

```json
{
    "code": 200,
    "message": "操作成功",
    "data": {}
}
```

## 数组操作

```go
import "github.com/lvjiaben/go-wheel/pkg/utils/datatype"

// 判断元素是否在数组中
exists := datatype.InArray("a", []string{"a", "b", "c"})

// 数组去重
unique := datatype.ArrayUnique([]int{1, 2, 2, 3})

// 数组差集
diff := datatype.ArrayDiff([]int{1, 2, 3}, []int{2, 3, 4})
```

## 字符串操作

```go
import "github.com/lvjiaben/go-wheel/pkg/utils/datatype"

// 生成随机字符串
str := datatype.RandomString(16)

// 驼峰转下划线
snake := datatype.CamelToSnake("UserName")  // user_name

// 下划线转驼峰
camel := datatype.SnakeToCamel("user_name")  // UserName
```

## 时间操作

```go
import "github.com/lvjiaben/go-wheel/pkg/utils/datatype"

// 格式化时间
str := datatype.FormatTime(time.Now())  // 2024-01-01 12:00:00

// 解析时间
t := datatype.ParseTime("2024-01-01 12:00:00")

// 获取今天开始时间
start := datatype.TodayStart()

// 获取今天结束时间
end := datatype.TodayEnd()
```

## 文件操作

```go
import "github.com/lvjiaben/go-wheel/pkg/utils/file"

// 判断文件是否存在
exists := file.Exists("/path/to/file")

// 创建目录
file.MkdirAll("/path/to/dir")

// 获取文件扩展名
ext := file.GetExt("image.png")  // png

// 获取文件大小
size := file.GetSize("/path/to/file")
```

## 获取上下文信息

```go
import "github.com/lvjiaben/go-wheel/pkg/utils/http"

// 获取语言
lang := http.GetLang(ctx)  // "zh-CN" 或 "en-US"

// 获取容器
container := http.GetContainer(ctx)
```

## 使用示例

### 用户注册

```go
func (s *UserService) Register(username, password string) error {
    // 生成盐值
    salt := crypto.GenerateSalt()
    
    // 加密密码
    hashedPassword, err := crypto.HashPassword(password, salt)
    if err != nil {
        return err
    }
    
    // 保存用户
    user := &model.User{
        Username: username,
        Password: hashedPassword,
        Salt:     salt,
    }
    
    return s.db.Create(user).Error
}
```

### 用户登录

```go
func (s *UserService) Login(username, password string) (*model.User, error) {
    var user model.User
    if err := s.db.Where("username = ?", username).First(&user).Error; err != nil {
        return nil, err
    }
    
    // 验证密码
    if !crypto.PasswordVerifyWithSalt(password, user.Salt, user.Password) {
        return nil, errors.New("密码错误")
    }
    
    return &user, nil
}
```

