package model

import (
	"encoding/json"
	"time"

	"github.com/lvjiaben/go-wheel/pkg/container"
	"go.uber.org/zap"
)

// AdminAuthGroup 角色表
type AdminAuthGroup struct {
	Id        int    `json:"id" gorm:"primaryKey"`             // 主键ID
	Pid       int    `json:"pid" gorm:"not null"`              // 上级角色ID
	Name      string `json:"name" gorm:"not null"`             // 角色名称
	Rules     string `json:"rules" gorm:"not null"`            // 角色权限
	CreatedAt int    `json:"created_at" gorm:"autoCreateTime"` // 创建时间
	UpdatedAt int    `json:"updated_at" gorm:"autoUpdateTime"` // 更新时间
}

// CustomAdminAuthGroup 自定义角色表
type CustomAdminAuthGroup struct {
	AdminAuthGroup
	Children []*CustomAdminAuthGroup `json:"children"` // 子角色列表
}

func (AdminAuthGroup) TableName() string {
	return "admin_auth_group"
}

// GetRules 获取角色权限
func (g *AdminAuthGroup) GetRules() ([]int, error) {
	var rules []int
	if err := json.Unmarshal([]byte(g.Rules), &rules); err != nil {
		return nil, err
	}
	return rules, nil
}

// SetRules 设置角色权限
func (g *AdminAuthGroup) SetRules(rules []int) error {
	data, err := json.Marshal(rules)
	if err != nil {
		return err
	}
	g.Rules = string(data)
	return nil
}

// BeforeCreate 创建前的钩子
func (g *AdminAuthGroup) BeforeCreate() error {
	g.CreatedAt = int(time.Now().Unix())
	g.UpdatedAt = int(time.Now().Unix())
	return nil
}

// BeforeUpdate 更新前的钩子
func (g *AdminAuthGroup) BeforeUpdate() error {
	g.UpdatedAt = int(time.Now().Unix())
	return nil
}

// GetAssocList 获取关联列表
func GetAssocList(c *container.Container) []CustomAdminAuthGroup {
	var groups []AdminAuthGroup
	if err := c.GetDB().Find(&groups).Error; err != nil {
		c.GetLogger().Error("获取角色列表失败", zap.Error(err))
		return nil
	}

	// 构建树形结构
	groupMap := make(map[int]*CustomAdminAuthGroup)
	var result []CustomAdminAuthGroup

	// 初始化map
	for _, group := range groups {
		customGroup := CustomAdminAuthGroup{
			AdminAuthGroup: group,
			Children:       make([]*CustomAdminAuthGroup, 0),
		}
		groupMap[group.Id] = &customGroup
		if group.Pid == 0 {
			result = append(result, customGroup)
		}
	}

	// 建立父子关系
	for _, group := range groups {
		if group.Pid != 0 {
			if parent, ok := groupMap[group.Pid]; ok {
				customGroup := groupMap[group.Id]
				parent.Children = append(parent.Children, customGroup)
			}
		}
	}

	return result
}

func (e *AdminAuthGroup) GetChildrenIds(c *container.Container, parentId *int, includeParentId bool, limitLevel int) ([]int, error) {
	var ids []int
	err := e.loadChildrenIds(c, parentId, includeParentId, limitLevel, 0, &ids)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func (e *AdminAuthGroup) loadChildrenIds(c *container.Container, parentId *int, includeParentId bool, limitLevel int, currentLevel int, ids *[]int) error {
	if limitLevel != -1 && currentLevel > limitLevel {
		return nil
	}

	query := c.GetDB()
	if parentId != nil {
		query = query.Where("pid = ?", *parentId)
	} else {
		query = query.Where("pid IS NULL")
	}

	var children []AdminAuthGroup
	if err := query.Find(&children).Error; err != nil {
		return err
	}

	if includeParentId && parentId != nil {
		*ids = append(*ids, *parentId)
	}

	for _, child := range children {
		*ids = append(*ids, child.Id)
		err := e.loadChildrenIds(c, &child.Id, includeParentId, limitLevel, currentLevel+1, ids)
		if err != nil {
			return err
		}
	}

	return nil
}
