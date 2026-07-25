package models

import "time"

type UserModel struct {
	Model
	Username       string `gorm:"size:32" json:"username"`
	Nickname       string `gorm:"size:32" json:"nickname"`
	Email          string `gorm:"size:128" json:"email"`
	Password       string `gorm:"size:64" json:"-"`
	Avatar         string `gorm:"size:256" json:"avatar"`
	Abstract       string `gorm:"size:256" json:"abstract"` // 简介
	RegisterSource int8   `json:"registerSource"`           // 注册来源
	CodeAge        int    `json:"codeAge"`                  // 码龄
	OpenID         string `gorm:"size:64" json:"openID"`    // 第三方登录ID
	Role           uint8  `json:"role"`                     // 角色 1管理员 2普通用户 3访客
}

func (UserModel) TableName() string {
	return "user"
}

type UserConfModel struct {
	UserID             uint       `gorm:"unique" json:"userID"`
	UserModel          UserModel  `gorm:"foreignKey:UserID;references:ID" json:"-"`
	LikeTags           []string   `gorm:"type:longtext;serializer:json" json:"likeTags"` // 兴趣标签
	UpdateUsernameDate *time.Time `json:"updateUsernameDate"`                            // 上次修改用户名时间
	OpenCollect        bool       `json:"openCollect"`                                   // 公开我的收藏
	OpenFollow         bool       `json:"openFollow"`                                    // 公开我的关注
	OpenFans           bool       `json:"openFans"`                                      // 公开我的粉丝
	HomeStyleID        uint       `json:"homeStyleID"`                                   // 主页样式ID
}

func (UserConfModel) TableName() string {
	return "user_conf"
}
