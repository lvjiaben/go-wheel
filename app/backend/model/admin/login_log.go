package admin

// LoginLog 管理员登录日志
type LoginLog struct {
	Id        int    `json:"id" gorm:"primaryKey;autoIncrement"` // 自增ID
	Username  string `json:"username" gorm:"size:64;not null"`   // 用户名
	Ip        string `json:"ip" gorm:"size:64;not null"`         // IP地址
	Status    int    `json:"status" gorm:"default:0"`            // 状态：0=失败，1=成功
	CreatedAt int64  `json:"created_at" gorm:"autoCreateTime"`   // 创建时间
}

// TableName 表名
func (LoginLog) TableName() string {
	return "admin_login_log"
}
