package log_api

import (
	"blogx_server/common"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/service/log_service"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type LogApi struct{}

type LogListRequest struct {
	common.PageInfo
	LogType     enum.LogType      `form:"logType"`
	Level       enum.LogLevelType `form:"level"`
	UserID      uint              `form:"userID"`
	IP          string            `form:"ip"`
	LoginStatus bool              `form:"loginStatus"`
	ServiceName string            `form:"serviceName"`
}

type LogListResponse struct {
	models.LogModel
	UserName string `json:"userName"`
	NickName string `json:"nickName"`
}

// LogListView 日志列表查询接口
func (LogApi) LogListView(c *gin.Context) {
	fmt.Println("--->", c.Request.URL)
	var cr LogListRequest
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		logrus.Errorf("查询日志列表参数绑定失败 %s", err.Error())
		res.FailWithError(err, c)
		return
	}
	fmt.Printf("入参 %+v\n", cr)

	list, count, err := common.ListQuery(models.LogModel{
		LogType:     cr.LogType,
		Level:       cr.Level,
		UserID:      cr.UserID,
		IP:          cr.IP,
		LoginStatus: cr.LoginStatus,
		ServiceName: cr.ServiceName,
	}, common.Options{
		PageInfo: cr.PageInfo,
		Likes:    []string{"title"},
		PreLoads: []string{"UserModel"},
		Debug:    true,
	})
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	_list := []LogListResponse{}
	for _, log := range list {
		_list = append(_list, LogListResponse{
			LogModel: log,
			UserName: log.UserModel.Username,
			NickName: log.UserModel.Nickname,
		})
	}

	res.OkWithList(_list, int(count), c)
}

// LogReadView 日志读取与更新接口
func (LogApi) LogReadView(c *gin.Context) {
	var cr models.IDRequest
	err := c.ShouldBindUri(&cr)
	if err != nil {
		logrus.Error("日志读取参数绑定失败：%s", err.Error())
		res.FailWithError(err, c)
		return
	}
	var log models.LogModel
	err = global.DB.Where("id = ?", cr.ID).Take(&log).Error
	if err != nil {
		res.FailWithError(err, c)
		logrus.Errorf("查询读取记录出错:%s", err.Error())
		return
	}
	if !log.IsRead {
		global.DB.Model(&log).Update("is_read", true)
	}
	res.OkWithMsg("更新成功", c)
}

// LogDeleteView 日志删除接口
func (LogApi) LogDeleteView(c *gin.Context) {
	actionLog := log_service.GetLogFromGinContext(c)
	actionLog.SetRequest()
	actionLog.SetTitle("删除日志")
	actionLog.SetLevel(enum.LogInfoLevel)

	var cr models.DeleteRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, c)
		logrus.Errorf("删除日志参数绑定失败: %s", err.Error())
		return
	}

	result := global.DB.Where("id in ?", cr.IDList).Delete(&models.LogModel{})
	if result.Error != nil {
		res.FailWithError(result.Error, c)
		logrus.Error("日志删除失败: %s", result.Error.Error())
		return
	}

	tip := fmt.Sprintf("共删除%d日志", result.RowsAffected)
	res.OkWithMsg(tip, c)
}
