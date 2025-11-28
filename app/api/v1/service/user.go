package service

import (
	"errors"
	"time"

	"github.com/lvjiaben/go-wheel/app/api/v1/validate"
	"github.com/lvjiaben/go-wheel/app/common/model"
	commonService "github.com/lvjiaben/go-wheel/app/common/service"
	"github.com/lvjiaben/go-wheel/pkg/container"
	"github.com/lvjiaben/go-wheel/pkg/jwt"
	"github.com/lvjiaben/go-wheel/pkg/utils/crypto"
	"github.com/lvjiaben/go-wheel/pkg/utils/datatype"
)

type UserService struct {
	container  *container.Container
	smsService *commonService.SmsService
}

func NewUserService(c *container.Container) *UserService {
	return &UserService{
		container:  c,
		smsService: commonService.NewSmsService(c),
	}
}

// LoginResponse 登录响应
type LoginResponse struct {
	Id          int    `json:"id"`
	AccessToken string `json:"accessToken"`
	Expires     string `json:"expires"`
	Avatar      string `json:"avatar"`
	Username    string `json:"username"`
	Mobile      string `json:"mobile"`
}

// Login 用户登录
func (s *UserService) Login(req *validate.LoginRequest) (*LoginResponse, error) {
	var user model.User
	// 查找用户（支持用户名或手机号登录）
	if err := s.container.GetDB().Where("(username = ? OR mobile = ?) AND status = 1", req.Username, req.Username).First(&user).Error; err != nil {
		return nil, errors.New("user_not_found")
	}

	// 验证密码
	if !crypto.PasswordVerifyWithSalt(req.Password, user.Salt, user.Password) {
		return nil, errors.New("password_incorrect")
	}

	// 生成JWT令牌
	token, err := jwt.GenerateToken(user.Id, user.Username,
		s.container.GetConfig().GetString("jwt.secret"),
		s.container.GetConfig().GetInt("jwt.expire_day"))
	if err != nil {
		return nil, errors.New("generate_token_failed")
	}

	// 如果启用单点登录，更新token
	if s.container.GetConfig().GetBool("api.login_sso") {
		s.container.GetDB().Model(&model.User{}).Where("id = ?", user.Id).Update("token", token)
	}

	return &LoginResponse{
		Id:          user.Id,
		AccessToken: token,
		Expires:     time.Now().AddDate(0, 0, s.container.GetConfig().GetInt("jwt.expire_day")).Format("2006/01/02 15:04:05"),
		Avatar:      user.Avatar,
		Username:    user.Username,
		Mobile:      user.Mobile,
	}, nil
}

// MobileLogin 手机号登录
func (s *UserService) MobileLogin(req *validate.MobileLoginRequest) (*LoginResponse, error) {
	// 验证手机验证码
	if !s.smsService.Verify(req.Mobile, req.Code, "login") {
		return nil, errors.New("sms_code_invalid")
	}

	// 验证通过后删除验证码
	s.smsService.Delete(req.Mobile, "login")

	var user model.User
	// 查找用户
	if err := s.container.GetDB().Where("mobile = ? AND status = 1", req.Mobile).First(&user).Error; err != nil {
		return nil, errors.New("user_not_found")
	}

	// 生成JWT令牌
	token, err := jwt.GenerateToken(user.Id, user.Username,
		s.container.GetConfig().GetString("jwt.secret"),
		s.container.GetConfig().GetInt("jwt.expire_day"))
	if err != nil {
		return nil, errors.New("generate_token_failed")
	}

	// 如果启用单点登录，更新token
	if s.container.GetConfig().GetBool("api.login_sso") {
		s.container.GetDB().Model(&model.User{}).Where("id = ?", user.Id).Update("token", token)
	}

	return &LoginResponse{
		Id:          user.Id,
		AccessToken: token,
		Expires:     time.Now().AddDate(0, 0, s.container.GetConfig().GetInt("jwt.expire_day")).Format("2006/01/02 15:04:05"),
		Avatar:      user.Avatar,
		Username:    user.Username,
		Mobile:      user.Mobile,
	}, nil
}

// Reg 用户注册
func (s *UserService) Reg(req *validate.RegRequest) (*LoginResponse, error) {
	// 验证手机验证码
	if !s.smsService.Verify(req.Mobile, req.Code, "register") {
		return nil, errors.New("sms_code_invalid")
	}

	// 验证通过后删除验证码
	s.smsService.Delete(req.Mobile, "register")

	// 检查手机号是否已注册
	var count int64
	s.container.GetDB().Model(&model.User{}).Where("mobile = ?", req.Mobile).Count(&count)
	if count > 0 {
		return nil, errors.New("mobile_exists")
	}

	// 处理邀请码
	var pid, tid int
	if req.InviteCode != "" {
		var inviter model.User
		if err := s.container.GetDB().Where("code = ?", req.InviteCode).First(&inviter).Error; err == nil {
			pid = inviter.Id
			if inviter.Tid > 0 {
				tid = inviter.Tid
			} else {
				tid = inviter.Id
			}
		}
	}

	// 生成盐值和密码
	salt, err := crypto.GenerateSalt()
	if err != nil {
		return nil, errors.New("generate_salt_failed")
	}
	hashedPassword, err := crypto.PasswordHashWithSalt(req.Password, salt)
	if err != nil {
		return nil, errors.New("hash_password_failed")
	}

	// 生成邀请码
	inviteCode := datatype.GenerateRandomString(8)

	// 创建用户
	user := model.User{
		Mobile:   req.Mobile,
		Username: req.Mobile,
		Password: hashedPassword,
		Salt:     salt,
		Code:     inviteCode,
		Pid:      pid,
		Tid:      tid,
		Status:   1,
	}

	if err := s.container.GetDB().Create(&user).Error; err != nil {
		return nil, errors.New("create_user_failed")
	}

	// 生成JWT令牌
	token, err := jwt.GenerateToken(user.Id, user.Username,
		s.container.GetConfig().GetString("jwt.secret"),
		s.container.GetConfig().GetInt("jwt.expire_day"))
	if err != nil {
		return nil, errors.New("generate_token_failed")
	}

	// 更新token到数据库
	s.container.GetDB().Model(&model.User{}).Where("id = ?", user.Id).Update("token", token)

	return &LoginResponse{
		Id:          user.Id,
		AccessToken: token,
		Expires:     time.Now().AddDate(0, 0, s.container.GetConfig().GetInt("jwt.expire_day")).Format("2006/01/02 15:04:05"),
		Avatar:      user.Avatar,
		Username:    user.Username,
		Mobile:      user.Mobile,
	}, nil
}

