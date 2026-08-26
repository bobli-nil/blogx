package router

import (
	"blogx_server/api"
	"blogx_server/api/focus_api"
	"blogx_server/middleware"

	"github.com/gin-gonic/gin"
)

func FocusRouter(r *gin.RouterGroup) {
	fr := r.Group("focus")
	focusApi := api.App.FocusApi

	fr.POST("", middleware.AuthMiddleware, middleware.BindJSONMiddleware[focus_api.FocusUserRequest], focusApi.FocusUserView)
	fr.DELETE("", middleware.AuthMiddleware, middleware.BindJSONMiddleware[focus_api.FocusUserRequest], focusApi.UnfocusUserView)
	fr.GET("my_focus", middleware.BindQueryMiddleware[focus_api.FocusUserListRequest], focusApi.FocusUserListView)
	fr.GET("my_fans", middleware.BindQueryMiddleware[focus_api.FansUserListRequest], focusApi.FansUserListView)
}
