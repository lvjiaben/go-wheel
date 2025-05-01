package model

type AdminAuthRule struct {
	Id        int    `json:"id" gorm:"column:id"`
	Sort      int    `json:"sort" gorm:"column:sort"`             // 排序
	Status    int    `json:"status" gorm:"column:status"`         // 状态
	Pid       int    `json:"pid" gorm:"column:pid"`               // 上级
	Route     string `json:"route" gorm:"column:route"`           // 路由
	Name      string `json:"name" gorm:"column:name"`             // 名称
	Alias     string `json:"alias" gorm:"column:alias"`           // 别名
	Icon      string `json:"icon" gorm:"column:icon"`             // 图标
	Type      string `json:"type" gorm:"column:type"`             // 类型
	Path      string `json:"path" gorm:"column:path"`             // 地址
	Enname    string `json:"enname" gorm:"column:enname"`         // 国际化
	Hide      int    `json:"hide" gorm:"column:hide"`             // 隐藏
	CreatedAt int    `json:"created_at" gorm:"column:created_at"` // 创建时间
	UpdatedAt int    `json:"updated_at" gorm:"column:updated_at"` // 更新时间
}

func (e *AdminAuthRule) TableName() string {
	return "admin_auth_rule"
}
