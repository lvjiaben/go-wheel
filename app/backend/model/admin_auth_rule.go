package model

import (
	"time"
)

// AdminAuthRule 权限规则表
type AdminAuthRule struct {
	Id        int    `json:"id" gorm:"primaryKey"`             // 主键ID
	Sort      int    `json:"sort" gorm:"not null"`             // 排序
	Status    int    `json:"status" gorm:"not null"`           // 状态：0=禁用，1=启用
	Pid       int    `json:"pid" gorm:"not null"`              // 上级权限ID
	Name      string `json:"name" gorm:"not null"`             // 权限名称
	Alias     string `json:"alias" gorm:"not null"`            // 权限别名
	Icon      string `json:"icon" gorm:"not null"`             // 图标
	Type      string `json:"type" gorm:"not null"`             // 类型：menu=菜单，button=按钮
	Path      string `json:"path" gorm:"not null"`             // 路径
	Method    string `json:"method" gorm:"not null"`           // 请求方法
	Enname    string `json:"enname" gorm:"not null"`           // 英文名称
	Hide      int    `json:"hide" gorm:"not null"`             // 是否隐藏：0=显示，1=隐藏
	CreatedAt int    `json:"created_at" gorm:"autoCreateTime"` // 创建时间
	UpdatedAt int    `json:"updated_at" gorm:"autoUpdateTime"` // 更新时间
}

func (AdminAuthRule) TableName() string {
	return "admin_auth_rule"
}

// BeforeCreate 创建前的钩子
func (r *AdminAuthRule) BeforeCreate() error {
	r.CreatedAt = int(time.Now().Unix())
	r.UpdatedAt = int(time.Now().Unix())
	return nil
}

// BeforeUpdate 更新前的钩子
func (r *AdminAuthRule) BeforeUpdate() error {
	r.UpdatedAt = int(time.Now().Unix())
	return nil
}
