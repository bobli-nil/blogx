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
	UserID      uint `form:"userID"` // 查某个用户的关注
}

type FocusUserListResponse struct {
	FocusUserID       uint       `json:"focusUserID"`
	FocusUserNickname string     `json:"focusUserNickname"`
	FocusUserAvatar   string     `json:"focusUserAvatar"`
	FocusUserAbstract string     `json:"focusUserAbstract"`
	CreatedAt         *time.Time `json:"createdAt"`
}

// FocusUserListView 我的关注/用户的关注
func (FocusApi) FocusUserListView(c *gin.Context) {
	cr := middleware.GetBind[FocusUserListRequest](c)

	if cr.UserID == 0 {
		claims, err := jwt.ParseTokenByGin(c)
		if err != nil || claims == nil {
			res.FailWithMsg("请登录", c)
			return
		}
		cr.UserID = claims.UserID
	} else {
		var userConf models.UserConfModel
		if err := global.DB.Take(&userConf, "user_id = ?", cr.UserID).Error; err != nil {
			res.FailWithMsg("用户配置不存在", c)
			return
		}
		if !userConf.OpenFollow {
			res.FailWithMsg("该用户未公开我的关注", c)
			return
		}
	}

	_list, count, _ := common.ListQuery(models.UserFocusModel{
		FocusUserID: cr.FocusUserID,
		UserID:      cr.UserID,
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

type FansUserListRequest struct {
	common.PageInfo
	FansUserID uint `form:"focusUserID"`
	UserID     uint `form:"userID"` // 查某个用户的粉丝
}

type FansUserListResponse struct {
	FansUserID       uint       `json:"fansUserID"`
	FansUserNickname string     `json:"fansUserNickname"`
	FansUserAvatar   string     `json:"fansUserAvatar"`
	FansUserAbstract string     `json:"fansUserAbstract"`
	CreatedAt        *time.Time `json:"createdAt"`
}

// FansUserListView 我的粉丝/用户的粉丝
func (FocusApi) FansUserListView(c *gin.Context) {
	cr := middleware.GetBind[FansUserListRequest](c)

	if cr.UserID == 0 {
		claims, err := jwt.ParseTokenByGin(c)
		if err != nil || claims == nil {
			res.OkWithMsg("请登录", c)
			return
		}
		cr.UserID = claims.UserID
	} else {
		var userConf models.UserConfModel
		if err := global.DB.Take(&userConf, "user_id = ?", cr.UserID).Error; err != nil {
			res.FailWithMsg("用户配置不存在", c)
			return
		}
		if !userConf.OpenFans {
			res.FailWithMsg("用户未公开我的粉丝", c)
			return
		}
	}

	_list, count, _ := common.ListQuery(models.UserFocusModel{
		FocusUserID: cr.UserID,
		UserID:      cr.FansUserID,
	}, common.Options{
		PageInfo: cr.PageInfo,
		PreLoads: []string{"UserModel"},
	})

	var list = make([]FansUserListResponse, 0)
	for _, model := range _list {
		list = append(list, FansUserListResponse{
			FansUserID:       model.UserID,
			FansUserNickname: model.UserModel.Nickname,
			FansUserAvatar:   model.UserModel.Avatar,
			FansUserAbstract: model.UserModel.Abstract,
			CreatedAt:        model.CreatedAt,
		})
	}

	res.OkWithList(list, count, c)
}

// UnfocusUserView 登录人取关
func (FocusApi) UnfocusUserView(c *gin.Context) {
	cr := middleware.GetBind[FocusUserRequest](c)
	claims := jwt.GetClaims(c)

	if claims.UserID == cr.FocusUserID {
		res.FailWithMsg("你不能取关自己", c)
		return
	}

	var user models.UserModel
	if err := global.DB.Take(&user, cr.FocusUserID).Error; err != nil {
		res.FailWithMsg("取关的用户不存在", c)
		return
	}

	var focusUser models.UserFocusModel
	err := global.DB.Take(&focusUser, "user_id = ? and focus_user_id = ?", claims.UserID, cr.FocusUserID).Error
	if err != nil {
		res.FailWithMsg("未关注此用户", c)
		return
	}

	global.DB.Delete(&focusUser)
	res.OkWithMsg("取消关注成功", c)
	return
}
