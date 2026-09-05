package middleware

import (
	"blogx_server/global"
	"fmt"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type CacheMiddlewarePrefix string

const (
	CacheBannerPrefix CacheMiddlewarePrefix = "cache_banner_"
)

type CacheOption struct {
	Prefix CacheMiddlewarePrefix `json:"prefix"`
	Time   time.Duration         `json:"time"`
	Params []string              `json:"params"`
}

type CacheResponseWriter struct {
	gin.ResponseWriter
	Body []byte
}

func (w *CacheResponseWriter) Write(data []byte) (int, error) {
	w.Body = append(w.Body, data...)
	return w.ResponseWriter.Write(data)
}

func NewBannerCacheOption() CacheOption {
	return CacheOption{
		Prefix: CacheBannerPrefix,
		Time:   time.Hour,
		Params: []string{"type"},
	}
}

func CacheMiddleware(option CacheOption) gin.HandlerFunc {
	return func(c *gin.Context) {
		values := url.Values{}
		for _, key := range option.Params {
			values.Add(key, c.Query(key))
		}
		key := fmt.Sprintf("%s%s", option.Prefix, values.Encode())
		// 请求部分
		val, err := global.Redis.Get(key).Result()
		if err == nil {
			c.Abort()
			fmt.Println("走缓存了")
			c.Header("Content-Type", "application/json; charset=utf-8")
			c.Writer.Write([]byte(val))
			return
		}

		w := &CacheResponseWriter{
			ResponseWriter: c.Writer,
		}
		c.Writer = w
		c.Next()

		// 响应
		body := string(w.Body)
		global.Redis.Set(key, body, option.Time)
	}
}

func CacheClose(prefix CacheMiddlewarePrefix) {
	keys, err := global.Redis.Keys(fmt.Sprintf("%s*", prefix)).Result()
	if err != nil {
		logrus.Errorf("找不到redis keys %s", prefix)
		return
	}
	if len(keys) > 0 {
		global.Redis.Del(keys...)
		logrus.Infof("删除前缀 %s 缓存，共 %d 条", prefix, len(keys))
	}
}
