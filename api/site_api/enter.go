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
	log := log_service.GetLogFromGinContext(c)
	log.SetRequest()
	log.SetTitle("这是SiteUpdateView的标题")
	log.SetLevel(enum.LogInfoLevel)
	log.SetItemInfo("a", "a")
	log.SetItemInfo("aa", 13)
	log.SetItemWarn("b", map[string]string{"b": "这是B"})
	log.SetItemError("c", []string{"这是切片1", "这是切片2"})

	c.JSON(200, gin.H{
		"code": 200,
		"msg":  "站点信息更新成功",
	})
}
