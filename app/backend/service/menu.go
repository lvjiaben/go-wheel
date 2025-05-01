package service

import (
	"github.com/lvjiaben/go-wheel/app/backend/model"
	"github.com/lvjiaben/go-wheel/pkg/actions"
	"github.com/lvjiaben/go-wheel/pkg/global"
)

type MenuService struct{}

type Menu struct {
	Component string  `json:"component"`
	Meta      Meta    `json:"meta"`
	Name      string  `json:"name"`
	Path      string  `json:"path"`
	Children  []*Menu `json:"children,omitempty"`
	Id        int     `json:"id"`
	Pid       int     `json:"pid"`
}

type Meta struct {
	Order      int    `json:"order,omitempty"`
	Title      string `json:"title"`
	HideInMenu bool   `json:"hideInMenu,omitempty"`
	Icon       string `json:"icon,omitempty"`
}

func GetMenusFromDB(userId int, lang string) []*Menu {
	service := AuthService{}
	ruleIds := service.GetRuleIds(userId)
	isSuper := false
	for _, id := range ruleIds {
		if id == "*" {
			isSuper = true
			break
		}
	}
	var authRules []model.AdminAuthRule
	query := global.DB.Table((&model.AdminAuthRule{}).TableName()).Where("status = 1 AND type != ?", "button")
	if !isSuper {
		query = query.Where("id in (?)", ruleIds)
	}
	query.Find(&authRules)
	// 转换 authRules 到 menus
	menus := make([]Menu, len(authRules))
	for i, rule := range authRules {
		var title string
		if lang == "zh" {
			title = rule.Name
		} else {
			title = rule.Enname
		}
		var component = "BasicLayout"
		if rule.Pid != 0 {
			component = rule.Path
		}
		menus[i] = Menu{
			Id:        rule.Id,
			Pid:       rule.Pid,
			Name:      title,
			Path:      rule.Path,
			Component: component,
			Meta: Meta{
				Title:      title,
				Order:      rule.Sort,
				Icon:       rule.Icon,
				HideInMenu: actions.IntToBool(rule.Hide),
			},
		}
	}
	return BuildTree(menus)
}

func BuildTree(tmpArr []Menu) []*Menu {
	result := make([]*Menu, 0)
	menuMap := make(map[int]*Menu)

	for i := range tmpArr {
		menuMap[tmpArr[i].Id] = &tmpArr[i]
	}

	// 遍历映射并构建树
	for _, menu := range menuMap {
		if parent, exists := menuMap[menu.Pid]; exists {
			parent.Children = append(parent.Children, menu)
		} else {
			result = append(result, menu)
		}
	}

	return result
}
