package models

type LogModel struct {
	Model
	Title     string    `gorm:"size:64" json:"title"`
	Content   string    `json:"content"`
	LogType   uint8     `json:"logType"` // 日志类型
	Level     uint8     `json:"level"`
	UserID    uint      `json:"userId"`
	UserModel UserModel `gorm:"foreignKey:UserID;references:ID" json:"-"`
	IP        string    `gorm:"size:32" json:"ip"`
	Addr      string    `gorm:"size:64" json:"addr"`
	IsRead    bool      `json:"isRead"`
}

func (LogModel) TableName() string {
	return "log"
}
