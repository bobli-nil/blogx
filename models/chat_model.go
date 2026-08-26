package models

import (
	"blogx_server/models/ctype"
	"blogx_server/models/enum/chat_msg_type"
)

type ChatModel struct {
	Model
	SendUserID    uint                  `json:"sendUserID"`
	SendUserModel UserModel             `gorm:"foreignKey:SendUserID" json:"-"`
	RevUserID     uint                  `json:"revUserID"`
	RevUserModel  UserModel             `gorm:"foreignKey:RevUserID" json:"-"`
	MsgType       chat_msg_type.MsgType `json:"msgType"`
	Msg           ctype.ChatMsg         `gorm:"type:longtext;serializer:json" json:"msg"`
}

func (ChatModel) TableName() string {
	return "chat"
}
