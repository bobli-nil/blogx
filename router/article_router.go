package router

import (
	"blogx_server/api"
	"blogx_server/api/article_api"
	"blogx_server/middleware"

	"github.com/gin-gonic/gin"
)

func ArticleRouter(r *gin.RouterGroup) {
	articleApi := api.App.ArticleApi
	ar := r.Group("article")
	ar.POST("", middleware.AuthMiddleware, middleware.BindJSONMiddleware[article_api.ArticleCreateReq], articleApi.ArticleCreateView)
}
