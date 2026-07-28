package res

import (
	"blogx_server/utils/validate"

	"github.com/gin-gonic/gin"
)

type Code int

func (c Code) String() string {
	switch c {
	case SuccessCode:
		return "成功"
	case FailValidCode:
		return "校验失败"
	case FailServiceCode:
		return "服务异常"
	default:
		return ""
	}
}

const (
	SuccessCode     Code = 0
	FailValidCode   Code = 1001
	FailServiceCode Code = 1002
)

type Response struct {
	Code Code   `json:"code"`
	Data any    `json:"data"`
	Msg  string `json:"msg"`
}

var empty = map[string]any{}

func (r *Response) Json(c *gin.Context) {
	c.JSON(200, r)
}

func Ok(data any, msg string, c *gin.Context) {
	res := &Response{SuccessCode, data, msg}
	res.Json(c)
}

func OkWithData(data any, c *gin.Context) {
	res := &Response{SuccessCode, data, SuccessCode.String()}
	res.Json(c)
}

func OkWithList(list any, count int, c *gin.Context) {
	res := &Response{
		SuccessCode,
		map[string]any{
			"list":  list,
			"count": count,
		},
		SuccessCode.String(),
	}
	res.Json(c)
}

func OkWithMsg(msg string, c *gin.Context) {
	res := &Response{SuccessCode, empty, msg}
	res.Json(c)
}

func FailWithData(data any, msg string, c *gin.Context) {
	res := &Response{FailServiceCode, data, msg}
	res.Json(c)
}

func FailWithMsg(msg string, c *gin.Context) {
	res := &Response{FailValidCode, empty, msg}
	res.Json(c)
}

func FailWithCode(code Code, c *gin.Context) {
	res := &Response{code, empty, code.String()}
	res.Json(c)
}

func FailWithError(err error, c *gin.Context) {
	msg := validate.ValidateErr(err)
	FailWithMsg(msg, c)
}
