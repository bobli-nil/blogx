package models

import (
	"blogx_server/global"
	"blogx_server/models/enum/relationship_enum"
)

type UserFocusModel struct {
	Model
	UserID         uint      `json:"userID"`
	UserModel      UserModel `gorm:"foreignKey:UserID" json:"-"`
	FocusUserID    uint      `json:"focusUserID"`
	FocusUserModel UserModel `gorm:"foreignKey:FocusUserID" json:"-"`
}

func (UserFocusModel) TableName() string {
	return "user_focus"
}

func CalcUserRelationship(A, B uint) relationship_enum.Relation {
	var userFocusList []UserFocusModel
	global.DB.Find(&userFocusList,
		"(user_id = ? and focus_user_id = ?) or (user_id = ? and focus_user_id = ?)",
		B, A,
		A, B,
	)
	if len(userFocusList) == 0 {
		return relationship_enum.RelationStranger
	}
	if len(userFocusList) == 2 {
		return relationship_enum.RelationFriend
	}
	if userFocusList[0].FocusUserID == A {
		return relationship_enum.RelationFans
	}
	return relationship_enum.RelationFocus
}
