package controller

import (
	"time"

	"github.com/lvjiaben/go-wheel/pkg/constants"
	"github.com/lvjiaben/go-wheel/pkg/monitor"
	"github.com/lvjiaben/go-wheel/pkg/utils/crypto"

	"github.com/lvjiaben/go-wheel/pkg/utils/http"

	"github.com/lvjiaben/go-wheel/app/backend/model/admin"

	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/app/backend/service"
	serviceAdmin "github.com/lvjiaben/go-wheel/app/backend/service/admin"
	"github.com/lvjiaben/go-wheel/app/backend/utils"
	"github.com/lvjiaben/go-wheel/app/backend/validate"
	"github.com/lvjiaben/go-wheel/pkg/captcha"
	"github.com/lvjiaben/go-wheel/pkg/container"
)

// AuthController 认证控制器
type AuthController struct {
	container      *container.Container
	authService    *service.AuthService
	roleService    *serviceAdmin.RoleService
	adminService   *serviceAdmin.AdminService
	menuService    *serviceAdmin.MenuService
	captchaService *captcha.CaptchaService
	authUtils      *utils.AuthUtils
}

// NewAuthController 创建认证控制器
func NewAuthController(c *container.Container) *AuthController {
	return &AuthController{
		container:      c,
		roleService:    serviceAdmin.NewRoleService(c),
		menuService:    serviceAdmin.NewMenuService(c),
		authService:    service.NewAuthService(c),
		adminService:   serviceAdmin.NewAdminService(c),
		captchaService: captcha.NewCaptchaService(c.GetRDB()),
		authUtils:      utils.NewAuthUtils(c),
	}
}

// Info 获取用户信息
func (c *AuthController) Info(ctx *gin.Context) {
	id := ctx.GetInt("admin_id")
	user, err := c.adminService.GetById(id)
	if err != nil {
		http.ErrorWithI18n(ctx, "common.error", nil)
		return
	}
	http.SuccessWithI18n(ctx, "common.success", gin.H{
		"id":       user.Id,
		"username": user.Username,
		"avatar":   user.Avatar,
		"email":    user.Email,
		"realName": user.Username,
	})
}

// Login 登录
func (c *AuthController) Login(ctx *gin.Context) {
	// 记录登录尝试（Prometheus）
	if metrics := c.getPrometheusMetrics(); metrics != nil {
		metrics.RecordLoginAttempt("admin")
	}

	// 使用 ValidateLogin 进行验证
	req, valid := validate.ValidateLogin(ctx)
	if !valid {
		if metrics := c.getPrometheusMetrics(); metrics != nil {
			metrics.RecordLoginFailure("admin", "invalid_request")
		}
		return
	}

	// 获取客户端IP地址
	clientIP := ctx.ClientIP()

	// 检查IP登录失败次数限制
	if blocked, remainingTime := c.checkLoginRateLimit(clientIP); blocked {
		if metrics := c.getPrometheusMetrics(); metrics != nil {
			metrics.RecordLoginFailure("admin", "rate_limited")
		}
		http.ErrorWithI18n(ctx, "backend.auth.too_many_attempts", gin.H{
			"remaining_seconds": remainingTime,
		})
		return
	}

	// 验证码验证
	if !c.captchaService.Verify(req.Captcha.ID, req.Captcha.Code) {
		if metrics := c.getPrometheusMetrics(); metrics != nil {
			metrics.RecordLoginFailure("admin", "invalid_captcha")
		}
		http.ErrorWithI18n(ctx, "backend.auth.captcha_invalid", nil)
		return
	}

	// 登录验证，传递IP地址
	user, err := c.authService.Login(req.Username, req.Password, clientIP)
	if err != nil {
		// 登录失败，记录失败次数
		c.recordLoginFailure(clientIP)
		if metrics := c.getPrometheusMetrics(); metrics != nil {
			metrics.RecordLoginFailure("admin", "invalid_credentials")
		}
		http.ErrorWithI18n(ctx, "backend.auth.login_failed", nil)
		return
	}

	// 登录成功，清除失败记录
	c.clearLoginFailure(clientIP)

	http.SuccessWithI18n(ctx, "backend.auth.login_success", gin.H{
		"id":          user.Id,
		"accessToken": user.Token,
		"roles":       []string{},
		"expires":     time.Now().AddDate(0, 0, c.container.GetConfig().GetInt("jwt.expire_day")).Format("2006/01/02 15:04:05"),
		"avatar":      user.Avatar,
		"username":    user.Username,
		"realName":    user.Username,
		"email":       user.Email,
	})
}

