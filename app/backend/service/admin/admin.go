package admin

import (
	"fmt"
	"time"

	"github.com/lvjiaben/go-wheel/pkg/utils/crypto"
	"github.com/lvjiaben/go-wheel/pkg/utils/datatype"

	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/app/backend/model/admin"
	"github.com/lvjiaben/go-wheel/app/backend/utils"
	"github.com/lvjiaben/go-wheel/pkg/container"
	"gorm.io/gorm"
)

// AdminService 管理员服务
type AdminService struct {
	container *container.Container
	authUtils *utils.AuthUtils
}

// NewAdminService 创建管理员服务
func NewAdminService(c *container.Container) *AdminService {
	return &AdminService{
		container: c,
		authUtils: utils.NewAuthUtils(c),
	}
}

// GetAll 获取所有管理员
func (s *AdminService) GetAll(ctx *gin.Context) ([]admin.AdminWithRoles, error) {
	currentAdminId := ctx.GetInt("admin_id")

	// 检查是否为超级管理员
	isSuper, err := s.authUtils.IsAdminSuper(currentAdminId)
	if err != nil {
		return nil, fmt.Errorf("backend.admin.check_super_failed")
	}

	var admins []admin.Admin

	// 如果是超级管理员，获取所有管理员（除了自己）
	if isSuper {
		if err := s.container.GetDB().Where("id != ?", currentAdminId).Order("id ASC").Find(&admins).Error; err != nil {
			return nil, fmt.Errorf("backend.admin.get_list_failed")
		}
	} else {
		// 获取当前管理员创建的所有管理员（包括无限级下级）
		adminIds, err := s.authUtils.GetAdminAllSubordinateIds(currentAdminId)
		if err != nil {
			return nil, err
		}

		if len(adminIds) == 0 {
			return []admin.AdminWithRoles{}, nil
		}

		if err := s.container.GetDB().Where("id IN ?", adminIds).Order("id ASC").Find(&admins).Error; err != nil {
			return nil, fmt.Errorf("backend.admin.get_list_failed")
		}
	}

	// 转换为 AdminWithRoles 并添加角色ID
	result := make([]admin.AdminWithRoles, 0, len(admins))
	for _, adminItem := range admins {
		roleIds, _ := s.authUtils.GetAdminDirectRoleIds(adminItem.Id)

		adminWithRoles := admin.AdminWithRoles{
			Id:            adminItem.Id,
			Pid:           adminItem.Pid,
			Username:      adminItem.Username,
			Realname:      adminItem.Realname,
			Avatar:        adminItem.Avatar,
			Email:         adminItem.Email,
			Mobile:        adminItem.Mobile,
			Status:        adminItem.Status,
			LastLoginTime: adminItem.LastLoginTime,
			CreatedAt:     adminItem.CreatedAt,
			UpdatedAt:     adminItem.UpdatedAt,
			RoleIds:       roleIds,
		}
		result = append(result, adminWithRoles)
	}

	return result, nil
}

// GetById 根据ID获取管理员
func (s *AdminService) GetById(id int) (*admin.Admin, error) {
	var adminItem admin.Admin
	if err := s.container.GetDB().First(&adminItem, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("backend.admin.not_found")
		}
		return nil, fmt.Errorf("backend.admin.get_detail_failed")
	}
	return &adminItem, nil
}

