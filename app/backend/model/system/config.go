package system

// Config 配置表
type Config struct {
	Id        uint   `json:"id" gorm:"primaryKey;autoIncrement"`                 // 主键ID
	Dir       string `json:"dir" gorm:"size:255;not null;comment:配置组"`           // 配置组
	Key       string `json:"key" gorm:"size:255;not null;comment:配置键"`           // 配置键
	Name      string `json:"name" gorm:"size:255;not null;comment:配置名称"`         // 配置名称
	Tip       string `json:"tip" gorm:"size:255;default:'';comment:提示说明"`        // 提示说明
	Type      string `json:"type" gorm:"size:255;default:'string';comment:配置类型"` // 配置类型
	Value     string `json:"value" gorm:"type:longtext;comment:配置值"`             // 配置值
	Variable  string `json:"variable" gorm:"size:255;default:'';comment:配置变量"`   // 配置变量
	CreatedAt uint   `json:"created_at" gorm:"not null;comment:创建时间"`            // 创建时间
}

// TableName 表名
func (Config) TableName() string {
	return "config"
}
