package controller

import (
	"strconv"

	"github.com/lvjiaben/go-wheel/pkg/utils/http"

	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/app/backend/service"
	"github.com/lvjiaben/go-wheel/pkg/container"
)

// HomeController 认证控制器
type HomeController struct {
	container   *container.Container
	homeService *service.HomeService
}

// NewHomeController 创建认证控制器
func NewHomeController(c *container.Container) *HomeController {
	return &HomeController{
		container:   c,
		homeService: service.NewHomeService(c),
	}
}

func (c *HomeController) Index(ctx *gin.Context) {
	timeStr := ctx.Query("time")
	time, err := strconv.Atoi(timeStr)
	if err != nil {
		http.ErrorWithI18n(ctx, "common.error", nil)
		return
	}
	if time != 7 && time != 30 {
		http.ErrorWithI18n(ctx, "common.error", nil)
		return
	}
	http.SuccessWithI18n(ctx, "common.success", c.homeService.Index(time))
}
