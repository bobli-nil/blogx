package chat_service

import (
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/ctype"
	"blogx_server/models/enum/chat_msg_type"
	"blogx_server/utils/xss"

	"github.com/sirupsen/logrus"
)

// ToChat A给B发消息
func ToChat(A, B uint, msgType chat_msg_type.MsgType, msg ctype.ChatMsg) {
	err := global.DB.Create(&models.ChatModel{
		SendUserID: A,
		RevUserID:  B,
		MsgType:    msgType,
		Msg:        msg,
	}).Error
	if err != nil {
		logrus.Errorf("创建对话失败 %s", err)
	}
}

func ToTextChat(A, B uint, content string) {
	ToChat(A, B, chat_msg_type.TextMsgType, ctype.ChatMsg{
		TextMsg: &ctype.TextMsg{
			Content: content,
		},
	})
}

func ToImageChat(A, B uint, src string) {
	ToChat(A, B, chat_msg_type.ImageMsgType, ctype.ChatMsg{
		ImageMsg: &ctype.ImageMsg{
			Src: src,
		},
	})
}

func ToMarkdownChat(A, B uint, content string) {
	filterContent, err := xss.XssFilter(content)
	if err != nil {
		logrus.Errorf("xss过滤markdown内容失败 %s", err)
	}
	ToChat(A, B, chat_msg_type.MarkdownMsgType, ctype.ChatMsg{
		MarkdownMsg: &ctype.MarkdownMsg{
			Content: filterContent,
		},
	})
}
