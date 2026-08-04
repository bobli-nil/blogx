package user_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/service/email_service"
	"blogx_server/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type UserApi struct{}

type SendEmailRequest struct {
	Type  int    `json:"type" binding:"oneof=1 2"`
	Email string `json:"email" binding:"required"`
}

type SendEmailResponse struct {
	EmailId string `json:"emailId"`
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
		err = email_service.SendResetPwdCode(cr.Email, code)
	}
	if err != nil {
		res.FailWithError(err, c)
		logrus.Errorf("邮件发送失败 %s", err.Error())
		return
	}

	global.Store.Set(uuidStr, code)

	res.OkWithData(SendEmailResponse{
		EmailId: uuidStr,
	}, c)
}
