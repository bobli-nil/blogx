package core

import (
	"blogx_server/global"

	"github.com/go-redis/redis"
	"github.com/sirupsen/logrus"
)

func InitRedis() *redis.Client {
	r := global.Conf.Redis
	redisDB := redis.NewClient(&redis.Options{
		Addr:     r.Addr,
		Password: r.Password,
		DB:       r.DB,
	})
	_, err := redisDB.Ping().Result()
	if err != nil {
		logrus.Fatalf("redis连接失败：%s", err.Error())
	}
	logrus.Info("redis连接成功")
	return redisDB
}
