package router

import (
	"blogx_server/api"
	"blogx_server/api/data_api"
	"blogx_server/middleware"

	"github.com/gin-gonic/gin"
)

func DataRouter(r *gin.RouterGroup) {
	dr := r.Group("data")
	dataApi := api.App.DataApi

	dr.GET("sum", middleware.AdminMiddleware, dataApi.SumView)
	dr.GET("article", middleware.AdminMiddleware, dataApi.ArticleYearDataView)
	dr.GET("growth", middleware.AdminMiddleware, middleware.BindQueryMiddleware[data_api.GrowthDataRequest], dataApi.GrowthData)
	dr.GET("computer", middleware.AdminMiddleware, dataApi.ComputerDataView)
}
