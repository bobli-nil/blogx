package article_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/ctype"
	"blogx_server/models/enum"
	"blogx_server/utils"
	"blogx_server/utils/jwt"

	"github.com/gin-gonic/gin"
)

func (ArticleApi) ArticleTagOptionsView(c *gin.Context) {
	claims := jwt.GetClaims(c)

	var articleList []models.ArticleModel
	global.DB.Model(&models.ArticleModel{}).Where("user_id = ? and Status = ?", claims.UserID, enum.ArticleStatusPublished).Find(&articleList)

	tagList := make(ctype.List, 0)
	for _, model := range articleList {
		tagList = append(tagList, model.TagList...)
	}
	tagList = utils.Unique(tagList)

	list := make([]models.OptionResponse[string], 0)
	for _, s := range tagList {
		list = append(list, models.OptionResponse[string]{
			Label: s,
			Value: s,
		})
	}

	res.OkWithData(list, c)
}
