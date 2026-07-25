package models

import "time"

type UserArticleCollectModel struct {
	UserID       uint         `gorm:"uniqueIndex:articleId_collectId_userId" json:"userId"`
	UserModel    UserModel    `gorm:"foreignKey:UserID;references:ID" json:"-"`
	ArticleID    uint         `gorm:"uniqueIndex:articleId_collectId_userId" json:"articleId"` // 文章ID
	ArticleModel ArticleModel `gorm:"foreignKey:ArticleID;references:ID" json:"-"`
	CollectID    uint         `gorm:"uniqueIndex:articleId_collectId_userId" json:"collectId"`
	CollectModel CollectModel `gorm:"foreignKey:CollectID;references:ID" json:"-"`
	CreatedAt    time.Time    `json:"createdAt"`
}

func (UserArticleCollectModel) TableName() string {
	return "user_article_collect"
}
