package chat_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/utils/jwt"

	"github.com/gin-gonic/gin"
)

func (ChatApi) UserChatDeleteView(c *gin.Context) {
	cr := middleware.GetBind[models.IDRequest](c)
	claims := jwt.GetClaims(c)

	// 消息是否存在
	var chat models.ChatModel
	if err := global.DB.Take(&chat, cr.ID).Error; err != nil {
		res.FailWithMsg("消息不存在", c)
		return
	}

	// 之前是否删过
	var chatAc models.UserChatActionModel
	err := global.DB.Take(&chatAc, "user_id = ? and chat_id = ?", claims.UserID, cr.ID).Error
	if err != nil {
		global.DB.Create(&models.UserChatActionModel{
			UserID:   claims.UserID,
			ChatID:   cr.ID,
			IsDelete: true,
		})
		res.OkWithMsg("消息删除成功", c)
		return
	}
	if chatAc.IsDelete {
		res.FailWithMsg("该消息已删除", c)
		return
	}

	global.DB.Model(&chatAc).Update("is_delete", true)
	res.OkWithMsg("消息删除成功", c)
}
