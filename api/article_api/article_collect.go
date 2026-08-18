package article_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/utils/jwt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ArticleCollectRequest struct {
	ArticleID uint `json:"articleID" binding:"required"`
	CollectID uint `json:"collectID"`
}

func (ArticleApi) ArticleCollectView(c *gin.Context) {
	cr := middleware.GetBind[ArticleCollectRequest](c)

	var article models.ArticleModel
	err := global.DB.Take(&article, "id = ? and status = ?", cr.ArticleID, enum.ArticleStatusPublished).Error
	if err != nil {
		res.FailWithMsg("文章不存在", c)
		return
	}

	claims := jwt.GetClaims(c)

	var collectModel models.CollectModel
	if cr.CollectID == 0 {
		// 默认收藏夹
		var err = global.DB.Take(&collectModel, "user_id = ? and is_default = ?", claims.UserID, 1).Error
		if err != nil {
			collectModel.Title = "默认收藏夹"
			collectModel.UserID = claims.UserID
			collectModel.IsDefault = true
			global.DB.Create(&collectModel)
		}
		cr.CollectID = collectModel.ID
	} else {
		// 判断该收藏夹是否存在且是自己的
		err = global.DB.Take(&collectModel, "id = ? and user_id = ?", cr.CollectID, claims.UserID).Error
		if err != nil {
			res.FailWithMsg("该收藏夹不存在", c)
			return
		}
	}

	// 判断该文章是否已收藏
	var articleCollect models.UserArticleCollectModel
	err = global.DB.Where(models.UserArticleCollectModel{
		UserID:    claims.UserID,
		ArticleID: cr.ArticleID,
		CollectID: cr.CollectID,
	}).Take(&articleCollect).Error
	if err != nil {
		// 未收藏，进行收藏
		err = global.DB.Create(&models.UserArticleCollectModel{
			UserID:    claims.UserID,
			ArticleID: cr.ArticleID,
			CollectID: cr.CollectID,
		}).Error
		if err != nil {
			res.FailWithMsg("收藏失败", c)
			return
		}
		global.DB.Model(&collectModel).Update("article_count", gorm.Expr("article_count + ?", 1))
		res.OkWithMsg("收藏成功", c)
		return
	}

	// 已收藏，进行取消收藏
	err = global.DB.Where(&models.UserArticleCollectModel{
		UserID:    claims.UserID,
		ArticleID: cr.ArticleID,
		CollectID: cr.CollectID,
	}).Delete(&articleCollect).Error
	if err != nil {
		res.FailWithMsg("取消收藏失败", c)
		return
	}
	global.DB.Model(&collectModel).Update("article_count", gorm.Expr("article_count - ?", 1))
	res.OkWithMsg("取消收藏成功", c)
	return
}
