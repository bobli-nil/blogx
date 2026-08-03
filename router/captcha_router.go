package router

import (
	"blogx_server/api"

	"github.com/gin-gonic/gin"
)

func CaptchaRouter(c *gin.RouterGroup) {
	captchaApi := api.App.CaptchaApi
	cr := c.Group("captcha")
	cr.GET("", captchaApi.CaptchaCreateView)
}
