package model

import (
	"gorm.io/gorm"
)

// User 用户表
type User struct {
	Id            int     `json:"id" gorm:"primaryKey;autoIncrement"`           // 主键ID
	Pid           int     `json:"pid" gorm:"default:0"`                         // 上级
	Tid           int     `json:"tid" gorm:"default:0"`                         // 顶级
	Status        int8    `json:"status" gorm:"default:1"`                      // 状态
	StatusText    string  `json:"status_text" gorm:"size:255"`                  // 状态信息
	Code          string  `json:"code" gorm:"size:255;index"`                   // 邀请码
	WechatUnionid string  `json:"wechat_unionid" gorm:"size:255"`               // UNIONID
	WechatOpenid  string  `json:"wechat_openid" gorm:"size:255"`                // OPENID
	Version       int     `json:"version" gorm:"default:1"`                     // 版本
	Avatar        string  `json:"avatar" gorm:"size:255"`                       // 头像
	Username      string  `json:"username" gorm:"size:32;not null"`             // 账号
	Password      string  `json:"password" gorm:"size:255;not null"`            // 密码（bcrypt需要60-72字节）
	Salt          string  `json:"salt" gorm:"size:32;not null"`                 // 密码盐值
	Email         string  `json:"email" gorm:"size:255"`                        // 邮箱
	Mobile        string  `json:"mobile" gorm:"size:11;uniqueIndex"`            // 手机号码
	Score         float64 `json:"score" gorm:"type:decimal(10,2);default:0.00"` // 积分
	Money         float64 `json:"money" gorm:"type:decimal(10,2);default:0.00"` // 余额
	Token         string  `json:"token" gorm:"type:text"`                       // TOKEN
	CreatedAt     int     `json:"created_at" gorm:"autoCreateTime"`             // 创建时间
	UpdatedAt     int     `json:"updated_at" gorm:"autoUpdateTime"`             // 更新时间
}

// TableName 表名
func (User) TableName() string {
	return "user"
}

// BeforeCreate 创建前的钩子
func (u *User) BeforeCreate(tx *gorm.DB) error {
	return nil
}

// BeforeUpdate 更新前的钩子
func (u *User) BeforeUpdate(tx *gorm.DB) error {
	return nil
}
