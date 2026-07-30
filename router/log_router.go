package router

import (
	"blogx_server/api"
	"blogx_server/middleware"

	"github.com/gin-gonic/gin"
)

func LogRouter(r *gin.RouterGroup) {
	logApi := api.App.LogApi
	lr := r.Group("log")
	lr.Use(middleware.AdminMiddleware)
	lr.GET("", logApi.LogListView)
	lr.DELETE("", logApi.LogDeleteView)
	lr.GET(":id", logApi.LogReadView)
}
