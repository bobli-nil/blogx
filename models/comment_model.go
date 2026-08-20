package models

type CommentModel struct {
	Model
	UserID         uint            `json:"userId"`
	UserModel      UserModel       `gorm:"foreignKey:UserID;references:ID" json:"-"`
	ArticleID      uint            `json:"articleId"`
	ArticleModel   ArticleModel    `gorm:"foreignKey:ArticleID;references:ID" json:"-"`
	Content        string          `gorm:"size:256" json:"content"`
	ParentID       *uint           `json:"parentId"` // 父级评论，可能没有，所以是指针
	ParentModel    *CommentModel   `gorm:"foreignKey:ParentID;references:ID" json:"-"`
	RootParentID   *uint           `json:"rootParentId"` // 根评论
	SubCommentList []*CommentModel `gorm:"foreignKey:ParentID;references:ID" json:"-"`
	DiggCount      int             `json:"diggCount"` // 评论点赞数
}

func (CommentModel) TableName() string {
	return "comment"
}
