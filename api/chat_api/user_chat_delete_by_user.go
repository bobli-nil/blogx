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

func (ChatApi) UserChatDeleteByUserView(c *gin.Context) {
	cr := middleware.GetBind[models.IDRequest](c)
	claims := jwt.GetClaims(c)

	var user models.UserModel
	if err := global.DB.Take(&user, cr.ID).Error; err != nil {
		res.FailWithMsg("用户不存在", c)
		return
	}

	// 找我和他产生了那些消息
	var chatList []models.ChatModel
	global.DB.
		Find(&chatList, "(send_user_id = ? and rev_user_id = ?) or (send_user_id = ? and rev_user_id = ?)", claims.UserID, cr.ID, cr.ID, claims.UserID)
	idList := make([]uint, 0)
	for _, model := range chatList {
		idList = append(idList, model.ID)
	}

	// TODO 这里的逻辑未测试
	chatAcMap := common.ScanMapV2(models.UserChatActionModel{}, common.ScanOption{
		// TODO 这里的逻辑和视频里不一样
		Where: global.DB.Where("(user_id = ? or user_id = ?) and chat_id in ?", claims.UserID, cr.ID, idList),
		Key:   "ChatID",
	})

	var addChatAc []models.UserChatActionModel
	var updateChatAcIDList []uint
	for _, model := range chatList {
		chatAc, ok := chatAcMap[model.ID]
		if !ok {
			addChatAc = append(addChatAc, models.UserChatActionModel{
				UserID:   model.SendUserID,
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
