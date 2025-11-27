package admin

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/app/backend/model/admin"
	"github.com/lvjiaben/go-wheel/app/backend/utils"

	"github.com/lvjiaben/go-wheel/pkg/container"
	"github.com/lvjiaben/go-wheel/pkg/utils/datatype"
	"gorm.io/gorm"
)

// RoleService 角色服务
type RoleService struct {
	container *container.Container
	authUtils *utils.AuthUtils
}

// NewRoleService 创建角色服务
func NewRoleService(c *container.Container) *RoleService {
	return &RoleService{
		container: c,
		authUtils: utils.NewAuthUtils(c),
	}
}

// GetAll 获取所有角色
func (s *RoleService) GetAll(ctx *gin.Context) ([]map[string]interface{}, error) {
	adminId := ctx.GetInt("admin_id")

	// 检查是否为超级管理员
	isSuper, err := s.authUtils.IsAdminSuper(adminId)
	if err != nil {
		return nil, fmt.Errorf("backend.role.check_super_failed")
	}

	var roles []admin.Role

	// 如果是超级管理员，获取所有角色
	if isSuper {
		if err := s.container.GetDB().Order("sort DESC, id ASC").Find(&roles).Error; err != nil {
			return nil, fmt.Errorf("backend.role.get_list_failed")
		}
	} else {
		// 获取当前管理员的所有可管理角色ID（包括无限级下级）
		adminRoleIds, err := s.authUtils.GetAdminAllRoleIds(adminId)
		if err != nil {
			return nil, err
		}

		if len(adminRoleIds) == 0 {
			return nil, nil
		}

		// 获取这些角色的所有无限级下级角色
		allRoleIds, err := s.authUtils.GetAllChildRoleIds(adminRoleIds)
		if err != nil {
			return nil, err
		}

		if err := s.container.GetDB().Where("id IN ?", allRoleIds).Order("sort DESC, id ASC").Find(&roles).Error; err != nil {
			return nil, fmt.Errorf("backend.role.get_list_failed")
		}
	}

	// 转换为 RoleWithMenus 并添加菜单ID
	result := make([]admin.RoleWithMenus, 0, len(roles))
	for _, role := range roles {
		menuIds, _ := s.authUtils.GetRoleMenuIds(role.Id)

		roleWithMenus := admin.RoleWithMenus{
			Id:          role.Id,
			Pid:         role.Pid,
			Name:        role.Name,
			Description: role.Description,
			IsSuper:     role.IsSuper,
			Status:      role.Status,
			Sort:        role.Sort,
			CreatedAt:   role.CreatedAt,
			UpdatedAt:   role.UpdatedAt,
			MenuIds:     menuIds,
		}
		result = append(result, roleWithMenus)
	}
	treeData := datatype.ToTreeAssocMap(result, "id", "pid", "children")
	return treeData, nil
}

// GetById 根据ID获取角色
func (s *RoleService) GetById(id int) (*admin.Role, error) {
	var role admin.Role
	if err := s.container.GetDB().First(&role, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("backend.role.not_found")
		}
		return nil, fmt.Errorf("backend.role.get_detail_failed")
	}
	return &role, nil
}

