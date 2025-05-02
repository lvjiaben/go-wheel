package validate

// AdminLogin 管理员登录验证
type AdminLogin struct {
	Username string `json:"username" binding:"required" label:"用户名"`
	Password string `json:"password" binding:"required" label:"密码"`
}

// AdminCreate 创建管理员验证
type AdminCreate struct {
	Username string `json:"username" binding:"required" label:"用户名"`
	Password string `json:"password" binding:"required" label:"密码"`
	Avatar   string `json:"avatar" label:"头像"`
	Failures int    `json:"failures" binding:"required" msg:"登陆失败次数有误"`
	Token    string `json:"token" binding:"-" msg:"TOKEN有误"`
}

// AdminUpdate 更新管理员验证
type AdminUpdate struct {
	Id       int    `json:"id" binding:"required" label:"ID"`
	Username string `json:"username" binding:"required" label:"用户名"`
	Password string `json:"password" label:"密码"`
	Avatar   string `json:"avatar" label:"头像"`
	Failures int    `json:"failures" binding:"required" msg:"登陆失败次数有误"`
	Token    string `json:"token" binding:"-" msg:"TOKEN有误"`
}

// AdminDelete 删除管理员验证
type AdminDelete struct {
	Id int `json:"id" binding:"required" label:"ID"`
}

type AdminSort struct {
}

type AdminAuthGroupCreate struct {
	Pid   int    `json:"pid" binding:"required" msg:"上级角色有误"`
	Name  string `json:"name" binding:"required" msg:"角色名称有误"`
	Rules string `json:"rules" binding:"required" msg:"角色权限有误"`
}

type AdminAuthGroupUpdate struct {
	Id    int    `json:"id" binding:"required" msg:"Id有误"`
	Pid   int    `json:"pid" binding:"required" msg:"上级角色有误"`
	Name  string `json:"name" binding:"required" msg:"角色名称有误"`
	Rules string `json:"rules" binding:"required" msg:"角色权限有误"`
}

type AdminAuthGroupDelete struct {
	Id int `json:"id" binding:"required" msg:"Id有误"`
}

type AdminAuthRuleCreate struct {
	Pid    int    `json:"pid" binding:"required" msg:"上级权限有误"`
	Name   string `json:"name" binding:"required" msg:"权限名称有误"`
	Enname string `json:"enname" binding:"required" msg:"英文名称有误"`
	Path   string `json:"path" binding:"required" msg:"权限路径有误"`
	Method string `json:"method" binding:"required" msg:"请求方法有误"`
	Type   string `json:"type" binding:"required" msg:"权限类型有误"`
	Status int    `json:"status" binding:"required" msg:"状态有误"`
	Sort   int    `json:"sort" binding:"required" msg:"排序有误"`
	Icon   string `json:"icon" binding:"-" msg:"图标有误"`
	Hide   int    `json:"hide" binding:"required" msg:"是否隐藏有误"`
	Alias  string `json:"alias" binding:"-" msg:"别名有误"`
}

type AdminAuthRuleUpdate struct {
	Id     int    `json:"id" binding:"required" msg:"Id有误"`
	Pid    int    `json:"pid" binding:"required" msg:"上级权限有误"`
	Name   string `json:"name" binding:"required" msg:"权限名称有误"`
	Enname string `json:"enname" binding:"required" msg:"英文名称有误"`
	Path   string `json:"path" binding:"required" msg:"权限路径有误"`
	Method string `json:"method" binding:"required" msg:"请求方法有误"`
	Type   string `json:"type" binding:"required" msg:"权限类型有误"`
	Status int    `json:"status" binding:"required" msg:"状态有误"`
	Sort   int    `json:"sort" binding:"required" msg:"排序有误"`
	Icon   string `json:"icon" binding:"-" msg:"图标有误"`
	Hide   int    `json:"hide" binding:"required" msg:"是否隐藏有误"`
	Alias  string `json:"alias" binding:"-" msg:"别名有误"`
}

type AdminAuthRuleDelete struct {
	Id int `json:"id" binding:"required" msg:"Id有误"`
}
