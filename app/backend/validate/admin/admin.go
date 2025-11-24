package admin

import (
	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/app/common/validator"
)

type AdminSave struct {
	Id       int    `json:"id" binding:"gte=0" label:"backend.admin.id"`
	Pid      int    `json:"pid" binding:"gte=0" label:"backend.admin.pid" msg:"backend.admin.pid_invalid"`
	Username string `json:"username" binding:"required,max=50" label:"backend.admin.username" msg:"backend.admin.username_required"`
	Password string `json:"password" label:"backend.admin.password" msg:"backend.admin.password_required"`
	Realname string `json:"realname" binding:"max=50" label:"backend.admin.realname"`
	Avatar   string `json:"avatar" binding:"max=255" label:"backend.admin.avatar"`
	Email    string `json:"email" binding:"omitempty,email,max=100" label:"backend.admin.email" msg:"backend.admin.email_invalid"`
	Mobile   string `json:"mobile" binding:"max=15" label:"backend.admin.mobile"`
	Status   int    `json:"status" binding:"oneof=0 1" label:"backend.admin.status" msg:"backend.admin.status_invalid"`
	RoleIds  []int  `json:"role_ids" label:"backend.admin.role_ids"`
}

func ValidateAdminSave(c *gin.Context) (*AdminSave, bool) {
	return validator.ValidateStructWithConvert[AdminSave](c)
}
