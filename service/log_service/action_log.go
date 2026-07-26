package log_service

import (
	"blogx_server/core"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ActionLog struct {
	c            *gin.Context
	level        enum.LogLevelType
	title        string
	requestBody  []byte
	responseBody []byte
	middleList   []string // 中间信息
	itemList     []string // 请求体、响应体
}

func (ac *ActionLog) SetTitle(title string) {
	ac.title = title
}

func (ac *ActionLog) SetLevel(level enum.LogLevelType) {
	ac.level = level
}

func (ac *ActionLog) SetRequest() {
	contentType := ac.c.GetHeader("Content-Type")
	var requestBody []byte
	if strings.Contains(contentType, "application/json") || strings.Contains(contentType, "text/plain") {
		byteData, err := io.ReadAll(ac.c.Request.Body)
		if err != nil {
			logrus.Errorf("获取请求体失败 %s", err.Error())
			return
		}
		ac.requestBody = byteData
		ac.c.Request.Body = io.NopCloser(bytes.NewReader(requestBody))
	} else {
		logrus.Warnf("请求体不是JSON")
	}
}

func (ac *ActionLog) SetResponse(responseBody []byte) {
	ac.responseBody = responseBody
}

func (ac *ActionLog) SetItem(key string, value any, logLevelType enum.LogLevelType) {
	var v string
	of := reflect.TypeOf(value)
	switch of.Kind() {
	case reflect.Struct, reflect.Map, reflect.Slice:
		byteData, err := json.Marshal(value)
		if err != nil {
			logrus.Errorf("设置信息转化JSON失败: %s", err)
			return
		}
		v = string(byteData)
	default:
		v = fmt.Sprintf("%v", value)
	}

	str := fmt.Sprintf("<div class=\"middle_info\"><span class=\"key\">%s</span><span class=\"value\">%s</span><span class=\"level\">%s</span></div>", key, v, logLevelType.String())
	ac.middleList = append(ac.middleList, str)
}

func (ac *ActionLog) SetItemInfo(key string, value any) {
	ac.SetItem(key, value, enum.LogInfoLevel)
}

func (ac *ActionLog) SetItemWarn(key string, value any) {
	ac.SetItem(key, value, enum.LogWarnLevel)
}

func (ac *ActionLog) SetItemError(key string, value any) {
	ac.SetItem(key, value, enum.LogErrorLevel)
}

func (ac *ActionLog) Save() {
	// 请求体
	if ac.requestBody != nil {
		str := "<div class=\"log_request\"><div class=\"log_request_head\"><span class=\"log_request_method\">%s</span><span class=\"log_request_path\">%s</span></div><div class=\"log_request_body\"><pre class=\"log_json_body\">%s</pre></div></div>"
		ac.itemList = append(ac.itemList, fmt.Sprintf(str, ac.c.Request.Method, ac.c.Request.URL.String(), string(ac.requestBody)))
	}

	// 中间信息
	ac.itemList = append(ac.itemList, ac.middleList...)

	// 响应体
	if ac.responseBody != nil {
		str := "<div class=\"log_response\"><pre class=\"log_json_body\">%s</pre></div>"
		ac.itemList = append(ac.itemList, fmt.Sprintf(str, string(ac.responseBody)))
	}

	ip := ac.c.ClientIP()
	addr := core.GetIPAddr(ip)
	// TODO 后面这里改
	userID := uint(1)

	err := global.DB.Create(&models.LogModel{
		LogType: enum.ActionLogType,
		Title:   ac.title,
		Content: strings.Join(ac.itemList, "\n"),
		Level:   ac.level,
		UserID:  userID,
		IP:      ip,
		Addr:    addr,
	}).Error
	if err != nil {
		logrus.Errorf("创建操作日志失败 %s", err)
		return
	}
}

func NewActionLogByGin(c *gin.Context) *ActionLog {
	return &ActionLog{c: c}
}

func GetLogFromGinContext(c *gin.Context) *ActionLog {
	_log, ok := c.Get("log")
	if !ok {
		logrus.Errorf("从Gin Context获取log对象失败")
		return nil
	}
	log, ok := _log.(*ActionLog)
	if !ok {
		logrus.Errorf("从Gin Context进行log断言失败")
		return nil
	}
	return log
}
