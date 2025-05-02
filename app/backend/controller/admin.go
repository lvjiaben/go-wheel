package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/app/backend/model"
	"github.com/lvjiaben/go-wheel/app/backend/service"
	"github.com/lvjiaben/go-wheel/app/backend/validate"
	"github.com/lvjiaben/go-wheel/pkg/container"
	"github.com/lvjiaben/go-wheel/pkg/jwt"
	"go.uber.org/zap"
)

// AdminController 管理员控制器
type AdminController struct {
	container    *container.Container
	adminService *service.AdminService
}

// NewAdminController 创建管理员控制器
func NewAdminController(c *container.Container) *AdminController {
	return &AdminController{
		container:    c,
		adminService: service.NewAdminService(c),
	}
}

// Login 登录
func (c *AdminController) Login(ctx *gin.Context) {
	var req validate.AdminLogin
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.container.GetLogger().Error("参数错误", zap.Error(err))
		ctx.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	admin, err := c.adminService.Login(req.Username, req.Password)
	if err != nil {
		c.container.GetLogger().Error("登录失败", zap.Error(err))
		ctx.JSON(400, gin.H{"code": 400, "msg": err.Error()})
		return
	}

	// 生成token
	token, err := jwt.GenerateToken(admin.Id, c.container.GetConfig().GetString("jwt.secret"), c.container.GetConfig().GetInt("jwt.expire_day"))
	if err != nil {
		c.container.GetLogger().Error("生成token失败", zap.Error(err))
		ctx.JSON(500, gin.H{"code": 500, "msg": "生成token失败"})
		return
	}

	// 更新token
	if err := c.adminService.UpdateToken(admin.Id, token); err != nil {
		c.container.GetLogger().Error("更新token失败", zap.Error(err))
		ctx.JSON(500, gin.H{"code": 500, "msg": "更新token失败"})
		return
	}

	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "登录成功",
		"data": gin.H{
			"token": token,
			"admin": admin,
		},
	})
}

// Logout 退出登录
func (c *AdminController) Logout(ctx *gin.Context) {
	adminId := ctx.GetInt("admin_id")
	if err := c.adminService.UpdateToken(adminId, ""); err != nil {
		c.container.GetLogger().Error("退出登录失败", zap.Error(err))
		ctx.JSON(500, gin.H{"code": 500, "msg": "退出登录失败"})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "msg": "退出登录成功"})
}

// List 获取管理员列表
func (c *AdminController) List(ctx *gin.Context) {
	admins, err := c.adminService.GetList()
	if err != nil {
		c.container.GetLogger().Error("获取管理员列表失败", zap.Error(err))
		ctx.JSON(500, gin.H{"code": 500, "msg": "获取管理员列表失败"})
		return
	}

	ctx.JSON(200, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": admins,
	})
}

// Create 创建管理员
func (c *AdminController) Create(ctx *gin.Context) {
	var req validate.AdminCreate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.container.GetLogger().Error("参数错误", zap.Error(err))
		ctx.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	admin := &model.Admin{
		Username: req.Username,
		Password: req.Password,
		Avatar:   req.Avatar,
	}

	if err := c.adminService.Create(admin); err != nil {
		c.container.GetLogger().Error("创建管理员失败", zap.Error(err))
		ctx.JSON(500, gin.H{"code": 500, "msg": "创建管理员失败"})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "msg": "创建成功"})
}

// Update 更新管理员
func (c *AdminController) Update(ctx *gin.Context) {
	var req validate.AdminUpdate
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.container.GetLogger().Error("参数错误", zap.Error(err))
		ctx.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	admin := &model.Admin{
		Id:       req.Id,
		Username: req.Username,
		Password: req.Password,
		Avatar:   req.Avatar,
	}

	if err := c.adminService.Update(admin); err != nil {
		c.container.GetLogger().Error("更新管理员失败", zap.Error(err))
		ctx.JSON(500, gin.H{"code": 500, "msg": "更新管理员失败"})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "msg": "更新成功"})
}

// Delete 删除管理员
func (c *AdminController) Delete(ctx *gin.Context) {
	var req validate.AdminDelete
	if err := ctx.ShouldBindJSON(&req); err != nil {
		c.container.GetLogger().Error("参数错误", zap.Error(err))
		ctx.JSON(400, gin.H{"code": 400, "msg": "参数错误"})
		return
	}

	if err := c.adminService.Delete(req.Id); err != nil {
		c.container.GetLogger().Error("删除管理员失败", zap.Error(err))
		ctx.JSON(500, gin.H{"code": 500, "msg": "删除管理员失败"})
		return
	}

	ctx.JSON(200, gin.H{"code": 200, "msg": "删除成功"})
}
