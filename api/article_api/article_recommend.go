package article_api

import (
	"blogx_server/common"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"

	"github.com/gin-gonic/gin"
)

type ArticleRecommendResponse struct {
	ID        uint64 `json:"id" gorm:"column:id"`
	Title     string `json:"title" gorm:"column:title"`
	LookCount int    `json:"lookCount" gorm:"column:look_count"`
}

func (ArticleApi) ArticleRecommendView(c *gin.Context) {
	cr := middleware.GetBind[common.PageInfo](c)

	list := make([]ArticleRecommendResponse, 0)
	global.DB.
		Model(&models.ArticleModel{}).
		Where("date(created_at) = date(now())").
		Order("look_count desc").
		Offset(cr.GetOffset()).
		Limit(cr.GetLimit()).
		Select("id", "title", "look_count").
		Scan(&list)

	res.OkWithData(list, c)
}
