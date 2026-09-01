package router

import (
	"blogx_server/api"
	"blogx_server/api/search_api"
	"blogx_server/middleware"

	"github.com/gin-gonic/gin"
)

func SearchRouter(r *gin.RouterGroup) {
	searchApi := api.App.SearchApi
	r.GET("article/search", middleware.BindQueryMiddleware[search_api.ArticleSearchRequest], searchApi.ArticleSearchView)
	r.GET("text/search", middleware.BindQueryMiddleware[search_api.TextSearchRequest], searchApi.TextSearchView)
}
