package middleware

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"bytes"
	"io"

	"github.com/gin-gonic/gin"
)

type CaptchaRequest struct {
	CaptchaId   string `json:"captchaId" binding:"required"`
	CaptchaCode string `json:"captchaCode" binding:"required"`
}

func CaptchaMiddleware(c *gin.Context) {
	if !global.Conf.Site.Login.Captcha {
		return
	}
	var cr CaptchaRequest
	byteData, _ := io.ReadAll(c.Request.Body)
	c.Request.Body = io.NopCloser(bytes.NewReader(byteData))
	err := c.ShouldBindJSON(&cr)
	c.Request.Body = io.NopCloser(bytes.NewReader(byteData))
	if err != nil {
		res.FailWithError(err, c)
		c.Abort()
		return
	}

	answer := global.Store.Get(cr.CaptchaId, true)
	if answer != cr.CaptchaCode {
		res.FailWithMsg("验证码错误", c)
		c.Abort()
		return
	}

	c.Next()
}
