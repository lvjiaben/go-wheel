package model

import (
	"time"

	"gorm.io/gorm"
)

// Admin 管理员表
type Admin struct {
	Id        int       `json:"id" gorm:"primaryKey"`             // 主键ID
	Username  string    `json:"username" gorm:"not null"`         // 用户名
	Password  string    `json:"password" gorm:"not null"`         // 密码
	Avatar    string    `json:"avatar" gorm:"not null"`           // 头像
	Failures  int       `json:"failures" gorm:"not null"`         // 登录失败次数
	Token     string    `json:"token" gorm:"not null"`            // Token
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"` // 创建时间
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"` // 更新时间
}

func (Admin) TableName() string {
	return "admin"
}

// BeforeCreate 创建前的钩子
func (a *Admin) BeforeCreate(tx *gorm.DB) error {
	return nil
}

// BeforeUpdate 更新前的钩子
func (a *Admin) BeforeUpdate(tx *gorm.DB) error {
	return nil
}
