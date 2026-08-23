package redis_comment

import (
	"blogx_server/global"
	"strconv"

	"github.com/sirupsen/logrus"
)

type CommentCacheType string

const (
	commentCacheApply CommentCacheType = "comment_cache_apply"
	commentCacheDigg  CommentCacheType = "comment_cache_digg"
)

func set(t CommentCacheType, commentID uint, n int) {
	num, _ := global.Redis.HGet(string(t), strconv.Itoa(int(commentID))).Int()
	num += n
	global.Redis.HSet(string(t), strconv.Itoa(int(commentID)), strconv.Itoa(num))
}

func SetCacheApply(commentID uint, n int) {
	set(commentCacheApply, commentID, n)
}

func SetCacheDigg(commentID uint, n int) {
	set(commentCacheDigg, commentID, n)
}

func get(t CommentCacheType, commentID uint) int {
	num, _ := global.Redis.HGet(string(t), strconv.Itoa(int(commentID))).Int()
	return num
}

func GetCacheApply(commentID uint) int {
	return get(commentCacheApply, commentID)
}

func GetCacheDigg(commentID uint) int {
	return get(commentCacheDigg, commentID)
}

func GetAll(t CommentCacheType) (mps map[uint]int) {
	res, err := global.Redis.HGetAll(string(t)).Result()
	if err != nil {
		return
	}
	mps = make(map[uint]int)
	for key, val := range res {
		k, err := strconv.Atoi(key)
		if err != nil {
			continue
		}
		v, err := strconv.Atoi(val)
		if err != nil {
			continue
		}
		mps[uint(k)] = v
	}
	return
}

func GetAllCacheApply() map[uint]int {
	return GetAll(commentCacheApply)
}

func GetAllCacheDigg() map[uint]int {
	return GetAll(commentCacheDigg)
}

func Clear() {
	if err := global.Redis.Del(string(commentCacheApply), string(commentCacheDigg)).Err(); err != nil {
		logrus.Error(err)
	} else {
		logrus.Info("redis评论数清空成功")
	}
}
