package i18n

import (
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

// I18n 国际化结构体
type I18n struct {
	mu       sync.RWMutex
	messages map[string]map[string]string
}

// New 创建新的国际化实例
func New() *I18n {
	return &I18n{
		messages: make(map[string]map[string]string),
	}
}

// Load 加载语言文件
func (i *I18n) Load(dir string) error {
	files, err := filepath.Glob(filepath.Join(dir, "*.yaml"))
	if err != nil {
		return err
	}

	for _, file := range files {
		lang := filepath.Base(file)
		lang = lang[:len(lang)-5] // 去掉.yaml后缀

		data, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		var messages map[string]map[string]string
		if err := yaml.Unmarshal(data, &messages); err != nil {
			return err
		}

		i.mu.Lock()
		i.messages[lang] = make(map[string]string)
		for section, items := range messages {
			for key, value := range items {
				i.messages[lang][section+"."+key] = value
			}
		}
		i.mu.Unlock()
	}

	return nil
}

// T 获取翻译
func (i *I18n) T(lang, key string) string {
	i.mu.RLock()
	defer i.mu.RUnlock()

	if messages, ok := i.messages[lang]; ok {
		if message, ok := messages[key]; ok {
			return message
		}
	}

	// 如果找不到翻译，返回key
	return key
}

// HasLang 判断语言是否存在
func (i *I18n) HasLang(lang string) bool {
	i.mu.RLock()
	defer i.mu.RUnlock()

	_, ok := i.messages[lang]
	return ok
}

// GetLangs 获取所有支持的语言
func (i *I18n) GetLangs() []string {
	i.mu.RLock()
	defer i.mu.RUnlock()

	langs := make([]string, 0, len(i.messages))
	for lang := range i.messages {
		langs = append(langs, lang)
	}
	return langs
}
