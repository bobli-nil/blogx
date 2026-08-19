package cron_service

import (
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/service/redis_service/redis_article"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func SyncArticle() {
	lookMap := redis_article.GetAllCacheLook()
	diggMap := redis_article.GetAllCacheDigg()
	collectMap := redis_article.GetAllCacheCollect()

	var list []models.ArticleModel
	global.DB.Find(&list)

	for _, model := range list {
		lookCount := lookMap[model.ID]
		diggCount := diggMap[model.ID]
		collectCount := collectMap[model.ID]
		if lookCount == 0 && diggCount == 0 && collectCount == 0 {
			continue
		}
		err := global.DB.Model(&model).Updates(map[string]any{
			"look_count":    gorm.Expr("look_count + ?", lookCount),
			"digg_count":    gorm.Expr("digg_count + ?", diggCount),
			"collect_count": gorm.Expr("collect_count + ?", collectCount),
		}).Error
		if err != nil {
			logrus.Errorf("%d 更新失败 %s", model.ID, err)
			continue
		}
		logrus.Infof("%d %s 更新成功", model.ID, model.Title)
	}

	// TODO 在清空之前获取同步期间增量数据

	redis_article.Clear()

	// TODO 清空之后再把同步期间的增量数据写回去
}
