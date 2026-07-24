package global

import (
	"blogx_server/conf"

	"gorm.io/gorm"
)

var (
	Conf *conf.Config
	DB   *gorm.DB
)
