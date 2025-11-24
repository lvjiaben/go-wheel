package service

import (
	"time"

	"github.com/lvjiaben/go-wheel/pkg/utils/datatype"

	"github.com/gin-gonic/gin"
	"github.com/lvjiaben/go-wheel/app/backend/model/admin"
	"github.com/lvjiaben/go-wheel/app/backend/model/system"
	"github.com/lvjiaben/go-wheel/app/common/model"
	"github.com/lvjiaben/go-wheel/pkg/container"
)

type HomeService struct {
	container *container.Container
}

func NewHomeService(c *container.Container) *HomeService {
	return &HomeService{
		container: c,
	}
}

func (s *HomeService) Index(t int) gin.H {
	todayStar, todayEnd := datatype.GetQueryTime()
	var userMoney float64
	s.container.GetDB().Model(&model.User{}).Select("COALESCE(SUM(money), 0)").Scan(&userMoney)
	var userScore float64
	s.container.GetDB().Model(&model.User{}).Select("COALESCE(SUM(score), 0)").Scan(&userScore)
	var userCount int64
	s.container.GetDB().Model(&model.User{}).Count(&userCount)
	var userToday int64
	s.container.GetDB().Model(&model.User{}).Where("created_at >= ? AND created_at < ?", todayStar, todayEnd).Count(&userToday)
	var userStatus int64
	s.container.GetDB().Model(&model.User{}).Where("status = ?", 0).Count(&userStatus)
	var adminCount int64
	s.container.GetDB().Model(&admin.Admin{}).Count(&adminCount)
	var uploadCount int64
	s.container.GetDB().Model(&system.Attachment{}).Count(&uploadCount)
	var uploadTodayCount int64
	s.container.GetDB().Model(&system.Attachment{}).Where("created_at >= ? AND created_at < ?", todayStar, todayEnd).Count(&uploadTodayCount)

	var xAxis []string
	var userRegLists []int64
	var userMoneyLists []float64
	var userScoreLists []float64
	for i := 0; i <= (t - 1); i++ {
		start, end := datatype.GetDayRange(i)
		xAxis = append(xAxis, time.Unix(int64(start), 0).Format("06-01-02"))
		var _userReg int64
		s.container.GetDB().Model(&model.User{}).Where("created_at >= ? AND created_at <= ?", start, end).Count(&_userReg)
		userRegLists = append(userRegLists, _userReg)
		var _userMoneyInc float64
		var _userMoneyDec float64
		s.container.GetDB().Model(&model.UserMoneyLog{}).Where("type=1").Where("created_at >= ? AND created_at <= ?", start, end).Select("COALESCE(SUM(money), 0)").Scan(&_userMoneyInc)
		s.container.GetDB().Model(&model.UserMoneyLog{}).Where("type=0").Where("created_at >= ? AND created_at <= ?", start, end).Select("COALESCE(SUM(money), 0)").Scan(&_userMoneyDec)
		userMoneyLists = append(userMoneyLists, _userMoneyInc-_userMoneyDec)
		var _userScoreInc float64
		var _userScoreDec float64
		s.container.GetDB().Model(&model.UserScoreLog{}).Where("type=1").Where("created_at >= ? AND created_at <= ?", start, end).Select("COALESCE(SUM(score), 0)").Scan(&_userScoreInc)
		s.container.GetDB().Model(&model.UserScoreLog{}).Where("type=0").Where("created_at >= ? AND created_at <= ?", start, end).Select("COALESCE(SUM(score), 0)").Scan(&_userScoreDec)
		userScoreLists = append(userScoreLists, _userScoreInc-_userScoreDec)
	}

	return gin.H{
		"user_money":         userMoney,
		"user_score":         userScore,
		"user_count":         userCount,
		"user_today":         userToday,
		"user_status":        userStatus,
		"admin_count":        adminCount,
		"upload_count":       uploadCount,
		"upload_today_count": uploadTodayCount,
		"line_chart": gin.H{
			"xAxis": xAxis,
			"yAxis": []interface{}{
				userRegLists,
				userMoneyLists,
				userScoreLists,
			},
		},
	}
}
