package system

import (
	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/app/backend/gen"
	"github.com/lvjiaben/go-wheel/app/common/validator"
)

// GenTableList 获取表列表
type GenTableList struct {
	Search string `form:"search" label:"搜索关键词"`
}

// GenTableInfo 获取表信息
type GenTableInfo struct {
	TableName string `form:"table_name" binding:"required" label:"表名"`
}

// GenPreview 预览代码
type GenPreview struct {
	Config gen.GenConfig `json:"config" binding:"required" label:"生成配置"`
}

// GenGenerate 生成代码
type GenGenerate struct {
	Config gen.GenConfig `json:"config" binding:"required" label:"生成配置"`
}

// GenDelete 删除生成的代码
type GenDelete struct {
	Id int `json:"id" binding:"required,min=1" label:"历史记录ID"`
}

// ValidateGenTableList 验证获取表列表
func ValidateGenTableList(c *gin.Context) (*GenTableList, bool) {
	return validator.ValidateStruct[GenTableList](c)
}

// ValidateGenTableInfo 验证获取表信息
func ValidateGenTableInfo(c *gin.Context) (*GenTableInfo, bool) {
	return validator.ValidateStruct[GenTableInfo](c)
}

// ValidateGenPreview 验证预览代码
func ValidateGenPreview(c *gin.Context) (*GenPreview, bool) {
	return validator.ValidateStruct[GenPreview](c)
}

// ValidateGenGenerate 验证生成代码
func ValidateGenGenerate(c *gin.Context) (*GenGenerate, bool) {
	return validator.ValidateStruct[GenGenerate](c)
}

// ValidateGenDelete 验证删除
func ValidateGenDelete(c *gin.Context) (*GenDelete, bool) {
	return validator.ValidateStruct[GenDelete](c)
}

