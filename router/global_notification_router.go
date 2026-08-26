package router

import (
	"blogx_server/api"
	"blogx_server/api/global_notification_api"
	"blogx_server/middleware"
	"blogx_server/models"

	"github.com/gin-gonic/gin"
)

func GlobalNotificationRouter(r *gin.RouterGroup) {
	gr := r.Group("global_notification")
	globalNotificationApi := api.App.GlobalNotificationApi
	gr.POST("", middleware.AdminMiddleware, middleware.BindJSONMiddleware[global_notification_api.CreateRequest], globalNotificationApi.CreateView)
	gr.GET("", middleware.AuthMiddleware, middleware.BindQueryMiddleware[global_notification_api.ListRequest], globalNotificationApi.ListView)
	gr.DELETE("", middleware.AdminMiddleware, middleware.BindJSONMiddleware[models.DeleteRequest], globalNotificationApi.RemoveAdminView)
	gr.POST("user", middleware.AuthMiddleware, middleware.BindJSONMiddleware[global_notification_api.UserMsgActionRequest], globalNotificationApi.UserMsgActionView)
}
