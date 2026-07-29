package router

import (
	"blogx_server/api"

	"github.com/gin-gonic/gin"
)

func LogRouter(r *gin.RouterGroup) {
	logApi := api.App.LogApi
	r.GET("log", logApi.LogListView)
	r.GET("log/:id", logApi.LogReadView)
	r.DELETE("log", logApi.LogDeleteView)
}
