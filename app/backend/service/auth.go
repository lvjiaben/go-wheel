package service

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/lvjiaben/go-wheel/app/backend/model"
	"github.com/lvjiaben/go-wheel/pkg/container"
	"github.com/lvjiaben/go-wheel/pkg/jwt"
	"github.com/lvjiaben/go-wheel/pkg/utils"
	"go.uber.org/zap"
)

type AuthService struct {
	container *container.Container
	ruleCache map[int][]int     // 用户ID -> 规则ID缓存
	cacheMu   sync.RWMutex      // 缓存锁
	cacheExp  map[int]time.Time // 缓存过期时间
}

func NewAuthService(c *container.Container) *AuthService {
	return &AuthService{
		container: c,
		ruleCache: make(map[int][]int),
		cacheExp:  make(map[int]time.Time),
	}
}

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

func (a *AuthService) GetGroups(userId int) []model.AdminAuthGroup {
	var groups []model.AdminAuthGroup
	var groupAccess []model.AdminAuthGroupAccess
	if err := a.container.GetDB().Where("uid = ?", userId).Find(&groupAccess).Error; err != nil {
		a.container.GetLogger().Error("获取用户角色失败", zap.Error(err))
		return nil
	}

	var gids []int
	for _, access := range groupAccess {
		gids = append(gids, access.Gid)
	}

	if err := a.container.GetDB().Where("id in (?)", gids).Find(&groups).Error; err != nil {
		a.container.GetLogger().Error("获取角色信息失败", zap.Error(err))
		return nil
	}

	return groups
}

