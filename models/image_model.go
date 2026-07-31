package models

type ImageModel struct {
	Model
	FileName string `gorm:"size:32" json:"fileName"`
	Path     string `gorm:"size:256" json:"path"`
	Size     int64  `json:"size"`
	Hash     string `gorm:"size:32" json:"hash"`
}

func (ImageModel) TableName() string {
	return "image"
}

func (i *ImageModel) WebPath() string {
	return "/" + i.Path
}
