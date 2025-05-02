package initialize

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
	"github.com/lvjiaben/go-wheel/pkg/config"
	"github.com/lvjiaben/go-wheel/pkg/container"
	"github.com/lvjiaben/go-wheel/pkg/types"
	"github.com/lvjiaben/go-wheel/routes"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	once   sync.Once
	DB     *gorm.DB
	RDB    *redis.Client
	Logger *zap.Logger
	Config *viper.Viper
)

// Initialize 统一初始化所有组件
func Initialize() {
	once.Do(func() {
		// 初始化配置
		initConfig()
		// 初始化日志
		initLogger()
		// 初始化数据库
		initMySQL()
		// 初始化Redis
		initRedis()
		// 初始化国际化
		initI18n()
		// 初始化队列
		initQueue()
		// 初始化定时任务
		initCron()
	})
}

func initConfig() {
	Config = viper.New()
	Config.SetConfigName("config")
	Config.SetConfigType("yaml")
	Config.AddConfigPath("./configs")

	if err := Config.ReadInConfig(); err != nil {
		panic("读取配置文件失败：" + err.Error())
	}
}

func initLogger() {
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{"stdout", "./logs/app.log"}

	var err error
	Logger, err = config.Build()
	if err != nil {
		panic("初始化日志失败：" + err.Error())
	}
}

func initMySQL() {
	cfg := Config
	if cfg == nil {
		panic("配置未初始化")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.GetString("mysql.user"),
		cfg.GetString("mysql.pass"),
		cfg.GetString("mysql.host"),
		cfg.GetInt("mysql.port"),
		cfg.GetString("mysql.dbname"),
		cfg.GetString("mysql.charset"),
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:  logger.Default.LogMode(logger.Info),            // 开启详细日志
		NowFunc: func() time.Time { return time.Now().Local() }, // 使用本地时间
	})
	if err != nil {
		panic("连接数据库失败: " + err.Error())
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic("获取数据库连接失败: " + err.Error())
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(cfg.GetInt("mysql.max_idle_conns"))
	sqlDB.SetMaxOpenConns(cfg.GetInt("mysql.max_open_conns"))

	// 设置连接生命周期
	if cfg.GetInt("mysql.max_lifetime") > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.GetInt("mysql.max_lifetime")) * time.Second)
	} else {
		sqlDB.SetConnMaxLifetime(time.Hour) // 默认1小时
	}

	// 设置空闲连接生命周期
	if cfg.GetInt("mysql.max_idle_time") > 0 {
		sqlDB.SetConnMaxIdleTime(time.Duration(cfg.GetInt("mysql.max_idle_time")) * time.Second)
	} else {
		sqlDB.SetConnMaxIdleTime(30 * time.Minute) // 默认30分钟
	}

	// 验证连接
	if err := sqlDB.Ping(); err != nil {
		panic("数据库连接验证失败: " + err.Error())
	}

	Logger.Info("数据库连接初始化成功",
		zap.Int("maxOpenConns", cfg.GetInt("mysql.max_open_conns")),
		zap.Int("maxIdleConns", cfg.GetInt("mysql.max_idle_conns")),
		zap.Int("maxLifetime", cfg.GetInt("mysql.max_lifetime")),
		zap.Int("maxIdleTime", cfg.GetInt("mysql.max_idle_time")))

	DB = db
}

func initRedis() {
	cfg := Config
	if cfg == nil {
		panic("配置未初始化")
	}

	if !cfg.GetBool("redis.state") {
		return
	}

	// Redis选项
	options := &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.GetString("redis.host"), cfg.GetInt("redis.port")),
		Password: cfg.GetString("redis.pass"),
		DB:       cfg.GetInt("redis.db"),
		PoolSize: cfg.GetInt("redis.pool_size"), // 连接池大小
	}

	// 设置最小空闲连接数
	if cfg.GetInt("redis.min_idle_conns") > 0 {
		options.MinIdleConns = cfg.GetInt("redis.min_idle_conns")
	}

	// 设置读写超时
	if cfg.GetInt("redis.read_timeout") > 0 {
		options.ReadTimeout = time.Duration(cfg.GetInt("redis.read_timeout")) * time.Second
	}

	if cfg.GetInt("redis.write_timeout") > 0 {
		options.WriteTimeout = time.Duration(cfg.GetInt("redis.write_timeout")) * time.Second
	}

	// 创建Redis客户端
	client := redis.NewClient(options)

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 测试连接
	if err := client.Ping(ctx).Err(); err != nil {
		panic("连接Redis失败: " + err.Error())
	}

	Logger.Info("Redis连接初始化成功",
		zap.Int("poolSize", cfg.GetInt("redis.pool_size")),
		zap.Int("minIdleConns", cfg.GetInt("redis.min_idle_conns")))

	RDB = client
}

