package model

type App struct {
	Name    string `mapstructure:"name"`
	Mode    string `mapstructure:"mode"`
	Port    int    `mapstructure:"port"`
	Version string `mapstructure:"version"`
}
type Log struct {
	Level      string `mapstructure:"level"`
	Filename   string `mapstructure:"filename"`
	MaxSize    int    `mapstructure:"max_size"`
	MaxAge     int    `mapstructure:"max_age"`
	MaxBackups int    `mapstructure:"max_backups"`
}
type Redis struct {
	State    bool   `mapstructure:"state"`
	Host     string `mapstructure:"host"`
	Pass     string `mapstructure:"pass"`
	Port     int    `mapstructure:"port"`
	Db       int    `mapstructure:"db"`
	PoolSize int    `mapstructure:"pool_size"`
}
type Mysql struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	User         string `mapstructure:"user"`
	Pass         string `mapstructure:"pass"`
	Dbname       string `mapstructure:"dbname"`
	Charset      string `mapstructure:"charset"`
	MaxOpenConns int    `mapstructure:"max_open_conns"`
	MaxIdleConns int    `mapstructure:"max_idle_conns"`
}
type Jwt struct {
	Secret    string `mapstructure:"secret"`
	ExpireDay int    `mapstructure:"expire_day"`
}
type Admin struct {
	LoginFailures       int  `mapstructure:"login_failures"`
	LoginFailuresSecond int  `mapstructure:"login_failures_second"`
	LoginSso            bool `mapstructure:"login_sso"`
}

type Config struct {
	App   App   `mapstructure:"app"`
	Log   Log   `mapstructure:"log"`
	Redis Redis `mapstructure:"redis"`
	Mysql Mysql `mapstructure:"mysql"`
	Jwt   Jwt   `mapstructure:"jwt"`
	Admin Admin `mapstructure:"admin"`
}
