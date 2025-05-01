package service

import (
	"strings"

	"github.com/lvjiaben/go-wheel/app/backend/model"
	"github.com/lvjiaben/go-wheel/pkg/actions"
	"github.com/lvjiaben/go-wheel/pkg/global"
)

type AuthService struct{}

type Group struct {
	List       []model.AdminAuthGroup
	AccessList []model.AdminAuthGroupAccess
}

func (s *AuthService) Check(rule string, userId int) bool {
	ruleList := s.GetRuleList(userId, true)
	for _, v := range ruleList {
		if v == "*" {
			return true
		} else if v == rule {
			return true
		}
	}
	return false
}

func (s *AuthService) GetGroups(userId int) []model.AdminAuthGroup {
	var data []model.AdminAuthGroup
	global.DB.Table("admin_auth_group").Where("id IN (SELECT gid FROM admin_auth_group_access WHERE uid = ?)", userId).Find(&data)
	return data
}

func (s *AuthService) GetGroupList(userId int) []model.CustomAdminAuthGroup {
	groups := s.GetGroups(userId)
	var groupIds []int
	for _, item := range groups {
		groupIds = append(groupIds, item.Id)
	}
	originGroupIds := groupIds
	for k, v := range groups {
		if actions.IsIntInSlice(v.Pid, originGroupIds) {
			groupIds = remove(groupIds, v.Id)
			groups = removeGroup(groups, k)
		}
	}
	var list []model.CustomAdminAuthGroup
	for _, v := range groups {
		if v.Rules == "*" {
			return model.GetAssocList()
		}
		var authGroup model.AdminAuthGroup
		global.DB.Where("id=?", v.Id).First(&authGroup)
		ids, err := authGroup.GetChildrenIds(nil, true, -1)
		if err != nil {
			return nil
		}
		result := model.GetAssocList("id IN (?)", ids)
		if result == nil {
			newItem := model.CustomAdminAuthGroup{
				AdminAuthGroup: authGroup,
				Children:       []model.CustomAdminAuthGroup{},
			}
			list = append(list, newItem)
		} else {
			list = append(list, result[0])
		}
	}
	return list
}

func (s *AuthService) GetChildrenAdminIds(userId int, withSelf bool) []int {
	var childrenIds []int
	if s.IsSuperAdmin(userId) {
		global.DB.Table((&model.Admin{}).TableName()).Select("id").Pluck("id", &childrenIds)
	} else {
		groupIds := s.GetChildrenGroupIds(userId, false)
		if len(groupIds) == 0 {
			childrenIds = []int{}
		} else {
			var list []model.AdminAuthGroupAccess
			global.DB.Table((&model.AdminAuthGroupAccess{}).TableName()).Where("gid in (?)", groupIds).Find(&list)
			for _, v := range list {
				childrenIds = append(childrenIds, v.Uid)
			}
		}
	}
	if withSelf {
		isIn := false
		for _, v := range childrenIds {
			if v == userId {
				isIn = true
			}
		}
		if !isIn {
			childrenIds = append(childrenIds, userId)
		}
	} else {
		childrenIds = remove(childrenIds, userId)
	}
	return childrenIds
}

func (s *AuthService) GetChildrenGroupIds(userId int, withSelf bool) []int {
	groups := s.GetGroups(userId)
	var groupIds []int
	for _, item := range groups {
		groupIds = append(groupIds, item.Id)
	}
	originGroupIds := groupIds
	for k, v := range groups {
		if actions.IsIntInSlice(v.Pid, originGroupIds) {
			groupIds = remove(groupIds, v.Id)
			groups = removeGroup(groups, k)
		}
	}
	var list []int
	for _, v := range groups {
		if v.Rules == "*" {
			var ids []int
			global.DB.Table((&model.AdminAuthGroup{}).TableName()).Select("id").Pluck("id", &ids)
			return ids
		}
		var authGroup model.AdminAuthGroup
		global.DB.Where("id=?", v.Id).First(&authGroup)
		ids, err := authGroup.GetChildrenIds(nil, true, -1)
		if err != nil {
			continue
		}
		list = append(list, ids...)
	}
	if !withSelf {
		list = arrayDiff(list, groupIds)
	}
	return list
}

func (s *AuthService) GetRuleIds(userId int) []string {
	groups := s.GetGroups(userId)
	var ids []string
	for _, item := range groups {
		ids = append(ids, strings.Split(item.Rules, ",")...)
	}
	return stringArrayUnique(ids)
}

func (s *AuthService) GetRuleList(userId int, check bool) map[int]string {
	ids := s.GetRuleIds(userId)
	if len(ids) == 0 {
		return map[int]string{}
	}
	if contains(ids, "*") && check {
		return map[int]string{0: "*"}
	}
	var authRule []model.AdminAuthRule
	if contains(ids, "*") {
		global.DB.Order("sort desc").Find(&authRule)
	} else {
		global.DB.Where("id IN (?)", ids).Order("sort desc").Find(&authRule)
	}
	list := make(map[int]string)
	for _, v := range authRule {
		list[v.Id] = v.Route
	}
	return list
}

func (s *AuthService) IsSuperAdmin(userId int) bool {
	ruleIds := s.GetRuleIds(userId)
	for _, v := range ruleIds {
		if v == "*" {
			return true
		}
	}
	return false
}

func contains(ids []string, value string) bool {
	for _, id := range ids {
		if id == value {
			return true
		}
	}
	return false
}

func stringArrayUnique(input []string) []string {
	uniqueMap := make(map[string]bool)
	var uniqueSlice []string
	for _, val := range input {
		if _, ok := uniqueMap[val]; !ok {
			uniqueMap[val] = true
			uniqueSlice = append(uniqueSlice, val)
		}
	}
	return uniqueSlice
}

func remove(slice []int, val int) []int {
	for i, v := range slice {
		if v == val {
			return append(slice[:i], slice[i+1:]...)
		}
	}
	return slice
}

func removeGroup(groups []model.AdminAuthGroup, index int) []model.AdminAuthGroup {
	return append(groups[:index], groups[index+1:]...)
}

func arrayDiff(list []int, groupIds []int) []int {
	diff := make([]int, 0)
	idMap := make(map[int]bool)

	for _, id := range groupIds {
		idMap[id] = true
	}

	for _, item := range list {
		if _, found := idMap[item]; !found {
			diff = append(diff, item)
		}
	}

	return diff
}
