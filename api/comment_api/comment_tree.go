package comment_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/service/comment_service"

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

	var comments []models.CommentModel
	var list = make([]*comment_service.CommentResponse, 0)
	global.DB.Where("article_id = ? and parent_id is null", article.ID).Find(&comments)
	if len(comments) > 0 {
		for _, comment := range comments {
			list = append(list, comment_service.GetCommentTreeV4(comment.ID))
		}
	}

	res.OkWithData(list, c)
}
