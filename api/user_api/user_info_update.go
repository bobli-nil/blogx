package user_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/utils/jwt"
	"blogx_server/utils/mps"
	"fmt"

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

	userInfo := mps.Struct2Map(cr, "s-u")
	fmt.Printf("%+v \n", userInfo)
	if len(userInfo) > 0 {
		var user models.UserModel
		err := global.DB.Preload("UserConfModel").Take(&user, claims.UserID).Error
		if err != nil {
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
	}
	userConfInfo := mps.Struct2Map(cr, "s-u-c")
	fmt.Printf("%+v \n", userConfInfo)
	if len(userConfInfo) > 0 {

	}
}
