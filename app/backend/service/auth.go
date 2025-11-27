package service

import (
	"errors"
	"fmt"
	"time"

	"github.com/lvjiaben/go-wheel/app/backend/utils"

	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/app/common/builder"

	"github.com/lvjiaben/go-wheel/app/backend/validate"

	"github.com/lvjiaben/go-wheel/app/backend/model/admin"
	serviceAdmin "github.com/lvjiaben/go-wheel/app/backend/service/admin"
	"github.com/lvjiaben/go-wheel/pkg/container"
	"github.com/lvjiaben/go-wheel/pkg/jwt"
	"github.com/lvjiaben/go-wheel/pkg/utils/crypto"
	"go.uber.org/zap"
)

type AuthService struct {
	container    *container.Container
	adminService *serviceAdmin.AdminService
	menuService  *serviceAdmin.MenuService
	roleService  *serviceAdmin.RoleService
	authUtils    *utils.AuthUtils
}

func NewAuthService(c *container.Container) *AuthService {
	return &AuthService{
		container:    c,
		adminService: serviceAdmin.NewAdminService(c),
		menuService:  serviceAdmin.NewMenuService(c),
		roleService:  serviceAdmin.NewRoleService(c),
		authUtils:    utils.NewAuthUtils(c),
	}
}

func (s *AuthService) Log(ctx *gin.Context) map[string]interface{} {
	db := s.container.GetDB().WithContext(ctx.Request.Context())
	return builder.NewCRUDBuilder[admin.LoginLog](db).WithFilter(false).WithContext(ctx).WithSearchFields("ip").List()
}

// GenerateToken 生成JWT令牌
func (s *AuthService) GenerateToken(adminId int) (string, error) {
	// 获取JWT相关配置
	secret := s.container.GetConfig().GetString("jwt.secret")
	expire := s.container.GetConfig().GetInt("jwt.expire_day")

	// 获取用户名
	var username string
	s.container.GetDB().Model(&admin.Admin{}).Where("id = ?", adminId).Pluck("username", &username)

	// 生成令牌
	return jwt.GenerateToken(adminId, username, secret, expire)
}

func (s *AuthService) Profile(req *validate.ProfileRequest, admin_id int) {
	s.container.GetDB().Model(&admin.Admin{}).Where("id = ?", admin_id).Updates(req)
}

// Login 登录
func (s *AuthService) Login(username, password string, ip string) (*admin.Admin, error) {
	var adminUser admin.Admin
	if err := s.container.GetDB().Where("username = ? AND status = 1", username).First(&adminUser).Error; err != nil {
		s.recordLoginLog(username, ip, 0, "用户不存在")
		return nil, errors.New("用户不存在")
	}
	if !crypto.PasswordVerifyWithSalt(password, adminUser.Salt, adminUser.Password) {
		s.recordLoginLog(username, ip, 0, "密码不正确")
		return nil, fmt.Errorf("password incorrect")
	}

	// 生成JWT令牌
	token, err := jwt.GenerateToken(adminUser.Id, adminUser.Username,
		s.container.GetConfig().GetString("jwt.secret"),
		s.container.GetConfig().GetInt("jwt.expire_day"))
	if err != nil {
		return nil, err
	}

	// 如果启用单点登录，更新token
	if s.container.GetConfig().GetBool("admin.login_sso") {
		s.container.GetDB().Model(&admin.Admin{}).Where("id = ?", adminUser.Id).Update("token", token)
	}

	// 记录登录成功日志
	s.recordLoginLog(username, ip, 1, "登录成功")

	adminUser.Token = token
	return &adminUser, nil
}

// recordLoginLog 记录登录日志
func (s *AuthService) recordLoginLog(username string, ip string, status int, message string) {
	// 创建登录日志记录
	loginLog := admin.LoginLog{
		Username:  username,
		Ip:        ip,
		Status:    status,
		CreatedAt: time.Now().Unix(),
	}

	// 保存日志
	if err := s.container.GetDB().Create(&loginLog).Error; err != nil {
		s.container.GetLogger().Error("记录登录日志失败",
			zap.String("username", username),
			zap.String("ip", ip),
			zap.Error(err))
	} else {
		s.container.GetLogger().Info("记录登录日志",
			zap.String("username", username),
			zap.String("ip", ip),
			zap.Int("status", status),
			zap.String("message", message))
	}
}

// GetUserMenus 获取用户菜单（Vben格式）
func (s *AuthService) GetUserMenus(ctx *gin.Context, adminId int) ([]serviceAdmin.VbenRoute, error) {
	// 如果是超级管理员，返回所有菜单
	if isSuper, _ := s.authUtils.IsAdminSuper(adminId); isSuper {
		return s.menuService.GetVbenRoutes(ctx, []int{}, "all")
	}

	// 获取用户角色
	roleIds, err := s.authUtils.GetAdminDirectRoleIds(adminId)
	if err != nil {
		return nil, err
	}
	// 获取角色菜单
	var allMenuIds []int
	for _, roleId := range roleIds {
		menuIds, err := s.authUtils.GetRoleMenuIds(roleId)
		if err != nil {
			continue
		}
		allMenuIds = append(allMenuIds, menuIds...)
	}

	// 去重
	menuIds := make([]int, 0)
	menuMap := make(map[int]bool)
	for _, id := range allMenuIds {
		if !menuMap[id] {
			menuMap[id] = true
			menuIds = append(menuIds, id)
		}
	}

	if len(menuIds) == 0 {
		menuIds = []int{0}
	}
	// 获取菜单树
	return s.menuService.GetVbenRoutes(ctx, menuIds, "all")
}
