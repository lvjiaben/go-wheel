package controller

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/app/backend/model"
	"github.com/lvjiaben/go-wheel/app/backend/service"
	"github.com/lvjiaben/go-wheel/app/backend/validate"
	"github.com/lvjiaben/go-wheel/pkg/container"
	"github.com/lvjiaben/go-wheel/pkg/utils"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct {
	authService *service.AuthService
	menuService *service.MenuService
	container   *container.Container
}

func NewAuthController(c *container.Container) *AuthController {
	return &AuthController{
		authService: service.NewAuthService(c),
		menuService: service.NewMenuService(c),
		container:   c,
	}
}

func (a *AuthController) User(c *gin.Context) {
	// 设置Container到上下文中，供响应函数使用
	c.Set("container", a.container)

	var user model.Admin
	if res := a.container.GetDB().First(&user, "id = ?", c.GetInt("admin_id")); res.Error != nil {
		utils.NotFoundErrorI18n(c, "user.not_found")
		return
	}

	isSuper := a.authService.IsSuperAdmin(c.GetInt("admin_id"))
	var codes []string

	// 获取用户权限码
	query := a.container.GetDB().Table((&model.AdminAuthRule{}).TableName()).Where("type = ? AND alias IS NOT NULL", "button")
	if !isSuper {
		var groupAccess model.AdminAuthGroupAccess
		if err := a.container.GetDB().First(&groupAccess, "uid = ?", user.Id).Error; err != nil {
			utils.NotFoundErrorI18n(c, "role.not_found")
			return
		}

		var group model.AdminAuthGroup
		if err := a.container.GetDB().First(&group, "id = ?", groupAccess.Gid).Error; err != nil {
			utils.NotFoundErrorI18n(c, "role.not_found")
			return
		}

		var ruleIds []int
		if err := json.Unmarshal([]byte(group.Rules), &ruleIds); err != nil {
			utils.ErrorWithI18n(c, "common.server_error", nil)
			return
		}

		query = query.Where("id in (?)", ruleIds)
	}

	if err := query.Pluck("alias", &codes).Error; err != nil {
		utils.ErrorWithI18n(c, "common.server_error", nil)
		return
	}

	utils.SuccessWithI18n(c, "common.success", gin.H{
		"id":       user.Id,
		"username": user.Username,
		"avatar":   user.Avatar,
		"codes":    codes,
	})
}

func (a *AuthController) Menus(c *gin.Context) {
	c.Set("container", a.container)
	utils.SuccessWithI18n(c, "common.success", a.menuService.GetMenusFromDB(c.GetInt("admin_id"), c.GetString("lang")))
}

func (a *AuthController) Codes(c *gin.Context) {
	c.Set("container", a.container)
	utils.SuccessWithI18n(c, "common.success", []string{"admin", "editor"})
}

// GetUserInfo方法已被User方法替代，保留此方法是为了兼容旧API
func (a *AuthController) GetUserInfo(c *gin.Context) {
	a.User(c)
}

func (a *AuthController) Login(c *gin.Context) {
	// 设置Container到上下文中，供响应函数使用
	c.Set("container", a.container)

	var req validate.AdminLogin
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidateErrorI18n(c, "common.bad_request")
		return
	}

	// 验证用户名和密码
	var user model.Admin
	if res := a.container.GetDB().First(&user, "username = ?", req.Username); res.Error != nil {
		utils.NotFoundErrorI18n(c, "user.not_found")
		return
	}

	// 使用bcrypt比较密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		a.container.GetLogger().Debug("密码不匹配",
			zap.String("input", req.Password),
			zap.String("stored", user.Password))
		utils.AuthErrorI18n(c, "auth.login_failed")
		return
	}

	// 生成 token
	token, err := a.authService.GenerateToken(user.Id)
	if err != nil {
		a.container.GetLogger().Error("生成token失败", zap.Error(err))
		utils.ErrorWithI18n(c, "common.server_error", nil)
		return
	}

	utils.SuccessWithI18n(c, "common.success", gin.H{
		"token": token,
	})
}

func (a *AuthController) Logout(c *gin.Context) {
	c.Set("container", a.container)
	utils.SuccessWithI18n(c, "common.success", nil)
}
