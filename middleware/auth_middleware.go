package middleware

import (
	"blogx_server/common/res"
	"blogx_server/models/enum"
	"blogx_server/service/redis_service/redis_jwt"
	"blogx_server/utils/jwt"
	"fmt"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware(c *gin.Context) {
	claims, err := jwt.ParseTokenByGin(c)
	if err != nil {
		res.FailWithError(err, c)
		c.Abort()
		return
	}
	blk, ok := redis_jwt.HasTokenBlackByGin(c)
	if ok {
		res.FailWithMsg(blk.Msg(), c)
		c.Abort()
		return
	}
	c.Set("claims", claims)
	fmt.Printf("AuthMiddleware里的claims %+v", claims)
	c.Next()
}

func AdminMiddleware(c *gin.Context) {
	claims, err := jwt.ParseTokenByGin(c)
	if err != nil {
		res.FailWithError(err, c)
		c.Abort()
		return
	}
	blk, ok := redis_jwt.HasTokenBlackByGin(c)
	if ok {
		res.FailWithMsg(blk.Msg(), c)
		c.Abort()
		return
	}
	if claims.Role != enum.AdminRole {
		res.FailWithMsg("权限错误", c)
		c.Abort()
		return
	}
	c.Set("claims", claims)
	c.Next()
}