// Save 创建或更新管理员
func (s *AdminService) Save(adminItem *admin.Admin, roleIds []int, ctx *gin.Context) error {
	currentAdminId := ctx.GetInt("admin_id")

	// 检查是否为超级管理员
	isSuper, err := s.authUtils.IsAdminSuper(currentAdminId)
	if err != nil {
		return fmt.Errorf("backend.admin.check_super_failed")
	}

	// 如果不是超级管理员，需要进行权限检查
	if !isSuper {
		// 创建时，设置父级为当前管理员
		if adminItem.Id == 0 {
			adminItem.Pid = currentAdminId
		} else {
			// 更新时，检查要修改的管理员是否在管理范围内
			if exists, err := s.authUtils.CheckAdminInScope(adminItem.Id, currentAdminId); err != nil {
				return fmt.Errorf("backend.admin.check_scope_failed")
			} else if !exists {
				return fmt.Errorf("backend.admin.not_in_scope")
			}

			// 不能修改自己
			if adminItem.Id == currentAdminId {
				return fmt.Errorf("backend.admin.cannot_modify_self")
			}
		}

		// 检查分配的角色是否在当前管理员的权限范围内
		currentAdminRoleIds, err := s.authUtils.GetAdminAllRoleIds(currentAdminId)
		if err != nil {
			return fmt.Errorf("backend.admin.get_roles_failed")
		}

		for _, roleId := range roleIds {
			if !datatype.ContainsInt(currentAdminRoleIds, roleId) {
				return fmt.Errorf("backend.admin.role_not_in_scope")
			}
		}
	}

	// 检查用户名是否重复
	if err := s.checkUsernameExists(adminItem.Username, adminItem.Id); err != nil {
		return err
	}

	// 检查邮箱是否重复
	if adminItem.Email != "" {
		if err := s.checkEmailExists(adminItem.Email, adminItem.Id); err != nil {
			return err
		}
	}

	// 设置时间戳
	now := int(time.Now().Unix())

	// 开始事务
	return s.container.GetDB().Transaction(func(tx *gorm.DB) error {
		if adminItem.Id > 0 {
			adminItem.UpdatedAt = now
			// 更新管理员（不更新密码）
			updateData := map[string]interface{}{
				"pid":        adminItem.Pid,
				"username":   adminItem.Username,
				"realname":   adminItem.Realname,
				"avatar":     adminItem.Avatar,
				"email":      adminItem.Email,
				"mobile":     adminItem.Mobile,
				"status":     adminItem.Status,
				"updated_at": adminItem.UpdatedAt,
			}

			// 如果密码不为空，则更新密码和盐值
			if adminItem.Password != "" {
				// 生成新的盐值
				salt, err := crypto.GenerateSalt()
				if err != nil {
					return fmt.Errorf("backend.admin.generate_salt_failed")
				}

				// 使用盐值加密密码
				hashedPassword, _ := crypto.PasswordHashWithSalt(adminItem.Password, salt)

				updateData["password"] = hashedPassword
				updateData["salt"] = salt
			}

			// 使用 Select 明确指定要更新的字段，包括零值字段
			fields := []string{"pid", "username", "realname", "avatar", "email", "mobile", "status", "updated_at"}
			if adminItem.Password != "" {
				fields = append(fields, "password", "salt")
			}

			if err := tx.Model(&admin.Admin{}).Where("id = ?", adminItem.Id).Select(fields).Updates(updateData).Error; err != nil {
				return fmt.Errorf("backend.admin.update_failed")
			}
		} else {
			adminItem.CreatedAt = now
			adminItem.UpdatedAt = now

			// 生成盐值
			salt, err := crypto.GenerateSalt()
			if err != nil {
				return fmt.Errorf("backend.admin.generate_salt_failed")
			}
			adminItem.Salt = salt

			// 使用盐值加密密码
			hashedPassword, _ := crypto.PasswordHashWithSalt(adminItem.Password, salt)
			adminItem.Password = hashedPassword

			// 创建管理员
			if err := tx.Create(adminItem).Error; err != nil {
				return fmt.Errorf("backend.admin.create_failed")
			}
		}

		// 更新角色关联
		if err := s.updateAdminRoles(adminItem.Id, roleIds, tx); err != nil {
			return err
		}

		return nil
	})
}

