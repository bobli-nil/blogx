package middleware

import (
	"bytes"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func LogMiddleware(c *gin.Context) {
	byteData, err := io.ReadAll(c.Request.Body)
	if err != nil {
		logrus.Errorf("获取请求体失败 %s", err.Error())
		return
	}
	fmt.Println("body: ", string(byteData))
	c.Request.Body = io.NopCloser(bytes.NewReader(byteData))

	c.Next()
}
