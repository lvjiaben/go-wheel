package global

type Config struct {
	App struct {
		Name    string
		Mode    string
		Port    int
		Version string
	}
	Log struct {
		Level      string
		Filename   string
		MaxSize    int
		MaxBackups int
		MaxAge     int
	}
	Mysql struct {
		Host         string
		Port         int
		User         string
		Pass         string
		Dbname       string
		Charset      string
		MaxOpenConns int
		MaxIdleConns int
	}
	Redis struct {
		State    bool
		Host     string
		Port     int
		Pass     string
		Db       int
		PoolSize int
	}
	Jwt struct {
		Secret    string
		ExpireDay int
	}
	Admin struct {
		LoginFailures       int
		LoginFailuresSecond int
		LoginSso            bool
	}
	Nacos struct {
		Host      string
		Port      int
		Namespace string
		DataId    string
		Group     string
	}
}
