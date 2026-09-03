package article_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/service/redis_service/redis_article"
	"blogx_server/utils/jwt"

	"github.com/gin-gonic/gin"
)

type ArticleDetailResponse struct {
	models.ArticleModel
	Username      string `json:"username"`
	Nickname      string `json:"nickname"`
	Avatar        string `json:"avatar"`
	CategoryTitle string `json:"categoryTitle"`
	IsCollect     bool   `json:"isCollect"`
	IsDigg        bool   `json:"isDigg"`
}

func (ArticleApi) ArticleDetailView(c *gin.Context) {
	cr := middleware.GetBind[models.IDRequest](c)

	var article models.ArticleModel
	if err := global.DB.Preload("UserModel").Preload("CategoryModel").Take(&article, cr.ID).Error; err != nil {
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

	data := ArticleDetailResponse{
		ArticleModel: article,
		Username:     article.UserModel.Username,
		Nickname:     article.UserModel.Nickname,
		Avatar:       article.UserModel.Avatar,
	}

	if err == nil && claims != nil {
		if claims.Role == enum.UserRole {
			if article.UserID != claims.UserID && article.Status != enum.ArticleStatusPublished {
				res.FailWithMsg("文章不存在", c)
				return
			}
		}
		// 查登录用户是否收藏该文章、是否点赞该文章
		var userDiggModel models.ArticleDiggModel
		err := global.DB.Take(&userDiggModel, "user_id = ? AND article_id = ?", claims.UserID, article.ID).Error
		if err == nil {
			data.IsDigg = true
		}
		var userCollectModel models.UserArticleCollectModel
		err = global.DB.Take(&userCollectModel, "user_id = ? AND article_id = ?", claims.UserID, article.ID).Error
		if err == nil {
			data.IsCollect = true
		}
	}

	// 已登录用户，能看到自己的所有文章，只能看到已发布的文章

	data.DiggCount = article.DiggCount + redis_article.GetCacheDigg(article.ID)
	data.LookCount = article.LookCount + redis_article.GetCacheLook(article.ID)
	data.CollectCount = article.CollectCount + redis_article.GetCacheCollect(article.ID)
	data.CommentCount = article.CommentCount + redis_article.GetCacheComment(article.ID)

	// TODO 增加CategoryTitle这个功能未测试
	if article.CategoryModel != nil {
		data.CategoryTitle = article.CategoryModel.Title
	}

	res.OkWithData(data, c)
}
