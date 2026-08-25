package models

import "time"

type UserCommentDiggModel struct {
	Model
	UserID       uint         `json:"userID"`
	UserModel    UserModel    `gorm:"foreignKey:UserID;references:ID" json:"-"`
	CommentID    uint         `json:"commentID"`
	CommentModel CommentModel `gorm:"foreignKey:CommentID;references:ID" json:"-"`
	CreatedAt    time.Time    `json:"createdAt"`
}

func (UserCommentDiggModel) TableName() string {
	return "user_comment_digg"
}
