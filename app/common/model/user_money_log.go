package model

import (
	"gorm.io/gorm"
)

// UserMoneyLog 用户资金日志表
type UserMoneyLog struct {
	Id        int     `json:"id" gorm:"primaryKey;autoIncrement"` // 主键ID
	UserId    int     `json:"user_id" gorm:"not null;index"`      // 用户ID
	Type      int     `json:"type" gorm:"default:1"`              // 类型：0=减少，1=增加
	Money     float64 `json:"money" gorm:"type:decimal(10,2)"`    // 金额
	Note      string  `json:"note" gorm:"size:255"`               // 备注
	Source    string  `json:"source" gorm:"size:255;index"`       // 来源
	CreatedAt int     `json:"created_at" gorm:"autoCreateTime"`   // 创建时间
}

// TableName 表名
func (UserMoneyLog) TableName() string {
	return "user_money_log"
}

// BeforeCreate 创建前的钩子
func (u *UserMoneyLog) BeforeCreate(tx *gorm.DB) error {
	return nil
}

// BeforeUpdate 更新前的钩子
func (u *UserMoneyLog) BeforeUpdate(tx *gorm.DB) error {
	return nil
}

// UserMoneyLogInfo 资金日志信息
type UserMoneyLogInfo struct {
	Id        int     `json:"id"`         // 主键ID
	UserId    int     `json:"user_id"`    // 用户ID
	Type      int     `json:"type"`       // 类型
	Money     float64 `json:"money"`      // 金额
	Note      string  `json:"note"`       // 备注
	Source    string  `json:"source"`     // 来源
	CreatedAt int     `json:"created_at"` // 创建时间
	Username  string  `json:"username"`   // 用户名（关联查询）
}

// UserMoneyLogWithUser 资金日志（包含用户信息）
type UserMoneyLogWithUser struct {
	UserMoneyLog
	Username string `json:"username"` // 用户名
	Mobile   string `json:"mobile"`   // 手机号
}
