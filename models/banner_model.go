package models

type BannerModel struct {
	Model
	Cover string `gorm:"size:256" json:"cover"`
	Href  string `gorm:"size:256" json:"href"`
	Show  *bool  `json:"show"`
	Type  int8   `json:"type"` // 1banner 2独家推广
}

func (BannerModel) TableName() string {
	return "banner"
}
