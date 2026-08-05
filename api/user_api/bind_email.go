package user_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/utils/jwt"

	"github.com/gin-gonic/gin"
)

// TODO 未测试
func (UserApi) BindEmailView(c *gin.Context) {
	myClaim := jwt.GetClaims(c)
	value, _ := c.Get("email")
	email := value.(string)

	var user models.UserModel
	err := global.DB.Take(&user, myClaim.UserID).Error
	if err != nil {
		res.FailWithMsg("用户不存在", c)
		return
	}

	err = global.DB.Model(&user).Update("email", email).Error
	if err != nil {
		res.FailWithMsg("绑定失败", c)
		return
	}

	res.OkWithMsg("绑定成功", c)
}
