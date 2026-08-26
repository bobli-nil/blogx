package models

type UserChatActionModel struct {
	Model
	ChatID   uint `json:"chatID"`
	UserID   uint `json:"userID"`
	IsRead   bool `json:"isRead"`
	IsDelete bool `json:"isDelete"`
}

func (UserChatActionModel) TableName() string {
	return "user_chat_action"
}
