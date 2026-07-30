package site_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models/enum"
	"blogx_server/utils/jwt"

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
		res.OkWithData(global.Conf.Site.SiteInfo, c)
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
		data = global.Conf.Email
	case "qq":
		data = global.Conf.QQ
	case "qiNiu":
		data = global.Conf.QiNiu
	case "ai":
		data = global.Conf.Ai
	default:
		res.FailWithMsg("不存在配置", c)
		return
	}

	res.OkWithData(data, c)
}

type SiteUpdateReq struct {
	Name string `json:"name" binding:"required" label:"名称"`
}

func (s *SiteApi) SiteUpdateView(c *gin.Context) {
	var requestBody SiteUpdateReq
	err := c.ShouldBindJSON(&requestBody)
	if err != nil {
		logrus.Errorf("参数绑定失败 %s", err)
		res.FailWithError(err, c)
		return
	}
	res.OkWithMsg("更新成功", c)

}
