package models

import (
	"os"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

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

func (i *ImageModel) BeforeDelete(tx *gorm.DB) error {
	err := os.Remove(i.Path)
	if err != nil {
		logrus.Warnf("删除图片二进制文件失败 %s", err)
	}
	return nil
}
