package router

import (
	"blogx_server/api"
	"blogx_server/api/ai_api"
	"blogx_server/middleware"

	"github.com/gin-gonic/gin"
)

func AiRouter(r *gin.RouterGroup) {
	aiApi := api.App.AiApi
	ar := r.Group("ai")

	ar.POST("analysis", middleware.AuthMiddleware, middleware.BindJSONMiddleware[ai_api.ArticleAnalysisRequest], aiApi.ArticleAnalysisView)
}
