package article_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/service/redis_service/redis_article"
	"blogx_server/utils/jwt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ArticleLookRequest struct {
	ArticleID  uint `json:"articleID" binding:"required"`
	TimeSecond int  `json:"timeSecond"` // 读文章一共用了多久
}

func (ArticleApi) ArticleLookView(c *gin.Context) {
	cr := middleware.GetBind[ArticleLookRequest](c)

	cliams, err := jwt.ParseTokenByGin(c)
	if err != nil {
		// 未登录
		res.OkWithMsg("未登录", c)
		return
	}

	var article models.ArticleModel
	err = global.DB.Take(&article, "id = ? and status = ?", cr.ArticleID, enum.ArticleStatusPublished).Error
	if err != nil {
		res.FailWithMsg("文章不存在", c)
		return
	}

	// 这里要引入缓存，如果改用户今天看过了这篇文章，直接返回，后面不用查表了
	if redis_article.GetUserArticleHistoryCache(cr.ArticleID, cliams.UserID) {
		res.OkWithMsg("成功", c)
		logrus.Info("在缓存里")
		return
	}

	// 查询今天有没有浏览过这个文章
	var history models.UserArticleReadHistoryModel
	err = global.DB.Take(&history,
		"article_id = ? and user_id = ? and created_at > ? and created_at < ?",
		cr.ArticleID,
		cliams.UserID,
		time.Now().Format("2006-01-02")+" 00:00:00",
		time.Now().Format("2006-01-02 15:04:05")).Error
	if err == nil {
		res.OkWithMsg("成功", c)
		return
	}
	// 没找到，创建
	err = global.DB.Create(&models.UserArticleReadHistoryModel{
		ArticleID: cr.ArticleID,
		UserID:    cliams.UserID,
	}).Error
	if err != nil {
		res.FailWithMsg("失败", c)
		return
	}

	res.OkWithMsg("成功", c)
	redis_article.SetCacheLook(cr.ArticleID, true)
	redis_article.SetUserArticleHistoryCache(cr.ArticleID, cliams.UserID)
	return

}
