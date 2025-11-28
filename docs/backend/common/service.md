# Common Service

项目在 `app/common/service` 中提供了多个通用服务。

## SmsService 短信服务

支持阿里云、腾讯云、云片、短信宝等多个短信平台。

### 初始化

```go
import "github.com/lvjiaben/go-wheel/app/common/service"

smsService := service.NewSmsService(container)
```

### 发送验证码

```go
// 发送验证码（自动生成）
err := smsService.Send(mobile, "", "login")

// 发送指定验证码
err := smsService.Send(mobile, "123456", "register")
```

事件类型（event）用于区分不同场景：
- `login` - 登录
- `register` - 注册
- `reset_password` - 重置密码
- `change_mobile` - 修改手机号

### 验证验证码

```go
if smsService.Verify(mobile, code, "login") {
    // 验证成功
}
```

### 删除验证码

```go
smsService.Delete(mobile, "login")
```

### 配置

在系统配置中设置：

| 配置键 | 说明 |
|--------|------|
| `sms_type` | 短信平台：aliyun/tencent/yunpian/smsbao |
| `sms_id` | Access Key ID |
| `sms_key` | Access Key Secret |
| `sms_token` | 签名名称 |
| `sms_template` | 默认模板 ID |
| `sms_template_login` | 登录模板 ID |

## UploadService 上传服务

支持本地存储、阿里云 OSS、腾讯云 COS、七牛云。

### 初始化

```go
uploadService := service.NewUploadService(container)
```

### 上传文件

```go
func (ctrl *Controller) Upload(ctx *gin.Context) {
    file, err := ctx.FormFile("file")
    if err != nil {
        http.ErrorWithI18n(ctx, "common.file_required", nil)
        return
    }

    result, err := uploadService.Upload(file)
    if err != nil {
        http.ErrorWithI18n(ctx, err.Error(), nil)
        return
    }

    http.SuccessWithI18n(ctx, "common.success", result)
}
```

返回结果：

```go
type UploadResult struct {
    Path      string // 存储路径
    Parent    string // 父级文件夹
    URL       string // 在线 HTTP 链接
    Filename  string // 文件名称
    Size      int64  // 文件大小
    MediaType string // 文件类型
    Extension string // 文件后缀
}
```

### 删除文件

```go
uploadService.Delete(path, "local")  // 本地
uploadService.Delete(path, "oss")    // 阿里云 OSS
uploadService.Delete(path, "cos")    // 腾讯云 COS
uploadService.Delete(path, "qiniu")  // 七牛云
```

### 配置

```yaml
upload:
  type: "local"           # local/oss/cos/qiniu
  base_url: "http://localhost:8801"
  upload_path: "uploads"
  max_size: 10485760      # 10MB
  allowed_extensions:
    - jpg
    - jpeg
    - png
    - gif
    - pdf
  allowed_types:
    - image/jpeg
    - image/png
    - image/gif
    - application/pdf
  
  # 阿里云 OSS
  oss:
    endpoint: ""
    access_key_id: ""
    access_key_secret: ""
    bucket_name: ""
  
  # 腾讯云 COS
  cos:
    bucket: ""
    region: ""
    secret_id: ""
    secret_key: ""
  
  # 七牛云
  qiniu:
    access_key: ""
    secret_key: ""
    bucket: ""
```

## CodeGeneratorService 验证码生成服务

生成各种类型的验证码和随机字符串。

```go
codeGenerator := service.NewCodeGeneratorService(container)

// 生成邀请码（10位大写字母+数字）
inviteCode := codeGenerator.GenerateInviteCode()

// 生成随机密码（12位字母+数字+特殊字符）
password := codeGenerator.GenerateRandomPassword()

// 生成纯数字验证码
numericCode := codeGenerator.GenerateNumericCode(6)

// 生成字母数字混合验证码
alphanumericCode := codeGenerator.GenerateAlphanumericCode(8, true)
```

