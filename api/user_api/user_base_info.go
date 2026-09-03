package user_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/service/redis_service/redis_user"

	"github.com/gin-gonic/gin"
)

type UserBaseInfoResponse struct {
	UserID       uint   `json:"userID"`
	CodeAge      int    `json:"codeAge"`
	Avatar       string `json:"avatar"`
	NickName     string `json:"nickName"`
	LookCount    int    `json:"lookCount"`
	ArticleCount int    `json:"articleCount"`
	FansCount    int    `json:"fansCount"`
	FollowCount  int    `json:"followCount"`
	Place        string `json:"place"`       // IP归属地
	OpenCollect  bool   `json:"openCollect"` // 公开我的收藏
	OpenFollow   bool   `json:"openFollow"`  // 公开我的关注
	OpenFans     bool   `json:"openFans"`    // 公开我的粉丝
	HomeStyleID  uint   `json:"homeStyleID"` // 主页样式ID
}

func (UserApi) UserBaseInfoView(c *gin.Context) {
	var cr models.IDRequest
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	var user models.UserModel
	err = global.DB.Preload("UserConfModel").Preload("ArticleList").Take(&user, cr.ID).Error
	if err != nil {
		res.FailWithMsg("用户不存在", c)
		return
	}

	data := UserBaseInfoResponse{
		UserID:       user.ID,
		CodeAge:      user.CodeAge(),
		Avatar:       user.Avatar,
		NickName:     user.Nickname,
		LookCount:    user.UserConfModel.LookCount + redis_user.GetCacheLook(cr.ID),
		ArticleCount: len(user.ArticleList),
		FansCount:    0,
		FollowCount:  0,
		Place:        user.Avatar,
		OpenCollect:  user.UserConfModel.OpenCollect,
		OpenFollow:   user.UserConfModel.OpenFollow,
		OpenFans:     user.UserConfModel.OpenFans,
		HomeStyleID:  user.UserConfModel.HomeStyleID,
	}

	var focusList []models.UserFocusModel
	global.DB.Find(&focusList, "user_id = ? or focus_user_id = ?", cr.ID, cr.ID)
	for _, model := range focusList {
		if model.UserID == cr.ID {
			data.FollowCount += 1
		}
		if model.FocusUserID == cr.ID {
			data.FansCount += 1
		}
	}

	redis_user.SetCacheLook(cr.ID, true)

	res.OkWithData(data, c)
}
