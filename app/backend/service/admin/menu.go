package admin

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/pkg/utils/datatype"

	"github.com/lvjiaben/go-wheel/app/backend/model/admin"

	"github.com/lvjiaben/go-wheel/pkg/container"
	"gorm.io/gorm"
)

// VbenRoute Vben Admin路由配置结构
type VbenRoute struct {
	Name      string      `json:"name"`                // 路由名称（英文）
	Path      string      `json:"path"`                // 路由路径
	Component string      `json:"component,omitempty"` // 组件路径
	Redirect  string      `json:"redirect,omitempty"`  // 路由重定向
	Meta      VbenMeta    `json:"meta"`                // 元信息
	Children  []VbenRoute `json:"children,omitempty"`  // 子路由
}

// VbenMeta Vben Admin路由元信息结构
type VbenMeta struct {
	Title      string   `json:"title"`                // 菜单标题
	Icon       string   `json:"icon,omitempty"`       // 菜单图标
	Order      int      `json:"order,omitempty"`      // 排序
	AffixTab   bool     `json:"affixTab,omitempty"`   // 固定标签页
	HideInTab  bool     `json:"hideInTab,omitempty"`  // 隐藏标签页
	HideInMenu bool     `json:"hideInMenu,omitempty"` // 隐藏菜单
	KeepAlive  bool     `json:"keepAlive,omitempty"`  // keepAlive
	IframeSrc  string   `json:"iframeSrc,omitempty"`  // iframe地址
	Link       string   `json:"link,omitempty"`       // 外部链接
	Authority  []string `json:"authority,omitempty"`  // 权限控制
}

// MenuService 菜单服务
type MenuService struct {
	container *container.Container
}

// NewMenuService 创建菜单服务
func NewMenuService(c *container.Container) *MenuService {
	return &MenuService{container: c}
}

// GetAll 获取所有菜单
func (s *MenuService) GetAll() ([]map[string]interface{}, error) {
	var menus []admin.Menu
	if err := s.container.GetDB().Order("sort DESC").Find(&menus).Error; err != nil {
		return nil, fmt.Errorf("backend.menu.get_list_failed")
	}
	treeData := datatype.ToTreeAssocMap(menus, "id", "pid", "children")
	return treeData, nil
}

// 创建/更新
func (s *MenuService) Save(menu *admin.Menu) error {
	// 检查名称唯一性
	if err := s.checkNameUnique(menu); err != nil {
		return err
	}

	if menu.Id > 0 {
		if err := s.container.GetDB().Save(menu).Error; err != nil {
			return fmt.Errorf("backend.menu.update_failed")
		}
	} else {
		s.container.GetDB().Model(&admin.Menu{}).Where("pid = ?", menu.Pid).Select("COALESCE(MAX(sort), 0)").Scan(&menu.Sort)
		menu.Sort++
		if err := s.container.GetDB().Create(menu).Error; err != nil {
			return fmt.Errorf("backend.menu.create_failed")
		}
	}
	return nil
}

// Delete 删除菜单（递归删除所有子菜单）
func (s *MenuService) Delete(id int) error {
	// 开始事务删除
	return s.container.GetDB().Transaction(func(tx *gorm.DB) error {
		// 递归收集所有需要删除的菜单ID（包括自身和所有子菜单）
		allMenuIds, err := s.collectAllChildMenuIds(tx, id)
		if err != nil {
			return err
		}

		// 删除所有菜单的角色关联
		if err := tx.Where("menu_id IN ?", allMenuIds).Delete(&admin.RoleMenu{}).Error; err != nil {
			return fmt.Errorf("backend.menu.delete_role_menu_failed")
		}

		// 删除所有菜单
		if err := tx.Where("id IN ?", allMenuIds).Delete(&admin.Menu{}).Error; err != nil {
			return fmt.Errorf("backend.menu.delete_failed")
		}

		return nil
	})
}

// collectAllChildMenuIds 递归收集所有子菜单ID（包括自身）
func (s *MenuService) collectAllChildMenuIds(tx *gorm.DB, parentId int) ([]int, error) {
	var allIds []int

	// 添加当前菜单ID
	allIds = append(allIds, parentId)

	// 查找直接子菜单
	var childMenus []admin.Menu
	if err := tx.Where("pid = ?", parentId).Find(&childMenus).Error; err != nil {
		return nil, fmt.Errorf("backend.menu.get_children_failed")
	}

	// 递归收集每个子菜单的ID
	for _, child := range childMenus {
		childIds, err := s.collectAllChildMenuIds(tx, child.Id)
		if err != nil {
			return nil, err
		}
		allIds = append(allIds, childIds...)
	}

	return allIds, nil
}