func initI18n() {
	// 创建 i18n 管理器
	i18nManager := NewI18n()

	// 加载语言文件
	langDir := filepath.Join("configs", "i18n")
	if err := i18nManager.LoadTranslations(langDir); err != nil {
		panic("加载语言文件失败: " + err.Error())
	}
}

func initQueue() {
	// 队列初始化逻辑
}

func initCron() {
	// 定时任务初始化逻辑
}

// ViperLoad 加载配置
func ViperLoad(c *container.Container) error {
	v := viper.New()
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AddConfigPath(".")
	v.AddConfigPath("./configs")

	if err := v.ReadInConfig(); err != nil {
		return fmt.Errorf("读取配置文件失败: %v", err)
	}

	var cfg config.Config
	if err := v.Unmarshal(&cfg); err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}

	c.SetConfig(&cfg)

	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		c.GetLogger().Info("配置文件发生变化",
			zap.String("name", e.Name),
			zap.String("op", e.Op.String()))

		var newCfg config.Config
		if err := v.Unmarshal(&newCfg); err != nil {
			c.GetLogger().Error("重新解析配置文件失败", zap.Error(err))
			return
		}

		c.SetConfig(&newCfg)
	})

	return nil
}

// ZapLoad 加载日志
func ZapLoad(c *container.Container) error {
	cfg := c.GetConfig()
	if cfg == nil {
		return fmt.Errorf("配置未初始化")
	}

	// 创建日志目录
	if err := os.MkdirAll(filepath.Dir(cfg.Log.Filename), 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %v", err)
	}

	// 配置日志
	config := zap.NewProductionConfig()
	config.OutputPaths = []string{cfg.Log.Filename, "stdout"}
	config.ErrorOutputPaths = []string{cfg.Log.Filename, "stderr"}
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err := config.Build()
	if err != nil {
		return fmt.Errorf("创建日志器失败: %v", err)
	}

	c.SetLogger(logger)
	return nil
}

