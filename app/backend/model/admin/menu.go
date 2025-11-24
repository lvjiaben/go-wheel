package admin

// Menu 菜单模型
type Menu struct {
	Id         int    `json:"id" gorm:"primaryKey;autoIncrement"`
	Pid        int    `json:"pid" gorm:"default:0"`
	Name       string `json:"name" gorm:"size:64;not null"`
	Enname     string `json:"enname" gorm:"size:64;not null"`
	Route      string `json:"route" gorm:"size:128;"`
	Component  string `json:"component" gorm:"size:128;"`
	Path       string `json:"path" gorm:"size:128;"`
	Icon       string `json:"icon" gorm:"size:64"`
	Sort       int    `json:"sort" gorm:"default:0"`
	Visible    int    `json:"visible" gorm:"default:1"`
	FixedTag   int    `json:"fixed_tag" gorm:"default:0"`
	ShowTag    int    `json:"show_tag" gorm:"default:0"`
	Iframe     string `json:"iframe" gorm:"size:128"`
	External   string `json:"external" gorm:"size:128"`
	Type       string `json:"type" gorm:"size:16;not null"`
	Permission string `json:"permission" gorm:"size:128"`
	CreatedAt  int    `json:"created_at" gorm:"not null"`
	UpdatedAt  int    `json:"updated_at" gorm:"not null"`
}

// TableName 表名
func (Menu) TableName() string {
	return "admin_menu"
}

// MenuTree 菜单树结构
type MenuTree struct {
	Id         int        `json:"id"`
	Pid        int        `json:"pid"`
	Name       string     `json:"name"`
	Enname     string     `json:"enname"`
	Route      string     `json:"route"`
	Component  string     `json:"component"`
	Path       string     `json:"path"`
	Icon       string     `json:"icon"`
	Sort       int        `json:"sort"`
	Visible    int        `json:"visible"`
	FixedTag   int        `json:"fixed_tag"`
	ShowTag    int        `json:"show_tag"`
	Iframe     string     `json:"iframe"`
	External   string     `json:"external"`
	Type       string     `json:"type"`
	Permission string     `json:"permission"`
	CreatedAt  int        `json:"created_at"`
	UpdatedAt  int        `json:"updated_at"`
	Children   []MenuTree `json:"children"`
}

// MenuItem 前端菜单项
type MenuItem struct {
	Name      string      `json:"name"`               // 菜单名称（英文）
	Path      string      `json:"path"`               // 路由路径
	Component string      `json:"component"`          // 组件路径
	Meta      MenuMeta    `json:"meta"`               // 菜单元数据
	Children  []*MenuItem `json:"children,omitempty"` // 子菜单
}

// MenuMeta 菜单元数据
type MenuMeta struct {
	Title      string `json:"title"`      // 菜单标题（中文）
	Icon       string `json:"icon"`       // 菜单图标
	HideInMenu bool   `json:"hideInMenu"` // 在菜单中隐藏
	Order      int    `json:"order"`      // 排序值
}
