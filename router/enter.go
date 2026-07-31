package router

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models/enum"
	"blogx_server/utils/jwt"

	"github.com/gin-gonic/gin"
)

func Run() {
	gin.SetMode(global.Conf.System.GinMode)

	r := gin.Default()
	r.Static("/uploads", "uploads")

	nr := r.Group("/api")
	nr.Use(middleware.LogMiddleware)
	SiteRouter(nr)
	LogRouter(nr)
	ImageRouter(nr)

	// TODO 临时用来生成token，后面删除
	r.GET("token/:role", func(c *gin.Context) {
		userType := c.Param("role")

		var userID uint = 1
		var userName string = "张三"
		var role enum.RoleType = enum.AdminRole
		if userType == "common" {
			userID = 2
			userName = "李四"
			role = enum.UserRole
		}

		token, err := jwt.GenerateToken(userID, userName, role)
		if err != nil {
			res.FailWithError(err, c)
			return
		}
		res.OkWithData(token, c)
	})

	addr := global.Conf.System.Addr()
	r.Run(addr)
}
