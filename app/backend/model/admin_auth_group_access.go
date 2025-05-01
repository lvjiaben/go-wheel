package model

type AdminAuthGroupAccess struct {
	Uid int `json:"uid" gorm:"column:uid"`
	Gid int `json:"gid" gorm:"column:gid"`
}

func (e *AdminAuthGroupAccess) TableName() string {
	return "admin_auth_group_access"
}
