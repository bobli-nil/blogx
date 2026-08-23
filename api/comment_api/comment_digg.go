package comment_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/service/redis_service/redis_comment"
	"blogx_server/utils/jwt"

	"github.com/gin-gonic/gin"
)

func (CommentApi) CommentDiggView(c *gin.Context) {
	cr := middleware.GetBind[models.IDRequest](c)

	var comment models.CommentModel
	if err := global.DB.Take(&comment, cr.ID).Error; err != nil {
		res.FailWithMsg("评论不存在", c)
		return
	}

	// 查一下之前有没有点过
	claims := jwt.GetClaims(c)
	var commentDigg models.UserCommentDiggModel
	if err := global.DB.Take(&commentDigg, "comment_id = ? and user_id = ?", comment.ID, claims.UserID).Error; err != nil {
		err = global.DB.Create(&models.UserCommentDiggModel{
			UserID:    claims.UserID,
			CommentID: comment.ID,
		}).Error
		if err != nil {
			res.FailWithMsg("点赞失败", c)
			return
		}
		redis_comment.SetCacheDigg(cr.ID, 1)

		res.OkWithMsg("点赞成功", c)
		return
	}
	// 点过就删除
	global.DB.Model(&models.UserCommentDiggModel{}).Delete("comment_id = ? and user_id = ?", commentDigg.CommentID, commentDigg.UserID)
	redis_comment.SetCacheDigg(cr.ID, -1)
	res.OkWithMsg("取消点赞成功", c)
	return
}
