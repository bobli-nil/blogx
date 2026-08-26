package site_msg_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum/message_type_enum"
	"blogx_server/utils/jwt"

	"github.com/gin-gonic/gin"
)

type UserMsgResponse struct {
	CommentMsgCount int `json:"commentMsgCount"`
	DiggMsgCount    int `json:"diggMsgCount"`
	PrivateMsgCount int `json:"privateMsgCount"`
	SystemMsgCount  int `json:"systemMsgCount"`
}

func (SiteMsgApi) UserMsgView(c *gin.Context) {
	claims := jwt.GetClaims(c)

	var msgList []models.MessageModel
	global.DB.Find(&msgList, "rev_user_id = ? and is_read = ?", claims.UserID, false)
	var data UserMsgResponse
	for _, model := range msgList {
		switch model.Type {
		case message_type_enum.CommentType, message_type_enum.ApplyType:
			data.CommentMsgCount++
		case message_type_enum.DiggCommentType, message_type_enum.DiggArticleType, message_type_enum.CollectArticleType:
			data.DiggMsgCount++
		case message_type_enum.SystemType:
			data.SystemMsgCount++
		}
	}

	// TODO 未读的私信

	// 算未读的全局消息
	var userReadIDList []uint
	global.DB.Model(&models.UserGlobalNotificationModel{}).
		Where("rev_user_id = ? and (is_read = ? or is_delete = ?)", claims.UserID, false, false).
		Select("notification_id").
		Scan(&userReadIDList)
	query := global.DB.Where("")
	if len(userReadIDList) > 0 {
		query.Where("id not in ?", userReadIDList)
	}

	var systemMsg []models.GlobalNotificationModel
	global.DB.Where(query).Find(&systemMsg)
	data.SystemMsgCount += len(systemMsg)

	res.OkWithData(data, c)
}
