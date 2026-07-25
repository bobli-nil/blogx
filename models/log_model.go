package models

import "blogx_server/models/enum"

type LogModel struct {
	Model
	Title       string         `gorm:"size:64" json:"title"`
	Content     string         `json:"content"`
	LogType     enum.LogType   `json:"logType"` // 日志类型
	Level       uint8          `json:"level"`
	UserID      uint           `json:"userId"`
	UserModel   UserModel      `gorm:"foreignKey:UserID;references:ID" json:"-"`
	IP          string         `gorm:"size:32" json:"ip"`
	Addr        string         `gorm:"size:64" json:"addr"`
	IsRead      bool           `json:"isRead"`
	LoginStatus bool           `json:"loginStatus"`             // 登录状态
	UserName    string         `gorm:"size:32" json:"userName"` // 登录日志的用户名
	Pwd         string         `gorm:"size:32" json:"pwd"`      // 登录日志的密码
	LoginType   enum.LoginType `json:"loginType"`               // 登录类型
}

func (LogModel) TableName() string {
	return "log"
}
