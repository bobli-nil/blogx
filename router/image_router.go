package router

import (
	"blogx_server/api"
	"blogx_server/middleware"

	"github.com/gin-gonic/gin"
)

func ImageRouter(r *gin.RouterGroup) {
	imageApi := api.App.ImageApi
	ir := r.Group("image")

	ir.POST("upload", middleware.AuthMiddleware, imageApi.UploadImageView)
}
