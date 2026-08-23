package models

import (
	"blogx_server/service/redis_service/redis_article"

	"github.com/gin-gonic/gin"
)

type UserArticleCollectModel struct {
	Model
	UserID       uint         `gorm:"uniqueIndex:articleId_collectId_userId" json:"userId"`
	UserModel    UserModel    `gorm:"foreignKey:UserID;references:ID" json:"-"`
	ArticleID    uint         `gorm:"uniqueIndex:articleId_collectId_userId" json:"articleId"` // 文章ID
	ArticleModel ArticleModel `gorm:"foreignKey:ArticleID;references:ID" json:"-"`
	CollectID    uint         `gorm:"uniqueIndex:articleId_collectId_userId" json:"collectId"`
	CollectModel CollectModel `gorm:"foreignKey:CollectID;references:ID" json:"-"`
}

func (UserArticleCollectModel) TableName() string {
	return "user_article_collect"
}

func (u UserArticleCollectModel) BeforeDelete(c *gin.Context) {
	redis_article.SetCacheCollect(u.ArticleID, false)
	return
}
