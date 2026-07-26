package log_service

import (
	"blogx_server/core"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum"
	"fmt"

	"github.com/gin-gonic/gin"
)

// NewLoginSuccess 登录成功日志记录
func NewLoginSuccess(c *gin.Context, loginType enum.LoginType) {
	ip := c.ClientIP()
	addr := core.GetIPAddr(ip)
	token := c.GetHeader("token")
	fmt.Println("token", token)
	// TODO 此处是模拟，userID和userName后续要从token中解析
	userID := uint(1)
	userName := ""

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
