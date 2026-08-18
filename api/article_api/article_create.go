package article_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/models/ctype"
	"blogx_server/models/enum"
	"blogx_server/utils/jwt"
	"blogx_server/utils/markdown"
	"blogx_server/utils/xss"
	"fmt"

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

	// 判断文章分类是不是自己创建的
	var category models.CategoryModel
	err = global.DB.Take(&category, "id = ? and user_id = ?", cr.CategoryID, user.ID).Error
	if err != nil {
		res.FailWithMsg("文章分类不存在", c)
		return
	}

	// 文章正文防止XSS攻击
	newContent, err := xss.XssFilter(cr.Content)
	if err != nil {
		res.FailWithMsg(err.Error(), c)
		return
	}
	cr.Content = newContent

	// 如果不传简介，从正文中取前30个字符
	if cr.Abstract == "" {
		abs, err := markdown.ExtractContent(cr.Content, 200)
		if err != nil {
			res.FailWithMsg(err.Error(), c)
			return
		}
		cr.Abstract = abs
	}

	// 正文内容图片转存

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
	if cr.Status == enum.ArticleStatusExamined && global.Conf.Site.Article.NoExamine {
		article.Status = enum.ArticleStatusPublished
	}

	err = global.DB.Debug().Create(&article).Error
	if err != nil {
		fmt.Println("---->", err)
		res.FailWithMsg("文章创建失败", c)
		return
	}
	res.OkWithMsg("文章创建成功", c)
}
