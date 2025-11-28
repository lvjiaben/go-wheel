# 日志系统

项目使用 [Zap](https://github.com/uber-go/zap) 作为日志库，支持日志分级、文件轮转、JSON 格式输出。

## 配置

在 `configs/config.yaml` 中配置：

```yaml
log:
  level: "debug"           # 日志级别: debug, info, warn, error
  format: "json"           # 输出格式: json, console
  output: "file"           # 输出方式: console, file, both
  file_path: "storage/logs/app.log"  # 日志文件路径
  max_size: 100            # 单个文件最大大小 (MB)
  max_backups: 10          # 保留的旧文件数量
  max_age: 30              # 保留天数
  compress: true           # 是否压缩旧文件
```

## 基本用法

### 获取 Logger

```go
// 从容器获取
logger := container.GetLogger()

// 在控制器中
func (c *MyController) MyMethod(ctx *gin.Context) {
    c.container.GetLogger().Info("处理请求")
}
```

### 日志级别

```go
logger.Debug("调试信息", zap.String("key", "value"))
logger.Info("普通信息", zap.Int("count", 100))
logger.Warn("警告信息", zap.Error(err))
logger.Error("错误信息", zap.Error(err))
```

### 结构化日志

```go
logger.Info("用户登录",
    zap.String("username", "admin"),
    zap.String("ip", "192.168.1.1"),
    zap.Int("user_id", 1),
    zap.Duration("duration", time.Second),
)
```

### 带上下文的日志

```go
// 创建带字段的 logger
userLogger := logger.With(
    zap.Int("user_id", userId),
    zap.String("request_id", requestId),
)

userLogger.Info("执行操作")
userLogger.Error("操作失败", zap.Error(err))
```

## 日志格式

### JSON 格式（生产环境推荐）

```json
{
  "level": "info",
  "ts": "2024-01-01T12:00:00.000Z",
  "caller": "controller/user.go:42",
  "msg": "用户登录",
  "username": "admin",
  "ip": "192.168.1.1"
}
```

### Console 格式（开发环境）

```
2024-01-01T12:00:00.000+0800	INFO	controller/user.go:42	用户登录	{"username": "admin", "ip": "192.168.1.1"}
```

## 日志轮转

使用 [lumberjack](https://github.com/natefinch/lumberjack) 实现日志轮转：

- 按文件大小轮转
- 自动压缩旧文件
- 自动删除过期文件

## 最佳实践

### 1. 使用结构化字段

```go
// ✅ 推荐
logger.Info("创建订单", zap.String("order_id", orderId), zap.Float64("amount", 99.9))

// ❌ 不推荐
logger.Info(fmt.Sprintf("创建订单 %s, 金额 %.2f", orderId, amount))
```

### 2. 错误日志包含堆栈

```go
logger.Error("数据库查询失败",
    zap.Error(err),
    zap.String("sql", sql),
    zap.Any("params", params),
)
```

### 3. 敏感信息脱敏

```go
// ✅ 脱敏处理
logger.Info("用户登录", zap.String("mobile", "138****1234"))

// ❌ 不要记录敏感信息
logger.Info("用户登录", zap.String("password", password))
```

### 4. 合理使用日志级别

| 级别 | 使用场景 |
|------|----------|
| Debug | 开发调试信息 |
| Info | 正常业务流程 |
| Warn | 可恢复的异常 |
| Error | 需要关注的错误 |

## 日志文件位置

默认日志文件位于 `storage/logs/` 目录：

```
storage/logs/
├── app.log           # 当前日志
├── app-2024-01-01.log.gz  # 历史日志（压缩）
└── app-2024-01-02.log.gz
```

