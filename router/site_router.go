package router

import (
	"blogx_server/api"

	"github.com/gin-gonic/gin"
)

func SiteRouter(r *gin.RouterGroup) {
	siteApi := api.App.SiteApi
	r.GET("site", siteApi.SiteInfoView)
	r.PUT("site", siteApi.SiteUpdateView)
}
