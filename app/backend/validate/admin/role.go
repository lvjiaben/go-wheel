package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/app/common/validator"
)

type RoleSave struct {
	Id          int    `json:"id" binding:"gte=0" label:"backend.role.id"`
	Pid         int    `json:"pid" binding:"gte=0" label:"backend.role.pid" msg:"backend.role.pid_invalid"`
	Name        string `json:"name" binding:"required,max=64" label:"backend.role.name" msg:"backend.role.name_required"`
	Description string `json:"description" binding:"max=255" label:"backend.role.description"`
	IsSuper     int    `json:"is_super" binding:"oneof=0 1" label:"backend.role.is_super" msg:"backend.role.is_super_invalid"`
	Status      int    `json:"status" binding:"oneof=0 1" label:"backend.role.status" msg:"backend.role.status_invalid"`
	Sort        int    `json:"sort" label:"backend.role.sort" msg:"backend.role.sort_invalid"`
	MenuIds     []int  `json:"menu_ids" label:"backend.role.menu_ids"` // 菜单权限
}

func ValidateRoleSave(c *gin.Context) (*RoleSave, bool) {
	return validator.ValidateStructWithConvert[RoleSave](c)
}
