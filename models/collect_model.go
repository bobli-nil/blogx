package models

type CollectModel struct {
	Model
	Title       string                    `gorm:"size:32" json:"title"`     // 名称
	Abstract    string                    `gorm:"size:256" json:"abstract"` // 简介
	Cover       string                    `gorm:"size:256" json:"cover"`    // 封面
	ArticleList []UserArticleCollectModel `gorm:"foreignKey:CollectID" json:"-"`
	UserID      uint                      `json:"userId"`    // 用户ID
	IsDefault   bool                      `json:"isDefault"` // 是否是默认收藏夹
	UserModel   UserModel                 `gorm:"foreignKey:UserID;references:ID" json:"-"`
}

func (CollectModel) TableName() string {
	return "collect"
}
