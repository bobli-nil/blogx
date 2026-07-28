package middleware

import (
	"blogx_server/service/log_service"
	"bytes"

	"github.com/gin-gonic/gin"
)

type ResponseBodyWriter struct {
	gin.ResponseWriter
	Body *bytes.Buffer
}

func (w *ResponseBodyWriter) Write(b []byte) (int, error) {
	w.Body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w *ResponseBodyWriter) WriteString(s string) (int, error) {
	w.Body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

func LogMiddleware(c *gin.Context) {
	// 自定义响应ResponseBodyWriter
	blw := &ResponseBodyWriter{
		ResponseWriter: c.Writer,
		Body:           bytes.NewBufferString(""),
	}
	c.Writer = blw

	log := log_service.NewActionLogByGin(c)
	c.Set("log", log) // 把log对象挂到 context 上

	c.Next()

	// 注意这里，部分路由不要Save
	value, exists := c.Get("GetLogFromGinContext")
	b, _ := value.(bool)
	if exists && b {
		log.SetResponse(blw.Body.Bytes())
		log.Save()
	}
}