// MysqlLoad 加载数据库
func MysqlLoad(c *container.Container) error {
	cfg := c.GetConfig()
	if cfg == nil {
		return fmt.Errorf("配置未初始化")
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=True&loc=Local",
		cfg.Mysql.User,
		cfg.Mysql.Pass,
		cfg.Mysql.Host,
		cfg.Mysql.Port,
		cfg.Mysql.Dbname,
		cfg.Mysql.Charset,
	)

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger:  logger.Default.LogMode(logger.Info),            // 开启详细日志
		NowFunc: func() time.Time { return time.Now().Local() }, // 使用本地时间
	})
	if err != nil {
		return fmt.Errorf("连接数据库失败: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("获取数据库连接失败: %v", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(cfg.Mysql.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Mysql.MaxOpenConns)

	// 设置连接生命周期
	if cfg.Mysql.MaxLifetime > 0 {
		sqlDB.SetConnMaxLifetime(time.Duration(cfg.Mysql.MaxLifetime) * time.Second)
	} else {
		sqlDB.SetConnMaxLifetime(time.Hour) // 默认1小时
	}

	// 设置空闲连接生命周期
	if cfg.Mysql.MaxIdleTime > 0 {
		sqlDB.SetConnMaxIdleTime(time.Duration(cfg.Mysql.MaxIdleTime) * time.Second)
	} else {
		sqlDB.SetConnMaxIdleTime(30 * time.Minute) // 默认30分钟
	}

	// 验证连接
	if err := sqlDB.Ping(); err != nil {
		return fmt.Errorf("数据库连接验证失败: %v", err)
	}

	c.GetLogger().Info("数据库连接初始化成功",
		zap.Int("maxOpenConns", cfg.Mysql.MaxOpenConns),
		zap.Int("maxIdleConns", cfg.Mysql.MaxIdleConns),
		zap.Int("maxLifetime", cfg.Mysql.MaxLifetime),
		zap.Int("maxIdleTime", cfg.Mysql.MaxIdleTime))

	c.SetDB(db)
	return nil
}

// RedisLoad 加载Redis
func RedisLoad(c *container.Container) error {
	cfg := c.GetConfig()
	if cfg == nil {
		return fmt.Errorf("配置未初始化")
	}

	if !cfg.Redis.State {
		return nil
	}

	// Redis选项
	options := &redis.Options{
		Addr:         fmt.Sprintf("%s:%d", cfg.Redis.Host, cfg.Redis.Port),
		Password:     cfg.Redis.Pass,
		DB:           cfg.Redis.Db,
		PoolSize:     cfg.Redis.PoolSize,     // 连接池大小
		MinIdleConns: cfg.Redis.MinIdleConns, // 最小空闲连接数
	}

	// 设置读写超时
	if cfg.Redis.ReadTimeout > 0 {
		options.ReadTimeout = time.Duration(cfg.Redis.ReadTimeout) * time.Second
	}

	if cfg.Redis.WriteTimeout > 0 {
		options.WriteTimeout = time.Duration(cfg.Redis.WriteTimeout) * time.Second
	}

	// 创建Redis客户端
	client := redis.NewClient(options)

	// 创建上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 测试连接
	if err := client.Ping(ctx).Err(); err != nil {
		return fmt.Errorf("连接Redis失败: %v", err)
	}

	c.GetLogger().Info("Redis连接初始化成功",
		zap.Int("poolSize", cfg.Redis.PoolSize),
		zap.Int("minIdleConns", cfg.Redis.MinIdleConns))

	c.SetRedis(&types.RedisWrapper{Client: client})
	return nil
}

// ValidateLoad 加载验证器
func ValidateLoad(c *container.Container) error {
	validate := validator.New()
	zh := zh.New()
	uni := ut.New(zh, zh)
	trans, _ := uni.GetTranslator("zh")
	if err := zh_translations.RegisterDefaultTranslations(validate, trans); err != nil {
		return fmt.Errorf("注册验证器翻译失败: %v", err)
	}

	c.SetValidator(validate)
	c.SetTranslator(trans)
	return nil
}

// NewI18n 创建国际化实例
func NewI18n() types.I18n {
	return types.NewI18nManager()
}

// NewMessageQueue 创建消息队列实例
type MessageQueueManager struct{}

func NewMessageQueue() types.MessageQueue {
	return &MessageQueueManager{}
}

func (m *MessageQueueManager) Push(ctx context.Context, topic string, message interface{}) error {
	return nil
}

func (m *MessageQueueManager) Pop(ctx context.Context, topic string) (interface{}, error) {
	return nil, nil
}

func (m *MessageQueueManager) Close() error {
	return nil
}

// NewDelayQueue 创建延迟队列实例
type DelayQueueManager struct{}

func NewDelayQueue() types.DelayQueue {
	return &DelayQueueManager{}
}

func (d *DelayQueueManager) Push(ctx context.Context, topic string, message interface{}, delay time.Duration) error {
	return nil
}

func (d *DelayQueueManager) Pop(ctx context.Context, topic string) (interface{}, error) {
	return nil, nil
}

func (d *DelayQueueManager) Close() error {
	return nil
}

// NewCronManager 创建定时任务管理器实例
type CronManagerImpl struct{}

func NewCronManager() types.CronManager {
	return &CronManagerImpl{}
}

func (c *CronManagerImpl) AddJob(spec string, cmd func()) error {
	return nil
}

func (c *CronManagerImpl) Start() {}

func (c *CronManagerImpl) Stop() {}

// RoutersLoad 加载路由
func RoutersLoad(c *container.Container) *gin.Engine {
	r := gin.New()

	// 使用中间件
	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	// 注册路由
	routes.RegisterRoutes(r, c)

	return r
}
