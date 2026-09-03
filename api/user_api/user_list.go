package user_api

import (
	"blogx_server/common"
	"blogx_server/common/res"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/models/enum"
	"time"

	"github.com/gin-gonic/gin"
)

type UserListRequest struct {
	common.PageInfo
}

type UserListResponse struct {
	ID            uint          `json:"id"`
	Nickname      string        `json:"nickname"`
	Username      string        `json:"username"`
	Avatar        string        `json:"avatar"`
	IP            string        `json:"ip"`
	Addr          string        `json:"addr"`
	ArticleCount  int           `json:"articleCount"`
	IndexCount    uint          `json:"indexCount"`    // 主页访问数
	CreatedAt     *time.Time    `json:"createdAt"`     // 注册时间
	LastLoginDate *time.Time    `json:"lastLoginDate"` // 最后登录时间
	Role          enum.RoleType `json:"role"`
}

func (UserApi) UserListView(c *gin.Context) {
	cr := middleware.GetBind[UserListRequest](c)

	_list, count, _ := common.ListQuery(models.UserModel{}, common.Options{
		PageInfo: cr.PageInfo,
		Likes:    []string{"nickname", "username"},
		PreLoads: []string{"ArticleList", "LoginList"},
	})

	var list = make([]UserListResponse, 0)
	for _, model := range _list {
		item := UserListResponse{
			ID:           model.ID,
			Nickname:     model.Nickname,
			Username:     model.Username,
			Avatar:       model.Avatar,
			ArticleCount: len(model.ArticleList),
			CreatedAt:    model.CreatedAt,
			Role:         model.Role,
		}
		if len(model.LoginList) > 0 {
			item.IP = model.LoginList[0].IP
			item.Addr = model.LoginList[0].Addr
			item.LastLoginDate = model.LoginList[len(model.LoginList)-1].CreatedAt
		}
		list = append(list, item)
	}

	res.OkWithList(list, count, c)
}
