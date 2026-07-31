package site_api

import (
	"blogx_server/common/res"
	"blogx_server/conf"
	"blogx_server/core"
	"blogx_server/global"
	"blogx_server/models/enum"
	"blogx_server/utils/jwt"

	"github.com/gin-gonic/gin"
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

func UpdateSite(Site conf.Site) error {
	return nil
}
