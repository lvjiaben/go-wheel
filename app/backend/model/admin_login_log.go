package model

import (
	"time"
)

// AdminLoginLog 登录日志表
type AdminLoginLog struct {
	Id        int    `json:"id" gorm:"primaryKey"`             // 主键ID
	Username  string `json:"username" gorm:"not null"`         // 用户名
	Ip        string `json:"ip" gorm:"not null"`               // IP地址
	Status    int    `json:"status" gorm:"not null"`           // 状态：0=失败，1=成功
	CreatedAt int    `json:"created_at" gorm:"autoCreateTime"` // 创建时间
}

func (AdminLoginLog) TableName() string {
	return "admin_login_log"
}

// BeforeCreate 创建前的钩子
func (l *AdminLoginLog) BeforeCreate() error {
	l.CreatedAt = int(time.Now().Unix())
	return nil
}