// Save 创建或更新角色（包含菜单权限分配）
func (s *RoleService) Save(role *admin.Role, menuIds []int, ctx *gin.Context) error {
	adminId := ctx.GetInt("admin_id")

	// 检查是否为超级管理员
	isSuper, err := s.authUtils.IsAdminSuper(adminId)
	if err != nil {
		return fmt.Errorf("backend.role.check_super_failed")
	}

	// 如果不是超级管理员，需要进行权限检查
	if !isSuper {
		// 普通管理员不能创建或修改超级管理员角色
		if role.IsSuper == 1 {
			return fmt.Errorf("backend.role.cannot_create_super_role")
		}

		if role.Pid == 0 {
			return fmt.Errorf("backend.role.parent_not_in_scope")
		}

		// 检查父级角色是否存在且在当前管理员的管理范围内
		if role.Pid > 0 {
			if exists, err := s.authUtils.CheckRoleInAdminScope(role.Pid, adminId); err != nil {
				return fmt.Errorf("backend.role.check_parent_failed")
			} else if !exists {
				return fmt.Errorf("backend.role.parent_not_in_scope")
			}
		}

		// 更新时，检查要修改的角色是否在管理范围内
		if role.Id > 0 {
			if exists, err := s.authUtils.CheckRoleInAdminScope(role.Id, adminId); err != nil {
				return fmt.Errorf("backend.role.check_scope_failed")
			} else if !exists {
				return fmt.Errorf("backend.role.not_in_scope")
			}

			// 不能修改自己直属的角色组
			if s.authUtils.IsAdminDirectRole(role.Id, adminId) {
				return fmt.Errorf("backend.role.cannot_modify_own_role")
			}

			// 检查原角色是否为超级管理员角色，普通管理员不能修改超级管理员角色
			var originalRole admin.Role
			if err := s.container.GetDB().First(&originalRole, role.Id).Error; err == nil {
				if originalRole.IsSuper == 1 {
					return fmt.Errorf("backend.role.cannot_modify_super_role")
				}
			}
		}
	}

	// 检查角色名称是否重复
	if err := s.checkRoleNameExists(role.Name, role.Id); err != nil {
		return err
	}

	// 使用事务保存角色和菜单权限
	return s.container.GetDB().Transaction(func(tx *gorm.DB) error {
		// 设置时间戳
		now := int(time.Now().Unix())
		if role.Id > 0 {
			role.UpdatedAt = now
			// 更新角色，使用 Select 明确指定要更新的字段，包括零值字段
			updateData := map[string]interface{}{
				"pid":         role.Pid,
				"name":        role.Name,
				"description": role.Description,
				"is_super":    role.IsSuper,
				"status":      role.Status,
				"sort":        role.Sort,
				"updated_at":  role.UpdatedAt,
			}

			fields := []string{"pid", "name", "description", "is_super", "status", "sort", "updated_at"}
			if err := tx.Model(&admin.Role{}).Where("id = ?", role.Id).Select(fields).Updates(updateData).Error; err != nil {
				return fmt.Errorf("backend.role.update_failed")
			}
		} else {
			role.CreatedAt = now
			role.UpdatedAt = now
			// 创建角色
			if err := tx.Create(role).Error; err != nil {
				return fmt.Errorf("backend.role.create_failed")
			}
		}

		// 分配菜单权限
		// 删除旧的菜单关联
		if err := tx.Where("role_id = ?", role.Id).Delete(&admin.RoleMenu{}).Error; err != nil {
			return fmt.Errorf("backend.role.delete_old_menus_failed")
		}

		// 添加新的菜单关联
		if len(menuIds) > 0 {
			roleMenus := make([]admin.RoleMenu, 0, len(menuIds))
			for _, menuId := range menuIds {
				roleMenus = append(roleMenus, admin.RoleMenu{
					RoleId:    role.Id,
					MenuId:    menuId,
					CreatedAt: now,
				})
			}

			if err := tx.Create(&roleMenus).Error; err != nil {
				return fmt.Errorf("backend.role.assign_menus_failed")
			}
		}

		return nil
	})
}

// Delete 删除角色
func (s *RoleService) Delete(id int, ctx *gin.Context) error {
	adminId := ctx.GetInt("admin_id")

	// 检查是否为超级管理员
	isSuper, err := s.authUtils.IsAdminSuper(adminId)
	if err != nil {
		return fmt.Errorf("backend.role.check_super_failed")
	}

	// 如果不是超级管理员，需要进行权限检查
	if !isSuper {
		// 检查要删除的角色是否在管理范围内
		if exists, err := s.authUtils.CheckRoleInAdminScope(id, adminId); err != nil {
			return fmt.Errorf("backend.role.check_scope_failed")
		} else if !exists {
			return fmt.Errorf("backend.role.not_in_scope")
		}

		// 不能删除自己直属的角色组
		if s.authUtils.IsAdminDirectRole(id, adminId) {
			return fmt.Errorf("backend.role.cannot_delete_own_role")
		}
	}

	// 开始事务删除（自动清理所有关联关系）
	return s.container.GetDB().Transaction(func(tx *gorm.DB) error {
		// 删除角色菜单关联
		if err := tx.Where("role_id = ?", id).Delete(&admin.RoleMenu{}).Error; err != nil {
			return fmt.Errorf("backend.role.delete_role_menu_failed")
		}

		// 删除管理员角色关联
		if err := tx.Where("role_id = ?", id).Delete(&admin.RoleAdmin{}).Error; err != nil {
			return fmt.Errorf("backend.role.delete_role_admin_failed")
		}

		// 删除角色
		if err := tx.Delete(&admin.Role{}, id).Error; err != nil {
			return fmt.Errorf("backend.role.delete_failed")
		}

		return nil
	})
}

