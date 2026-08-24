package api

import (
	"blogx_server/common/res"
	"blogx_server/models/enum"
	"blogx_server/utils/jwt"

	"github.com/gin-gonic/gin"
)

func GenerateTempToken(r *gin.Engine) {
	r.GET("token/:role", func(c *gin.Context) {
		userType := c.Param("role")

		var userID uint
		var userName string
		var role enum.RoleType
		switch userType {
		case "admin":
			userID = 1
			userName = "zhangsan"
			role = enum.AdminRole
		case "common":
			userID = 2
			userName = "lisi"
			role = enum.UserRole
		case "wangwu":
			userID = 6
			userName = "wangwu"
			role = enum.UserRole
		}

		token, err := jwt.GenerateToken(userID, userName, role)
		if err != nil {
			res.FailWithError(err, c)
			return
		}
		res.OkWithData(token, c)
	})
}