// GetVbenRoutes 获取Vben Admin路由配置
// menuIds: 菜单ID列表，如果为空则获取所有菜单
// menuType: 菜单类型，如果为空或"all"则获取所有类型
func (s *MenuService) GetVbenRoutes(ctx *gin.Context, menuIds []int, menuType string) ([]VbenRoute, error) {
	// 检查容器是否为空
	if s.container == nil {
		return nil, fmt.Errorf("容器为空")
	}

	// 检查数据库连接是否为空
	db := s.container.GetDB()
	if db == nil {
		return nil, fmt.Errorf("数据库连接为空")
	}

	// 查询菜单
	var menus []admin.Menu
	query := db.Order("sort DESC")

	if menuType != "all" && menuType != "" {
		query = query.Where("type = ?", menuType)
	} else {
		query = query.Where("type != ?", "button")
	}

	if len(menuIds) > 0 {
		query = query.Where("id IN ?", menuIds)
	}
	if err := query.Find(&menus).Error; err != nil {
		return nil, fmt.Errorf("backend.menu.get_list_failed")
	}

	isChinese := ctx.Value("isCn").(bool)
	// 构建菜单树
	menuMap := make(map[int][]admin.Menu)
	for _, menu := range menus {
		menuMap[menu.Pid] = append(menuMap[menu.Pid], menu)
	}

	// 构建Vben路由配置
	var routes []VbenRoute
	if rootMenus, ok := menuMap[0]; ok {
		for _, rootMenu := range rootMenus {
			route := s.convertMenuToVbenRoute(rootMenu, menuMap, isChinese)
			routes = append(routes, route)
		}
	}

	return routes, nil
}

// 将菜单转换为Vben路由配置
func (s *MenuService) convertMenuToVbenRoute(menu admin.Menu, menuMap map[int][]admin.Menu, isChinese bool) VbenRoute {
	// 确定显示的标题（根据语言）
	title := menu.Name
	if !isChinese && menu.Enname != "" {
		title = menu.Enname
	}

	// 创建基本Vben路由配置
	route := VbenRoute{
		Name: menu.Enname, // 路由名称使用英文名
		Path: menu.Path,
		Meta: VbenMeta{
			Title: title, // 标题根据语言选择
			Icon:  menu.Icon,
			// 不返回Order字段，后端已按sort DESC排序
		},
	}

	// 根据菜单类型设置组件和特殊属性
	switch menu.Type {
	case "menu":
		if menu.Pid == 0 {
			// 顶级菜单，使用BasicLayout
			route.Component = "BasicLayout"
		} else {
			// 子菜单，使用具体的组件路径
			route.Component = menu.Component
		}
	case "iframe":
		route.Meta.IframeSrc = menu.Iframe
		route.Component = "IFrameView"
	case "link":
		route.Meta.Link = menu.External
	}

	route.Meta.KeepAlive = true

	// 处理显示控制
	if menu.Visible == 0 {
		route.Meta.HideInMenu = true
	}

	// 处理标签页控制
	route.Meta.HideInTab = menu.ShowTag == 1

	// 处理固定标签页
	if menu.FixedTag == 1 {
		route.Meta.AffixTab = true
	}

	// 处理子菜单
	if children, ok := menuMap[menu.Id]; ok && len(children) > 0 {
		// 递归转换所有子菜单
		for _, child := range children {
			childRoute := s.convertMenuToVbenRoute(child, menuMap, isChinese)
			route.Children = append(route.Children, childRoute)
		}

		// 如果有子菜单且没有设置重定向，设置重定向到第一个子菜单
		if route.Redirect == "" && len(route.Children) > 0 {
			route.Redirect = route.Children[0].Path
		}
	}

	return route
}

// checkNameUnique 检查菜单名称唯一性
func (s *MenuService) checkNameUnique(menu *admin.Menu) error {
	var count int64

	// 检查中文名称唯一性
	query := s.container.GetDB().Model(&admin.Menu{}).Where("name = ?", menu.Name)
	if menu.Id > 0 {
		query = query.Where("id != ?", menu.Id)
	}
	if err := query.Count(&count).Error; err != nil {
		return fmt.Errorf("backend.menu.check_name_failed")
	}
	if count > 0 {
		return fmt.Errorf("backend.menu.name_exists")
	}

	// 检查英文名称唯一性
	query = s.container.GetDB().Model(&admin.Menu{}).Where("enname = ?", menu.Enname)
	if menu.Id > 0 {
		query = query.Where("id != ?", menu.Id)
	}
	if err := query.Count(&count).Error; err != nil {
		return fmt.Errorf("backend.menu.check_enname_failed")
	}
	if count > 0 {
		return fmt.Errorf("backend.menu.enname_exists")
	}

	return nil
}
