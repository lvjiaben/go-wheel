package validate

type AdminLogin struct {
	Username string `json:"username" binding:"required" msg:"管理员账号有误"`
	Password string `json:"password" binding:"required" msg:"管理员密码有误"`
}

type AdminCreate struct {
	Pid      int    `json:"pid" binding:"required" msg:"上级管理员有误"`
	Username string `json:"username" binding:"required" msg:"管理员账号有误"`
	Password string `json:"password" binding:"required" msg:"管理员密码有误"`
	Avatar   string `json:"avatar" binding:"-" msg:"管理员头像有误"`
	Failures int    `json:"failures" binding:"required" msg:"登陆失败次数有误"`
	Token    string `json:"token" binding:"-" msg:"TOKEN有误"`
}

type AdminUpdate struct {
	Id       int    `json:"id" binding:"required" msg:"Id有误"`
	Pid      int    `json:"pid" binding:"required" msg:"上级管理员有误"`
	Username string `json:"username" binding:"required" msg:"管理员账号有误"`
	Password string `json:"password" binding:"required" msg:"管理员密码有误"`
	Avatar   string `json:"avatar" binding:"-" msg:"管理员头像有误"`
	Failures int    `json:"failures" binding:"required" msg:"登陆失败次数有误"`
	Token    string `json:"token" binding:"-" msg:"TOKEN有误"`
}

type AdminDelete struct {
	Id int `json:"id" binding:"required" msg:"Id有误"`
}

type AdminSort struct {
}
