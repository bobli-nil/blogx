package log_api

import (
	"blogx_server/common"
	"blogx_server/common/res"
	"blogx_server/models"
	"blogx_server/models/enum"
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

	//if cr.Page < 0 {
	//	cr.Page = 1
	//}
	//if cr.Limit <= 0 {
	//	cr.Limit = 10
	//}
	//offset := (cr.Page - 1) * cr.Limit
	//
	//var list []models.LogModel
	//condition := models.LogModel{
	//	LogType:     cr.LogType,
	//	Level:       cr.Level,
	//	UserID:      cr.UserID,
	//	IP:          cr.IP,
	//	LoginStatus: cr.LoginStatus,
	//	ServiceName: cr.ServiceName,
	//}
	//like := global.DB.Where("title LIKE ?", "%"+cr.Keyword+"%")
	//global.DB.Preload("UserModel").Debug().Where(condition).Where(like).Offset(offset).Limit(cr.Limit).Find(&list)
	//var count int64
	//global.DB.Debug().Where(condition).Where(like).Model(models.LogModel{}).Count(&count)

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
