package middleware

import (
	"blogx_server/global"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type CacheMiddlewarePrefix string

const (
	CacheBannerPrefix       CacheMiddlewarePrefix = "cache_banner_"
	CacheDataComputerPrefix CacheMiddlewarePrefix = "cache_data_computer_"
)

type CacheOption struct {
	Prefix      CacheMiddlewarePrefix     `json:"prefix"`
	Time        time.Duration             `json:"time"`
	Params      []string                  `json:"params"`
	NoCacheFunc func(c *gin.Context) bool `json:"-"`
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
		NoCacheFunc: func(c *gin.Context) bool {
			referer := c.GetHeader("referer")
			return strings.Contains(referer, "admin")
		},
	}
}

func NewDataComputerCacheOption() CacheOption {
	return CacheOption{
		Prefix: CacheDataComputerPrefix,
		Time:   time.Minute,
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
		// 如果请求头referer中包含admin路径，就不走缓存
		if (err == nil) && (option.NoCacheFunc == nil || !option.NoCacheFunc(c)) {
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
