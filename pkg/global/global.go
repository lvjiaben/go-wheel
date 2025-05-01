package global

import (
	"admin/pkg/types"

	"github.com/go-redis/redis"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	CONFIG Config
	DB     *gorm.DB
	LOG    *zap.Logger
	RDB    *redis.Client
	I18N   types.I18n
	MQ     types.MessageQueue
	DQ     types.DelayQueue
	CRON   types.CronManager
)
