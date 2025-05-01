package controller

import (
	"crypto/md5"
	"encoding/hex"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/app/backend/model"
	"github.com/lvjiaben/go-wheel/app/backend/service"
	"github.com/lvjiaben/go-wheel/app/backend/validate"
	"github.com/lvjiaben/go-wheel/pkg/global"
	"github.com/lvjiaben/go-wheel/pkg/http"
	"github.com/lvjiaben/go-wheel/pkg/jwt"
	"github.com/spf13/viper"
)

type AuthController struct{}

var authService = service.AuthService{}

func (a *AuthController) User(c *gin.Context) {
	var user model.Admin
	if res := global.DB.First(&user, "id = ?", c.GetInt("admin_id")); res.Error != nil {
		http.Error(c, "账户不存在")
		return
	}
	http.Success(c, "", map[string]any{
		"id":       user.Id,
		"realName": user.Username,
		"username": user.Username,
		"avatar":   user.Avatar,
		"lang":     c.GetString("lang"),
	})
}

func (a *AuthController) Menus(c *gin.Context) {
	http.Success(c, "", service.GetMenusFromDB(c.GetInt("admin_id"), c.GetString("lang")))
}

func (a *AuthController) Codes(c *gin.Context) {
	ruleIds := authService.GetRuleIds(c.GetInt("admin_id"))
	isSuper := false
	for _, id := range ruleIds {
		if id == "*" {
			isSuper = true
			break
		}
	}
	var codes []string
	query := global.DB.Table((&model.AdminAuthRule{}).TableName()).Where("type = ? AND alias IS NOT NULL", "button")
	if !isSuper {
		query = query.Where("id in (?)", ruleIds)
	}
	query.Pluck("alias", &codes)
	http.Success(c, "", codes)
}

func (a *AuthController) Login(c *gin.Context) {
	var u validate.AdminLogin
	if err := c.ShouldBind(&u); err != nil {
		http.Error(c, http.Translate(err, u))
		return
	}
	var user model.Admin
	if res := global.DB.First(&user, "username = ?", u.Username); res.Error != nil {
		http.Error(c, "账户或密码错误")
		return
	}
	log := &model.AdminLoginLog{
		Username: u.Username,
		Status:   0,
		Ip:       c.ClientIP(),
	}
	global.DB.Create(&log)
	currentTime := int(time.Now().Unix())
	if user.Failures > viper.GetInt("admin.login_failures") && currentTime-user.UpdatedAt <= viper.GetInt("admin.login_failures_second") {
		http.Error(c, "密码重试次数过多，请稍后再试")
		return
	}
	hash := md5.New()
	hash.Write([]byte(u.Password))
	md5Hash := hex.EncodeToString(hash.Sum(nil))
	if md5Hash != user.Password {
		global.DB.Model(&user).Where("id = ?", user.Id).Update("failures", user.Failures+1)
		http.Error(c, "账户或密码错误")
		return
	}
	token, err := jwt.GenToken(int64(user.Id))
	if err != nil {
		http.Error(c, "TOKEN系统出错，请联系技术同学")
		return
	}
	global.DB.Model(&user).Where("id = ?", user.Id).Updates(model.Admin{Failures: 0, Token: token})
	var loginLog model.AdminLoginLog
	global.DB.Model(&loginLog).Where("id = ?", log.Id).Update("status", 1)
	http.Success(c, "登陆成功", map[string]any{
		"accessToken": token,
		"id":          user.Id,
		"username":    user.Username,
		"avatar":      user.Avatar,
	}, nil)
	return
}
