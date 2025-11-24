package system

import (
	"fmt"

	"github.com/lvjiaben/go-wheel/app/backend/model/system"
	"github.com/lvjiaben/go-wheel/pkg/container"
)

// ConfigService 附件服务
type ConfigService struct {
	container *container.Container
}

// NewConfigService 创建附件服务
func NewConfigService(c *container.Container) *ConfigService {
	return &ConfigService{
		container: c,
	}
}

func (s *ConfigService) List() map[string]interface{} {
	// 第一步：查询所有不同的 dir 分组
	type DirGroup struct {
		Id  uint   `json:"id"`
		Dir string `json:"dir"`
	}
	var dirGroups []DirGroup
	s.container.GetDB().Model(&system.Config{}).
		Select("MIN(id) as id, dir").
		Group("dir").
		Order("MIN(id) ASC").
		Find(&dirGroups)

	// 第二步：遍历每个 dir，查询该分组下的所有配置项
	result := make([]map[string]interface{}, 0)
	for _, dirGroup := range dirGroups {
		// 查询该 dir 下的所有配置项
		var configs []system.Config
		s.container.GetDB().Where("dir = ?", dirGroup.Dir).
			Order("id ASC").
			Find(&configs)

		// 构建返回结构
		item := map[string]interface{}{
			"id":       dirGroup.Id,
			"dir":      dirGroup.Dir,
			"children": configs,
		}
		result = append(result, item)
	}

	return map[string]interface{}{
		"list": result,
	}
}

func (s *ConfigService) Delete(id int) {
	var config system.Config
	// 数据库直接删除
	s.container.GetDB().Where("id = ?", id).Delete(&config)
}

func (s *ConfigService) Create(form *system.Config) error {
	if err := s.container.GetDB().Create(form).Error; err != nil {
		return fmt.Errorf("common.fail")
	}
	return nil
}

// UpdateByKey 根据key更新配置值
func (s *ConfigService) UpdateByKey(key string, value string) error {
	// 查询配置是否存在
	var config system.Config
	if err := s.container.GetDB().Where("`key` = ?", key).First(&config).Error; err != nil {
		return err
	}

	// 更新配置值
	if err := s.container.GetDB().Model(&config).Update("value", value).Error; err != nil {
		return err
	}

	return nil
}
