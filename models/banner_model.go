package models

type BannerModel struct {
	Model
	Cover string `gorm:"size:256" json:"cover"`
	Href  string `gorm:"size:256" json:"href"`
	Show  *bool  `json:"show"`
}

func (BannerModel) TableName() string {
	return "banner"
}
