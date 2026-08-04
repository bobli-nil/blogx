package global

import (
	"blogx_server/conf"

	"github.com/go-redis/redis"
	"github.com/mojocn/base64Captcha"
	"gorm.io/gorm"
)

const Version = "1.0.1"

var (
	Conf  *conf.Config
	DB    *gorm.DB
	Redis *redis.Client
	Store = base64Captcha.DefaultMemStore
)
