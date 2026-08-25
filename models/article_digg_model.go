package models

import "time"

type ArticleDiggModel struct {
	Model
	UserID       uint         `gorm:"uniqueIndex:idx_name" json:"userId"`
	UserModel    UserModel    `gorm:"foreignKey:UserID;references:ID" json:"-"`
	ArticleID    uint         `gorm:"uniqueIndex:idx_name" json:"articleId"`
	ArticleModel ArticleModel `gorm:"foreignKey:ArticleID;references:ID" json:"-"`
	CreatedAt    time.Time    `json:"createdAt"`
}

func (ArticleDiggModel) TableName() string {
	return "article_digg"
}
