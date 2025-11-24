package system

import (
	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/app/common/validator"
)

type ConfigUpdate struct {
	Id int `json:"id" binding:"required,min=1" label:"Id"`
	ConfigCreate
}

type ConfigCreate struct {
	Name     string `json:"name" binding:"required,min=1,max=255" label:"Name"`
	Key      string `json:"key" binding:"required,min=1,max=255" label:"Key"`
	Dir      string `json:"dir" binding:"required,min=1,max=255" label:"Dir"`
	Tip      string `json:"tip" binding:"max=255" label:"Tip"`
	Type     string `json:"type" binding:"required,min=1,max=255" label:"Type"`
	Value    string `json:"value" label:"Value"`
	Variable string `json:"variable" label:"Variable"`
}

func ValidateConfigUpdate(c *gin.Context) (*ConfigUpdate, bool) {
	return validator.ValidateStruct[ConfigUpdate](c)
}

func ValidateConfigCreate(c *gin.Context) (*ConfigCreate, bool) {
	return validator.ValidateStruct[ConfigCreate](c)
}
