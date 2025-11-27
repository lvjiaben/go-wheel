package types

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"gopkg.in/yaml.v3"
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
}

// RedisClient Redis客户端接口
type RedisClient struct {
	*redis.Client
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

// I18n 国际化接口
type I18n interface {
	Get(key string, lang string) string
	Set(key string, value string, lang string) error
	LoadTranslations(dir string) error
}

// I18nManager 国际化管理器
type I18nManager struct {
	Translations map[string]map[string]string
}

func NewI18nManager() *I18nManager {
	return &I18nManager{
		Translations: make(map[string]map[string]string),
	}
}

func (i *I18nManager) Get(key string, lang string) string {
	// 检查映射是否存在对应语言
	if translations, ok := i.Translations[lang]; ok {
		if value, ok := translations[key]; ok {
			return value
		}
	}

	// 尝试切分语言代码，如zh-CN -> zh
	langBase := strings.Split(lang, "-")[0]

	// 尝试使用基础语言代码查找
	for transLang, translations := range i.Translations {
		// 如果找到基础语言代码匹配的语言
		if strings.HasPrefix(transLang, langBase) {
			if value, ok := translations[key]; ok {
				return value
			}
			break
		}
	}

	// 如果找不到，回退到默认语言（中文）
	if lang != "zh-CN" && lang != "zh" {
		if translations, ok := i.Translations["zh-CN"]; ok {
			if value, ok := translations[key]; ok {
				return value
			}
		}
	}

	// 都找不到，返回原键
	return key
}

func (i *I18nManager) Set(key string, value string, lang string) error {
	if _, ok := i.Translations[lang]; !ok {
		i.Translations[lang] = make(map[string]string)
	}
	i.Translations[lang][key] = value
	return nil
}

func (i *I18nManager) LoadTranslations(dir string) error {
	// 获取目录下的所有 yaml 文件
	files, err := filepath.Glob(filepath.Join(dir, "*.yaml"))
	if err != nil {
		return fmt.Errorf("获取语言文件失败: %v", err)
	}

	for _, file := range files {
		// 获取语言代码（文件名）
		lang := filepath.Base(file)
		lang = lang[:len(lang)-5] // 去掉 .yaml 后缀

		// 读取语言文件
		data, err := os.ReadFile(file)
		if err != nil {
			return fmt.Errorf("读取语言文件 %s 失败: %v", file, err)
		}

		// 使用yaml解析器处理文件
		var nested map[string]interface{}
		if err := yaml.Unmarshal(data, &nested); err != nil {
			return fmt.Errorf("解析语言文件 %s 失败: %v", file, err)
		}

		// 将嵌套的结构扁平化
		translations := make(map[string]string)
		flattenMap("", nested, translations)
		i.Translations[lang] = translations
	}

	return nil
}

// flattenMap 递归处理嵌套的map结构，将其扁平化为点分隔的键
func flattenMap(prefix string, nested map[string]interface{}, result map[string]string) {
	for k, v := range nested {
		key := k
		if prefix != "" {
			key = prefix + "." + k
		}

		switch value := v.(type) {
		case string:
			result[key] = value
		case map[string]interface{}:
			flattenMap(key, value, result)
		case map[interface{}]interface{}:
			// 将map[interface{}]interface{}转换为map[string]interface{}
			stringMap := make(map[string]interface{})
			for mk, mv := range value {
				if mkString, ok := mk.(string); ok {
					stringMap[mkString] = mv
				}
			}
			flattenMap(key, stringMap, result)
		}
	}
}

// MessageQueue 消息队列接口
type MessageQueue interface {
	Push(ctx context.Context, topic string, message interface{}) error
	Pop(ctx context.Context, topic string) (interface{}, error)
	Close() error
}

// DelayQueue 延迟队列接口
type DelayQueue interface {
	Push(ctx context.Context, topic string, message interface{}, delay time.Duration) error
	Pop(ctx context.Context, topic string) (interface{}, error)
	Close() error
}

// CronManager 定时任务管理器接口
type CronManager interface {
	AddJob(spec string, name string, cmd func()) error
	AddJobWithContext(spec string, name string, cmd func(context.Context)) error
	RemoveJob(name string)
	Start()
	Stop()
	GetJobs() []string
	GetJobCount() int
	IsRunning() bool
}
