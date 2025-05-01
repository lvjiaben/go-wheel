package model

type Admin struct {
	Id        int    `json:"id" gorm:"column:id"`
	Pid       int    `json:"pid" gorm:"column:pid"`               // 上级管理员
	Username  string `json:"username" gorm:"column:username"`     // 管理员账号
	Password  string `json:"password" gorm:"column:password"`     // 管理员密码
	Avatar    string `json:"avatar" gorm:"column:avatar"`         // 管理员头像
	Failures  int    `json:"failures" gorm:"column:failures"`     // 登陆失败次数
	Token     string `json:"token" gorm:"column:token"`           // TOKEN
	CreatedAt int    `json:"created_at" gorm:"column:created_at"` // 创建时间
	UpdatedAt int    `json:"updated_at" gorm:"column:updated_at"` // 更新时间
}

func (e *Admin) TableName() string {
	return "admin"
}
