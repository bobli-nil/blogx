package article_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/utils/jwt"
	"fmt"

	"github.com/gin-gonic/gin"
)

func (ArticleApi) ArticleDiggView(c *gin.Context) {
	cr := middleware.GetBind[models.IDRequest](c)

	var article models.ArticleModel
	if err := global.DB.Take(&article, "id = ? and status = ?", cr.ID, enum.ArticleStatusPublished).Error; err != nil {
		res.FailWithMsg("文章不存在", c)
		return
	}

	// 查一下之前有没有点过
	claims := jwt.GetClaims(c)
	var articleDigg models.ArticleDiggModel
	if err := global.DB.Take(&articleDigg, "article_id = ? and user_id = ?", article.ID, claims.UserID).Error; err != nil {
		err = global.DB.Create(&models.ArticleDiggModel{
			UserID:    claims.UserID,
			ArticleID: article.ID,
		}).Error
		if err != nil {
			res.FailWithMsg("点赞失败", c)
			return
		}
		res.OkWithMsg("点赞成功", c)
		return
	}
	fmt.Printf("%+v\n", articleDigg)
	global.DB.Model(&models.ArticleDiggModel{}).Delete("article_id = ? and user_id = ?", articleDigg.ArticleID, articleDigg.UserID)
	res.OkWithMsg("取消点赞成功", c)
	return
}
