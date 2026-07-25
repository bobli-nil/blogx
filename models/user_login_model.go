package models

type UserLoginModel struct {
	Model
	UserID    uint      `json:"userId"`
	UserModel UserModel `gorm:"foreignKey:UserID;references:ID" json:"-"`
	IP        string    `gorm:"size:32" json:"ip"`
	Addr      string    `gorm:"size:64" json:"addr"`
	UA        string    `gorm:"size:128" json:"ua"`
}

func (UserLoginModel) TableName() string {
	return "user_login"
}
