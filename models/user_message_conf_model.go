package models

type UserMessageConfModel struct {
	UserID             uint      `gorm:"primary_key"`
	UserModel          UserModel `gorm:"foreignkey:UserID" json:"-"`
	OpenCommentMessage bool      `json:"openCommentMessage"` // 是否开启回复和评论
	OpenDiggMessage    bool      `json:"openDiggMessage"`    // 是否开启赞和收藏
	OpenPrivateChat    bool      `json:"openPrivateChat"`    // 是否开启私聊
}

func (UserMessageConfModel) TableName() string {
	return "user_message_conf"
}
