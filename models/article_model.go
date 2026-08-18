package models

import (
	"blogx_server/models/ctype"
	"blogx_server/models/enum"
	_ "embed"
)

type ArticleModel struct {
	Model
	Title        string             `gorm:"size:32" json:"title"`
	Abstract     string             `gorm:"size:256" json:"abstract"`
	Content      string             `json:"content,omitempty"`
	CategoryID   *uint              `json:"categoryId"`
	TagList      ctype.List         `gorm:"type:longtext" json:"tagList"`
	Cover        string             `gorm:"size:256" json:"cover"`
	UserID       uint               `json:"userId"`
	UserModel    UserModel          `gorm:"foreignKey:UserID;references:ID" json:"-"`
	LookCount    int                `json:"lookCount"`
	DiggCount    int                `json:"diggCount"`
	CommentCount int                `json:"commentCount"`
	CollectCount int                `json:"collectCount"`
	OpenComment  bool               `json:"openComment"`
	Status       enum.ArticleStatus `json:"status"` // 状态：草稿 审核中 已发布
}

func (ArticleModel) TableName() string {
	return "article"
}

//go:embed mappings/article_mapping.json
var articleMapping string

func (ArticleModel) Mapping() string {
	return articleMapping
}

func (ArticleModel) Index() string {
	return "article_index"
}
