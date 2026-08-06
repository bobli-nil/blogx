package user_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/utils/jwt"
	"blogx_server/utils/mps"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type UserInfoUpdateRequest struct {
	Username    *string   `json:"username" s-u:"username"`
	Nickname    *string   `json:"nickname" s-u:"nickname"`
	Avatar      *string   `json:"avatar" s-u:"avatar"`
	Abstract    *string   `json:"abstract" s-u:"abstract"`
	LikeTags    *[]string `json:"likeTags" s-u-c:"like_tags"`
	OpenCollect *bool     `json:"openCollect" s-u-c:"open_collect"`
	OpenFollow  *bool     `json:"openFollow" s-u-c:"open_follow"`
	OpenFans    *bool     `json:"openFans" s-u-c:"open_fans"`
	HomeStyleID *uint     `json:"homeStyleId" s-u-c:"home_style_id"`
}

func (UserApi) UserInfoUpdateView(c *gin.Context) {
	var cr UserInfoUpdateRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	claims := jwt.GetClaims(c)

	userMap := mps.Struct2Map(cr, "s-u")
	userConfMap := mps.Struct2Map(cr, "s-u-c")
	fmt.Printf("%+v \n", userMap)
	fmt.Printf("%+v \n", userConfMap)
	if len(userMap) > 0 {
		var user models.UserModel
		if err := global.DB.Preload("UserConfModel").Take(&user, claims.UserID).Error; err != nil {
			res.FailWithMsg("用户不存在", c)
			return
		}

		// 判断新的用户名不能重复
		if cr.Username != nil {
			var userCount int64
			global.DB.Model(&models.UserModel{}).
				Where("username = ? and id <> ?", *cr.Username, claims.UserID).
				Count(&userCount)
			if userCount > 0 {
				res.FailWithMsg("用户名已被使用", c)
				return
			}
		}

		// 用户名30天内只能更新一次
		if usernameUpdateTime := user.UserConfModel.UpdateUsernameDate; usernameUpdateTime != nil {
			if time.Now().Sub(*usernameUpdateTime).Hours() < 720 {
				res.FailWithMsg("用户名30天内只能更新一次", c)
				return
			}
		}
		userConfMap["update_username_date"] = time.Now()

		// QQ用户不能修改昵称和头像
		if (cr.Avatar != nil || cr.Nickname != nil) && user.RegisterSource == enum.RegisterQQSourceType {
			res.FailWithMsg("QQ用户不能修改昵称和头像", c)
			return
		}

		if err = global.DB.Model(&user).Updates(userMap).Error; err != nil {
			res.FailWithMsg("更新失败", c)
			return
		}

	}

	if len(userConfMap) > 0 {
		var userConf models.UserConfModel
		if err := global.DB.Take(&userConf, "user_id = ?", claims.UserID).Error; err != nil {
			res.FailWithMsg("用户配置信息不存在", c)
			return
		}
		if err = global.DB.Model(&userConf).Updates(userConfMap).Error; err != nil {
			res.FailWithMsg("更新失败", c)
			return
		}
	}

	res.OkWithMsg("更新成功", c)
}
