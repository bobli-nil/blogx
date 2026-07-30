package router

import (
	"blogx_server/api"
	"blogx_server/middleware"

	"github.com/gin-gonic/gin"
)

func SiteRouter(r *gin.RouterGroup) {
	sr := r.Group("/site")
	sr.Use(middleware.AuthMiddleware)

	siteApi := api.App.SiteApi
	sr.GET(":name", siteApi.SiteInfoView)
	sr.PUT("", siteApi.SiteUpdateView)
}
