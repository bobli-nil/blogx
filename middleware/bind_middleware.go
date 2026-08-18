package middleware

import (
	"blogx_server/common/res"

	"github.com/gin-gonic/gin"
)

func BindJSONMiddleware[T any](c *gin.Context) {
	var cr T
	if err := c.ShouldBindJSON(&cr); err != nil {
		res.FailWithError(err, c)
		c.Abort()
		return
	}
	c.Set("request", cr)
	c.Next()
}

func BindQueryMiddleware[T any](c *gin.Context) {
	var cr T
	if err := c.ShouldBindQuery(&cr); err != nil {
		res.FailWithError(err, c)
		c.Abort()
		return
	}
	c.Set("request", cr)
	c.Next()
}

func BindUriMiddleware[T any](c *gin.Context) {
	var cr T
	if err := c.ShouldBindUri(&cr); err != nil {
		res.FailWithError(err, c)
		c.Abort()
		return
	}
	c.Set("request", cr)
	c.Next()
}

func GetBind[T any](c *gin.Context) (cr T) {
	return c.MustGet("request").(T)
}
