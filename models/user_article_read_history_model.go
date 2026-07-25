package models

import "time"

type UserArticleReadHistoryModel struct {
	UserID       uint         `json:"userId"`
	UserModel    UserModel    `gorm:"foreignKey:UserID;references:ID" json:"-"`
	ArticleID    uint         `json:"articleId"`
	ArticleModel ArticleModel `gorm:"foreignKey:ArticleID;references:ID" json:"-"`
	CreatedAt    time.Time    `json:"createdAt"`
}

func (UserArticleReadHistoryModel) TableName() string {
	return "user_article_read_history"
}
