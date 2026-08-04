package user_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/utils"
	"blogx_server/utils/email_store"
	"blogx_server/utils/jwt"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type RegisterEmailRequest struct {
	EmailID string `json:"emailID" binding:"required"`
	Code    string `json:"code" binding:"required"`
	Pwd     string `json:"pwd" binding:"required"`
}

type RegisterEmailResponse struct {
	Token string `json:"token"`
}

func (UserApi) RegisterEmailView(c *gin.Context) {
	var cr RegisterEmailRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	value, ok := global.EmailVerifyStore.Load(cr.EmailID)
	global.EmailVerifyStore.Delete(cr.EmailID)
	if !ok {
		res.FailWithMsg("邮箱验证失败", c)
		return
	}
	emailInfo, ok := value.(email_store.EmailStoreInfo)
	if !ok {
		res.FailWithMsg("邮箱验证失败", c)
		return
	}
	if emailInfo.Code != cr.Code {
		res.FailWithMsg("邮箱验证失败", c)
		return
	}

	// 入库
	uname := utils.GenerateRandomDigitsSimple(4)
	var user = models.UserModel{
		Username:       fmt.Sprintf("b_%s", uname),
		Nickname:       "邮箱用户",
		RegisterSource: enum.RegisterEmailSourceType,
		Password:       "",
		Email:          emailInfo.Email,
		Role:           enum.UserRole,
	}
	err = global.DB.Create(&user).Error
	if err != nil {
		res.FailWithMsg("注册失败", c)
		logrus.Errorf("邮箱注册失败 %s", err)
		return
	}

	// 颁发token
	token, err := jwt.GenerateToken(user.ID, user.Username, user.Role)
	if err != nil {
		res.FailWithMsg("邮箱登录失败", c)
		return
	}
	res.OkWithData(RegisterEmailResponse{Token: token}, c)
}