func (s *AuthService) GetGroupList(userId int) []model.CustomAdminAuthGroup {
	groups := s.GetGroups(userId)
	var groupIds []int
	for _, item := range groups {
		groupIds = append(groupIds, item.Id)
	}
	originGroupIds := groupIds
	for k, v := range groups {
		if v.Pid != 0 && utils.ContainsInt(originGroupIds, v.Pid) {
			groupIds = remove(groupIds, v.Id)
			groups = removeGroup(groups, k)
		}
	}
	var list []model.CustomAdminAuthGroup
	for _, v := range groups {
		if v.Rules == "*" {
			return model.GetAssocList(s.container)
		}
		var authGroup model.AdminAuthGroup
		s.container.GetDB().Where("id=?", v.Id).First(&authGroup)
		_, err := authGroup.GetChildrenIds(s.container, nil, true, -1)
		if err != nil {
			return nil
		}
		result := model.GetAssocList(s.container)
		if result == nil {
			newItem := model.CustomAdminAuthGroup{
				AdminAuthGroup: authGroup,
				Children:       make([]*model.CustomAdminAuthGroup, 0),
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
		s.container.GetDB().Table((&model.Admin{}).TableName()).Select("id").Pluck("id", &childrenIds)
	} else {
		groupIds := s.GetChildrenGroupIds(userId, false)
		if len(groupIds) == 0 {
			childrenIds = []int{}
		} else {
			var list []model.AdminAuthGroupAccess
			s.container.GetDB().Table((&model.AdminAuthGroupAccess{}).TableName()).Where("gid in (?)", groupIds).Find(&list)
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
		if v.Pid != 0 && utils.ContainsInt(originGroupIds, v.Pid) {
			groupIds = remove(groupIds, v.Id)
			groups = removeGroup(groups, k)
		}
	}
	var list []int
	for _, v := range groups {
		if v.Rules == "*" {
			var ids []int
			s.container.GetDB().Table((&model.AdminAuthGroup{}).TableName()).Select("id").Pluck("id", &ids)
			return ids
		}
		var authGroup model.AdminAuthGroup
		s.container.GetDB().Where("id=?", v.Id).First(&authGroup)
		ids, err := authGroup.GetChildrenIds(s.container, nil, true, -1)
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

// 从缓存获取规则，如果缓存不存在或已过期则查询数据库
func (a *AuthService) GetRuleIdsWithCache(userId int) []int {
	// 检查缓存
	a.cacheMu.RLock()
	if rules, ok := a.ruleCache[userId]; ok {
		// 检查是否过期
		if expiry, exists := a.cacheExp[userId]; exists && time.Now().Before(expiry) {
			a.cacheMu.RUnlock()
			return rules
		}
	}
	a.cacheMu.RUnlock()

	// 缓存不存在或已过期，从数据库查询
	rules := a.GetRuleIds(userId)

	// 更新缓存
	a.cacheMu.Lock()
	a.ruleCache[userId] = rules
	a.cacheExp[userId] = time.Now().Add(5 * time.Minute) // 5分钟缓存
	a.cacheMu.Unlock()

	return rules
}

// 清除用户规则缓存
func (a *AuthService) ClearUserRuleCache(userId int) {
	a.cacheMu.Lock()
	delete(a.ruleCache, userId)
	delete(a.cacheExp, userId)
	a.cacheMu.Unlock()
}

func (a *AuthService) GetRuleIds(userId int) []int {
	var user model.Admin
	if err := a.container.GetDB().First(&user, "id = ?", userId).Error; err != nil {
		a.container.GetLogger().Error("获取用户信息失败", zap.Error(err))
		return nil
	}

	if user.Username == "admin" {
		return []int{-1} // 特殊标记，表示超级管理员
	}

	var groupAccess model.AdminAuthGroupAccess
	if err := a.container.GetDB().First(&groupAccess, "uid = ?", userId).Error; err != nil {
		a.container.GetLogger().Error("获取用户角色失败", zap.Error(err))
		return nil
	}

	var group model.AdminAuthGroup
	if err := a.container.GetDB().First(&group, "id = ?", groupAccess.Gid).Error; err != nil {
		a.container.GetLogger().Error("获取角色信息失败", zap.Error(err))
		return nil
	}

	var ruleIds []int
	if err := json.Unmarshal([]byte(group.Rules), &ruleIds); err != nil {
		a.container.GetLogger().Error("解析角色权限失败", zap.Error(err))
		return nil
	}

	return ruleIds
}

func (s *AuthService) GetRuleList(userId int, check bool) map[int]string {
	ids := s.GetRuleIdsWithCache(userId) // 使用缓存版本
	if len(ids) == 0 {
		return map[int]string{}
	}
	if ids[0] == -1 && check { // 超级管理员
		return map[int]string{0: "*"}
	}
	var authRule []model.AdminAuthRule
	if ids[0] == -1 { // 超级管理员
		s.container.GetDB().Order("sort desc").Find(&authRule)
	} else {
		s.container.GetDB().Where("id IN (?)", ids).Order("sort desc").Find(&authRule)
	}
	list := make(map[int]string)
	for _, v := range authRule {
		list[v.Id] = v.Path
	}
	return list
}

func (s *AuthService) IsSuperAdmin(userId int) bool {
	ruleIds := s.GetRuleIdsWithCache(userId) // 使用缓存版本
	return len(ruleIds) > 0 && ruleIds[0] == -1
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

	for _, id := range list {
		if !idMap[id] {
			diff = append(diff, id)
		}
	}

	return diff
}

func (a *AuthService) GetRules(userId int) []model.AdminAuthRule {
	ruleIds := a.GetRuleIds(userId)
	if len(ruleIds) == 0 {
		return nil
	}

	var rules []model.AdminAuthRule
	query := a.container.GetDB().Table((&model.AdminAuthRule{}).TableName())
	if ruleIds[0] != -1 { // 不是超级管理员
		query = query.Where("id in (?)", ruleIds)
	}
	if err := query.Find(&rules).Error; err != nil {
		a.container.GetLogger().Error("获取权限列表失败", zap.Error(err))
		return nil
	}

	return rules
}

func (a *AuthService) CheckPermission(userId int, path string, method string) bool {
	ruleIds := a.GetRuleIds(userId)
	if len(ruleIds) == 0 {
		return false
	}

	if ruleIds[0] == -1 { // 超级管理员
		return true
	}

	var rule model.AdminAuthRule
	if err := a.container.GetDB().Where("path = ? AND method = ?", path, method).First(&rule).Error; err != nil {
		a.container.GetLogger().Error("获取权限信息失败", zap.Error(err))
		return false
	}

	for _, id := range ruleIds {
		if id == rule.Id {
			return true
		}
	}

	return false
}

func (a *AuthService) GenerateToken(userId int) (string, error) {
	// 生成token
	token, err := jwt.GenerateToken(userId, a.container.GetConfig().GetString("jwt.secret"), a.container.GetConfig().GetInt("jwt.expire_day"))
	if err != nil {
		return "", err
	}

	// 更新用户token
	if err := a.container.GetDB().Model(&model.Admin{}).Where("id = ?", userId).Update("token", token).Error; err != nil {
		return "", err
	}

	return token, nil
}

// RefreshToken 刷新用户Token
func (a *AuthService) RefreshToken(oldToken string) (string, error) {
	// 验证旧token
	claims, err := jwt.ParseToken(oldToken, a.container.GetConfig().GetString("jwt.secret"))
	if err != nil {
		return "", fmt.Errorf("无效的token: %v", err)
	}

	// 检查用户存在性
	var user model.Admin
	if err := a.container.GetDB().First(&user, "id = ?", claims.Id).Error; err != nil {
		return "", fmt.Errorf("用户不存在: %v", err)
	}

	// 检查token是否一致
	if user.Token != oldToken {
		return "", fmt.Errorf("令牌不匹配")
	}

	// 生成新token
	newToken, err := jwt.GenerateToken(claims.Id, a.container.GetConfig().GetString("jwt.secret"), a.container.GetConfig().GetInt("jwt.expire_day"))
	if err != nil {
		return "", fmt.Errorf("生成token失败: %v", err)
	}

	// 更新用户token
	if err := a.container.GetDB().Model(&model.Admin{}).Where("id = ?", claims.Id).Update("token", newToken).Error; err != nil {
		return "", fmt.Errorf("更新token失败: %v", err)
	}

	return newToken, nil
}
