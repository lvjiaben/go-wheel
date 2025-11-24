package user

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/app/backend/validate"
	"github.com/lvjiaben/go-wheel/app/common/builder"
	commonModel "github.com/lvjiaben/go-wheel/app/common/model"
	commonService "github.com/lvjiaben/go-wheel/app/common/service"
	"github.com/lvjiaben/go-wheel/pkg/container"
	"github.com/lvjiaben/go-wheel/pkg/utils/crypto"
	"github.com/lvjiaben/go-wheel/pkg/utils/datatype"
	"gorm.io/gorm"
)

// UserService 用户服务
type UserService struct {
	container     *container.Container
	crudService   *builder.CRUDBuilder[commonModel.User]
	codeGenerator *commonService.CodeGeneratorService
}

// NewUserService 创建用户服务
func NewUserService(c *container.Container) *UserService {
	return &UserService{
		container: c,
		crudService: builder.NewCRUDBuilderWithProvider[commonModel.User](
			func(ctx *gin.Context) *gorm.DB {
				return c.GetDB().WithContext(ctx.Request.Context())
			},
		).WithSearchFields("username", "email", "mobile", "code"),
		codeGenerator: commonService.NewCodeGeneratorService(c),
	}
}

// List 获取用户列表
func (s *UserService) List(ctx *gin.Context) map[string]interface{} {
	return s.crudService.WithContext(ctx).List()
}

// Create 创建用户
func (s *UserService) Create(ctx *gin.Context, form *validate.UserCreate) (*commonModel.User, error) {
	// 生成邀请码（如果没有提供）
	if form.Code == "" {
		form.Code = s.codeGenerator.GenerateUniqueInviteCode("user", "code", 0)
	}
	// 生成随机密码（如果没有提供）
	if form.Password == "" {
		form.Password = s.codeGenerator.GenerateRandomPassword()
	}
	// 使用 bcrypt+salt 加密密码
	salt, err := crypto.GenerateSalt()
	if err != nil {
		return nil, fmt.Errorf("生成盐值失败: %v", err)
	}
	hashedPassword, err := crypto.HashPassword(form.Password, salt)
	if err != nil {
		return nil, fmt.Errorf("密码加密失败: %v", err)
	}
	form.Password = hashedPassword
	form.Salt = salt // 保存盐值

	return s.crudService.WithContext(ctx).Create(form)
}

// Update 更新用户
func (s *UserService) Update(ctx *gin.Context, form *validate.UserUpdate) (*commonModel.User, error) {
	// 如果有密码，先加密
	if form.Password != "" {
		salt, err := crypto.GenerateSalt()
		if err != nil {
			return nil, fmt.Errorf("生成盐值失败: %v", err)
		}
		hashedPassword, err := crypto.HashPassword(form.Password, salt)
		if err != nil {
			return nil, fmt.Errorf("密码加密失败: %v", err)
		}
		form.Password = hashedPassword
		form.Salt = salt // 同时更新盐值
	}
	return s.crudService.WithContext(ctx).Update(form.Id, form)
}

// Delete 删除用户
func (s *UserService) Delete(ctx *gin.Context, ids []int) error {
	return s.crudService.WithContext(ctx).Delete(ids)
}

// UpdateMoney 更新用户余额（自增或自减）
func (s *UserService) UpdateMoney(id int, operationType string, money float64, note, source string) error {
	return s.container.GetDB().Transaction(func(tx *gorm.DB) error {
		var logType int
		var model commonModel.User
		if operationType == "add" {
			tx.Model(&model).Where("id = ?", id).Update("money", gorm.Expr("money + ?", money))
			logType = 1 // 增加
		} else {
			tx.Model(&model).Where("id = ?", id).Update("money", gorm.Expr("money - ?", money))
			logType = 0 // 减少
		}
		// 创建余额日志
		log := commonModel.UserMoneyLog{
			UserId: id,
			Type:   logType,
			Money:  money,
			Note:   note,
			Source: source,
		}
		tx.Create(&log)
		return nil
	})
}

// UpdateScore 更新用户积分（自增或自减）
func (s *UserService) UpdateScore(id int, operationType string, score float64, note, source string) error {
	return s.container.GetDB().Transaction(func(tx *gorm.DB) error {
		var model commonModel.User
		var logType int
		if operationType == "add" {
			tx.Model(&model).Where("id = ?", id).Update("score", gorm.Expr("score + ?", score))
			logType = 1 // 增加
		} else {
			tx.Model(&model).Where("id = ?", id).Update("score", gorm.Expr("score - ?", score))
			logType = 0 // 减少
		}
		// 创建积分日志
		log := commonModel.UserScoreLog{
			UserId: id,
			Type:   logType,
			Score:  score,
			Note:   note,
			Source: source,
		}
		tx.Create(&log)

		return nil
	})
}

// Operate 直接操作字段
func (s *UserService) Operate(ids []int, field string, value int8) error {
	// 检查字段是否允许操作
	if !datatype.Contains([]string{"status"}, field) {
		return fmt.Errorf("common.server_error")
	}
	s.container.GetDB().Model(&commonModel.User{}).Where("id IN ?", ids).Update(field, value)
	return nil
}
