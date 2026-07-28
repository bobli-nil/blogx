package log_service

import (
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum"
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/sirupsen/logrus"
)

type RuntimeDateType int8

const (
	RuntimeDateHour  RuntimeDateType = 1 // 按小时
	RuntimeDateDay   RuntimeDateType = 2 // 按天
	RuntimeDateWeek  RuntimeDateType = 3 // 按周
	RuntimeDateMonth RuntimeDateType = 4 // 按月
)

func (r RuntimeDateType) SqlTime() string {
	switch r {
	case RuntimeDateHour:
		return "date_sub(now(), interval 1 hour)"
	case RuntimeDateDay:
		return "date_sub(now(), interval 1 day)"
	case RuntimeDateWeek:
		return "date_sub(now(), interval 7 day)"
	case RuntimeDateMonth:
		return "date_sub(now(), interval 1 month)"
	}
	return "date_sub(now(), interval 1 day)"
}

type RuntimeLog struct {
	serviceName     string
	level           enum.LogLevelType
	title           string
	itemList        []string
	runtimeDateType RuntimeDateType
}

func (r *RuntimeLog) SetTitle(title string) {
	r.title = title
}

func (r *RuntimeLog) SetLevel(level enum.LogLevelType) {
	r.level = level
}

func (r *RuntimeLog) SetItem(key string, value any, logLevelType enum.LogLevelType) {
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
	r.itemList = append(r.itemList, str)
}

func (r *RuntimeLog) SetItemInfo(key string, value any) {
	r.SetItem(key, value, enum.LogInfoLevel)
}

func (r *RuntimeLog) SetItemWarn(key string, value any) {
	r.SetItem(key, value, enum.LogWarnLevel)
}

func (r *RuntimeLog) SetNowTime() {
	str := fmt.Sprintf("<div class=\"now-time\">%s</div>", time.Now().Format("2006-01-02 15:04:05"))
	r.itemList = append(r.itemList, str)
}

func (r *RuntimeLog) Save() {
	r.SetNowTime()

	var log models.LogModel
	conditionStr := fmt.Sprintf("where service_name = ? and log_type = ? and created_at > %s", r.runtimeDateType.SqlTime())
	global.DB.Find(&log, conditionStr, r.serviceName, enum.RuntimeLogType)
	if log.ID != 0 {
		// 已经存在，需要更新
		c := strings.Join(r.itemList, "\n")
		newContent := log.Content + "\n" + c
		global.DB.Model(&models.LogModel{}).Updates(map[string]any{
			"content": newContent,
		})
		r.itemList = []string{}
		return
	}

	err := global.DB.Save(&models.LogModel{
		LogType:     enum.RuntimeLogType,
		Title:       r.title,
		Content:     strings.Join(r.itemList, "\n"),
		Level:       r.level,
		ServiceName: r.serviceName,
	}).Error
	if err != nil {
		logrus.Errorf("创建运行日志失败 %s", err)
		return
	}
	r.itemList = []string{}
}

func NewRuntimeLog(serviceName string, dateType RuntimeDateType) *RuntimeLog {
	return &RuntimeLog{
		serviceName:     serviceName,
		runtimeDateType: dateType,
	}
}
