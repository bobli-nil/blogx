package router

import (
	"blogx_server/api"
	"blogx_server/api/chat_api"
	"blogx_server/middleware"
	"blogx_server/models"

	"github.com/gin-gonic/gin"
)

func ChatRouter(r *gin.RouterGroup) {
	cr := r.Group("chat")
	chatApi := api.App.ChatApi

	cr.GET("", middleware.AuthMiddleware, middleware.BindQueryMiddleware[chat_api.ChatListRequest], chatApi.ChatListView)
	cr.GET("session", middleware.AuthMiddleware, middleware.BindQueryMiddleware[chat_api.SessionListRequest], chatApi.SessionListView)
	cr.DELETE("", middleware.AuthMiddleware, middleware.BindJSONMiddleware[models.DeleteRequest], chatApi.UserChatDeleteView)
}
