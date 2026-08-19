package article_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/utils/jwt"

	"github.com/gin-gonic/gin"
)

func (ArticleApi) ArticleRemoveUserView(c *gin.Context) {
	cr := middleware.GetBind[models.IDRequest](c)
	claims := jwt.GetClaims(c)

	var model models.ArticleModel
	err := global.DB.Take(&model, "id = ? and user_id = ?", cr.ID, claims.UserID).Error
	if err != nil {
		res.FailWithMsg("文章不存在", c)
		return
	}

	if err := global.DB.Delete(&model).Error; err != nil {
		res.FailWithMsg("文章删除失败", c)
		return
	}

	res.OkWithMsg("文章删除成功", c)
}
