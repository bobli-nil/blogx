package log_service

import (
	"blogx_server/core"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/utils/jwt"
	"fmt"

	"github.com/gin-gonic/gin"
)

// NewLoginSuccess 登录成功日志记录
func NewLoginSuccess(c *gin.Context, loginType enum.LoginType) {
	ip := c.ClientIP()
	addr := core.GetIPAddr(ip)
	token := c.GetHeader("token")
	fmt.Println("token", token)
	userID := uint(0)
	userName := ""
	// 从token中解析用户ID和用户名
	myClaims, err := jwt.ParseTokenByGin(c)
	if err == nil && myClaims != nil {
		userID = myClaims.UserID
		userName = myClaims.UserName
	}

	global.DB.Create(&models.LogModel{
		Title:       "用户登录成功",
		Content:     "",
		LogType:     enum.LoginLogType,
		UserID:      userID,
		LoginStatus: true,
		IP:          ip,
		Addr:        addr,
		UserName:    userName,
		Pwd:         "-",
		LoginType:   loginType,
	})
}

// NewLoginFail 登录失败的日志记录
func NewLoginFail(c *gin.Context, loginType enum.LoginType, msg string, userName string, pwd string) {
	ip := c.ClientIP()
	addr := core.GetIPAddr(ip)
	global.DB.Create(&models.LogModel{
		Title:       "用户登录失败",
		Content:     msg,
		LogType:     enum.LoginLogType,
		LoginStatus: true,
		IP:          ip,
		Addr:        addr,
		UserName:    userName,
		Pwd:         pwd,
		LoginType:   loginType,
	})
}
