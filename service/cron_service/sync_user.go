package cron_service

import (
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/service/redis_service/redis_user"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func SyncUser() {
	lookMap := redis_user.GetAllCacheLook()

	var list []models.UserConfModel
	global.DB.Find(&list)

	for _, model := range list {
		look := lookMap[model.UserID]
		if look == 0 {
			continue
		}
		err := global.DB.Model(&model).Update("look_count", gorm.Expr("look_count + ?", look)).Error
		if err != nil {
			logrus.Errorf("更新%d look_count字段失败", model.UserID)
			continue
		}
		logrus.Infof("%d 更新成功", model.UserID)
	}

	redis_user.Clear()
}
