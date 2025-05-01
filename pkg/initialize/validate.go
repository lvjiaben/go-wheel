package initialize

import (
	"admin/pkg/global"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/locales/zh"
	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	zh_translations "github.com/go-playground/validator/v10/translations/zh"
)

func ValidateLoad() {
	zh := zh.New()
	global.VALID_UNI = ut.New(zh, zh)
	global.VALID_TRANS, _ = global.VALID_UNI.GetTranslator("zh")
	//获取gin的校验器
	validate := binding.Validator.Engine().(*validator.Validate)
	//注册翻译器
	zh_translations.RegisterDefaultTranslations(validate, global.VALID_TRANS)
}
