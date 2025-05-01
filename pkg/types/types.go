package types

import (
	"context"
	"time"
)

// Config 配置结构体
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
		MaxIdleConns int
		MaxOpenConns int
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
		Secret string
		Expire int
		Issuer string
	}
	Admin struct {
		Username string
		Password string
	}
	Nacos struct {
		Host      string
		Port      int
		Namespace string
		DataId    string
		Group     string
	}
}

// RedisClient Redis客户端接口
type RedisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
	Del(ctx context.Context, keys ...string) error
	LPush(ctx context.Context, key string, values ...interface{}) error
	BRPop(ctx context.Context, timeout time.Duration, keys ...string) ([]string, error)
	ZAdd(ctx context.Context, key string, members ...interface{}) error
	ZRangeByScore(ctx context.Context, key string, min, max string) ([]string, error)
	ZRem(ctx context.Context, key string, members ...interface{}) error
	Close() error
}

// Message 消息结构体
type Message struct {
	ID      string
	Topic   string
	Content string
	Time    time.Time
}

// DelayTask 延迟任务结构体
type DelayTask struct {
	ID        string
	Topic     string
	Content   string
	ExecuteAt time.Time
}

// Task 定时任务结构体
type Task struct {
	ID       string
	Name     string
	Spec     string
	Handler  func()
	NextTime time.Time
}

// MessageQueue 消息队列接口
type MessageQueue interface {
	Publish(ctx context.Context, msg *Message) error
	Subscribe(ctx context.Context, topic string, handler func(*Message)) error
	Close() error
}

// DelayQueue 延迟队列接口
type DelayQueue interface {
	AddTask(ctx context.Context, task *DelayTask) error
	Process(ctx context.Context, handler func(*DelayTask)) error
	Close() error
}

// CronManager 定时任务管理器接口
type CronManager interface {
	AddTask(task *Task) error
	RemoveTask(id string) error
	Start()
	Stop()
}

// I18n 国际化接口
type I18n interface {
	LoadTranslations(dir string) error
	Translate(lang, key string, args ...interface{}) string
	GetAvailableLanguages() []string
}