// Profile 修改资料
func (c *AuthController) Profile(ctx *gin.Context) {
	req, valid := validate.ValidateProfile(ctx)
	if !valid {
		return
	}
	c.authService.Profile(req, ctx.GetInt("admin_id"))
	http.SuccessWithI18n(ctx, "common.success", nil)
}

// checkLoginRateLimit 检查登录频率限制
// 返回值：(是否被限制, 剩余封禁秒数)
func (c *AuthController) checkLoginRateLimit(ip string) (bool, int64) {
	redisKey := "login_fail:" + ip
	redisClient := c.container.GetRDB()
	ctx := c.container.GetContext()

	// 获取失败次数
	failCount, err := redisClient.Get(ctx, redisKey).Int64()
	if err != nil {
		// 如果key不存在，返回0
		return false, 0
	}

	// 从配置获取最大失败次数（默认使用常量）
	maxFailures := c.container.GetConfig().GetInt("admin.login_failures")
	if maxFailures == 0 {
		maxFailures = constants.DefaultMaxLoginFailures
	}

	// 如果失败次数超过限制
	if failCount >= int64(maxFailures) {
		// 获取剩余过期时间
		ttl, err := redisClient.TTL(ctx, redisKey).Result()
		if err != nil {
			return false, 0
		}
		return true, int64(ttl.Seconds())
	}

	return false, 0
}

// recordLoginFailure 记录登录失败
func (c *AuthController) recordLoginFailure(ip string) {
	redisKey := "login_fail:" + ip
	redisClient := c.container.GetRDB()
	ctx := c.container.GetContext()

	// 增加失败次数
	failCount, err := redisClient.Incr(ctx, redisKey).Result()
	if err != nil {
		c.container.GetLogger().Error("记录登录失败次数失败: " + err.Error())
		return
	}

	// 如果是第一次失败，设置过期时间
	if failCount == 1 {
		redisClient.Expire(ctx, redisKey, constants.LoginFailureRecordDuration)
	}

	// 从配置获取最大失败次数
	maxFailures := c.container.GetConfig().GetInt("admin.login_failures")
	if maxFailures == 0 {
		maxFailures = constants.DefaultMaxLoginFailures
	}

	// 如果达到最大失败次数，延长封禁时间
	if failCount >= int64(maxFailures) {
		redisClient.Expire(ctx, redisKey, constants.LoginFailureLockDuration)
		c.container.GetLogger().Warn("IP " + ip + " 登录失败次数过多，已被封禁")
	}
}

// clearLoginFailure 清除登录失败记录
func (c *AuthController) clearLoginFailure(ip string) {
	redisKey := "login_fail:" + ip
	redisClient := c.container.GetRDB()
	ctx := c.container.GetContext()

	err := redisClient.Del(ctx, redisKey).Err()
	if err != nil {
		c.container.GetLogger().Error("清除登录失败记录失败: " + err.Error())
	}
}

// getPrometheusMetrics 获取 Prometheus 指标收集器（辅助方法）
func (c *AuthController) getPrometheusMetrics() *monitor.PrometheusMetrics {
	// 从 container 获取（需要类型断言）
	if metricsInterface := c.container.Get("prometheus_metrics"); metricsInterface != nil {
		if metrics, ok := metricsInterface.(*monitor.PrometheusMetrics); ok {
			return metrics
		}
	}
	return nil
}

