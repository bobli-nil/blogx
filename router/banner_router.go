package router

import (
	"blogx_server/api"
	"blogx_server/middleware"

	"github.com/gin-gonic/gin"
)

func BannerRouter(r *gin.RouterGroup) {
	bannerApi := api.App.BannerApi
	br := r.Group("banner")
	br.POST("", middleware.AdminMiddleware, bannerApi.BannerCreateView)
	br.GET("", middleware.AdminMiddleware, bannerApi.BannerListView)
	br.PUT(":id", middleware.AdminMiddleware, bannerApi.BannerUpdateView)
	br.DELETE("", middleware.AdminMiddleware, bannerApi.BannerRemoveView)
}
