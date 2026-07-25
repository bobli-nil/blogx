package models

type GlobalNotificationModel struct {
	Model
	Title   string `gorm:"size:64" json:"title"`
	Content string `gorm:"size:64" json:"content"`
	Icon    string `gorm:"size:256" json:"icon"`
	Href    string `gorm:"size:256" json:"href"`
}

func (GlobalNotificationModel) TableName() string {
	return "global_notification"
}
