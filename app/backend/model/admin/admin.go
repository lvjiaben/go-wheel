package admin

// Admin 管理员表
type Admin struct {
	Id            int    `json:"id" gorm:"primaryKey"`             // 主键ID
	Pid           int    `json:"pid" gorm:"default:0"`             // 上级管理员ID
	Username      string `json:"username" gorm:"not null"`         // 用户名
	Password      string `json:"password" gorm:"not null"`         // 密码
	Salt          string `json:"salt" gorm:"default:''"`           // 密码盐
	Realname      string `json:"realname" gorm:"default:''"`       // 真实姓名
	Avatar        string `json:"avatar" gorm:"default:''"`         // 头像
	Email         string `json:"email" gorm:"default:''"`          // 邮箱
	Mobile        string `json:"mobile" gorm:"default:''"`         // 手机号
	Failures      int    `json:"failures" gorm:"default:0"`        // 登录失败次数
	Status        int    `json:"status" gorm:"default:1"`          // 状态：0=禁用，1=启用
	Token         string `json:"token" gorm:"default:''"`          // Token
	LastLoginTime int    `json:"last_login_time" gorm:"default:0"` // 最后登录时间
	CreatedAt     int    `json:"created_at" gorm:"autoCreateTime"` // 创建时间
	UpdatedAt     int    `json:"updated_at" gorm:"autoUpdateTime"` // 更新时间
}

// TableName 表名
func (Admin) TableName() string {
	return "admin"
}

// AdminInfo 管理员信息
type AdminInfo struct {
	Id          int          `json:"id"`          // 主键ID
	Username    string       `json:"username"`    // 用户名
	Realname    string       `json:"realname"`    // 真实姓名
	Avatar      string       `json:"avatar"`      // 头像
	Email       string       `json:"email"`       // 邮箱
	Mobile      string       `json:"mobile"`      // 手机号
	Status      int          `json:"status"`      // 状态
	IsSuper     bool         `json:"is_super"`    // 是否超级管理员
	Roles       []SimpleRole `json:"roles"`       // 角色列表
	Permissions []string     `json:"permissions"` // 权限列表
}

// SimpleRole 角色简要信息
type SimpleRole struct {
	Id   int    `json:"id"`   // 角色ID
	Name string `json:"name"` // 角色名称
}

// AdminWithRoles 管理员信息（包含角色ID列表）
type AdminWithRoles struct {
	Id            int    `json:"id"`              // 主键ID
	Pid           int    `json:"pid"`             // 上级管理员ID
	Username      string `json:"username"`        // 用户名
	Realname      string `json:"realname"`        // 真实姓名
	Avatar        string `json:"avatar"`          // 头像
	Email         string `json:"email"`           // 邮箱
	Mobile        string `json:"mobile"`          // 手机号
	Status        int    `json:"status"`          // 状态：0=禁用，1=启用
	LastLoginTime int    `json:"last_login_time"` // 最后登录时间
	CreatedAt     int    `json:"created_at"`      // 创建时间
	UpdatedAt     int    `json:"updated_at"`      // 更新时间
	RoleIds       []int  `json:"role_ids"`        // 角色ID列表
}
