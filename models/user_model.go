package models

import (
	"blogx_server/models/enum"
	"math"
	"time"

	"gorm.io/gorm"
)

type UserModel struct {
	Model
	Username       string              `gorm:"size:32" json:"username"`
	Nickname       string              `gorm:"size:32" json:"nickname"`
	Email          string              `gorm:"size:128" json:"email"`
	Password       string              `gorm:"size:64" json:"-"`
	Avatar         string              `gorm:"size:256" json:"avatar"`
	Abstract       string              `gorm:"size:256" json:"abstract"` // 简介
	RegisterSource enum.RegisterSource `json:"registerSource"`           // 注册来源
	OpenID         string              `gorm:"size:64" json:"openID"`    // 第三方登录ID
	Role           enum.RoleType       `json:"role"`                     // 角色 1管理员 2普通用户 3访客
	UserConfModel  *UserConfModel      `gorm:"foreignKey:UserID" json:"-"`
	IP             string              `gorm:"size:64" json:"ip"`
	Address        string              `gorm:"size:128" json:"address"`
}

func (u UserModel) GetID() uint {
	return u.ID
}

func (*UserModel) TableName() string {
	return "user"
}

func (u *UserModel) CodeAge() int {
	sub := time.Now().Sub(*u.CreatedAt)
	return int(math.Ceil(sub.Hours() / 24 / 365))
}

func (u *UserModel) AfterCreate(tx *gorm.DB) error {
	err := tx.Create(&UserConfModel{
		UserID:      u.ID,
		OpenCollect: true,
		OpenFollow:  true,
		OpenFans:    true,
		HomeStyleID: 1,
	}).Error
	err = tx.Create(&UserMessageConfModel{
		UserID:             u.ID,
		OpenCommentMessage: true,
		OpenDiggMessage:    true,
		OpenPrivateChat:    true,
	}).Error
	return err
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
