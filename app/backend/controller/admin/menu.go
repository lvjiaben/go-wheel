package admin

import (
	"github.com/lvjiaben/go-wheel/pkg/utils/http"

	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/app/backend/model/admin"
	"github.com/lvjiaben/go-wheel/app/backend/service"
	serviceAdmin "github.com/lvjiaben/go-wheel/app/backend/service/admin"
	validateAdmin "github.com/lvjiaben/go-wheel/app/backend/validate/admin"
	"github.com/lvjiaben/go-wheel/pkg/container"
)

// MenuController 菜单控制器
type MenuController struct {
	container   *container.Container
	menuService *serviceAdmin.MenuService
	authService *service.AuthService
}

// NewMenuController 创建菜单控制器
func NewMenuController(c *container.Container) *MenuController {
	return &MenuController{
		container:   c,
		menuService: serviceAdmin.NewMenuService(c),
		authService: service.NewAuthService(c),
	}
}

// List 获取菜单列表
func (m *MenuController) List(c *gin.Context) {

	// 获取菜单树
	menus, err := m.menuService.GetAll()
	if err != nil {
		http.ErrorWithI18n(c, "common.server_error", nil)
		return
	}

	http.SuccessWithI18n(c, "common.success", menus)
}

// Save 保存
func (m *MenuController) Save(c *gin.Context) {
	// 使用ValidateMenuSave进行验证
	menuForm, valid := validateAdmin.ValidateMenuSave(c)
	if !valid {
		return
	}
	// 创建菜单对象
	menu := admin.Menu{
		Id:         menuForm.Id,
		Pid:        menuForm.Pid,
		Sort:       menuForm.Sort,
		Name:       menuForm.Name,
		Enname:     menuForm.Enname,
		Route:      menuForm.Route,
		Component:  menuForm.Component,
		Path:       menuForm.Path,
		Icon:       menuForm.Icon,
		FixedTag:   menuForm.FixedTag,
		ShowTag:    menuForm.ShowTag,
		Visible:    menuForm.Visible,
		Iframe:     menuForm.Iframe,
		External:   menuForm.External,
		Type:       menuForm.Type,
		Permission: menuForm.Permission,
	}
	if err := m.menuService.Save(&menu); err != nil {
		http.ErrorWithI18n(c, err.Error(), nil)
		return
	}
	http.SuccessWithI18n(c, "common.success", menu)
}

// Delete 删除菜单
func (m *MenuController) Delete(c *gin.Context) {
	// 使用ValidateMenuSave进行验证
	menuForm, valid := validateAdmin.ValidateMenuDelete(c)
	if !valid {
		return
	}
	// 删除菜单（会自动删除角色菜单关联）
	if err := m.menuService.Delete(menuForm.Id); err != nil {
		http.ErrorWithI18n(c, err.Error(), nil)
		return
	}
	http.SuccessWithI18n(c, "common.success", nil)
}
