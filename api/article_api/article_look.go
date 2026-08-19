package article_api

import (
	"blogx_server/common"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/service/redis_service/redis_article"
	"blogx_server/utils/jwt"
	"fmt"
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

type ArticleLookListRequest struct {
	common.PageInfo
	UserID uint `form:"userID"`
	Type   int8 `form:"type" binding:"required,oneof=1 2"`
}

type ArticleLookListResponse struct {
	ID        uint       `json:"id"`
	CreatedAt *time.Time `json:"createAt"`
	ArticleID uint       `json:"articleID"`
	Title     string     `json:"title"`
	Cover     string     `json:"cover"`
	Nickname  string     `json:"nickname"`
	Avatar    string     `json:"avatar"`
	UserID    uint       `json:"userID"`
}

func (ArticleApi) ArticleLookListView(c *gin.Context) {
	cr := middleware.GetBind[ArticleLookListRequest](c)
	claims := jwt.GetClaims(c)

	switch cr.Type {
	case 1:
		cr.UserID = claims.UserID
	}

	_list, count, _ := common.ListQuery(models.UserArticleReadHistoryModel{
		UserID: cr.UserID,
	}, common.Options{
		Debug:    true,
		PageInfo: cr.PageInfo,
		PreLoads: []string{"UserModel", "ArticleModel"},
	})

	var list = make([]ArticleLookListResponse, 0)
	for _, model := range _list {
		list = append(list, ArticleLookListResponse{
			ID:        model.ID,
			CreatedAt: model.CreatedAt,
			ArticleID: model.ArticleID,
			Title:     model.ArticleModel.Title,
			Cover:     model.ArticleModel.Cover,
			Nickname:  model.UserModel.Nickname,
			Avatar:    model.UserModel.Avatar,
			UserID:    model.UserModel.ID,
		})
	}

	res.OkWithList(list, count, c)
}

func (ArticleApi) ArticleLookRemoveView(c *gin.Context) {
	cr := middleware.GetBind[models.DeleteRequest](c)
	claims := jwt.GetClaims(c)

	var list []models.UserArticleReadHistoryModel
	global.DB.Find(&list, "user_id = ? and article_id in ?", claims.UserID, cr.IDList)
	if len(list) > 0 {
		if err := global.DB.Delete(&list).Error; err != nil {
			res.FailWithMsg("足迹删除失败", c)
			return
		}
	}

	res.OkWithMsg(fmt.Sprintf("足迹删除成功 删除%d条", len(list)), c)
}
