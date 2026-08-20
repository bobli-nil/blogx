package router

import (
	"blogx_server/api"
	"blogx_server/api/comment_api"
	"blogx_server/middleware"

	"github.com/gin-gonic/gin"
)

func CommentRouter(r *gin.RouterGroup) {
	cr := r.Group("comment")
	commentApi := api.App.CommentApi
	cr.POST("", middleware.AuthMiddleware, middleware.BindJSONMiddleware[comment_api.CommentCreateRequest], commentApi.CommentCreateView)
}
