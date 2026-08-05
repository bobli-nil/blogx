package user_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	pwd2 "blogx_server/utils/pwd"

	"github.com/gin-gonic/gin"
)

type ResetPasswordRequest struct {
	Pwd string `json:"pwd" binding:"required"`
}

func (UserApi) ResetPassword(c *gin.Context) {
	var cr ResetPasswordRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	value, _ := c.Get("email")
	email := value.(string)

	var user models.UserModel
	err = global.DB.Take(&user, "email = ?", email).Error
	if err != nil {
		res.FailWithMsg("邮箱不存在", c)
		return
	}

	hashPwd, _ := pwd2.GenerateFromPassword(cr.Pwd)
	err = global.DB.Model(&user).Update("password", hashPwd).Error
	if err != nil {
		res.FailWithMsg("重置密码失败", c)
		return
	}

	res.OkWithMsg("重置密码成功", c)

}
