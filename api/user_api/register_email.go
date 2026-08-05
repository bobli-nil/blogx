package user_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/service/user_service"
	"blogx_server/utils"
	"blogx_server/utils/jwt"
	pwd2 "blogx_server/utils/pwd"
	"fmt"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RegisterEmailRequest struct {
	Pwd string `json:"pwd" binding:"required"`
}

type RegisterEmailResponse struct {
	Token string `json:"token"`
}

func (UserApi) RegisterEmailView(c *gin.Context) {
	var cr RegisterEmailRequest
	if err := c.ShouldBindJSON(&cr); err != nil {
		res.FailWithError(err, c)
		return
	}

	if !global.Conf.Site.Login.EmailPwdLogin {
		res.FailWithMsg("邮箱密码登录未开启", c)
		return
	}

	// 再 userModel和userConfModel里分别创建一条数据
	uname := utils.GenerateRandomDigitsSimple(4)
	pwd, _ := pwd2.GenerateFromPassword(cr.Pwd)
	value, ok := c.Get("email")
	if !ok {
		res.FailWithMsg("获取email失败", c)
		return
	}
	email, ok := value.(string)
	if !ok {
		res.FailWithMsg("获取email失败", c)
		return
	}

	var user = models.UserModel{
		Username:       fmt.Sprintf("b_%s", uname),
		Nickname:       "邮箱用户",
		RegisterSource: enum.RegisterEmailSourceType,
		Password:       pwd,
		Email:          email,
		Role:           enum.UserRole,
	}
	err := global.DB.Transaction(func(tx *gorm.DB) error {
		if err := global.DB.Create(&user).Error; err != nil {
			res.FailWithMsg("注册失败", c)
			return err
		}
		userConf := models.UserConfModel{
			UserID:             user.ID,
			UpdateUsernameDate: user.CreatedAt,
		}
		if err := global.DB.Create(&userConf).Error; err != nil {
			res.FailWithMsg("注册失败", c)
			return err
		}
		return nil
	})
	if err != nil {
		res.FailWithMsg("注册失败", c)
		return
	}

	// 颁发token
	token, err := jwt.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		res.FailWithMsg("邮箱登录失败", c)
		return
	}
	// 记录登录日志
	user_service.NewUserService(user).UserLogin(c)

	res.OkWithData(RegisterEmailResponse{Token: token}, c)
}
