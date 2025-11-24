package constants

import "time"

// 时间相关常量
const (
	// JWTExpireDays JWT 默认过期天数
	JWTExpireDays = 7
	
	// LoginFailureLockDuration 登录失败锁定时长
	LoginFailureLockDuration = 30 * time.Minute
	
	// LoginFailureRecordDuration 登录失败记录保存时长
	LoginFailureRecordDuration = 15 * time.Minute
	
	// CaptchaTTL 验证码有效期
	CaptchaTTL = 5 * time.Minute
	
	// GracefulShutdownTimeout 优雅关闭超时时间
	GracefulShutdownTimeout = 5 * time.Second
)

// 登录相关常量
const (
	// DefaultMaxLoginFailures 默认最大登录失败次数
	DefaultMaxLoginFailures = 5
)

// 密码相关常量
const (
	// BcryptCost bcrypt 加密成本
	BcryptCost = 10
	
	// SaltLength 盐值长度（字节）
	SaltLength = 16
	
	// RandomPasswordLength 随机密码长度
	RandomPasswordLength = 12
	
	// InviteCodeLength 邀请码长度
	InviteCodeLength = 10
)

// 文件上传相关常量
const (
	// DefaultMaxFileSize 默认最大文件大小（10MB）
	DefaultMaxFileSize = 10 * 1024 * 1024
	
	// FilePermission 文件权限
	FilePermission = 0755
)

// 验证码相关常量
const (
	// CaptchaHeight 验证码高度
	CaptchaHeight = 80
	
	// CaptchaWidth 验证码宽度
	CaptchaWidth = 240
	
	// CaptchaLength 验证码长度
	CaptchaLength = 4
)

// 重试相关常量
const (
	// MaxInviteCodeRetries 邀请码生成最大重试次数
	MaxInviteCodeRetries = 3

	// DefaultMaxRetries 默认最大重试次数
	DefaultMaxRetries = 3

	// DefaultInitialDelay 默认初始延迟
	DefaultInitialDelay = 5 * time.Second

	// DefaultMaxRetryDelay 默认最大重试延迟
	DefaultMaxRetryDelay = 30 * time.Second
)

// 熔断器相关常量
const (
	// DefaultCircuitBreakerThreshold 默认熔断器阈值（失败次数）
	DefaultCircuitBreakerThreshold = 3

	// DefaultCircuitBreakerTimeout 默认熔断器超时时间
	DefaultCircuitBreakerTimeout = 30 * time.Second

	// DefaultCircuitBreakerCheckInterval 默认熔断器检查间隔
	DefaultCircuitBreakerCheckInterval = 5 * time.Second
)

// 配置热重载相关常量
const (
	// DefaultGracePeriod 默认优雅关闭等待时间
	DefaultGracePeriod = 30 * time.Second
)

// 监控相关常量
const (
	// DefaultPrometheusCollectInterval 默认Prometheus指标收集间隔
	DefaultPrometheusCollectInterval = 15 * time.Second

	// DefaultResourceMonitorInterval 默认资源监控间隔
	DefaultResourceMonitorInterval = 30 * time.Second

	// DefaultConfigCacheInterval 默认配置缓存刷新间隔
	DefaultConfigCacheInterval = 30 * time.Second
)

// 超时相关常量
const (
	// DefaultDBTimeout 默认数据库操作超时
	DefaultDBTimeout = 5 * time.Second

	// DefaultRedisTimeout 默认Redis操作超时
	DefaultRedisTimeout = 5 * time.Second
)

