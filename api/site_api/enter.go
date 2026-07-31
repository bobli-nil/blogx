package site_api

import (
	"blogx_server/common/res"
	"blogx_server/conf"
	"blogx_server/core"
	"blogx_server/global"
	"blogx_server/models/enum"
	"blogx_server/utils/jwt"
	"errors"
	"fmt"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type SiteApi struct{}

type SiteInfoRequest struct {
	Name string `uri:"name"`
}

func (s *SiteApi) SiteInfoView(c *gin.Context) {
	cr := SiteInfoRequest{}
	err := c.ShouldBindUri(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	if cr.Name == "site" {
		global.Conf.Site.About.Version = global.Version
		res.OkWithData(global.Conf.Site, c)
		return
	}

	value, ok := c.Get("claims")
	claims, assertOK := value.(*jwt.MyClaims)
	if !assertOK || !ok || claims.Role != enum.AdminRole {
		res.FailWithMsg("没有权限", c)
		return
	}

	var data any

	switch cr.Name {
	case "email":
		rep := global.Conf.Email
		rep.AuthCode = "******"
		data = rep
	case "qq":
		rep := global.Conf.QQ
		rep.AppKey = "******"
		data = rep
	case "qiNiu":
		rep := global.Conf.QiNiu
		rep.SecretKey = "******"
		data = rep
	case "ai":
		rep := global.Conf.Ai
		rep.SecretKey = "******"
		data = rep
	default:
		res.FailWithMsg("不存在配置", c)
		return
	}

	res.OkWithData(data, c)
}

func (s *SiteApi) SiteInfoQQView(c *gin.Context) {
	res.OkWithData(global.Conf.QQ.Url(), c)
}

type SiteUpdateReq struct {
	Name string `json:"name" binding:"required" label:"名称"`
}

func (s *SiteApi) SiteUpdateView(c *gin.Context) {
	var cr SiteInfoRequest
	err := c.ShouldBindUri(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	var rep any
	switch cr.Name {
	case "site":
		var data conf.Site
		err = c.ShouldBindJSON(&data)
		rep = data
	case "email":
		var data conf.Email
		err = c.ShouldBindJSON(&data)
		rep = data
	case "qq":
		var data conf.QQ
		err = c.ShouldBindJSON(&data)
		rep = data
	case "qiNiu":
		var data conf.QiNiu
		err = c.ShouldBindJSON(&data)
		rep = data
	case "ai":
		var data conf.Ai
		err = c.ShouldBindJSON(&data)
		rep = data
	default:
		res.FailWithMsg("不存在的配置", c)
		return
	}

	if err != nil {
		res.FailWithError(err, c)
		return
	}

	switch s := rep.(type) {
	case conf.Site:
		err := UpdateSite(s)
		if err != nil {
			res.FailWithError(err, c)
			return
		}
		global.Conf.Site = s
	case conf.Email:
		if s.AuthCode == "******" {
			s.AuthCode = global.Conf.Email.AuthCode
		}
		global.Conf.Email = s
	case conf.QQ:
		if s.AppKey == "******" {
			s.AppKey = global.Conf.QQ.AppKey
		}
		global.Conf.QQ = s
	case conf.QiNiu:
		if s.SecretKey == "******" {
			s.SecretKey = global.Conf.QiNiu.SecretKey
		}
		global.Conf.QiNiu = s
	case conf.Ai:
		if s.SecretKey == "******" {
			s.SecretKey = global.Conf.Ai.SecretKey
		}
		global.Conf.Ai = s
	default:
		res.FailWithMsg("没有这个分类", c)
		return
	}

	core.SetConf()

	res.OkWithMsg("更新成功", c)

}

// TODO 未测试
func UpdateSite(Site conf.Site) error {
	project := Site.Project
	seo := Site.Seo
	if project.Icon == "" && project.Title == "" && project.WebPath == "" && seo.Keywords == "" && seo.Description == "" {
		return errors.New("项目和SEO相关配置不能为空")
	}
	if project.WebPath == "" {
		return errors.New("前端地址不能为空")
	}
	file, err := os.Open(project.WebPath)
	defer file.Close()
	if err != nil {
		return errors.New(fmt.Sprintf("%s 文件不存在", project.WebPath))
	}
	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		return errors.New(fmt.Sprintf("%s 文件解析失败", project.WebPath))
	}

	if project.Title != "" {
		doc.Find("title").SetText(project.Title)
	}
	if project.Icon != "" {
		iconSelection := doc.Find(`link[rel="icon"]`)
		fmt.Println(iconSelection.Length())
		if iconSelection.Length() > 0 {
			iconSelection.SetAttr("href", project.Icon)
		} else {
			doc.Find("head").AppendHtml(fmt.Sprintf("<link rel=\"icon\" href=\"%s\">", project.Icon))
		}
	}
	if seo.Keywords != "" {
		keywordSelection := doc.Find("meta[name='keywords']")
		if keywordSelection.Length() > 0 {
			keywordSelection.SetAttr("content", seo.Keywords)
		} else {
			doc.Find("head").AppendHtml(fmt.Sprintf("<meta name=\"keywords\" content=\"%s\">", seo.Keywords))
		}
	}
	if seo.Description != "" {
		descriptionSelection := doc.Find("meta[name='description']")
		if descriptionSelection.Length() > 0 {
			descriptionSelection.SetAttr("content", seo.Description)
		} else {
			doc.Find("head").AppendHtml(fmt.Sprintf("<meta name=\"description\" content=\"%s\">", seo.Description))
		}
	}

	htmlStr, err := doc.Html()
	if err != nil {
		return errors.New("转化失败")
	}

	err = os.WriteFile(project.WebPath, []byte(htmlStr), 0666)
	if err != nil {
		logrus.Errorf("文件写入失败 %s", err)
		return errors.New("文件写入失败")
	}

	return nil
}
