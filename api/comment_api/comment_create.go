package comment_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/service/comment_service"
	"blogx_server/service/redis_service/redis_article"
	"blogx_server/service/redis_service/redis_comment"
	"blogx_server/utils/jwt"

	"github.com/gin-gonic/gin"
)

type CommentCreateRequest struct {
	ArticleID uint   `json:"articleID" binding:"required"`
	Content   string `json:"content" binding:"required"`
	ParentID  *uint  `json:"parentID"`
}

func (CommentApi) CommentCreateView(c *gin.Context) {
	cr := middleware.GetBind[CommentCreateRequest](c)
	claims := jwt.GetClaims(c)

	var article models.ArticleModel
	err := global.DB.Take(&article, "id = ? and status = ?", cr.ArticleID, enum.ArticleStatusPublished).Error
	if err != nil {
		res.FailWithMsg("文章不存在", c)
		return
	}

	model := models.CommentModel{
		ArticleID: cr.ArticleID,
		Content:   cr.Content,
		UserID:    claims.UserID,
		ParentID:  cr.ParentID,
	}

	if cr.ParentID != nil {
		parentList := comment_service.GetParents(*cr.ParentID)
		if len(parentList) >= global.Conf.Site.Article.CommentLine {
			res.FailWithMsg("评论层级达到限制", c)
			return
		}

		if len(parentList) > 0 {
			model.RootParentID = parentList[len(parentList)-1].ParentID
			for _, commentModel := range parentList {
				redis_comment.SetCacheApply(commentModel.ID, 1)
			}
		}

	}

	err = global.DB.Create(&model).Error
	if err != nil {
		res.FailWithMsg("发布评论失败", c)
		return
	}
	redis_article.SetCacheComment(cr.ArticleID, 1)
	res.OkWithMsg("发布评论成功", c)
}
