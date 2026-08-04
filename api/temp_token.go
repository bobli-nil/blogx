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

		var userID uint = 13
		var userName = "zhangsan"
		var role = enum.AdminRole
		if userType == "common" {
			userID = 14
			userName = "lisi"
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