// Password 修改密码
func (c *AuthController) Password(ctx *gin.Context) {
	currentAdminId := ctx.GetInt("admin_id")
	var changeData struct {
		OldPassword string `json:"old_password" binding:"required"`
		NewPassword string `json:"new_password" binding:"required,min=6"`
	}
	if err := ctx.ShouldBindJSON(&changeData); err != nil {
		http.ErrorWithI18n(ctx, "common.invalid_params", nil)
		return
	}
	adminItem, err := c.adminService.GetById(currentAdminId)
	if err != nil {
		http.ErrorWithI18n(ctx, err.Error(), nil)
		return
	}
	if !crypto.PasswordVerifyWithSalt(changeData.OldPassword, adminItem.Salt, adminItem.Password) {
		http.ErrorWithI18n(ctx, "backend.admin.old_password_incorrect", nil)
		return
	}
	// 更新密码
	salt, _ := crypto.GenerateSalt()
	hashedPassword, _ := crypto.PasswordHashWithSalt(changeData.NewPassword, salt)
	if err := c.container.GetDB().Model(&admin.Admin{}).Where("id = ?", currentAdminId).Updates(map[string]interface{}{
		"password": hashedPassword,
		"salt":     salt,
	}).Error; err != nil {
		http.ErrorWithI18n(ctx, "backend.admin.update_password_failed", nil)
	}
	http.SuccessWithI18n(ctx, "backend.admin.change_password_success", nil)
}

// Log 日志
func (c *AuthController) Log(ctx *gin.Context) {
	http.SuccessWithI18n(ctx, "common.success", c.authService.Log(ctx))
}

// Logout 登出
func (c *AuthController) Logout(ctx *gin.Context) {
	if c.container.GetConfig().GetBool("admin.login_sso") {
		c.container.GetDB().Model(&admin.Admin{}).Where("id = ?", ctx.GetInt("admin_id")).Update("token", "")
	}
	http.SuccessWithI18n(ctx, "backend.auth.logout_success", nil)
}

// Permission 获取用户权限
func (c *AuthController) Permission(ctx *gin.Context) {
	adminId := ctx.GetInt("admin_id")
	db := c.container.GetDBWithContext(ctx.Request.Context())
	// 获取权限
	var permissions []string
	if isSuper, _ := c.authUtils.IsAdminSuper(adminId); isSuper {
		// 超级管理员获取所有按钮权限
		db.Model(&admin.Menu{}).Where("type = ?", "button").Pluck("permission", &permissions)
	} else {
		// 普通管理员获取角色权限
		roleIds, err := c.authUtils.GetAdminDirectRoleIds(adminId)
		if err == nil && len(roleIds) > 0 {
			// 收集所有角色的菜单ID
			var allMenuIds []int
			for _, roleId := range roleIds {
				menuIds, err := c.authUtils.GetRoleMenuIds(roleId)
				if err != nil {
					continue
				}
				allMenuIds = append(allMenuIds, menuIds...)
			}

			// 去重菜单ID
			uniqueMenuIds := make(map[int]bool)
			var finalMenuIds []int
			for _, id := range allMenuIds {
				if !uniqueMenuIds[id] {
					uniqueMenuIds[id] = true
					finalMenuIds = append(finalMenuIds, id)
				}
			}

			// 使用 IN 查询一次性获取所有菜单（优化 N+1 查询）
			if len(finalMenuIds) > 0 {
				db.Model(&admin.Menu{}).
					Where("id IN ? AND type = ?", finalMenuIds, "button").
					Pluck("permission", &permissions)
			}
		}
	}

	http.SuccessWithI18n(ctx, "common.success", gin.H{
		"permissions": permissions,
	})
}

// Menus 获取用户菜单
func (c *AuthController) Menus(ctx *gin.Context) {
	adminId := ctx.GetInt("admin_id")
	// 获取菜单（authService 内部会使用 container.GetDB()）
	menus, err := c.authService.GetUserMenus(adminId)
	if err != nil {
		http.ErrorWithI18n(ctx, "common.server_error", nil)
		return
	}
	http.SuccessWithI18n(ctx, "common.success", menus)
}
