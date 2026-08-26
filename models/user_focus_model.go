package models

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
