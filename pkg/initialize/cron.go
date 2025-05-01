package initialize

import (
	"admin/pkg/cron"
)

func NewCronManager() *cron.CronManager {
	return cron.NewCronManager()
}
