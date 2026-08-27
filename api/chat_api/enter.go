package chat_api

import (
	"blogx_server/common"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/utils/jwt"

	"github.com/gin-gonic/gin"
)

type ChatApi struct{}

type ChatListRequest struct {
	common.PageInfo
	SendUserID uint `form:"sendUserID"`
	RevUserID  uint `form:"revUserID" binding:"required"`
	Type       int8 `form:"type" binding:"required,oneof=1 2"` // 1用户 2管理员
}

type ChatListResponse struct {
	models.ChatModel
	SendUserNickname string `json:"sendUserNickname"`
	SendUserAvatar   string `json:"sendUserAvatar"`
	RevUserNickname  string `json:"revUserNickname"`
	RevUserAvatar    string `json:"revUserAvatar"`
	IsMe             bool   `json:"isMe"`
	IsRead           bool   `json:"isRead"` // 消息是否已读
}

func (ChatApi) ChatListView(c *gin.Context) {
	cr := middleware.GetBind[ChatListRequest](c)
	claims := jwt.GetClaims(c)
	cr.Order = "created_at desc"

	var deletedIDList []uint
	var userChatActionList []models.UserChatActionModel
	var chatReadMap = map[uint]bool{}
	global.DB.Find(&userChatActionList, "user_id = ? and (is_delete = ? or is_delete is null)", cr.RevUserID, 0)
	for _, model := range userChatActionList {
		chatReadMap[model.ChatID] = true
	}

	switch cr.Type {
	case 1:
		cr.SendUserID = claims.UserID
		global.DB.
			Model(&models.UserChatActionModel{}).
			Where("user_id = ? and is_delete = ?", claims.UserID, true).
			Select("chat_id").
			Scan(&deletedIDList)
	case 2:
		if claims.Role != enum.AdminRole {
			res.FailWithMsg("权限错误", c)
			return
		}
		if cr.SendUserID == 0 {
			res.FailWithMsg("发送人必填", c)
			return
		}
	}

	query := global.DB.
		Where("(send_user_id = ? and rev_user_id = ?) or (send_user_id = ? and rev_user_id = ?)", cr.SendUserID, cr.RevUserID, cr.RevUserID, cr.SendUserID)

	if len(deletedIDList) > 0 {
		query.Where("id not in ?", deletedIDList)
	}

	_list, count, _ := common.ListQuery(models.ChatModel{}, common.Options{
		PageInfo: cr.PageInfo,
		PreLoads: []string{"SendUserModel", "RevUserModel"},
		Where:    query,
	})

	list := make([]ChatListResponse, 0)
	for _, model := range _list {
		item := ChatListResponse{
			ChatModel:        model,
			SendUserNickname: model.SendUserModel.Nickname,
			SendUserAvatar:   model.SendUserModel.Avatar,
			RevUserNickname:  model.RevUserModel.Nickname,
			RevUserAvatar:    model.RevUserModel.Avatar,
			IsRead:           chatReadMap[model.ID],
		}
		if claims.UserID == model.SendUserID {
			item.IsMe = true
		}
		list = append(list, item)
	}

	res.OkWithList(list, count, c)
}
