package site_msg_api

import (
	"blogx_server/common"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/models/enum/message_type_enum"
	"blogx_server/utils/jwt"

	"github.com/gin-gonic/gin"
)

type SiteMsgListRequest struct {
	common.PageInfo
	T int8 `form:"t" binding:"required,oneof=1 2 3"` // 1评论和回复 2赞和收藏 3系统
}

func (SiteMsgApi) SiteMsgListView(c *gin.Context) {
	cr := middleware.GetBind[SiteMsgListRequest](c)
	claims := jwt.GetClaims(c)

	typeList := make([]message_type_enum.Type, 0)
	switch cr.T {
	case 1:
		typeList = append(typeList, message_type_enum.CommentType, message_type_enum.ApplyType)
	case 2:
		typeList = append(typeList, message_type_enum.DiggArticleType, message_type_enum.DiggCommentType, message_type_enum.CollectArticleType)
	case 3:
		typeList = append(typeList, message_type_enum.SystemType)
	}

	list, count, _ := common.ListQuery(models.MessageModel{
		RevUserID: claims.UserID,
	}, common.Options{
		Debug:    true,
		PageInfo: cr.PageInfo,
		Where:    global.DB.Where("type in ?", typeList),
	})

	res.OkWithList(list, count, c)
}
