package system

// GenHistory 代码生成历史记录
type GenHistory struct {
	Id              int    `json:"id" gorm:"primaryKey;autoIncrement"`     // 主键ID
	GenTableName    string `json:"table_name" gorm:"column:table_name;size:255;not null"` // 表名
	TableComment    string `json:"table_comment" gorm:"size:255"`          // 表注释
	ModuleName      string `json:"module_name" gorm:"size:255;not null"`   // 模块名
	StructName      string `json:"struct_name" gorm:"size:255;not null"`   // 结构体名
	PackageName     string `json:"package_name" gorm:"size:255;not null"`  // 包名
	FrontendSrcPath string `json:"frontend_src_path" gorm:"size:500"`      // 前端 src 路径
	Config          string `json:"config" gorm:"type:text"`                // 完整配置（JSON）
	CreatedAt       int    `json:"created_at" gorm:"autoCreateTime"`       // 创建时间
	UpdatedAt       int    `json:"updated_at" gorm:"autoUpdateTime"`       // 更新时间
}

// TableName 表名
func (GenHistory) TableName() string {
	return "gen_history"
}

