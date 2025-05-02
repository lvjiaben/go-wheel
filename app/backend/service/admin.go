package service

import (
	"github.com/lvjiaben/go-wheel/app/backend/model"
	"github.com/lvjiaben/go-wheel/pkg/container"
	"golang.org/x/crypto/bcrypt"
)

// AdminService 管理员服务
type AdminService struct {
	container *container.Container
}

// NewAdminService 创建管理员服务
func NewAdminService(c *container.Container) *AdminService {
	return &AdminService{
		container: c,
	}
}

// Login 登录
func (s *AdminService) Login(username, password string) (*model.Admin, error) {
	var admin model.Admin
	if err := s.container.GetDB().Where("username = ?", username).First(&admin).Error; err != nil {
		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(admin.Password), []byte(password)); err != nil {
		return nil, err
	}

	return &admin, nil
}

// UpdateToken 更新token
func (s *AdminService) UpdateToken(id int, token string) error {
	return s.container.GetDB().Model(&model.Admin{}).Where("id = ?", id).Update("token", token).Error
}

// GetList 获取管理员列表
func (s *AdminService) GetList() ([]model.Admin, error) {
	var admins []model.Admin
	if err := s.container.GetDB().Find(&admins).Error; err != nil {
		return nil, err
	}
	return admins, nil
}

// Create 创建管理员
func (s *AdminService) Create(admin *model.Admin) error {
	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	admin.Password = string(hashedPassword)

	return s.container.GetDB().Create(admin).Error
}

// Update 更新管理员
func (s *AdminService) Update(admin *model.Admin) error {
	// 如果密码不为空，则加密密码
	if admin.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(admin.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		admin.Password = string(hashedPassword)
	}

	return s.container.GetDB().Model(admin).Updates(admin).Error
}

// Delete 删除管理员
func (s *AdminService) Delete(id int) error {
	return s.container.GetDB().Delete(&model.Admin{}, id).Error
}
