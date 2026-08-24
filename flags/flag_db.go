package flags

import (
	"blogx_server/global"
	"blogx_server/models"

	"github.com/sirupsen/logrus"
)

func FlagDB() {
	err := global.DB.AutoMigrate(
		&models.UserModel{},
		&models.UserConfModel{},
		&models.ArticleModel{},
		&models.CategoryModel{},
		&models.ArticleDiggModel{},
		&models.CollectModel{},
		&models.UserArticleCollectModel{},
		&models.UserTopArticleModel{},
		&models.UserArticleReadHistoryModel{},
		&models.CommentModel{},
		&models.BannerModel{},
		&models.LogModel{},
		&models.ImageModel{},
		&models.UserLoginModel{},
		&models.UserCommentDiggModel{},
		&models.MessageModel{},
		&models.UserMessageConfModel{},
	)
	if err != nil {
		logrus.Errorf("数据库迁移失败 %s", err)
		return
	}
	logrus.Infof("数据库迁移成功")
}
