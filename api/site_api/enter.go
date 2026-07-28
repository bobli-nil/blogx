package site_api

import (
	"blogx_server/common/res"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type SiteApi struct{}

func (s *SiteApi) SiteInfoView(c *gin.Context) {
	res.OkWithData("data", c)
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
