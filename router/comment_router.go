package router

import (
	"blogx_server/api"
	"blogx_server/api/comment_api"
	"blogx_server/middleware"
	"blogx_server/models"

	"github.com/gin-gonic/gin"
)

func CommentRouter(r *gin.RouterGroup) {
	cr := r.Group("comment")
	commentApi := api.App.CommentApi
	cr.POST("", middleware.AuthMiddleware, middleware.BindJSONMiddleware[comment_api.CommentCreateRequest], commentApi.CommentCreateView)
	cr.GET("tree/:id", middleware.AuthMiddleware, middleware.BindUriMiddleware[models.IDRequest], commentApi.CommentTreeView)
	cr.GET("", middleware.AuthMiddleware, middleware.BindQueryMiddleware[comment_api.CommentListRequest], commentApi.CommentListView)
}
