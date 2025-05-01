package http

import (
	"net/http"
	"reflect"
	"time"

	ut "github.com/go-playground/universal-translator"

	"github.com/go-playground/validator/v10"

	"github.com/gin-gonic/gin"
)

func Success(c *gin.Context, msg string, data ...interface{}) {
	if len(data) < 1 {
		data = []any{nil}
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    200,
		"message": msg,
		"data":    data[0],
		"time":    time.Now().Unix(),
	})
}

func Error(c *gin.Context, msg string, data ...interface{}) {
	if len(data) < 1 {
		data = []any{nil}
	}
	c.JSON(http.StatusOK, map[string]interface{}{
		"code":    500,
		"message": msg,
		"data":    data[0],
		"time":    time.Now().Unix(),
	})
}

func ValidateError(errs validator.ValidationErrors, r interface{}) string {
	s := reflect.TypeOf(r)
	for _, fieldError := range errs {
		filed, _ := s.FieldByName(fieldError.Field())
		errTag := fieldError.Tag() + "_msg"
		errTagText := filed.Tag.Get(errTag)
		errText := filed.Tag.Get("msg")
		if errTagText != "" {
			return errTagText
		}
		if errText != "" {
			return errText
		}
	}
	return ""
}

func Translate(err error, r interface{}) string {
	var trans ut.Translator
	var result string
	self := ValidateError(err.(validator.ValidationErrors), r)
	if self != "" {
		return self
	}
	for _, err := range err.(validator.ValidationErrors) {
		result = result + err.Translate(trans)
	}
	return result
}
