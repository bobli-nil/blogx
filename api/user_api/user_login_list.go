package user_api

import (
	"blogx_server/common"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/utils/jwt"
	"time"

	"github.com/gin-gonic/gin"
)

type UserLoginListRequest struct {
	common.PageInfo
	IP        string `form:"ip"`
	StartTime string `form:"startTime"` // 时间戳 秒级
	EndTime   string `form:"endTime"`   // 时间戳 秒级
}

type UserLoginListResponse struct {
	models.UserLoginModel
	UserNickName string `json:"userNickname"`
	UserAvatar   string `json:"userAvatar"`
}

func (UserApi) UserLoginList(c *gin.Context) {
	var cr UserLoginListRequest
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	claims := jwt.GetClaims(c)
	if claims == nil {
		res.FailWithMsg("获取用户信息失败", c)
		return
	}
	role := claims.Role
	query := global.DB.Where("")
	if role == enum.UserRole {
		query = global.DB.Where("user_id = ?", claims.UserID)
	}
	loc, _ := time.LoadLocation("Asia/Shanghai")
	if cr.StartTime != "" {
		startStr, err := time.ParseInLocation("2006-01-02 15:04:05", cr.StartTime, loc)
		if err != nil {
			res.FailWithMsg("开始时间格式错误", c)
			return
		}
		query = query.Where("created_at >= ?", startStr)
	}
	if cr.EndTime != "" {
		endStr, err := time.ParseInLocation("2006-01-02 15:04:05", cr.EndTime, loc)
		if err != nil {
			res.FailWithMsg("结束时间格式错误", c)
			return
		}
		query = query.Where("created_at <= ?", endStr)
	}

	logList, count, err := common.ListQuery(&models.UserLoginModel{
		IP: cr.IP,
	}, common.Options{
		PageInfo: cr.PageInfo,
		Debug:    true,
		Likes:    []string{"addr"},
		Where:    query,
		PreLoads: []string{"UserModel"},
	})
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	var list []UserLoginListResponse
	for _, log := range logList {
		list = append(list, UserLoginListResponse{
			UserLoginModel: *log,
			UserNickName:   log.UserModel.Nickname,
			UserAvatar:     log.UserModel.Avatar,
		})
	}

	res.OkWithList(list, count, c)
}
