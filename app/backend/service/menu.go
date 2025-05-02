package service

import (
	"encoding/json"

	"github.com/lvjiaben/go-wheel/app/backend/model"
	"github.com/lvjiaben/go-wheel/pkg/container"
	"go.uber.org/zap"
)

type Menu struct {
	Id       int     `json:"id"`
	Pid      int     `json:"pid"`
	Name     string  `json:"name"`
	Path     string  `json:"path"`
	Icon     string  `json:"icon"`
	Type     string  `json:"type"`
	Children []*Menu `json:"children"`
}

type MenuService struct {
	container   *container.Container
	authService *AuthService
}

func NewMenuService(c *container.Container) *MenuService {
	return &MenuService{
		container:   c,
		authService: NewAuthService(c),
	}
}

func (m *MenuService) GetMenuList() ([]model.AdminAuthRule, error) {
	var rules []model.AdminAuthRule
	query := m.container.GetDB().Table((&model.AdminAuthRule{}).TableName()).Where("status = 1 AND type != ?", "button")
	if err := query.Find(&rules).Error; err != nil {
		m.container.GetLogger().Error("获取菜单列表失败", zap.Error(err))
		return nil, err
	}
	return rules, nil
}

func (m *MenuService) GetMenusFromDB(userId int, lang string) []*Menu {
	var user model.Admin
	if err := m.container.GetDB().First(&user, "id = ?", userId).Error; err != nil {
		m.container.GetLogger().Error("获取用户信息失败", zap.Error(err))
		return nil
	}

	isSuper := m.authService.IsSuperAdmin(userId)
	var ruleIds []int
	if !isSuper {
		var groupAccess model.AdminAuthGroupAccess
		if err := m.container.GetDB().First(&groupAccess, "uid = ?", userId).Error; err != nil {
			m.container.GetLogger().Error("获取用户角色失败", zap.Error(err))
			return nil
		}

		var group model.AdminAuthGroup
		if err := m.container.GetDB().First(&group, "id = ?", groupAccess.Gid).Error; err != nil {
			m.container.GetLogger().Error("获取角色信息失败", zap.Error(err))
			return nil
		}

		if err := json.Unmarshal([]byte(group.Rules), &ruleIds); err != nil {
			m.container.GetLogger().Error("解析角色权限失败", zap.Error(err))
			return nil
		}
	}

	var authRules []model.AdminAuthRule
	query := m.container.GetDB().Table((&model.AdminAuthRule{}).TableName()).Where("status = 1 AND type != ?", "button")
	if !isSuper {
		query = query.Where("id in (?)", ruleIds)
	}
	if err := query.Find(&authRules).Error; err != nil {
		m.container.GetLogger().Error("获取菜单列表失败", zap.Error(err))
		return nil
	}

	var menus []*Menu
	for _, rule := range authRules {
		if lang != "zh-CN" {
			rule.Name = rule.Enname
		}
		if rule.Pid == 0 {
			menu := &Menu{
				Id:       rule.Id,
				Pid:      rule.Pid,
				Name:     rule.Name,
				Path:     rule.Path,
				Icon:     rule.Icon,
				Type:     rule.Type,
				Children: make([]*Menu, 0),
			}
			menus = append(menus, menu)
		}
	}

	for _, menu := range menus {
		for _, rule := range authRules {
			if lang != "zh-CN" {
				rule.Name = rule.Enname
			}
			if rule.Pid == menu.Id {
				child := &Menu{
					Id:       rule.Id,
					Pid:      rule.Pid,
					Name:     rule.Name,
					Path:     rule.Path,
					Icon:     rule.Icon,
					Type:     rule.Type,
					Children: make([]*Menu, 0),
				}
				menu.Children = append(menu.Children, child)
			}
		}
	}

	return menus
}
