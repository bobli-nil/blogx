package redis_article

import (
	"blogx_server/global"
	"blogx_server/utils/date"
	"fmt"
	"strconv"

	"github.com/sirupsen/logrus"
)

type ArticleCacheType string

const (
	articleCacheLook    ArticleCacheType = "article_look_key"
	articleCacheDigg    ArticleCacheType = "article_digg_key"
	articleCacheCollect ArticleCacheType = "article_collect_key"
)

func set(t ArticleCacheType, articleID uint, increase bool) {
	num, _ := global.Redis.HGet(string(t), strconv.Itoa(int(articleID))).Int()
	if increase {
		num += 1
	} else {
		num -= 1
	}
	global.Redis.HSet(string(t), strconv.Itoa(int(articleID)), num)
}

func SetCacheLook(articleID uint, increase bool) {
	set(articleCacheLook, articleID, increase)
}

func SetCacheDigg(articleID uint, increase bool) {
	set(articleCacheDigg, articleID, increase)
}

func SetCacheCollect(articleID uint, increase bool) {
	set(articleCacheCollect, articleID, increase)
}

func get(t ArticleCacheType, articleID uint) int {
	num, _ := global.Redis.HGet(string(t), strconv.Itoa(int(articleID))).Int()
	return num
}

func GetCacheLook(articleID uint) int {
	return get(articleCacheLook, articleID)
}

func GetCacheDigg(articleID uint) int {
	return get(articleCacheDigg, articleID)
}

func GetCacheCollect(articleID uint) int {
	return get(articleCacheCollect, articleID)
}

func getAll(t ArticleCacheType) (mps map[uint]int) {
	res, err := global.Redis.HGetAll(string(t)).Result()
	if err != nil {
		return
	}
	mps = make(map[uint]int)
	for k, v := range res {
		keyInt, err := strconv.Atoi(k)
		if err != nil {
			continue
		}
		valueInt, err := strconv.Atoi(v)
		if err != nil {
			continue
		}
		mps[uint(keyInt)] = valueInt
	}
	return
}

func GetAllCacheLook() (mps map[uint]int) {
	return getAll(articleCacheLook)
}

func GetAllCacheDigg() (mps map[uint]int) {
	return getAll(articleCacheDigg)
}

func GetAllCacheCollect() (mps map[uint]int) {
	return getAll(articleCacheCollect)
}

func SetUserArticleHistoryCache(articleID uint, userID uint) {
	key := fmt.Sprintf("history_%d", userID)
	field := fmt.Sprintf("article_%d", articleID)
	endTime := date.GetNowAfter()
	if err := global.Redis.HSet(key, field, "").Err(); err != nil {
		logrus.Error(err)
		return
	}
	if err := global.Redis.ExpireAt(key, endTime).Err(); err != nil {
		logrus.Error(err)
		return
	}
}

func GetUserArticleHistoryCache(articleID uint, userID uint) bool {
	key := fmt.Sprintf("history_%d", userID)
	field := fmt.Sprintf("article_%d", articleID)
	err := global.Redis.HGet(key, field).Err()
	return err == nil
}
