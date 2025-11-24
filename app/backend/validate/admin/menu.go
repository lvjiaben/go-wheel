package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/app/common/validator"
)

type MenuSave struct {
	Id         int    `json:"id" binding:"gte=0"`
	Sort       int    `json:"sort" binding:"gte=0" label:"backend.menu.sort" msg:"backend.menu.sort_invalid"`
	Pid        int    `json:"pid" binding:"gte=0" label:"backend.menu.pid" msg:"backend.menu.pid_invalid"`
	Name       string `json:"name" binding:"required,max=64" label:"backend.menu.name" msg:"backend.menu.name_required"`
	Enname     string `json:"enname" binding:"required,max=64" label:"backend.menu.enname" msg:"backend.menu.enname_required"`
	Route      string `json:"route" binding:"max=128" label:"backend.menu.route" msg:"backend.menu.route_required"`
	Component  string `json:"component" binding:"max=128" label:"backend.menu.component" msg:"backend.menu.component_required"`
	Path       string `json:"path" binding:"max=128" label:"backend.menu.path" msg:"backend.menu.path_required"`
	Icon       string `json:"icon" binding:"max=64" label:"backend.menu.icon"`
	Visible    int    `json:"visible" binding:"oneof=0 1" label:"backend.menu.visible" msg:"backend.menu.visible_invalid"`
	FixedTag   int    `json:"fixed_tag" binding:"oneof=0 1" label:"backend.menu.fixed_tag" msg:"backend.menu.fixed_tag_invalid"`
	ShowTag    int    `json:"show_tag" binding:"oneof=0 1" label:"backend.menu.show_tag" msg:"backend.menu.show_tag_invalid"`
	Iframe     string `json:"iframe" binding:"max=128" label:"backend.menu.iframe" msg:"backend.menu.iframe_invalid"`
	External   string `json:"external" binding:"max=128" label:"backend.menu.external" msg:"backend.menu.external_invalid"`
	Type       string `json:"type" binding:"required,oneof=menu button iframe link" label:"backend.menu.type" msg:"backend.menu.type_required"`
	Permission string `json:"permission" binding:"max=128" label:"backend.menu.permission"`
}

type MenuDelete struct {
	Id int `json:"id" binding:"required" label:"backend.menu.id" msg:"backend.menu.id_required"`
}

func ValidateMenuSave(c *gin.Context) (*MenuSave, bool) {
	return validator.ValidateStructWithConvert[MenuSave](c)
}

func ValidateMenuDelete(c *gin.Context) (*MenuDelete, bool) {
	return validator.ValidateStruct[MenuDelete](c)
}
