package cron_service

import (
	"time"

	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
)

func Cron() {
	timezone, _ := time.LoadLocation("Asia/Shanghai")
	crontab := cron.New(cron.WithSeconds(), cron.WithLocation(timezone))
	_, err := crontab.AddFunc("* * 2 * * *", SyncArticle)
	if err != nil {
		logrus.Fatal(err)
	}
	crontab.Start()
}
