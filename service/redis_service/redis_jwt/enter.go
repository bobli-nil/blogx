package redis_jwt

import (
	"blogx_server/global"
	"blogx_server/utils/jwt"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type BlackType int8

const (
	UserBlackType   BlackType = 1 // 用户注销登录
	AdminBlackType  BlackType = 2 // 管理员注销登录
	DeviceBlackType BlackType = 3 // 其他设备挤下来
)

func (b BlackType) String() string {
	return fmt.Sprintf("%d", b)
}

func (b BlackType) Msg() string {
	switch b {
	case UserBlackType:
		return "已注销"
	case AdminBlackType:
		return "禁止登录"
	case DeviceBlackType:
		return "设备下线"
	}
	return "已注销"
}

func ParseBlackType(str string) BlackType {
	switch str {
	case "1":
		return UserBlackType
	case "2":
		return AdminBlackType
	case "3":
		return DeviceBlackType
	default:
		return UserBlackType
	}
}

func TokenBlack(token string, value BlackType) {
	key := fmt.Sprintf("token_black_%s", token)

	claims, err := jwt.ParseToken(token)
	if err != nil || claims == nil {
		logrus.Errorf("token解析失败：%s", err)
		return
	}
	seconds := claims.ExpiresAt.Unix() - time.Now().Unix()

	_, err1 := global.Redis.Set(key, value.String(), time.Duration(seconds)*time.Second).Result()
	if err1 != nil {
		logrus.Errorf("设置redis值失败:%s", err1.Error())
		return
	}
}

func HasTokenBlack(token string) (blk BlackType, ok bool) {
	key := fmt.Sprintf("token_black_%s", token)
	val, err := global.Redis.Get(key).Result()
	if err != nil {
		logrus.Errorf("从redis获取token失败: %s", err)
		return
	}
	blk = ParseBlackType(val)
	return blk, true
}

func HasTokenBlackByGin(c *gin.Context) (blk BlackType, ok bool) {
	token := c.GetHeader("token")
	if token == "" {
		token = c.Query("token")
	}
	return HasTokenBlack(token)
}
