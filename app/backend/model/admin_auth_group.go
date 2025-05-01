package model

import "github.com/lvjiaben/go-wheel/frame/global"

type AdminAuthGroup struct {
	Id        int    `json:"id" gorm:"column:id"`
	Pid       int    `json:"pid" gorm:"column:pid"`               // 父亲
	Name      string `json:"name" gorm:"column:name"`             // 名称
	Rules     string `json:"rules" gorm:"column:rules"`           // 规则
	CreatedAt int    `json:"created_at" gorm:"column:created_at"` // 创建时间
	UpdatedAt int    `json:"updated_at" gorm:"column:updated_at"` // 更新时间
}

type CustomAdminAuthGroup struct {
	AdminAuthGroup
	Children []CustomAdminAuthGroup `json:"children,omitempty"`
}

func (e *AdminAuthGroup) TableName() string {
	return "admin_auth_group"
}

func GetAssocList(query ...interface{}) []CustomAdminAuthGroup {
	var list []AdminAuthGroup
	if len(query) > 0 {
		if err := global.DB.Find(&list, query...).Error; err != nil {
			return nil
		}
	} else {
		if err := global.DB.Find(&list).Error; err != nil {
			return nil
		}
	}
	var rootList []CustomAdminAuthGroup
	for _, item := range list {
		if item.Pid == 0 {
			newItem := CustomAdminAuthGroup{
				AdminAuthGroup: item,
				Children:       []CustomAdminAuthGroup{},
			}
			loadChildren(&newItem, &list)
			rootList = append(rootList, newItem)
		}
	}
	return rootList
}

func loadChildren(parent *CustomAdminAuthGroup, categories *[]AdminAuthGroup) {
	var children []CustomAdminAuthGroup
	for i := 0; i < len(*categories); {
		if (*categories)[i].Pid == parent.Id {
			child := (*categories)[i]
			childItem := CustomAdminAuthGroup{
				AdminAuthGroup: child,
				Children:       []CustomAdminAuthGroup{},
			}
			loadChildren(&childItem, categories)
			children = append(children, childItem)
			*categories = append((*categories)[:i], (*categories)[i+1:]...)
		} else {
			i++
		}
	}
	parent.Children = children
}

func (e *AdminAuthGroup) GetChildrenIds(parentId *int, includeParentId bool, limitLevel int) ([]int, error) {
	var ids []int
	err := e.loadChildrenIds(parentId, includeParentId, limitLevel, 0, &ids)
	if err != nil {
		return nil, err
	}
	return ids, nil
}

func (e *AdminAuthGroup) loadChildrenIds(parentId *int, includeParentId bool, limitLevel int, currentLevel int, ids *[]int) error {
	if limitLevel != -1 && currentLevel > limitLevel {
		return nil
	}

	query := global.DB
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
		err := e.loadChildrenIds(&child.Id, includeParentId, limitLevel, currentLevel+1, ids)
		if err != nil {
			return err
		}
	}

	return nil
}
