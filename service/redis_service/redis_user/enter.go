package redis_user

import (
	"blogx_server/global"
	"strconv"

	"github.com/sirupsen/logrus"
)

type userCacheType string

const (
	userCacheLook userCacheType = "user_look_count"
)

func set(t userCacheType, userID uint, n int) {
	num, _ := global.Redis.HGet(string(t), strconv.Itoa(int(userID))).Int()
	num += n
	global.Redis.HSet(string(t), strconv.Itoa(int(userID)), num)
}

func SetCacheLook(userID uint, increase bool) {
	n := 1
	if !increase {
		n = -1
	}
	set(userCacheLook, userID, n)
}

func get(t userCacheType, userID uint) int {
	num, _ := global.Redis.HGet(string(t), strconv.Itoa(int(userID))).Int()
	return num
}

func GetCacheLook(userID uint) int {
	return get(userCacheLook, userID)
}

func GetAll(t userCacheType) (mp map[uint]int) {
	res, err := global.Redis.HGetAll(string(t)).Result()
	if err != nil {
		return nil
	}
	mp = make(map[uint]int)
	for key, value := range res {
		iK, err := strconv.Atoi(key)
		if err != nil {
			continue
		}
		iV, err := strconv.Atoi(value)
		if err != nil {
			continue
		}
		mp[uint(iK)] = iV
	}
	return
}

func GetAllCacheLook() map[uint]int {
	return GetAll(userCacheLook)
}

func Clear() {
	err := global.Redis.Del(string(userCacheLook))
	if err != nil {
		logrus.Error(err)
	}
}
