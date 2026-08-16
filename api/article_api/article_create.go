package article_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/models/ctype"
	"blogx_server/models/enum"
	"blogx_server/utils/jwt"

	"github.com/gin-gonic/gin"
)

type ArticleCreateReq struct {
	Title       string             `json:"title" binding:"required"`
	Abstract    string             `json:"abstract"`
	Content     string             `json:"content" binding:"required"`
	CategoryID  *uint              `json:"categoryID"`
	TagList     ctype.List         `json:"tagList"`
	Cover       string             `json:"cover"`
	OpenComment bool               `json:"openComment"`
	Status      enum.ArticleStatus `json:"status" binding:"oneof=1 2"`
}

func (ArticleApi) ArticleCreateView(c *gin.Context) {
	cr := middleware.GetBind[ArticleCreateReq](c)

	user, err := jwt.GetClaims(c).GetUser()
	if err != nil {
		res.FailWithMsg("用户不存在", c)
		return
	}

	article := models.ArticleModel{
		Title:       cr.Title,
		Abstract:    cr.Abstract,
		Content:     cr.Content,
		CategoryID:  cr.CategoryID,
		TagList:     cr.TagList,
		Cover:       cr.Cover,
		OpenComment: cr.OpenComment,
		UserID:      user.ID,
		Status:      cr.Status,
	}
	if global.Conf.Site.Article.NoExamine {
		article.Status = enum.ArticleStatusPublished
	}

	err = global.DB.Create(&article).Error
	if err != nil {
		res.FailWithMsg("文章创建失败", c)
		return
	}
	res.OkWithMsg("文章创建成功", c)
}
