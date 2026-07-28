package router

import (
	"blogx_server/global"
	"blogx_server/middleware"

	"github.com/gin-gonic/gin"
)

func Run() {
	gin.SetMode(global.Conf.System.GinMode)

	r := gin.Default()
	r.Static("/uploads", "uploads")

	nr := r.Group("/api")
	nr.Use(middleware.LogMiddleware)
	SiteRouter(nr)
	LogRouter(nr)

	addr := global.Conf.System.Addr()
	r.Run(addr)
}
