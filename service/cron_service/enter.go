package cron_service

import (
	"time"

	"github.com/robfig/cron/v3"
)

func Cron() {
	timezone, _ := time.LoadLocation("Asia/Shanghai")
	crontab := cron.New(cron.WithSeconds(), cron.WithLocation(timezone))
	crontab.AddFunc("0 0 2 * * *", SyncArticle)
	crontab.AddFunc("0 0 3 * * *", SyncComment)
	crontab.AddFunc("0 59 23 * * *", SyncSiteFlow)
	crontab.Start()
}
