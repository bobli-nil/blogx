package site_msg_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/models/enum/message_type_enum"
	"blogx_server/utils/jwt"
	"fmt"

	"github.com/gin-gonic/gin"
)

type SiteMsgReadRequest struct {
	ID uint `json:"id"`
	T  int8 `json:"t"` // 1评论和回复 2赞和收藏 3系统  指定类型一键已读
}

func (SiteMsgApi) SiteMsgReadView(c *gin.Context) {
	cr := middleware.GetBind[SiteMsgReadRequest](c)
	claims := jwt.GetClaims(c)

	// 单个
	if cr.ID != 0 {
		var msg models.MessageModel
		err := global.DB.Take(&msg, "id = ? and rev_user_id = ?", cr.ID, claims.UserID).Error
		if err != nil {
			res.FailWithMsg("消息不存在", c)
			return
		}
		if msg.IsRead {
			res.FailWithMsg("该消息已读", c)
			return
		}
		global.DB.Model(&msg).Update("is_read", true)

		res.OkWithMsg("消息读取成功", c)
		return
	}

	if cr.T == 0 {
		res.FailWithMsg("需传入参数t", c)
		return
	}
	typeList := make([]message_type_enum.Type, 0)
	switch cr.T {
	case 1:
		typeList = append(typeList, message_type_enum.CommentType, message_type_enum.ApplyType)
	case 2:
		typeList = append(typeList, message_type_enum.DiggArticleType, message_type_enum.DiggCommentType, message_type_enum.CollectArticleType)
	case 3:
		typeList = append(typeList, message_type_enum.SystemType)
	}
	var msgList []models.MessageModel
	global.DB.Model(&models.MessageModel{}).Find(&msgList, "type in ? and rev_user_id = ? and is_read = ?", typeList, claims.UserID, false)
	if len(msgList) > 0 {
		global.DB.Model(&msgList).Update("is_read", true)
	}
	msg := fmt.Sprintf("批量读取%d条消息成功", len(msgList))
	res.OkWithData(msg, c)
}
