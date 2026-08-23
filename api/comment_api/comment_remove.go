package comment_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/service/comment_service"
	"blogx_server/service/redis_service/redis_comment"
	"blogx_server/utils/jwt"
	"fmt"

	"github.com/gin-gonic/gin"
)

func (CommentApi) CommentRemoveView(c *gin.Context) {
	cr := middleware.GetBind[models.IDRequest](c)
	claims := jwt.GetClaims(c)

	var comment models.CommentModel
	if err := global.DB.Preload("ArticleModel").Take(&comment, cr.ID).Error; err != nil {
		res.FailWithMsg("评论不存在", c)
		return
	}

	// 普通用户：只能删除自己发布的评论 或 自己发布的文章的评论
	// 管理员：都可以删
	if claims.Role != enum.AdminRole {
		if !(comment.UserID == claims.UserID || comment.ArticleModel.UserID == claims.UserID) {
			res.FailWithMsg("权限错误", c)
			return
		}
	}
	// 找所有子评论（包括自己）删除，找所有父评论更新回复数
	subList := comment_service.GetCommentOneDimensional(comment.ID)
	if comment.ParentID != nil {
		parentList := comment_service.GetParents(*comment.ParentID)
		for _, commentModel := range parentList {
			redis_comment.SetCacheApply(commentModel.ID, -len(subList))
		}
	}
	global.DB.Delete(&subList)

	msg := fmt.Sprintf("删除成功，共删除%d条评论", len(subList))
	res.OkWithMsg(msg, c)
}
