package system

import (
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/app/backend/model/system"
	serviceSystem "github.com/lvjiaben/go-wheel/app/backend/service/system"
	validateSystem "github.com/lvjiaben/go-wheel/app/backend/validate/system"
	commonService "github.com/lvjiaben/go-wheel/app/common/service"
	"github.com/lvjiaben/go-wheel/pkg/container"
	"github.com/lvjiaben/go-wheel/pkg/utils/http"
)

// ConfigController 配置控制器
type ConfigController struct {
	container     *container.Container
	configService *serviceSystem.ConfigService
	configCache   *commonService.ConfigCacheService
}

// NewConfigController 创建配置控制器
func NewConfigController(c *container.Container) *ConfigController {
	return &ConfigController{
		container:     c,
		configService: serviceSystem.NewConfigService(c),
		configCache:   commonService.NewConfigCacheService(c),
	}
}

// List 获取配置列表
func (c *ConfigController) List(ctx *gin.Context) {
	http.SuccessWithI18n(ctx, "common.success", c.configService.List())
}

// Create 创建
func (m *ConfigController) Create(c *gin.Context) {
	// 使用ValidateConfigCreate进行验证
	form, valid := validateSystem.ValidateConfigCreate(c)
	if !valid {
		return
	}
	// 创建菜单对象
	data := system.Config{
		Dir:      form.Dir,
		Key:      form.Key,
		Name:     form.Name,
		Tip:      form.Tip,
		Type:     form.Type,
		Value:    form.Value,
		Variable: form.Variable,
	}
	if err := m.configService.Create(&data); err != nil {
		http.ErrorWithI18n(c, err.Error(), nil)
		return
	}

	// 刷新配置缓存
	m.configCache.Refresh()

	http.SuccessWithI18n(c, "common.success", nil)
}

// Update 更新配置（批量更新）
func (m *ConfigController) Update(c *gin.Context) {
	// 获取POST的所有参数
	var postData map[string]interface{}
	if err := c.ShouldBindJSON(&postData); err != nil {
		http.ErrorWithI18n(c, "common.invalid_params", nil)
		return
	}
	// 遍历所有参数，key为配置的key，value为配置的值
	updateCount := 0
	failedKeys := []string{}

	for key, value := range postData {
		// 将value转换为字符串
		var valueStr string
		switch v := value.(type) {
		case string:
			valueStr = v
		case float64:
			valueStr = strconv.FormatFloat(v, 'f', -1, 64)
		case int:
			valueStr = strconv.Itoa(v)
		case bool:
			valueStr = strconv.FormatBool(v)
		default:
			valueStr = ""
		}
		// 调用service更新配置
		if err := m.configService.UpdateByKey(key, valueStr); err != nil {
			failedKeys = append(failedKeys, key)
			continue
		}
		updateCount++
	}

	// 刷新配置缓存
	m.configCache.Refresh()

	http.SuccessWithI18n(c, "common.success", gin.H{
		"updated_count": updateCount,
		"failed_count":  len(failedKeys),
		"total_count":   len(postData),
		"failed_keys":   failedKeys,
	})
}

// Delete 删除配置
func (c *ConfigController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.ErrorWithI18n(ctx, "common.invalid_params", nil)
		return
	}
	c.configService.Delete(id)

	// 刷新配置缓存
	c.configCache.Refresh()

	http.SuccessWithI18n(ctx, "common.success", nil)
}
