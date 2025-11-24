package user

import (
	"github.com/gin-gonic/gin"
	serviceUser "github.com/lvjiaben/go-wheel/app/backend/service/user"
	validateUser "github.com/lvjiaben/go-wheel/app/backend/validate"
	"github.com/lvjiaben/go-wheel/pkg/container"
	"github.com/lvjiaben/go-wheel/pkg/utils/http"
)

type UserController struct {
	container   *container.Container
	userService *serviceUser.UserService
}

func NewUserController(c *container.Container) *UserController {
	return &UserController{
		container:   c,
		userService: serviceUser.NewUserService(c),
	}
}

func (ctrl *UserController) List(ctx *gin.Context) {
	http.SuccessWithI18n(ctx, "common.success", ctrl.userService.List(ctx))
}

func (ctrl *UserController) Create(ctx *gin.Context) {
	form, valid := validateUser.ValidateUserCreate(ctx)
	if !valid {
		return
	}
	res, err := ctrl.userService.Create(ctx, form)
	if err != nil {
		http.ErrorWithI18n(ctx, err.Error(), nil)
		return
	}
	http.SuccessWithI18n(ctx, "common.success", res)
}

func (ctrl *UserController) Update(ctx *gin.Context) {
	form, valid := validateUser.ValidateUserUpdate(ctx)
	if !valid {
		return
	}
	res, err := ctrl.userService.Update(ctx, form)
	if err != nil {
		http.ErrorWithI18n(ctx, err.Error(), nil)
		return
	}
	http.SuccessWithI18n(ctx, "common.success", res)
}

func (ctrl *UserController) Delete(ctx *gin.Context) {
	form, valid := validateUser.ValidateUserDelete(ctx)
	if !valid {
		return
	}
	if err := ctrl.userService.Delete(ctx, form.Ids); err != nil {
		http.ErrorWithI18n(ctx, err.Error(), nil)
		return
	}
	http.SuccessWithI18n(ctx, "common.success", nil)
}

// UpdateMoney 更新用户余额
func (ctrl *UserController) UpdateMoney(ctx *gin.Context) {
	// 验证参数
	form, valid := validateUser.ValidateUserUpdateMoney(ctx)
	if !valid {
		return
	}
	// 调用服务层
	if err := ctrl.userService.UpdateMoney(form.Id, form.Type, form.Money, form.Note, form.Source); err != nil {
		http.ErrorWithI18n(ctx, err.Error(), nil)
		return
	}
	http.SuccessWithI18n(ctx, "common.success", nil)
}

// UpdateScore 更新用户积分
func (ctrl *UserController) UpdateScore(ctx *gin.Context) {
	// 验证参数
	form, valid := validateUser.ValidateUserUpdateScore(ctx)
	if !valid {
		return
	}
	// 调用服务层
	if err := ctrl.userService.UpdateScore(form.Id, form.Type, form.Score, form.Note, form.Source); err != nil {
		http.ErrorWithI18n(ctx, err.Error(), nil)
		return
	}
	http.SuccessWithI18n(ctx, "common.success", nil)
}

// Operate 操作用户字段（status等开关字段）- 支持单个或批量
func (ctrl *UserController) Operate(ctx *gin.Context) {
	// 验证参数
	form, valid := validateUser.ValidateUserOperate(ctx)
	if !valid {
		return
	}
	// 调用服务层
	if err := ctrl.userService.Operate(form.Ids, form.Field, form.Value); err != nil {
		http.ErrorWithI18n(ctx, err.Error(), nil)
		return
	}
	http.SuccessWithI18n(ctx, "common.success", nil)
}
