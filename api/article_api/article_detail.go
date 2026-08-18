package article_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/utils/jwt"

	"github.com/gin-gonic/gin"
)

type ArticleDetailResponse struct {
	models.ArticleModel
	Username string `json:"username"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
}

func (ArticleApi) ArticleDetailView(c *gin.Context) {
	cr := middleware.GetBind[models.IDRequest](c)

	var article models.ArticleModel
	if err := global.DB.Preload("UserModel").Take(&article, cr.ID).Error; err != nil {
		res.FailWithMsg("文章不存在", c)
		return
	}

	claims, err := jwt.ParseTokenByGin(c)
	if err != nil {
		// 未登录的用户，只能看已发布的文章
		if article.Status != enum.ArticleStatusPublished {
			res.FailWithMsg("文章不存在", c)
			return
		}
	}

	// 已登录用户，能看到自己的所有文章，只能看到已发布的文章
	if claims != nil && claims.Role == enum.UserRole {
		if article.UserID != claims.UserID && article.Status != enum.ArticleStatusPublished {
			res.FailWithMsg("文章不存在", c)
			return
		}
	}

	// 管理员，能看所有人的所有文章

	// TODO 从缓存里获取点赞数和浏览量
	res.OkWithData(ArticleDetailResponse{
		ArticleModel: article,
		Username:     article.UserModel.Username,
		Nickname:     article.UserModel.Nickname,
		Avatar:       article.UserModel.Avatar,
	}, c)
}
