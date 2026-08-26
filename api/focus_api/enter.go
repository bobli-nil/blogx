package focus_api

import (
	"blogx_server/common"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/utils/jwt"
	"time"

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

type FocusUserListRequest struct {
	common.PageInfo
	FocusUserID uint `form:"focusUserID"`
}

type FocusUserListResponse struct {
	FocusUserID       uint       `json:"focusUserID"`
	FocusUserNickname string     `json:"focusUserNickname"`
	FocusUserAvatar   string     `json:"focusUserAvatar"`
	FocusUserAbstract string     `json:"focusUserAbstract"`
	CreatedAt         *time.Time `json:"createdAt"`
}

func (FocusApi) FocusUserListView(c *gin.Context) {
	cr := middleware.GetBind[FocusUserListRequest](c)
	claims := jwt.GetClaims(c)

	_list, count, _ := common.ListQuery(models.UserFocusModel{
		FocusUserID: cr.FocusUserID,
		UserID:      claims.UserID,
	}, common.Options{
		PageInfo: cr.PageInfo,
		PreLoads: []string{"FocusUserModel"},
	})

	var list = make([]FocusUserListResponse, 0)
	for _, model := range _list {
		list = append(list, FocusUserListResponse{
			FocusUserID:       model.FocusUserID,
			FocusUserNickname: model.FocusUserModel.Nickname,
			FocusUserAvatar:   model.FocusUserModel.Avatar,
			FocusUserAbstract: model.FocusUserModel.Abstract,
			CreatedAt:         model.CreatedAt,
		})
	}

	res.OkWithList(list, count, c)
}
