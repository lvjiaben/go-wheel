package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func SetLang() func(c *gin.Context) {
	return func(c *gin.Context) {
		acceptLanguage := c.GetHeader("Accept-Language")
		if strings.Contains(acceptLanguage, "en") {
			c.Set("lang", "en")
		} else {
			c.Set("lang", "zh")
		}
		c.Next()
	}
}