// GetMyMenus 获取当前管理员的菜单列表
func (s *RoleService) GetMyMenus(ctx *gin.Context) ([]admin.MenuTree, error) {
	adminId := ctx.GetInt("admin_id")

	// 检查是否为超级管理员
	isSuper, err := s.authUtils.IsAdminSuper(adminId)
	if err != nil {
		return nil, fmt.Errorf("backend.role.check_super_failed")
	}

	var menus []admin.Menu

	// 如果是超级管理员，返回所有菜单
	if isSuper {
		if err := s.container.GetDB().Order("sort DESC").Find(&menus).Error; err != nil {
			return nil, fmt.Errorf("backend.menu.get_list_failed")
		}
	} else {
		// 获取管理员的所有菜单权限
		menuIds, err := s.authUtils.GetAdminMenuIds(adminId)
		if err != nil {
			return nil, err
		}

		if len(menuIds) == 0 {
			return []admin.MenuTree{}, nil
		}

		if err := s.container.GetDB().Where("id IN ?", menuIds).Order("sort DESC").Find(&menus).Error; err != nil {
			return nil, fmt.Errorf("backend.menu.get_list_failed")
		}
	}

	// 构建菜单树
	return s.buildMenuTree(menus), nil
}

// ===== 私有辅助方法 =====

// checkRoleNameExists 检查角色名称是否重复
func (s *RoleService) checkRoleNameExists(name string, excludeId int) error {
	var count int64
	query := s.container.GetDB().Model(&admin.Role{}).Where("name = ?", name)
	if excludeId > 0 {
		query = query.Where("id != ?", excludeId)
	}

	if err := query.Count(&count).Error; err != nil {
		return fmt.Errorf("backend.role.check_name_failed")
	}

	if count > 0 {
		return fmt.Errorf("backend.role.name_exists")
	}

	return nil
}

// buildMenuTree 构建菜单树
func (s *RoleService) buildMenuTree(menus []admin.Menu) []admin.MenuTree {
	menuMap := make(map[int][]admin.Menu)

	// 构建映射
	for _, menu := range menus {
		menuMap[menu.Pid] = append(menuMap[menu.Pid], menu)
	}

	// 构建树
	return s.getMenuChildren(0, menuMap)
}

// getMenuChildren 递归获取子菜单
func (s *RoleService) getMenuChildren(pid int, menuMap map[int][]admin.Menu) []admin.MenuTree {
	children := menuMap[pid]
	if len(children) == 0 {
		return []admin.MenuTree{}
	}

	result := make([]admin.MenuTree, 0, len(children))
	for _, child := range children {
		tree := admin.MenuTree{
			Id:         child.Id,
			Pid:        child.Pid,
			Name:       child.Name,
			Enname:     child.Enname,
			Route:      child.Route,
			Component:  child.Component,
			Path:       child.Path,
			Icon:       child.Icon,
			Sort:       child.Sort,
			Visible:    child.Visible,
			FixedTag:   child.FixedTag,
			ShowTag:    child.ShowTag,
			Iframe:     child.Iframe,
			External:   child.External,
			Type:       child.Type,
			Permission: child.Permission,
			CreatedAt:  child.CreatedAt,
			UpdatedAt:  child.UpdatedAt,
			Children:   s.getMenuChildren(child.Id, menuMap),
		}
		result = append(result, tree)
	}

	return result
}
