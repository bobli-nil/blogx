package user_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/service/email_service"
	"blogx_server/utils"
	"blogx_server/utils/email_store"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type SendEmailRequest struct {
	Type  int    `json:"type" binding:"oneof=1 2 3"` // 1注册 2重置 3绑定邮箱
	Email string `json:"email" binding:"required"`
}

type SendEmailResponse struct {
	EmailID string `json:"emailID"`
}

func (UserApi) SendEmailView(c *gin.Context) {
	var cr SendEmailRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	code := utils.GenerateRandomDigitsSimple(4)
	uuidStr := uuid.New().String()

	switch cr.Type {
	case 1:
		err = global.DB.Take(&models.UserModel{}, "email = ?", cr.Email).Error
		if err == nil {
			res.FailWithMsg("该邮箱已存在", c)
			return
		}
		err = email_service.SendRegisterCode(cr.Email, code)
	case 2:
		var user models.UserModel
		err = global.DB.Take(&user, "email = ?", cr.Email).Error
		if err != nil {
			res.FailWithMsg("该邮箱不存在", c)
			return
		}
		if user.RegisterSource != enum.RegisterEmailSourceType {
			res.FailWithMsg("非邮箱注册无法重置密码", c)
			return
		}

		err = email_service.SendResetPwdCode(cr.Email, code)
	case 3:
		err = global.DB.Take(&models.UserModel{}, "email = ?", cr.Email).Error
		if err == nil {
			res.FailWithMsg("该邮箱已使用", c)
			return
		}
		err = email_service.SendBindEmailCode(cr.Email, code)
	}

	if err != nil {
		res.FailWithError(err, c)
		logrus.Errorf("邮件发送失败 %s", err.Error())
		return
	}

	email_store.Set(uuidStr, cr.Email, code)

	res.OkWithData(SendEmailResponse{
		EmailID: uuidStr,
	}, c)
}
