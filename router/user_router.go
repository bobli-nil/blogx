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
	ur.GET("login_list", middleware.AuthMiddleware, userApi.UserLoginList)
	ur.PUT("update_password", middleware.AuthMiddleware, userApi.UpdatePassword)
	ur.PUT("reset_password", middleware.EmailVerifyMiddleware, userApi.ResetPassword)
	ur.PUT("bind_email", middleware.AuthMiddleware, middleware.EmailVerifyMiddleware, userApi.BindEmailView)
	ur.PUT("", middleware.AuthMiddleware, userApi.UserInfoUpdateView)
	ur.PUT("admin", middleware.AdminMiddleware, userApi.AdminUserInfoUpdateView)
}
