package models

import (
	"blogx_server/global"
	"blogx_server/utils/file"
	"fmt"
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
	qiNiu := global.Conf.QiNiu
	suffix, _ := file.ImageSuffixJudge(i.FileName)
	return fmt.Sprintf("%s/%s/%s.%s", qiNiu.Uri, qiNiu.Prefix, i.Hash, suffix)
}

func (i *ImageModel) BeforeDelete(tx *gorm.DB) error {
	err := os.Remove(i.Path)
	if err != nil {
		logrus.Warnf("删除图片二进制文件失败 %s", err)
	}
	return nil
}
