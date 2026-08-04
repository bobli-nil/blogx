package middleware

import (
	"blogx_server/common/res"
	"blogx_server/utils/email_store"
	"bytes"
	"io"

	"github.com/gin-gonic/gin"
)

type EmailVerifyRequest struct {
	EmailID string `json:"emailID" binding:"required"`
	Code    string `json:"code" binding:"required"`
}

func EmailVerifyMiddleware(c *gin.Context) {
	var cr EmailVerifyRequest
	byteData, _ := io.ReadAll(c.Request.Body)
	c.Request.Body = io.NopCloser(bytes.NewReader(byteData))
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, c)
		c.Abort()
		return
	}
	c.Request.Body = io.NopCloser(bytes.NewReader(byteData))

	info, ok := email_store.Verify(cr.EmailID, cr.Code)
	if !ok {
		res.FailWithMsg("邮箱验证失败", c)
		c.Abort()
		return
	}

	c.Set("email", info.Email)

	c.Next()
}
