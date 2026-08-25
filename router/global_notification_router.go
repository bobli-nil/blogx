package router

import (
	"blogx_server/api"
	"blogx_server/api/global_notification_api"
	"blogx_server/middleware"

	"github.com/gin-gonic/gin"
)

func GlobalNotificationRouter(r *gin.RouterGroup) {
	gr := r.Group("global_notification")
	globalNotificationApi := api.App.GlobalNotificationApi
	gr.POST("", middleware.AdminMiddleware, middleware.BindJSONMiddleware[global_notification_api.CreateRequest], globalNotificationApi.CreateView)
}
