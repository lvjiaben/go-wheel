package model

// AdminAuthGroupAccess 角色访问表
type AdminAuthGroupAccess struct {
	Uid int `json:"uid" gorm:"primaryKey"` // 用户ID
	Gid int `json:"gid" gorm:"primaryKey"` // 角色ID
}

func (AdminAuthGroupAccess) TableName() string {
	return "admin_auth_group_access"
}
