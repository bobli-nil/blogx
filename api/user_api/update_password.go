package user_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models/enum"
	"blogx_server/utils/jwt"
	pwd2 "blogx_server/utils/pwd"

	"github.com/gin-gonic/gin"
)

type UpdatePasswordRequest struct {
	OldPassword string `json:"oldPassword"`
	NewPassword string `json:"newPassword"`
}

func (UserApi) UpdatePassword(c *gin.Context) {
	claims := jwt.GetClaims(c)
	if claims == nil {
		res.FailWithMsg("获取登录状态失败", c)
		return
	}

	var cr UpdatePasswordRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	user, err := claims.GetUser()
	if err != nil {
		res.FailWithMsg("用户不存在", c)
		return
	}

	if !(user.RegisterSource == enum.RegisterEmailSourceType || user.Email != "") {
		res.FailWithMsg("仅支持邮箱注册或绑定邮箱的账号进行密码修改", c)
		return
	}

	valid := pwd2.CompareHashAndPassword(user.Password, cr.OldPassword)
	if !valid {
		res.FailWithMsg("原密码不正确", c)
		return
	}

	newPassword, _ := pwd2.GenerateFromPassword(cr.NewPassword)
	err = global.DB.Model(&user).Update("password", newPassword).Error
	if err != nil {
		res.FailWithMsg("更新失败", c)
		return
	}

	res.OkWithMsg("更新成功", c)

}
