package config

// Config 配置结构体
type Config struct {
	App      App      `mapstructure:"app"`
	Mysql    Mysql    `mapstructure:"mysql"`
	Redis    Redis    `mapstructure:"redis"`
	Log      Log      `mapstructure:"log"`
	Jwt      Jwt      `mapstructure:"jwt"`
	Admin    Admin    `mapstructure:"admin"`
	Nacos    Nacos    `mapstructure:"nacos"`
	RabbitMQ RabbitMQ `mapstructure:"rabbitmq"`
}

// App 应用配置
type App struct {
	Name    string `mapstructure:"name"`
	Port    int    `mapstructure:"port"`
	Mode    string `mapstructure:"mode"`
	Version string `mapstructure:"version"`
}

// Mysql 数据库配置
type Mysql struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Pass         string `mapstructure:"pass"`
	Dbname       string `mapstructure:"dbname"`
	Charset      string `mapstructure:"charset"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxLifetime  int    `mapstructure:"max_lifetime"`
	MaxIdleTime  int    `mapstructure:"max_idle_time"`
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
	Expire int    `mapstructure:"expire"`
	Issuer string `mapstructure:"issuer"`
}

// Admin 管理员配置
type Admin struct {
	LoginFailures int    `mapstructure:"login_failures"`
	LoginSso      bool   `mapstructure:"login_sso"`
	Username      string `mapstructure:"username"`
	Password      string `mapstructure:"password"`
}

// Nacos 配置
type Nacos struct {
	Host      string `mapstructure:"host"`
	Port      int    `mapstructure:"port"`
	Namespace string `mapstructure:"namespace"`
	DataId    string `mapstructure:"data_id"`
	Group     string `mapstructure:"group"`
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
