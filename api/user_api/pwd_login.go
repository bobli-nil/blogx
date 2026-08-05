package user_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/service/user_service"
	"blogx_server/utils/jwt"
	pwd2 "blogx_server/utils/pwd"

	"github.com/gin-gonic/gin"
)

type PwdLoginRequest struct {
	Val      string `json:"val" binding:"required"` // 用户名或者邮箱
	Password string `json:"password" binding:"required"`
}

func (UserApi) PwdLogin(c *gin.Context) {
	// 参数绑定
	var cr PwdLoginRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	// 查数据库看有没有该用户，没有就报错
	var user models.UserModel
	err = global.DB.Take(&user, "(username = ? or email = ?) and password is not null and password <> ''", cr.Val, cr.Val).Error
	if err != nil {
		res.FailWithMsg("用户名/邮箱或密码错误", c)
		return
	}

	// 校验加密的密码
	pwdRight := pwd2.CompareHashAndPassword(user.Password, cr.Password)
	if !pwdRight {
		res.FailWithMsg("用户名/邮箱或密码错误", c)
		return
	}

	// 站点是否启用用户名密码登录
	if !global.Conf.Site.Login.UsernamePwdLogin || !global.Conf.Site.Login.EmailPwdLogin {
		res.FailWithMsg("站点未完全启用邮箱密码及用户名密码登录", c)
		return
	}

	// 通过就颁发token
	token, err := jwt.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		res.FailWithMsg("获取token失败", c)
		return
	}

	// 记录登录日志
	user_service.NewUserService(user).UserLogin(c)

	res.OkWithData(token, c)

}
