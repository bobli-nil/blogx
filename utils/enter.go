package utils

import (
	"math/rand"
	"time"
)

func InList[T comparable](key T, list []T) bool {
	for _, v := range list {
		if v == key {
			return true
		}
	}
	return false
}

// GenerateRandomDigitsSimple 生成固定长度随机数
func GenerateRandomDigitsSimple(length int) string {
	if length <= 0 {
		return ""
	}

	// 使用当前时间纳秒作为种子，保证每次执行结果不同
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	result := make([]byte, length)
	for i := 0; i < length; i++ {
		result[i] = byte('0' + rng.Intn(10))
	}
	return string(result)
}

func Unique[T comparable](list []T) []T {
	size := len(list)
	if size == 0 {
		return []T{}
	}
	mp := make(map[T]bool)
	for _, v := range list {
		if _, ok := mp[v]; !ok {
			mp[v] = true
		}
	}
	newSlice := make([]T, 0)
	for key, _ := range mp {
		newSlice = append(newSlice, key)
	}
	return newSlice
}
