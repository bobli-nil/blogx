package router

import (
	"blogx_server/api"
	"blogx_server/middleware"

	"github.com/gin-gonic/gin"
)

func UserRouter(c *gin.RouterGroup) {
	userApi := api.App.UserApi
	ur := c.Group("user")
	ur.POST("send_email", middleware.CaptchaMiddleware, userApi.SendEmailView)
	ur.POST("email", middleware.EmailVerifyMiddleware, userApi.RegisterEmailView)
	ur.POST("pwd_login", middleware.CaptchaMiddleware, userApi.PwdLogin)
	ur.GET("detail", middleware.AuthMiddleware, userApi.UserDetailView)
	ur.GET("base", userApi.UserBaseInfoView)
}
