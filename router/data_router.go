package router

import (
	"blogx_server/api"
	"blogx_server/middleware"

	"github.com/gin-gonic/gin"
)

func DataRouter(r *gin.RouterGroup) {
	dr := r.Group("data")
	dataApi := api.App.DataApi

	dr.GET("sum", middleware.AdminMiddleware, dataApi.SumView)
	dr.GET("article", middleware.AdminMiddleware, dataApi.ArticleDataView)
}
