package comment_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/models/enum/relationship_enum"
	"blogx_server/service/comment_service"
	"blogx_server/service/focus_service"
	"blogx_server/utils"
	"blogx_server/utils/jwt"

	"github.com/gin-gonic/gin"
)

func (CommentApi) CommentTreeView(c *gin.Context) {
	cr := middleware.GetBind[models.IDRequest](c)

	var article models.ArticleModel
	err := global.DB.Take(&article, "id = ? and status = ?", cr.ID, enum.ArticleStatusPublished).Error
	if err != nil {
		res.FailWithMsg("文章不存在", c)
		return
	}

	// TODO 登录用户是否点赞文章下的评论；文章下的评论者与登录用户关系，功能未测试
	var userDiggCommentMap = map[uint]bool{}
	var userRelationMap = map[uint]relationship_enum.Relation{}
	claims, err := jwt.ParseTokenByGin(c)
	if err == nil && claims != nil {
		var commentList []models.CommentModel
		global.DB.Find(&commentList, "article_id = ?", cr.ID)
		if len(commentList) > 0 {
			var commentIDList []uint
			var userIDList []uint
			for _, model := range commentList {
				commentIDList = append(commentIDList, model.ID)
				userIDList = append(userIDList, model.UserID)
			}
			userIDList = utils.Unique(userIDList) // 去重
			userRelationMap = focus_service.CalcUserPatchRelationship2(claims.UserID, userIDList)
			var commentDiggList []models.UserCommentDiggModel
			global.DB.Find(&commentDiggList, "user_id = ? and comment_id in ?", claims.UserID, commentIDList)
			for _, model := range commentDiggList {
				userDiggCommentMap[model.CommentID] = true
			}
		}
	}

	var comments []models.CommentModel
	var list = make([]*comment_service.CommentResponse, 0)
	global.DB.Where("article_id = ? and parent_id is null", article.ID).Find(&comments)
	if len(comments) > 0 {
		for _, comment := range comments {
			list = append(list, comment_service.GetCommentTreeV4(comment.ID, userDiggCommentMap, userRelationMap))
		}
	}

	res.OkWithData(list, c)
}
