package res

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

func SSEOk(data any, c *gin.Context) {
	byteData, _ := json.Marshal(Response{SuccessCode, data, SuccessCode.String()})
	c.SSEvent("", string(byteData))
	c.Writer.Flush()
}

func SSEFail(msg string, c *gin.Context) {
	byteData, _ := json.Marshal(Response{FailValidCode, msg, FailValidCode.String()})
	c.SSEvent("", string(byteData))
	c.Writer.Flush()
}
