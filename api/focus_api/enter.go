package focus_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/utils/jwt"

	"github.com/gin-gonic/gin"
)

type FocusApi struct{}

type FocusUserRequest struct {
	FocusUserID uint `json:"focusUserID"`
}

func (FocusApi) FocusUserView(c *gin.Context) {
	cr := middleware.GetBind[FocusUserRequest](c)
	claims := jwt.GetClaims(c)

	// 不能自己关注自己
	if claims.UserID == cr.FocusUserID {
		res.FailWithMsg("不能关注自己", c)
		return
	}

	// 关注的用户是否存在
	var focusUser models.UserModel
	err := global.DB.Take(&focusUser, cr.FocusUserID).Error
	if err != nil {
		res.FailWithMsg("用户不存在", c)
		return
	}

	// 是否已经关注了他
	var focus models.UserFocusModel
	err = global.DB.Take(&focus, "user_id = ? and focus_user_id = ?", claims.UserID, cr.FocusUserID).Error
	if err == nil {
		res.FailWithMsg("已关注该用户", c)
		return
	}

	global.DB.Create(&models.UserFocusModel{
		UserID:      claims.UserID,
		FocusUserID: cr.FocusUserID,
	})

	res.OkWithMsg("关注成功", c)
	return
}
