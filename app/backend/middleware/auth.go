package middleware

import (
	"net/http"
	"strings"

	"github.com/lvjiaben/go-wheel/app/backend/service"

	"github.com/lvjiaben/go-wheel/frame/global"

	"github.com/spf13/viper"

	"github.com/lvjiaben/go-wheel/app/backend/model"

	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/utils/jwt"
)

var authService = service.AuthService{}

func PermissionCheck() func(c *gin.Context) {
	return func(c *gin.Context) {
		userId := c.GetInt("admin_id")
		if userId <= 0 {
			c.JSON(http.StatusOK, gin.H{
				"code":    401,
				"message": "请先登陆",
			})
			c.Abort()
			return
		}
		path := c.Request.URL.Path
		if authService.Check(path, userId) == false {
			c.JSON(http.StatusOK, gin.H{
				"code":    401,
				"message": "无权访问",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func JWTAuthCheck() func(c *gin.Context) {
	return func(c *gin.Context) {
		authHeader := c.Request.Header.Get("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusOK, gin.H{
				"code":    401,
				"message": "请先登陆",
			})
			c.Abort()
			return
		}
		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusOK, gin.H{
				"code":    401,
				"message": "请先登陆",
			})
			c.Abort()
			return
		}
		mc, err := jwt.ParseToken(parts[1])
		if err != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    401,
				"message": "请先登陆",
			})
			c.Abort()
			return
		}
		var user model.Admin
		if res := global.DB.First(&user, "id = ?", mc.ID); res.Error != nil {
			c.JSON(http.StatusOK, gin.H{
				"code":    401,
				"message": "请先登陆",
			})
			c.Abort()
			return
		}
		if viper.GetBool("admin.login_sso") == true && user.Token != parts[1] {
			c.JSON(http.StatusOK, gin.H{
				"code":    401,
				"message": "登陆失效",
			})
			c.Abort()
			return
		}
		c.Set("admin_id", user.Id)
		c.Next()
	}
}
