package global

import (
	"blogx_server/conf"

	"github.com/go-redis/redis"
	"gorm.io/gorm"
)

var (
	Conf  *conf.Config
	DB    *gorm.DB
	Redis *redis.Client
)
