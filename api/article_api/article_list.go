package article_api

import (
	"blogx_server/common"
	"blogx_server/common/res"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/utils/jwt"

	"github.com/gin-gonic/gin"
)

type ArticleListRequest struct {
	common.PageInfo
	Type       int8               `form:"type" binding:"required,oneof=1 2 3"` //1用户查别人的 2查自己的 3管理员查所有
	UserID     uint               `form:"userID"`
	CategoryID *uint              `form:"categoryID"`
	Status     enum.ArticleStatus `form:"status" binding:"oneof=1 2 3"`
}

type ArticleListResponse struct {
	models.ArticleModel
	UserTop  bool `json:"userTop"`  // 是否是用户置顶
	AdminTop bool `json:"adminTop"` // 是否是管理员置顶
}

func (ArticleApi) ArticleListView(c *gin.Context) {
	cr := middleware.GetBind[ArticleListRequest](c)
	switch cr.Type {
	case 1:
		// 查别人，用户参数就是必填的
		if cr.Type == 0 {
			res.FailWithMsg("用户ID必填", c)
			return
		}
		cr.Status = 0
	case 2:
		claims, err := jwt.ParseTokenByGin(c)
		if err != nil {
			res.FailWithMsg("请登录", c)
			return
		}
		cr.UserID = claims.UserID
	case 3:
		claims, err := jwt.ParseTokenByGin(c)
		if !(err == nil && claims.Role == enum.AdminRole) {
			res.FailWithMsg("角色错误", c)
			return
		}
	}

	_list, count, err := common.ListQuery(models.ArticleModel{
		UserID:     cr.UserID,
		CategoryID: cr.CategoryID,
		Status:     cr.Status,
	}, common.Options{
		Debug:    true,
		Likes:    []string{"title"},
		PageInfo: cr.PageInfo,
	})
	if err != nil {
		res.FailWithMsg(err.Error(), c)
		return
	}

	var list = make([]ArticleListResponse, 0)
	for _, model := range _list {
		model.Content = ""
		list = append(list, ArticleListResponse{
			ArticleModel: model,
			UserTop:      false,
			AdminTop:     false,
		})
	}

	res.OkWithList(list, count, c)
}
