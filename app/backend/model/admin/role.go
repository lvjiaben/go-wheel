package admin

// Role 管理员角色
type Role struct {
	Id          int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Pid         int    `json:"pid" gorm:"default:0;comment:父ID"`
	Name        string `json:"name" gorm:"size:64;not null;comment:角色名称"`
	Description string `json:"description" gorm:"size:255;comment:角色描述"`
	IsSuper     int    `json:"is_super" gorm:"default:0;comment:是否超级管理员 0否 1是"`
	Status      int    `json:"status" gorm:"default:1;comment:状态 0禁用 1启用"`
	Sort        int    `json:"sort" gorm:"default:50;comment:排序"`
	CreatedAt   int    `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt   int    `json:"updated_at" gorm:"comment:更新时间"`
}

// RoleMenu 角色-菜单关联表
type RoleMenu struct {
	Id        int `gorm:"primaryKey" json:"id"` // 主键ID
	RoleId    int `json:"role_id"`              // 角色ID
	MenuId    int `json:"menu_id"`              // 菜单ID
	CreatedAt int `json:"created_at"`           // 创建时间
}

// RoleAdmin 管理员-角色关联表
type RoleAdmin struct {
	Id        int `gorm:"primaryKey" json:"id"` // 主键ID
	AdminId   int `json:"admin_id"`             // 管理员ID
	RoleId    int `json:"role_id"`              // 角色ID
	CreatedAt int `json:"created_at"`           // 创建时间
}

// RoleWithMenus 角色信息（包含菜单ID列表）
type RoleWithMenus struct {
	Id          int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Pid         int    `json:"pid" gorm:"default:0;comment:父ID"`
	Name        string `json:"name" gorm:"size:64;not null;comment:角色名称"`
	Description string `json:"description" gorm:"size:255;comment:角色描述"`
	IsSuper     int    `json:"is_super" gorm:"default:0;comment:是否超级管理员 0否 1是"`
	Status      int    `json:"status" gorm:"default:1;comment:状态 0禁用 1启用"`
	Sort        int    `json:"sort" gorm:"default:50;comment:排序"`
	CreatedAt   int    `json:"created_at" gorm:"comment:创建时间"`
	UpdatedAt   int    `json:"updated_at" gorm:"comment:更新时间"`
	MenuIds     []int  `json:"menu_ids"` // 菜单ID列表
}

// TableName 设置表名
func (Role) TableName() string {
	return "admin_role"
}

// TableName 表名
func (RoleMenu) TableName() string {
	return "admin_role_menu"
}

// TableName 表名
func (RoleAdmin) TableName() string {
	return "admin_role_admin"
}
