package user_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/utils/jwt"
	"time"

	"github.com/gin-gonic/gin"
)

type UserDetailResponse struct {
	ID             uint                  `json:"id"`
	CreatedAt      *time.Time            `json:"createdAt"`
	Username       string                `json:"username"`
	Nickname       string                `json:"nickname"`
	Avatar         string                `json:"avatar"`
	Abstract       string                `json:"abstract"`
	RegisterSource enum.RegisterSource   `json:"registerSource"`
	LikeTags       []string              `json:"likeTags"`
	UserConf       *models.UserConfModel `json:"userConf"`
}

func (UserApi) UserDetailView(c *gin.Context) {
	claims := jwt.GetClaims(c)
	id := claims.UserID

	var user models.UserModel
	err := global.DB.Debug().Preload("UserConfModel").Where("id = ?", id).Take(&user).Error
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	data := UserDetailResponse{
		ID:             user.ID,
		CreatedAt:      user.CreatedAt,
		Username:       user.Username,
		Nickname:       user.Nickname,
		Avatar:         user.Avatar,
		Abstract:       user.Abstract,
		RegisterSource: user.RegisterSource,
	}
	if user.UserConfModel != nil {
		data.UserConf = user.UserConfModel
	}

	res.OkWithData(data, c)

}
