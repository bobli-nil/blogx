package models

import "time"

type UserTopArticleModel struct {
	UserID       uint         `gorm:"uniqueIndex:userId_articleId" json:"userId"`
	UserModel    UserModel    `gorm:"foreignKey:UserID;references:ID" json:"-"`
	ArticleID    uint         `gorm:"uniqueIndex:userId_articleId" json:"articleId"` // 文章ID
	ArticleModel ArticleModel `gorm:"foreignKey:ArticleID;references:ID" json:"-"`
	CreatedAt    time.Time    `json:"createdAt"`
}

func (UserTopArticleModel) TableName() string {
	return "user_top_article"
}
