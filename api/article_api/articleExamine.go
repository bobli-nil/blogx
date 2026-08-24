package article_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/service/message_service"
	"fmt"

	"github.com/gin-gonic/gin"
)

type ArticleExamineRequest struct {
	ArticleID uint               `json:"articleID" binding:"required"`
	Status    enum.ArticleStatus `json:"status" binding:"required,oneof=3 4"`
	Msg       string             `json:"msg"`
}

func (ArticleApi) ArticleExamineView(c *gin.Context) {
	cr := middleware.GetBind[ArticleExamineRequest](c)

	var article models.ArticleModel
	if err := global.DB.Take(&article, cr.ArticleID).Error; err != nil {
		res.FailWithMsg("文章不存在", c)
		return
	}

	global.DB.Model(&article).Update("status", cr.Status)

	switch cr.Status {
	case 3:
		// TODO 这里的跳转地址看前端
		message_service.InsertSystemMessage(article.UserID, "管理员审核了你的文章", "审核通过", article.Title, "")
	case 4:
		message_service.InsertSystemMessage(article.UserID, "管理员审核了你的文章", fmt.Sprintf("审核不通过，原因：%s", cr.Msg), article.Title, "")
	}

	res.OkWithData(article, c)
}
