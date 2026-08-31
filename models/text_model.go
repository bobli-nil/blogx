package models

type TextModel struct {
	Model
	ArticleID uint   `json:"articleID"`
	Head      string `json:"head"`
	Body      string `json:"body"`
}

func (TextModel) TableName() string {
	return "text"
}
