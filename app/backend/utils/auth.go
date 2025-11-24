package utils

import (
	"github.com/lvjiaben/go-wheel/app/backend/model/admin"
	"github.com/lvjiaben/go-wheel/pkg/container"
	"github.com/lvjiaben/go-wheel/pkg/utils/datatype"
)

// AuthUtils 认证相关的工具方法
type AuthUtils struct {
	container *container.Container
}

// NewAuthUtils 创建认证工具实例
func NewAuthUtils(c *container.Container) *AuthUtils {
	return &AuthUtils{container: c}
}

// IsAdminSuper 检查管理员是否为超级管理员
func (u *AuthUtils) IsAdminSuper(adminId int) (bool, error) {
	var count int64
	err := u.container.GetDB().Table("admin_role_admin ara").
		Joins("JOIN admin_role ar ON ara.role_id = ar.id").
		Where("ara.admin_id = ? AND ar.is_super = 1", adminId).
		Count(&count).Error

	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetAdminDirectRoleIds 获取管理员的直接角色ID
func (u *AuthUtils) GetAdminDirectRoleIds(adminId int) ([]int, error) {
	var roleIds []int
	err := u.container.GetDB().Model(&admin.RoleAdmin{}).
		Where("admin_id = ?", adminId).
		Pluck("role_id", &roleIds).Error

	return roleIds, err
}

// GetAdminAllRoleIds 获取管理员的所有角色ID（包括无限级下级）
func (u *AuthUtils) GetAdminAllRoleIds(adminId int) ([]int, error) {
	// 先获取直接角色ID
	directRoleIds, err := u.GetAdminDirectRoleIds(adminId)
	if err != nil {
		return nil, err
	}

	if len(directRoleIds) == 0 {
		return []int{}, nil
	}

	// 获取所有下级角色ID
	return u.GetAllChildRoleIds(directRoleIds)
}

// GetAllChildRoleIds 递归获取角色的所有下级角色ID
func (u *AuthUtils) GetAllChildRoleIds(parentIds []int) ([]int, error) {
	if len(parentIds) == 0 {
		return []int{}, nil
	}

	allIds := make([]int, 0)
	visited := make(map[int]bool)

	// 添加父级ID
	for _, id := range parentIds {
		if !visited[id] {
			allIds = append(allIds, id)
			visited[id] = true
		}
	}

	// 递归获取下级
	currentLevel := parentIds
	for len(currentLevel) > 0 {
		var childIds []int
		err := u.container.GetDB().Model(&admin.Role{}).
			Where("pid IN ?", currentLevel).
			Pluck("id", &childIds).Error

		if err != nil {
			return nil, err
		}

		if len(childIds) == 0 {
			break
		}

		nextLevel := make([]int, 0)
		for _, id := range childIds {
			if !visited[id] {
				allIds = append(allIds, id)
				visited[id] = true
				nextLevel = append(nextLevel, id)
			}
		}
		currentLevel = nextLevel
	}

	return allIds, nil
}

// GetRoleMenuIds 获取角色的菜单ID列表
func (u *AuthUtils) GetRoleMenuIds(roleId int) ([]int, error) {
	var menuIds []int
	err := u.container.GetDB().Model(&admin.RoleMenu{}).
		Where("role_id = ?", roleId).
		Pluck("menu_id", &menuIds).Error

	return menuIds, err
}

// GetAdminMenuIds 获取管理员的所有菜单ID
func (u *AuthUtils) GetAdminMenuIds(adminId int) ([]int, error) {
	var menuIds []int
	err := u.container.GetDB().Table("admin_role_admin ara").
		Joins("JOIN admin_role_menu arm ON ara.role_id = arm.role_id").
		Where("ara.admin_id = ?", adminId).
		Distinct("arm.menu_id").
		Pluck("arm.menu_id", &menuIds).Error

	return menuIds, err
}

// CheckRoleInAdminScope 检查角色是否在管理员的管理范围内
func (u *AuthUtils) CheckRoleInAdminScope(roleId, adminId int) (bool, error) {
	adminRoleIds, err := u.GetAdminAllRoleIds(adminId)
	if err != nil {
		return false, err
	}

	return datatype.ContainsInt(adminRoleIds, roleId), nil
}

// IsAdminDirectRole 检查是否为管理员的直属角色
func (u *AuthUtils) IsAdminDirectRole(roleId, adminId int) bool {
	directRoleIds, err := u.GetAdminDirectRoleIds(adminId)
	if err != nil {
		return false
	}

	return datatype.ContainsInt(directRoleIds, roleId)
}

// GetAdminRoles 获取管理员的角色列表
func (u *AuthUtils) GetAdminRoles(adminId int) ([]admin.SimpleRole, error) {
	var roles []admin.SimpleRole
	err := u.container.GetDB().Table("admin_role ar").
		Select("ar.id, ar.name").
		Joins("JOIN admin_role_admin ara ON ar.id = ara.role_id").
		Where("ara.admin_id = ?", adminId).
		Scan(&roles).Error

	return roles, err
}

// GetAdminAllSubordinateIds 获取管理员的所有下级ID（包括无限级）
func (u *AuthUtils) GetAdminAllSubordinateIds(adminId int) ([]int, error) {
	allIds := make([]int, 0)
	visited := make(map[int]bool)

	// 递归获取下级
	currentLevel := []int{adminId}
	for len(currentLevel) > 0 {
		var childIds []int
		err := u.container.GetDB().Model(&admin.Admin{}).
			Where("pid IN ?", currentLevel).
			Pluck("id", &childIds).Error

		if err != nil {
			return nil, err
		}

		if len(childIds) == 0 {
			break
		}

		nextLevel := make([]int, 0)
		for _, id := range childIds {
			if !visited[id] {
				allIds = append(allIds, id)
				visited[id] = true
				nextLevel = append(nextLevel, id)
			}
		}
		currentLevel = nextLevel
	}

	return allIds, nil
}

// CheckAdminInScope 检查管理员是否在指定管理员的管理范围内
func (u *AuthUtils) CheckAdminInScope(targetAdminId, currentAdminId int) (bool, error) {
	subordinateIds, err := u.GetAdminAllSubordinateIds(currentAdminId)
	if err != nil {
		return false, err
	}

	return datatype.ContainsInt(subordinateIds, targetAdminId), nil
}
