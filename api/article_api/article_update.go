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
	"bytes"
	"fmt"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
)

type ArticleUpdateReq struct {
	ID          uint       `json:"id" binding:"required"`
	Title       string     `json:"title" binding:"required"`
	Abstract    string     `json:"abstract"`
	Content     string     `json:"content" binding:"required"`
	CategoryID  *uint      `json:"categoryID"`
	TagList     ctype.List `json:"tagList"`
	Cover       string     `json:"cover"`
	OpenComment bool       `json:"openComment"`
}

func (ArticleApi) ArticleUpdateView(c *gin.Context) {
	cr := middleware.GetBind[ArticleUpdateReq](c)

	user, err := jwt.GetClaims(c).GetUser()
	if err != nil {
		res.FailWithMsg("用户不存在", c)
		return
	}

	var article models.ArticleModel
	if err = global.DB.Take(&article, cr.ID).Error; err != nil {
		res.FailWithMsg("该文章不存在", c)
		return
	}
	if article.UserID != user.ID {
		res.FailWithMsg("只能更新自己的文章", c)
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
	contentDoc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(cr.Content)))
	if err != nil {
		res.FailWithMsg("正文解析错误", c)
		return
	}
	contentDoc.Find("script").Remove()
	contentDoc.Find("image").Remove()
	contentDoc.Find("iframe").Remove()
	contentDoc.Find("video").Remove()
	contentDoc.Find("audio").Remove()
	cr.Content = contentDoc.Text()

	// 如果不传简介，从正文中取前30个字符
	if cr.Abstract == "" {
		htmlStr := markdown.MdToHTML(cr.Content)
		doc, err := goquery.NewDocumentFromReader(bytes.NewReader([]byte(htmlStr)))
		if err != nil {
			fmt.Println(err)
			return
		}
		htmlText := doc.Text()
		cr.Abstract = string([]rune(htmlText)[:200])
	}

	// 要更新的数据
	mps := map[string]any{
		"title":        cr.Title,
		"abstract":     cr.Abstract,
		"content":      cr.Content,
		"category_id":  cr.CategoryID,
		"tag_list":     cr.TagList,
		"cover":        cr.Cover,
		"open_comment": cr.OpenComment,
	}
	if article.Status == enum.ArticleStatusPublished && !global.Conf.Site.Article.NoExamine {
		mps["status"] = enum.ArticleStatusExamined
	}

	err = global.DB.Debug().Model(&article).Updates(mps).Error
	if err != nil {
		res.FailWithMsg(err.Error(), c)
		return
	}

	res.OkWithMsg("文章更新成功", c)
}
