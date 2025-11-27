package validate

import (
	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/app/common/validator"
)

// UserCreate 创建用户
type UserCreate struct {
	Pid        int    `json:"pid" binding:"min=0" label:"上级ID"`
	Tid        int    `json:"tid" binding:"min=0" label:"顶级ID"`
	Status     int8   `json:"status" binding:"oneof=0 1" label:"状态"`
	StatusText string `json:"status_text" binding:"max=255" label:"状态信息"`
	Code       string `json:"code" binding:"max=255" label:"邀请码"`
	Avatar     string `json:"avatar" binding:"max=255" label:"头像"`
	Username   string `json:"username" binding:"required,min=1,max=32" label:"账号"`
	Password   string `json:"password" binding:"max=72" label:"密码"`
	Salt       string `json:"salt" binding:"max=32" label:"盐值"`
	Email      string `json:"email" binding:"omitempty,email,max=255" label:"邮箱"`
	Mobile     string `json:"mobile" binding:"omitempty,len=11" label:"手机号码"`
}

// UserUpdate 更新用户
type UserUpdate struct {
	Id         int    `json:"id" binding:"required,min=1" label:"用户ID"`
	Pid        int    `json:"pid" binding:"min=0" label:"上级ID"`
	Tid        int    `json:"tid" binding:"min=0" label:"顶级ID"`
	Status     int8   `json:"status" binding:"oneof=0 1" label:"状态"`
	StatusText string `json:"status_text" binding:"max=255" label:"状态信息"`
	Code       string `json:"code" binding:"max=255" label:"邀请码"`
	Avatar     string `json:"avatar" binding:"max=255" label:"头像"`
	Username   string `json:"username" binding:"required,min=1,max=32" label:"账号"`
	Password   string `json:"password" binding:"max=72" label:"密码"`
	Salt       string `json:"salt" binding:"max=32" label:"盐值"`
	Email      string `json:"email" binding:"omitempty,email,max=255" label:"邮箱"`
	Mobile     string `json:"mobile" binding:"omitempty,len=11" label:"手机号码"`
}

// UserUpdateMoney 更新用户余额
type UserUpdateMoney struct {
	Id     int     `json:"id" binding:"required,min=1" label:"用户ID"`
	Type   string  `json:"type" binding:"required,oneof=add sub" label:"操作类型"`
	Money  float64 `json:"money" binding:"required,gt=0" label:"金额"`
	Note   string  `json:"note" binding:"max=255" label:"备注"`
	Source string  `json:"source" binding:"max=255" label:"来源"`
}

// UserUpdateScore 更新用户积分
type UserUpdateScore struct {
	Id     int     `json:"id" binding:"required,min=1" label:"用户ID"`
	Type   string  `json:"type" binding:"required,oneof=add sub" label:"操作类型"`
	Score  float64 `json:"score" binding:"required,gt=0" label:"积分"`
	Note   string  `json:"note" binding:"max=255" label:"备注"`
	Source string  `json:"source" binding:"max=255" label:"来源"`
}

// UserOperate 操作用户字段（status等开关字段）
type UserOperate struct {
	Ids   []int  `json:"ids" binding:"required,dive,min=1" label:"用户ID列表"` // 批量ID（可选）
	Field string `json:"field" binding:"required" label:"字段名"`
	Value int    `json:"value" binding:"" label:"字段值"`
}

// UserDelete 删除
type UserDelete struct {
	Ids []int `json:"ids" binding:"required,min=1,dive,min=1" label:"用户ID列表"`
}

// ValidateUserCreate 验证创建用户
func ValidateUserCreate(c *gin.Context) (*UserCreate, bool) {
	return validator.ValidateStructWithConvert[UserCreate](c)
}

// ValidateUserUpdate 验证更新用户
func ValidateUserUpdate(c *gin.Context) (*UserUpdate, bool) {
	return validator.ValidateStructWithConvert[UserUpdate](c)
}

// ValidateUserUpdateMoney 验证更新余额
func ValidateUserUpdateMoney(c *gin.Context) (*UserUpdateMoney, bool) {
	return validator.ValidateStructWithConvert[UserUpdateMoney](c)
}

// ValidateUserUpdateScore 验证更新积分
func ValidateUserUpdateScore(c *gin.Context) (*UserUpdateScore, bool) {
	return validator.ValidateStructWithConvert[UserUpdateScore](c)
}

// ValidateUserOperate 验证操作字段
func ValidateUserOperate(c *gin.Context) (*UserOperate, bool) {
	return validator.ValidateStructWithConvert[UserOperate](c)
}

// ValidateUserDelete 验证删除
func ValidateUserDelete(c *gin.Context) (*UserDelete, bool) {
	return validator.ValidateStruct[UserDelete](c)
}
