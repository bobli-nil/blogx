package cron_service

import (
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/service/redis_service/redis_comment"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

func SyncComment() {
	// 同步评论回复数
	mps := redis_comment.GetAllCacheApply()
	keys := GetKeys(mps)
	var comments []models.CommentModel
	global.DB.Find(&comments, "id in ?", keys)
	if len(comments) > 0 {
		for _, comment := range comments {
			err := global.DB.Model(&comment).Update("apply_count", gorm.Expr("apply_count + ?", mps[comment.ID])).Error
			if err != nil {
				logrus.Error(err)
				continue
			}
			logrus.Infof("%d 评论回复数同步成功", comment.ID)
		}
	}
	// 同步评论点赞数
	mps1 := redis_comment.GetAllCacheDigg()
	keys1 := GetKeys(mps1)
	var comments1 []models.CommentModel
	global.DB.Find(&comments1, "id in ?", keys1)
	if len(comments1) > 0 {
		for _, comment := range comments1 {
			err := global.DB.Model(&comment).Update("digg_count", gorm.Expr("digg_count + ?", mps1[comment.ID])).Error
			if err != nil {
				logrus.Error(err)
				continue
			}
			logrus.Infof("%d 评论点赞数同步成功", comment.ID)
		}
	}

	redis_comment.Clear()
}

func GetKeys(mps map[uint]int) (list []uint) {
	list = make([]uint, 0)
	for key, _ := range mps {
		list = append(list, key)
	}
	return
}
