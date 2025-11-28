package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/pkg/container"
	"github.com/lvjiaben/go-wheel/pkg/utils/http"
)

// MenuController 菜单控制器
type MenuController struct {
	container *container.Container
}

// NewMenuController 创建菜单控制器
func NewMenuController(c *container.Container) *MenuController {
	return &MenuController{
		container: c,
	}
}

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
	Title      string `json:"title"`                // 菜单标题
	Icon       string `json:"icon,omitempty"`       // 菜单图标
	Order      int    `json:"order,omitempty"`      // 排序
	AffixTab   bool   `json:"affixTab,omitempty"`   // 固定标签页
	HideInTab  bool   `json:"hideInTab,omitempty"`  // 隐藏标签页
	HideInMenu bool   `json:"hideInMenu,omitempty"` // 隐藏菜单
	KeepAlive  bool   `json:"keepAlive,omitempty"`  // keepAlive
}

// All 获取所有菜单（返回固定JSON）
func (c *MenuController) All(ctx *gin.Context) {
	// 根据语言返回对应的标题
	isCn := ctx.GetBool("isCn")

	// 标题映射
	homeTitle := "Home"
	userCenterTitle := "User Center"
	settingsTitle := "Settings"
	if isCn {
		homeTitle = "首页"
		userCenterTitle = "个人中心"
		settingsTitle = "账号设置"
	}

	// 固定菜单配置
	menus := []VbenRoute{
		{
			Name:      "Home",
			Path:      "/home",
			Component: "BasicLayout",
			Redirect:  "/home/index",
			Meta: VbenMeta{
				Title: homeTitle,
				Icon:  "lucide:house",
				Order: 1,
			},
			Children: []VbenRoute{
				{
					Name:      "HomeIndex",
					Path:      "/home/index",
					Component: "/home/index",
					Meta: VbenMeta{
						Title:     homeTitle,
						Icon:      "lucide:house",
						AffixTab:  true,
						KeepAlive: true,
					},
				},
			},
		},
		{
			Name:      "Userinfo",
			Path:      "/userinfo",
			Component: "BasicLayout",
			Redirect:  "/userinfo",
			Meta: VbenMeta{
				Title: userCenterTitle,
				Icon:  "lucide:user",
				Order: 2,
			},
			Children: []VbenRoute{
				{
					Name:      "UserSettings",
					Path:      "/userinfo",
					Component: "/home/userinfo",
					Meta: VbenMeta{
						Title:     settingsTitle,
						Icon:      "lucide:settings",
						KeepAlive: true,
					},
				},
			},
		},
	}

	http.SuccessWithI18n(ctx, "common.success", menus)
}
