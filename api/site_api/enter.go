package site_api

import (
	"blogx_server/models/enum"
	"blogx_server/service/log_service"

	"github.com/gin-gonic/gin"
)

type SiteApi struct{}

func (s *SiteApi) SiteInfoView(c *gin.Context) {
	log_service.NewLoginSuccess(c, enum.UserNamePwdLoginType)
	log_service.NewLoginFail(c, enum.UserNamePwdLoginType, "用户不存在", "lisi", "password")
	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
	})
	return
}

func (s *SiteApi) SiteUpdateView(c *gin.Context) {
	ac := log_service.NewActionLogByGin(c)
	ac.Save()

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "ok",
	})
}
