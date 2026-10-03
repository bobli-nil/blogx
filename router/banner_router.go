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
	// TODO 这里是首页的，不应该走Admin权限
	//br.GET("", middleware.AdminMiddleware, middleware.CacheMiddleware(middleware.NewBannerCacheOption()), bannerApi.BannerListView)
	br.GET("", middleware.CacheMiddleware(middleware.NewBannerCacheOption()), bannerApi.BannerListView)
	br.PUT(":id", middleware.AdminMiddleware, bannerApi.BannerUpdateView)
	br.DELETE("", middleware.AdminMiddleware, bannerApi.BannerRemoveView)
}
