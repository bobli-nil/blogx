package article_api

import (
	"blogx_server/common"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/models/enum/relationship_enum"
	"blogx_server/service/focus_service"
	"blogx_server/utils/jwt"

	"github.com/gin-gonic/gin"
)

type AuthorRecommendResponse struct {
	UserID       uint   `json:"userID"`
	UserNickname string `json:"userNickname"`
	UserAvatar   string `json:"userAvatar"`
	UserAbstract string `json:"userAbstract"`
}

func (ArticleApi) AuthorRecommendView(c *gin.Context) {
	cr := middleware.GetBind[common.PageInfo](c)

	var count int64
	var userIDList []uint
	global.DB.
		Model(&models.ArticleModel{}).
		Group("user_id").
		Select("count(*) as count").
		Count(&count)
	global.DB.
		Model(&models.ArticleModel{}).
		Group("user_id").
		Offset(cr.GetOffset()).
		Limit(cr.GetLimit()).
		Select("user_id").
		Scan(&userIDList)

	claims, err := jwt.ParseTokenByGin(c)
	if err == nil && claims != nil {
		relationMap := focus_service.CalcUserPatchRelationship2(claims.UserID, userIDList)
		userIDList = []uint{}
		for userID, relation := range relationMap {
			if relation == relationship_enum.RelationStranger || relation == relationship_enum.RelationFans {
				userIDList = append(userIDList, userID)
			}
		}
	}

	var userList []models.UserModel
	global.DB.Find(&userList, "id in ?", userIDList)

	var list = make([]AuthorRecommendResponse, 0)
	for _, model := range userList {
		list = append(list, AuthorRecommendResponse{
			UserID:       model.ID,
			UserNickname: model.Nickname,
			UserAvatar:   model.Avatar,
			UserAbstract: model.Abstract,
		})
	}

	res.OkWithList(list, int(count), c)
}
