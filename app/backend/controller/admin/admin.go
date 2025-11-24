package admin

import (
	"strconv"

	"github.com/lvjiaben/go-wheel/pkg/utils/http"

	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/app/backend/model/admin"
	"github.com/lvjiaben/go-wheel/app/backend/service"
	serviceAdmin "github.com/lvjiaben/go-wheel/app/backend/service/admin"
	validateAdmin "github.com/lvjiaben/go-wheel/app/backend/validate/admin"
	"github.com/lvjiaben/go-wheel/pkg/container"
)

// AdminController 管理员控制器
type AdminController struct {
	container    *container.Container
	adminService *serviceAdmin.AdminService
	authService  *service.AuthService
}

// NewAdminController 创建管理员控制器
func NewAdminController(c *container.Container) *AdminController {
	return &AdminController{
		container:    c,
		adminService: serviceAdmin.NewAdminService(c),
		authService:  service.NewAuthService(c),
	}
}

// List 获取管理员列表
func (c *AdminController) List(ctx *gin.Context) {
	admins, err := c.adminService.GetAll(ctx)
	if err != nil {
		http.ErrorWithI18n(ctx, err.Error(), nil)
		return
	}

	http.SuccessWithI18n(ctx, "common.success", admins)
}

// Save 创建或更新管理员
func (c *AdminController) Save(ctx *gin.Context) {
	// 验证请求数据
	adminSave, valid := validateAdmin.ValidateAdminSave(ctx)
	if !valid {
		return
	}

	// 创建时密码必填，更新时密码可选
	if adminSave.Id == 0 && adminSave.Password == "" {
		http.ErrorWithI18n(ctx, "backend.admin.password_required", nil)
		return
	}

	// 创建管理员对象
	adminItem := &admin.Admin{
		Id:       adminSave.Id,
		Pid:      adminSave.Pid,
		Username: adminSave.Username,
		Realname: adminSave.Realname,
		Avatar:   adminSave.Avatar,
		Email:    adminSave.Email,
		Mobile:   adminSave.Mobile,
		Status:   adminSave.Status,
	}

	// 如果密码不为空，设置密码
	if adminSave.Password != "" {
		adminItem.Password = adminSave.Password
	}

	// 保存管理员
	if err := c.adminService.Save(adminItem, adminSave.RoleIds, ctx); err != nil {
		http.ErrorWithI18n(ctx, err.Error(), nil)
		return
	}

	// 不返回密码
	adminItem.Password = ""

	http.SuccessWithI18n(ctx, "common.success", adminItem)
}

// Delete 删除管理员
func (c *AdminController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.ErrorWithI18n(ctx, "common.invalid_params", nil)
		return
	}

	if err := c.adminService.Delete(id, ctx); err != nil {
		http.ErrorWithI18n(ctx, err.Error(), nil)
		return
	}

	http.SuccessWithI18n(ctx, "common.success", nil)
}
