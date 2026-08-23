package router

import (
	"blogx_server/api"
	"blogx_server/api/article_api"
	"blogx_server/middleware"
	"blogx_server/models"

	"github.com/gin-gonic/gin"
)

func ArticleRouter(r *gin.RouterGroup) {
	articleApi := api.App.ArticleApi
	ar := r.Group("article")
	ar.POST("", middleware.AuthMiddleware, middleware.BindJSONMiddleware[article_api.ArticleCreateReq], articleApi.ArticleCreateView)
	ar.GET("", middleware.BindQueryMiddleware[article_api.ArticleListRequest], articleApi.ArticleListView)
	ar.PUT("", middleware.AuthMiddleware, middleware.BindJSONMiddleware[article_api.ArticleUpdateReq], articleApi.ArticleUpdateView)
	ar.GET(":id", middleware.BindUriMiddleware[models.IDRequest], articleApi.ArticleDetailView)

	ar.POST("examine", middleware.AdminMiddleware, middleware.BindJSONMiddleware[article_api.ArticleExamineRequest], articleApi.ArticleExamineView)
	ar.GET("digg/:id", middleware.AuthMiddleware, middleware.BindUriMiddleware[models.IDRequest], articleApi.ArticleDiggView)
	ar.POST("collect", middleware.AuthMiddleware, middleware.BindJSONMiddleware[article_api.ArticleCollectRequest], articleApi.ArticleCollectView)

	ar.POST("history", middleware.BindJSONMiddleware[article_api.ArticleLookRequest], articleApi.ArticleLookView)
	ar.GET("history", middleware.AuthMiddleware, middleware.BindQueryMiddleware[article_api.ArticleLookListRequest], articleApi.ArticleLookListView)
	ar.DELETE("history", middleware.AuthMiddleware, middleware.BindJSONMiddleware[models.DeleteRequest], articleApi.ArticleLookRemoveView)

	ar.DELETE(":id", middleware.AuthMiddleware, middleware.BindUriMiddleware[models.IDRequest], articleApi.ArticleRemoveUserView)
	ar.DELETE("", middleware.AdminMiddleware, middleware.BindJSONMiddleware[models.DeleteRequest], articleApi.ArticleRemoveView)

	ar.POST("category", middleware.AuthMiddleware, middleware.BindJSONMiddleware[article_api.CategoryCreateRequest], articleApi.ArticleCategoryCreate)
	ar.GET("category", middleware.BindQueryMiddleware[article_api.CategoryListRequest], articleApi.CategoryListView)
	ar.DELETE("category", middleware.AuthMiddleware, middleware.BindJSONMiddleware[models.DeleteRequest], articleApi.CategoryRemoveView)

	ar.POST("collect/create", middleware.AuthMiddleware, middleware.BindJSONMiddleware[article_api.CollectCreateRequest], articleApi.CollectCreate)
	ar.GET("collect/list", middleware.BindQueryMiddleware[article_api.CollectListRequest], articleApi.CollectListView)
	ar.DELETE("collect/delete", middleware.AuthMiddleware, middleware.BindJSONMiddleware[models.DeleteRequest], articleApi.CollectRemoveView)

	ar.GET("category/options", middleware.AuthMiddleware, articleApi.CategoryOptionsView)
	ar.GET("tag/options", middleware.AuthMiddleware, articleApi.ArticleTagOptionsView)
}
