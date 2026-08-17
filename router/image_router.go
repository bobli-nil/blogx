package router

import (
	"blogx_server/api"
	"blogx_server/api/image_api"
	"blogx_server/middleware"

	"github.com/gin-gonic/gin"
)

func ImageRouter(r *gin.RouterGroup) {
	imageApi := api.App.ImageApi
	ir := r.Group("image")

	ir.POST("upload", middleware.AuthMiddleware, imageApi.UploadImageView)
	ir.GET("list", middleware.AuthMiddleware, imageApi.ImageListView)
	ir.DELETE("", middleware.AdminMiddleware, imageApi.ImageRemoveView)
	ir.POST("token", middleware.AuthMiddleware, imageApi.QiNiuGenToken)
	ir.POST("transfer_deposit", middleware.AuthMiddleware, middleware.BindJSONMiddleware[image_api.TransferDepositRequest], imageApi.TransferDepositView)
}
