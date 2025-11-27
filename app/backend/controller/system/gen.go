package system

import (
	"github.com/gin-gonic/gin"
	serviceSystem "github.com/lvjiaben/go-wheel/app/backend/service/system"
	validateSystem "github.com/lvjiaben/go-wheel/app/backend/validate/system"
	"github.com/lvjiaben/go-wheel/pkg/container"
	"github.com/lvjiaben/go-wheel/pkg/utils/http"
)

type GenController struct {
	container  *container.Container
	genService *serviceSystem.GenService
}

func NewGenController(c *container.Container) *GenController {
	return &GenController{
		container:  c,
		genService: serviceSystem.NewGenService(c),
	}
}

// TableList 获取数据库表列表
func (ctrl *GenController) TableList(ctx *gin.Context) {
	form, valid := validateSystem.ValidateGenTableList(ctx)
	if !valid {
		return
	}
	
	tables, err := ctrl.genService.GetTableList(form.Search)
	if err != nil {
		http.ErrorWithI18n(ctx, err.Error(), nil)
		return
	}
	
	http.SuccessWithI18n(ctx, "common.success", tables)
}

// TableInfo 获取表详细信息
func (ctrl *GenController) TableInfo(ctx *gin.Context) {
	form, valid := validateSystem.ValidateGenTableInfo(ctx)
	if !valid {
		return
	}
	
	tableInfo, err := ctrl.genService.GetTableInfo(form.TableName)
	if err != nil {
		http.ErrorWithI18n(ctx, err.Error(), nil)
		return
	}
	
	http.SuccessWithI18n(ctx, "common.success", tableInfo)
}

// TableConfig 获取表的默认配置
func (ctrl *GenController) TableConfig(ctx *gin.Context) {
	form, valid := validateSystem.ValidateGenTableInfo(ctx)
	if !valid {
		return
	}
	
	config, err := ctrl.genService.GetTableConfig(form.TableName)
	if err != nil {
		http.ErrorWithI18n(ctx, err.Error(), nil)
		return
	}
	
	http.SuccessWithI18n(ctx, "common.success", config)
}

// Preview 预览生成的代码
func (ctrl *GenController) Preview(ctx *gin.Context) {
	form, valid := validateSystem.ValidateGenPreview(ctx)
	if !valid {
		return
	}
	
	code, err := ctrl.genService.PreviewCode(&form.Config)
	if err != nil {
		http.ErrorWithI18n(ctx, err.Error(), nil)
		return
	}
	
	http.SuccessWithI18n(ctx, "common.success", code)
}

// Generate 生成代码
func (ctrl *GenController) Generate(ctx *gin.Context) {
	form, valid := validateSystem.ValidateGenGenerate(ctx)
	if !valid {
		return
	}
	
	if err := ctrl.genService.GenerateCode(&form.Config); err != nil {
		http.ErrorWithI18n(ctx, err.Error(), nil)
		return
	}
	
	http.SuccessWithI18n(ctx, "common.success", nil)
}

// History 获取生成历史
func (ctrl *GenController) History(ctx *gin.Context) {
	histories, err := ctrl.genService.GetHistory()
	if err != nil {
		http.ErrorWithI18n(ctx, err.Error(), nil)
		return
	}
	
	http.SuccessWithI18n(ctx, "common.success", histories)
}

// Delete 删除生成的代码
func (ctrl *GenController) Delete(ctx *gin.Context) {
	form, valid := validateSystem.ValidateGenDelete(ctx)
	if !valid {
		return
	}
	
	if err := ctrl.genService.DeleteGenerated(form.Id); err != nil {
		http.ErrorWithI18n(ctx, err.Error(), nil)
		return
	}
	
	http.SuccessWithI18n(ctx, "common.success", nil)
}

// Download 下载生成的代码
func (ctrl *GenController) Download(ctx *gin.Context) {
	form, valid := validateSystem.ValidateGenPreview(ctx)
	if !valid {
		return
	}
	
	zipData, err := ctrl.genService.DownloadCode(&form.Config)
	if err != nil {
		http.ErrorWithI18n(ctx, err.Error(), nil)
		return
	}
	
	// 设置响应头
	ctx.Header("Content-Type", "application/zip")
	ctx.Header("Content-Disposition", "attachment; filename=generated_code.zip")
	ctx.Data(200, "application/zip", zipData)
}

