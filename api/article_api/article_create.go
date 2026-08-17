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

type ArticleCreateReq struct {
	Title       string             `json:"title" binding:"required"`
	Abstract    string             `json:"abstract"`
	Content     string             `json:"content" binding:"required"`
	CategoryID  *uint              `json:"categoryID"`
	TagList     ctype.List         `json:"tagList gorm:"type:varchar(255)"`
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

	err = global.DB.Create(&article).Error
	if err != nil {
		fmt.Println("---->", err)
		res.FailWithMsg("文章创建失败", c)
		return
	}
	res.OkWithMsg("文章创建成功", c)
}
