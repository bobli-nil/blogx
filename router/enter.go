package router

import (
	"blogx_server/api"
	"blogx_server/global"
	"blogx_server/middleware"

	"github.com/gin-gonic/gin"
)

func Run() {
	gin.SetMode(global.Conf.System.GinMode)

	r := gin.Default()
	r.Static("/uploads", "uploads")

	// TODO 临时用来生成token，后面删除
	api.GenerateTempToken(r)

	nr := r.Group("/api")
	nr.Use(middleware.LogMiddleware)
	SiteRouter(nr)
	LogRouter(nr)
	ImageRouter(nr)
	BannerRouter(nr)
	CaptchaRouter(nr)
	UserRouter(nr)
	ArticleRouter(nr)
	CommentRouter(nr)
	SiteMsgRouter(nr)
	GlobalNotificationRouter(nr)
	FocusRouter(nr)

	addr := global.Conf.System.Addr()
	r.Run(addr)
}
