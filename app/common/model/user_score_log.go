package model

import (
	"gorm.io/gorm"
)

// UserScoreLog 用户积分日志表
type UserScoreLog struct {
	Id        int     `json:"id" gorm:"primaryKey;autoIncrement"` // 主键ID
	UserId    int     `json:"user_id" gorm:"not null;index"`      // 用户ID
	Type      int     `json:"type" gorm:"default:1"`              // 类型：0=减少，1=增加
	Score     float64 `json:"score" gorm:"type:decimal(10,2)"`    // 积分数量
	Note      string  `json:"note" gorm:"size:255"`               // 备注
	Source    string  `json:"source" gorm:"size:255;index"`       // 来源
	CreatedAt int     `json:"created_at" gorm:"autoCreateTime"`   // 创建时间
}

// TableName 表名
func (UserScoreLog) TableName() string {
	return "user_score_log"
}

// BeforeCreate 创建前的钩子
func (u *UserScoreLog) BeforeCreate(tx *gorm.DB) error {
	return nil
}

// BeforeUpdate 更新前的钩子
func (u *UserScoreLog) BeforeUpdate(tx *gorm.DB) error {
	return nil
}

// UserScoreLogInfo 积分日志信息
type UserScoreLogInfo struct {
	Id        int     `json:"id"`         // 主键ID
	UserId    int     `json:"user_id"`    // 用户ID
	Type      int     `json:"type"`       // 类型
	Score     float64 `json:"score"`      // 积分数量
	Note      string  `json:"note"`       // 备注
	Source    string  `json:"source"`     // 来源
	CreatedAt int     `json:"created_at"` // 创建时间
	Username  string  `json:"username"`   // 用户名（关联查询）
}

// UserScoreLogWithUser 积分日志（包含用户信息）
type UserScoreLogWithUser struct {
	UserScoreLog
	Username string `json:"username"` // 用户名
	Mobile   string `json:"mobile"`   // 手机号
}
