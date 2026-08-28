package chat_api

import (
	"blogx_server/common"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/utils/jwt"

	"github.com/gin-gonic/gin"
)

func (ChatApi) UserChatDeleteView(c *gin.Context) {
	cr := middleware.GetBind[models.DeleteRequest](c)
	claims := jwt.GetClaims(c)

	// 消息是否存在
	var chatList []models.ChatModel
	global.DB.Find(&chatList, "id in ?", cr.IDList)

	// 之前是否删过
	chatAcMap := common.ScanMapV2(models.UserChatActionModel{}, common.ScanOption{
		Key:   "ChatID",
		Where: global.DB.Where("user_id = ? and chat_id in ?", claims.UserID, cr.IDList),
	})

	var addChatAc []models.UserChatActionModel
	var updateChatAcIDList []uint
	for _, model := range chatList {
		chatAc, ok := chatAcMap[model.ID]
		if !ok {
			addChatAc = append(addChatAc, models.UserChatActionModel{
				UserID:   claims.UserID,
				ChatID:   model.ID,
				IsDelete: true,
			})
			continue
		}
		if chatAc.IsDelete {
			continue
		}
		updateChatAcIDList = append(updateChatAcIDList, chatAc.ID)
	}

	if len(addChatAc) > 0 {
		if err := global.DB.Create(&addChatAc).Error; err != nil {
			res.FailWithMsg("删除消息失败", c)
			return
		}
	}
	if len(updateChatAcIDList) > 0 {
		global.DB.Model(&models.UserChatActionModel{}).Where("id in ?", updateChatAcIDList).Update("is_delete", true)
	}

	res.OkWithMsg("删除消息成功", c)
}