// Logout 退出登录（拉黑token）
func (s *UserService) Logout(userId int, token string) error {
	// 将token加入黑名单
	ctx := s.container.GetContext()
	redisKey := "token_blacklist:" + token
	expireDays := s.container.GetConfig().GetInt("jwt.expire_day")
	if expireDays <= 0 {
		expireDays = 7
	}
	expiration := time.Duration(expireDays) * 24 * time.Hour

	if err := s.container.GetRDB().Set(ctx, redisKey, "1", expiration).Err(); err != nil {
		return errors.New("logout_failed")
	}

	// 清除数据库中的token
	s.container.GetDB().Model(&model.User{}).Where("id = ?", userId).Update("token", "")

	return nil
}

// ResetPwd 重置密码
func (s *UserService) ResetPwd(req *validate.ResetPwdRequest) error {
	// 验证手机验证码
	if !s.smsService.Verify(req.Mobile, req.Code, "resetpwd") {
		return errors.New("sms_code_invalid")
	}

	// 验证通过后删除验证码
	s.smsService.Delete(req.Mobile, "resetpwd")

	// 查找用户
	var user model.User
	if err := s.container.GetDB().Where("mobile = ?", req.Mobile).First(&user).Error; err != nil {
		return errors.New("user_not_found")
	}

	// 生成新的盐值和密码
	salt, err := crypto.GenerateSalt()
	if err != nil {
		return errors.New("generate_salt_failed")
	}
	hashedPassword, err := crypto.PasswordHashWithSalt(req.NewPassword, salt)
	if err != nil {
		return errors.New("hash_password_failed")
	}

	// 更新密码
	if err := s.container.GetDB().Model(&model.User{}).Where("id = ?", user.Id).Updates(map[string]any{
		"password": hashedPassword,
		"salt":     salt,
		"token":    "", // 清除token，强制重新登录
	}).Error; err != nil {
		return errors.New("update_password_failed")
	}

	return nil
}

// ChangeMobile 修改手机号
func (s *UserService) ChangeMobile(userId int, req *validate.ChangeMobileRequest) error {
	// 验证手机验证码（验证新手机号）
	if !s.smsService.Verify(req.NewMobile, req.Code, "changemobile") {
		return errors.New("sms_code_invalid")
	}

	// 验证通过后删除验证码
	s.smsService.Delete(req.NewMobile, "changemobile")

	// 验证当前手机号是否属于该用户
	var user model.User
	if err := s.container.GetDB().Where("id = ? AND mobile = ?", userId, req.Mobile).First(&user).Error; err != nil {
		return errors.New("mobile_not_match")
	}

	// 检查新手机号是否已被注册
	var count int64
	s.container.GetDB().Model(&model.User{}).Where("mobile = ? AND id != ?", req.NewMobile, userId).Count(&count)
	if count > 0 {
		return errors.New("new_mobile_exists")
	}

	// 更新手机号
	if err := s.container.GetDB().Model(&model.User{}).Where("id = ?", userId).Update("mobile", req.NewMobile).Error; err != nil {
		return errors.New("update_mobile_failed")
	}

	return nil
}

// GetById 根据ID获取用户
func (s *UserService) GetById(id int) (*model.User, error) {
	var user model.User
	if err := s.container.GetDB().Where("id = ?", id).First(&user).Error; err != nil {
		return nil, errors.New("user_not_found")
	}
	return &user, nil
}

// ChangePwd 修改密码
func (s *UserService) ChangePwd(userId int, req *validate.ChangePwdRequest) error {
	// 查找用户
	var user model.User
	if err := s.container.GetDB().Where("id = ?", userId).First(&user).Error; err != nil {
		return errors.New("user_not_found")
	}

	// 验证旧密码
	if !crypto.PasswordVerifyWithSalt(req.OldPassword, user.Salt, user.Password) {
		return errors.New("old_password_incorrect")
	}

	// 生成新的盐值和密码
	salt, err := crypto.GenerateSalt()
	if err != nil {
		return errors.New("generate_salt_failed")
	}
	hashedPassword, err := crypto.PasswordHashWithSalt(req.NewPassword, salt)
	if err != nil {
		return errors.New("hash_password_failed")
	}

	// 更新密码
	if err := s.container.GetDB().Model(&model.User{}).Where("id = ?", userId).Updates(map[string]any{
		"password": hashedPassword,
		"salt":     salt,
	}).Error; err != nil {
		return errors.New("update_password_failed")
	}

	return nil
}
