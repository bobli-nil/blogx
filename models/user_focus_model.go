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

func CalcUserPatchRelationship(A uint, BList []uint) (m map[uint]relationship_enum.Relation) {
	m = make(map[uint]relationship_enum.Relation)

	var userFocusList []UserFocusModel
	global.DB.
		Find(&userFocusList, "(user_id = ? and focus_user_id in ?) or (user_id in ? and focus_user_id = ?)",
			A, BList,
			BList, A,
		)

	for _, B := range BList {
		count := 0
		for _, model := range userFocusList {
			m[B] = relationship_enum.RelationStranger
			if model.FocusUserID == B {
				m[B] = relationship_enum.RelationFocus
				count++
			}
			if model.UserID == B {
				m[B] = relationship_enum.RelationFans
				count++
			}
		}
		if count == 2 {
			m[B] = relationship_enum.RelationFriend
		}
	}

	return
}

func CalcUserPatchRelationship2(A uint, BList []uint) (m map[uint]relationship_enum.Relation) {
	m = make(map[uint]relationship_enum.Relation)

	var userFocusList []UserFocusModel
	global.DB.
		Find(&userFocusList, "(user_id = ? and focus_user_id in ?) or (user_id in ? and focus_user_id = ?)",
			A, BList,
			BList, A,
		)

	mp := make(map[uint]UserFocusModel)
	for _, model := range userFocusList {
		mp[model.UserID] = model
	}

	for _, B := range BList {
		if mp[B].FocusUserID == A {
			m[B] = relationship_enum.RelationFans
		}
		if mp[A].FocusUserID == B {
			m[B] = relationship_enum.RelationFocus
		}
		if mp[B].FocusUserID == A && mp[A].FocusUserID == B {
			m[B] = relationship_enum.RelationFriend
		}
	}

	return
}
