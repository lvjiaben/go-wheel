package model

type AdminLoginLog struct {
	Id        int    `json:"id" gorm:"column:id"`
	Username  string `json:"username" gorm:"column:username"`     // 账号
	Ip        string `json:"ip" gorm:"column:ip"`                 // IP
	Status    int    `json:"status" gorm:"column:status"`         // 状态
	CreatedAt int    `json:"created_at" gorm:"column:created_at"` // 时间
}

func (e *AdminLoginLog) TableName() string {
	return "admin_login_log"
}
