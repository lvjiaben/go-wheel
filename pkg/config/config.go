package config

// Config 配置结构体
type Config struct {
	App      App      `mapstructure:"app"`
	Database Database `mapstructure:"database"`
	Redis    Redis    `mapstructure:"redis"`
	Log      Log      `mapstructure:"log"`
	Jwt      Jwt      `mapstructure:"jwt"`
	Admin    Admin    `mapstructure:"admin"`
	RabbitMQ RabbitMQ `mapstructure:"rabbitmq"`
	Upload   Upload   `mapstructure:"upload"`
}

type Upload struct {
	Type              string      `mapstructure:"type"`
	BaseUrl           string      `mapstructure:"base_url"`
	UploadPath        string      `mapstructure:"upload_path"`
	MaxSize           int64       `mapstructure:"max_size"`
	AllowedTypes      []string    `mapstructure:"allowed_types"`
	AllowedExtensions []string    `mapstructure:"allowed_extensions"`
	Cos               CosConfig   `mapstructure:"cos"`
	Oss               OssConfig   `mapstructure:"oss"`
	Qiniu             QiniuConfig `mapstructure:"qiniu"`
}

type CosConfig struct {
	SecretId  string `mapstructure:"secret_id"`
	SecretKey string `mapstructure:"secret_key"`
	Region    string `mapstructure:"region"`
	Bucket    string `mapstructure:"bucket"`
}

type OssConfig struct {
	Endpoint        string `mapstructure:"endpoint"`
	AccessKeyID     string `mapstructure:"access_key_id"`
	AccessKeySecret string `mapstructure:"access_key_secret"`
	BucketName      string `mapstructure:"bucket_name"`
}

type QiniuConfig struct {
	Region    string `mapstructure:"region"`
	AccessKey string `mapstructure:"access_key"`
	SecretKey string `mapstructure:"secret_key"`
	Bucket    string `mapstructure:"bucket"`
}

// App 应用配置
type App struct {
	Name           string `mapstructure:"name"`
	Port           int    `mapstructure:"port"`
	Mode           string `mapstructure:"mode"`
	Version        string `mapstructure:"version"`
	MaxRequestBody int64  `mapstructure:"max_request_body"` // 最大请求体大小（MB）
}

// Database 数据库配置（支持 MySQL 和 PostgreSQL）
type Database struct {
	Driver       string `mapstructure:"driver"`        // 数据库驱动：mysql 或 postgres
	Host         string `mapstructure:"host"`          // 主机地址
	Port         int    `mapstructure:"port"`          // 端口号
	User         string `mapstructure:"user"`          // 用户名
	Pass         string `mapstructure:"pass"`          // 密码
	Dbname       string `mapstructure:"dbname"`        // 数据库名
	Charset      string `mapstructure:"charset"`       // 字符集（仅 MySQL）
	SSLMode      string `mapstructure:"sslmode"`       // SSL 模式（仅 PostgreSQL）
	Timezone     string `mapstructure:"timezone"`      // 时区
	MaxIdleConns int    `mapstructure:"max_idle_conns"` // 最大空闲连接数
	MaxOpenConns int    `mapstructure:"max_open_conns"` // 最大打开连接数
	MaxLifetime  int    `mapstructure:"max_lifetime"`   // 连接最大生命周期（秒）
	MaxIdleTime  int    `mapstructure:"max_idle_time"`  // 空闲连接最大生命周期（秒）
}

// Redis 配置
type Redis struct {
	State         bool   `mapstructure:"state"`
	Host          string `mapstructure:"host"`
	Port          int    `mapstructure:"port"`
	Pass          string `mapstructure:"pass"`
	Db            int    `mapstructure:"db"`
	PoolSize      int    `mapstructure:"pool_size"`
	MinIdleConns  int    `mapstructure:"min_idle_conns"`
	MaxConnAge    int    `mapstructure:"max_conn_age"`
	IdleTimeout   int    `mapstructure:"idle_timeout"`
	PoolTimeout   int    `mapstructure:"pool_timeout"`
	ReadTimeout   int    `mapstructure:"read_timeout"`
	WriteTimeout  int    `mapstructure:"write_timeout"`
	EnableMetrics bool   `mapstructure:"enable_metrics"`
}

// RabbitMQ 配置
type RabbitMQ struct {
	State               bool   `mapstructure:"state"`
	Host                string `mapstructure:"host"`
	Port                int    `mapstructure:"port"`
	VirtualHost         string `mapstructure:"virtual_host"`
	User                string `mapstructure:"user"`
	Pass                string `mapstructure:"pass"`
	QueueName           string `mapstructure:"queue_name"`
	DelayQueueName      string `mapstructure:"delay_queue_name"`
	Exchange            string `mapstructure:"exchange"`
	DelayExchange       string `mapstructure:"delay_exchange"`
	RetryCount          int    `mapstructure:"retry_count"`
	ReconnectInterval   int    `mapstructure:"reconnect_interval"`
	HeartbeatInterval   int    `mapstructure:"heartbeat_interval"`
	ConnectionTimeout   int    `mapstructure:"connection_timeout"`
	EnableConfirmation  bool   `mapstructure:"enable_confirmation"`
	PrefetchCount       int    `mapstructure:"prefetch_count"`
	PrefetchSize        int    `mapstructure:"prefetch_size"`
	DefaultConsumerName string `mapstructure:"default_consumer_name"`
}

// Log 日志配置
type Log struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxBackups int    `mapstructure:"max_backups"`
	MaxAge     int    `mapstructure:"max_age"`
}

// Jwt 配置
type Jwt struct {
	Secret string `mapstructure:"secret"`
	Expire int    `mapstructure:"expire_day"`
	Issuer string `mapstructure:"issuer"`
}

// Admin 管理员配置
type Admin struct {
	LoginFailures int    `mapstructure:"login_failures"`
	LoginSso      bool   `mapstructure:"login_sso"`
	Username      string `mapstructure:"username"`
	Password      string `mapstructure:"password"`
}

// GetString 获取字符串配置
func (c *Config) GetString(key string) string {
	switch key {
	case "jwt.secret":
		return c.Jwt.Secret
	default:
		return ""
	}
}

// GetInt 获取整数配置
func (c *Config) GetInt(key string) int {
	switch key {
	case "admin.login_failures":
		return c.Admin.LoginFailures
	case "jwt.expire_day":
		return c.Jwt.Expire
	default:
		return 0
	}
}

// GetBool 获取布尔配置
func (c *Config) GetBool(key string) bool {
	switch key {
	case "admin.login_sso":
		return c.Admin.LoginSso
	default:
		return false
	}
}
