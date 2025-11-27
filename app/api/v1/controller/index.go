package controller

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	commonService "github.com/lvjiaben/go-wheel/app/common/service"
	"github.com/lvjiaben/go-wheel/pkg/container"
)

// IndexController 首页控制器
type IndexController struct {
	container   *container.Container
	configCache *commonService.ConfigCacheService
}

// NewIndexController 创建首页控制器
func NewIndexController(c *container.Container) *IndexController {
	return &IndexController{
		container:   c,
		configCache: commonService.NewConfigCacheService(c),
	}
}

// Index 首页方法
func (c *IndexController) Index(ctx *gin.Context) {

	// 准备模板数据
	data := gin.H{
		"title":       "欢迎页面",
		"message":     "Hello World ," + c.configCache.Get("site_name"),
		"currentTime": time.Now().Format("2006-01-02 15:04:05"),
	}

	// 渲染模板
	ctx.HTML(http.StatusOK, "index.html", data)
}
