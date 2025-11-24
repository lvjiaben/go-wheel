package system

// Attachment 附件表
type Attachment struct {
	Id        int    `json:"id" gorm:"primaryKey"`             // 主键ID
	Type      string `json:"type" gorm:"default:'local'"`      // 存储类型
	AdminId   int    `json:"admin_id" gorm:"default:0;index"`  // 管理员ID
	UserId    int    `json:"user_id" gorm:"default:0;index"`   // 用户ID
	Path      string `json:"path" gorm:"not null;index"`       // 存储路径
	Parent    string `json:"parent" gorm:"index"`              // 父级文件夹
	URL       string `json:"url" gorm:"not null;index"`        // 在线HTTP链接
	Filename  string `json:"filename" gorm:"not null;index"`   // 文件名称
	Size      int64  `json:"size" gorm:"not null"`             // 文件大小
	MediaType string `json:"mediatype" gorm:"not null"`        // 文件类型
	Extension string `json:"extension" gorm:"not null"`        // 文件后缀
	CreatedAt int    `json:"created_at" gorm:"autoCreateTime"` // 创建时间
	UpdatedAt int    `json:"updated_at" gorm:"autoUpdateTime"` // 更新时间
}

// TableName 设置表名
func (Attachment) TableName() string {
	return "attachment"
}

// AttachmentWithCount 带统计的附件结构
type AttachmentWithCount struct {
	Attachment
	Count int64 `json:"count,omitempty"` // 子文件数量（用于文件夹）
}