// Delete 删除管理员
func (s *AdminService) Delete(id int, ctx *gin.Context) error {
	currentAdminId := ctx.GetInt("admin_id")

	// 不能删除自己
	if id == currentAdminId {
		return fmt.Errorf("backend.admin.cannot_delete_self")
	}

	// 检查是否为超级管理员
	isSuper, err := s.authUtils.IsAdminSuper(currentAdminId)
	if err != nil {
		return fmt.Errorf("backend.admin.check_super_failed")
	}

	// 如果不是超级管理员，需要进行权限检查
	if !isSuper {
		// 检查要删除的管理员是否在管理范围内
		if exists, err := s.authUtils.CheckAdminInScope(id, currentAdminId); err != nil {
			return fmt.Errorf("backend.admin.check_scope_failed")
		} else if !exists {
			return fmt.Errorf("backend.admin.not_in_scope")
		}
	}

	// 检查是否有下级管理员
	if hasSubordinates, err := s.hasSubordinateAdmins(id); err != nil {
		return fmt.Errorf("backend.admin.check_subordinates_failed")
	} else if hasSubordinates {
		return fmt.Errorf("backend.admin.has_subordinates")
	}

	// 开始事务删除
	return s.container.GetDB().Transaction(func(tx *gorm.DB) error {
		// 删除管理员角色关联
		if err := tx.Where("admin_id = ?", id).Delete(&admin.RoleAdmin{}).Error; err != nil {
			return fmt.Errorf("backend.admin.delete_roles_failed")
		}

		// 删除管理员
		if err := tx.Delete(&admin.Admin{}, id).Error; err != nil {
			return fmt.Errorf("backend.admin.delete_failed")
		}

		return nil
	})
}

// ===== 私有辅助方法 =====

// checkUsernameExists 检查用户名是否重复
func (s *AdminService) checkUsernameExists(username string, excludeId int) error {
	var count int64
	query := s.container.GetDB().Model(&admin.Admin{}).Where("username = ?", username)
	if excludeId > 0 {
		query = query.Where("id != ?", excludeId)
	}

	if err := query.Count(&count).Error; err != nil {
		return fmt.Errorf("backend.admin.check_username_failed")
	}

	if count > 0 {
		return fmt.Errorf("backend.admin.username_exists")
	}

	return nil
}

// checkEmailExists 检查邮箱是否重复
func (s *AdminService) checkEmailExists(email string, excludeId int) error {
	var count int64
	query := s.container.GetDB().Model(&admin.Admin{}).Where("email = ?", email)
	if excludeId > 0 {
		query = query.Where("id != ?", excludeId)
	}

	if err := query.Count(&count).Error; err != nil {
		return fmt.Errorf("backend.admin.check_email_failed")
	}

	if count > 0 {
		return fmt.Errorf("backend.admin.email_exists")
	}

	return nil
}

// hasSubordinateAdmins 检查是否有下级管理员
func (s *AdminService) hasSubordinateAdmins(adminId int) (bool, error) {
	var count int64
	err := s.container.GetDB().Model(&admin.Admin{}).Where("pid = ?", adminId).Count(&count).Error
	return count > 0, err
}

// updateAdminRoles 更新管理员角色关联
func (s *AdminService) updateAdminRoles(adminId int, roleIds []int, tx *gorm.DB) error {
	// 删除现有角色关联
	if err := tx.Where("admin_id = ?", adminId).Delete(&admin.RoleAdmin{}).Error; err != nil {
		return fmt.Errorf("backend.admin.delete_old_roles_failed")
	}

	// 添加新的角色关联
	if len(roleIds) > 0 {
		now := int(time.Now().Unix())
		roleAdmins := make([]admin.RoleAdmin, 0, len(roleIds))
		for _, roleId := range roleIds {
			roleAdmins = append(roleAdmins, admin.RoleAdmin{
				AdminId:   adminId,
				RoleId:    roleId,
				CreatedAt: now,
			})
		}

		if err := tx.Create(&roleAdmins).Error; err != nil {
			return fmt.Errorf("backend.admin.assign_roles_failed")
		}
	}

	return nil
}
