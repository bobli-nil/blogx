package router

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models/enum"
	"blogx_server/utils/jwt"

	"github.com/gin-gonic/gin"
)

func Run() {
	gin.SetMode(global.Conf.System.GinMode)

	r := gin.Default()
	r.Static("/uploads", "uploads")

	nr := r.Group("/api")
	nr.Use(middleware.LogMiddleware)
	nr.Use(middleware.AdminMiddleware)
	SiteRouter(nr)
	LogRouter(nr)

	// TODO 临时用来生成token
	r.GET("token", func(c *gin.Context) {
		token, err := jwt.GenerateToken(1000, "admin", enum.AdminRole)
		if err != nil {
			res.FailWithError(err, c)
			return
		}
		res.OkWithData(token, c)
	})

	addr := global.Conf.System.Addr()
	r.Run(addr)
}
