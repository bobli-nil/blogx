package models

type CategoryModel struct {
	Model
	Title       string         `gorm:"size:32" json:"title"`
	UserID      uint           `json:"userId"`
	UserModel   UserModel      `gorm:"foreignKey:UserID;references:ID" json:"-"`
	ArticleList []ArticleModel `gorm:"foreignKey:CategoryID;references:ID" json:"-"`
}

func (CategoryModel) TableName() string {
	return "category"
}
