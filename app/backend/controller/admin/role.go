package admin

import (
	"strconv"

	"github.com/lvjiaben/go-wheel/pkg/utils/http"

	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/app/backend/model/admin"
	"github.com/lvjiaben/go-wheel/app/backend/service"
	serviceAdmin "github.com/lvjiaben/go-wheel/app/backend/service/admin"
	"github.com/lvjiaben/go-wheel/app/backend/utils"
	validateAdmin "github.com/lvjiaben/go-wheel/app/backend/validate/admin"
	"github.com/lvjiaben/go-wheel/pkg/container"
)

// RoleController 角色控制器
type RoleController struct {
	container   *container.Container
	roleService *serviceAdmin.RoleService
	authService *service.AuthService
	authUtils   *utils.AuthUtils
}

// NewRoleController 创建角色控制器
func NewRoleController(c *container.Container) *RoleController {
	return &RoleController{
		container:   c,
		roleService: serviceAdmin.NewRoleService(c),
		authService: service.NewAuthService(c),
		authUtils:   utils.NewAuthUtils(c),
	}
}

// List 获取角色列表
func (c *RoleController) List(ctx *gin.Context) {
	roles, err := c.roleService.GetAll(ctx)
	if err != nil {
		http.ErrorWithI18n(ctx, err.Error(), nil)
		return
	}

	http.SuccessWithI18n(ctx, "common.success", roles)
}

// Save 创建或更新角色
func (c *RoleController) Save(ctx *gin.Context) {
	// 验证请求数据
	roleSave, valid := validateAdmin.ValidateRoleSave(ctx)
	if !valid {
		return
	}

	if roleSave.Id == 1 {
		http.ErrorWithI18n(ctx, "common.error", nil)
		return
	}

	// 提前进行权限检查，避免不必要的验证
	adminId := ctx.GetInt("admin_id")
	isSuper, err := c.authUtils.IsAdminSuper(adminId)
	if err != nil {
		http.ErrorWithI18n(ctx, "backend.role.check_super_failed", nil)
		return
	}

	// 如果不是超级管理员且试图创建/修改超级管理员角色，直接拒绝
	if !isSuper && roleSave.IsSuper == 1 {
		http.ErrorWithI18n(ctx, "backend.role.cannot_create_super_role", nil)
		return
	}

	// 如果是更新操作，检查原角色是否为超级管理员角色
	if !isSuper && roleSave.Id > 0 {
		originalRole, err := c.roleService.GetById(roleSave.Id)
		if err == nil && originalRole.IsSuper == 1 {
			http.ErrorWithI18n(ctx, "backend.role.cannot_modify_super_role", nil)
			return
		}
	}

	// 创建角色对象
	role := &admin.Role{
		Id:          roleSave.Id,
		Pid:         roleSave.Pid,
		Name:        roleSave.Name,
		Description: roleSave.Description,
		IsSuper:     roleSave.IsSuper,
		Status:      roleSave.Status,
		Sort:        roleSave.Sort,
	}

	// 保存角色（包含菜单权限）
	if err := c.roleService.Save(role, roleSave.MenuIds, ctx); err != nil {
		http.ErrorWithI18n(ctx, err.Error(), nil)
		return
	}

	if role.Id > 0 {
		http.SuccessWithI18n(ctx, "common.success", role)
	} else {
		http.SuccessWithI18n(ctx, "common.success", role)
	}
}

// Delete 删除角色
func (c *RoleController) Delete(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.ErrorWithI18n(ctx, "common.invalid_params", nil)
		return
	}

	if id == 1 {
		http.ErrorWithI18n(ctx, "common.error", nil)
		return
	}

	if err := c.roleService.Delete(id, ctx); err != nil {
		http.ErrorWithI18n(ctx, err.Error(), nil)
		return
	}

	http.SuccessWithI18n(ctx, "common.success", nil)
}

// GetMyMenus 获取当前管理员的菜单列表
func (c *RoleController) GetMyMenus(ctx *gin.Context) {
	menus, err := c.roleService.GetMyMenus(ctx)
	if err != nil {
		http.ErrorWithI18n(ctx, err.Error(), nil)
		return
	}

	http.SuccessWithI18n(ctx, "common.success", menus)
}
