package router

import (
	"blogx_server/api"
	"blogx_server/api/site_msg_api"
	"blogx_server/middleware"

	"github.com/gin-gonic/gin"
)

func SiteMsgRouter(r *gin.RouterGroup) {
	sr := r.Group("site_msg")
	siteMsgApi := api.App.SiteMsgApi
	sr.GET("", middleware.AuthMiddleware, middleware.BindQueryMiddleware[site_msg_api.SiteMsgListRequest], siteMsgApi.SiteMsgListView)
	sr.POST("", middleware.AuthMiddleware, middleware.BindJSONMiddleware[site_msg_api.SiteMsgReadRequest], siteMsgApi.SiteMsgReadView)
	sr.GET("conf", middleware.AuthMiddleware, siteMsgApi.UserMessageConfView)
	sr.PUT("conf", middleware.AuthMiddleware, middleware.BindJSONMiddleware[site_msg_api.UserMessageConfRequest], siteMsgApi.UserMessageConfUpdateView)
}
