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

type ChatApi struct{}

type ChatListRequest struct {
	common.PageInfo
	UserID uint `form:"userID" binding:"required"`
}

type ChatListResponse struct {
	models.ChatModel
	SendUserNickname string `json:"sendUserNickname"`
	SendUserAvatar   string `json:"sendUserAvatar"`
	RevUserNickname  string `json:"revUserNickname"`
	RevUserAvatar    string `json:"revUserAvatar"`
	IsMe             bool   `json:"isMe"`
}

func (ChatApi) ChatListView(c *gin.Context) {
	cr := middleware.GetBind[ChatListRequest](c)
	claims := jwt.GetClaims(c)
	cr.Order = "created_at desc"

	query := global.DB.
		Where("(send_user_id = ? and rev_user_id = ?) or (send_user_id = ? and rev_user_id = ?)", cr.UserID, claims.UserID, claims.UserID, cr.UserID)

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
		}
		if claims.UserID == model.SendUserID {
			item.IsMe = true
		}
		list = append(list, item)
	}

	res.OkWithList(list, count, c)
}
