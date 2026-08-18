package article_api

import (
	"blogx_server/common"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/service/redis_service/redis_article"
	"blogx_server/utils/jwt"
	"blogx_server/utils/sql"
	"fmt"

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

	var orderColumnMap = map[string]bool{
		"look_count desc":    true,
		"look_count asc":     true,
		"digg_count desc":    true,
		"digg_count asc":     true,
		"comment_count desc": true,
		"comment_count asc":  true,
		"collect_count desc": true,
		"collect_count asc":  true,
	}

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

	_, ok := orderColumnMap[cr.Order]
	if !ok && cr.Order != "" {
		res.FailWithMsg("该字段不支持排序", c)
		return
	}

	// 处理置顶
	var userTopMap = map[uint]bool{}
	var adminTopMap = map[uint]bool{}
	var topArticleIDList []uint
	if cr.UserID != 0 {
		var userTopArticle []models.UserTopArticleModel
		// TODO 此处有bug，这里查出来的只有自己置顶的，还有管理员置顶的没查出，后面修复
		if err := global.DB.Debug().Preload("UserModel").Order("created_at desc").Find(&userTopArticle, cr.UserID).Error; err != nil {
			res.FailWithMsg(err.Error(), c)
			return
		}

		for _, model := range userTopArticle {
			topArticleIDList = append(topArticleIDList, model.ArticleID)
			if model.UserModel.Role == enum.AdminRole {
				adminTopMap[model.ArticleID] = true
			}
			userTopMap[model.ArticleID] = true
		}

	}

	var options = common.Options{
		Debug:        true,
		Likes:        []string{"title"},
		PageInfo:     cr.PageInfo,
		DefaultOrder: "created_at desc",
	}
	if len(topArticleIDList) > 0 {
		options.DefaultOrder = fmt.Sprintf("%s, created_at desc", sql.ConvertSliceOrderSql(topArticleIDList))
	}

	_list, count, err := common.ListQuery(models.ArticleModel{
		UserID:     cr.UserID,
		CategoryID: cr.CategoryID,
		Status:     cr.Status,
	}, options)
	if err != nil {
		res.FailWithMsg(err.Error(), c)
		return
	}

	lookMap := redis_article.GetAllCacheLook()
	diggMap := redis_article.GetAllCacheDigg()
	collectMap := redis_article.GetAllCacheCollect()

	var list = make([]ArticleListResponse, 0)
	for _, model := range _list {
		model.Content = ""
		model.CollectCount = model.CollectCount + collectMap[model.ID]
		model.LookCount = model.LookCount + lookMap[model.ID]
		model.DiggCount = model.DiggCount + diggMap[model.ID]
		list = append(list, ArticleListResponse{
			ArticleModel: model,
			UserTop:      userTopMap[model.ID],
			AdminTop:     adminTopMap[model.ID],
		})
	}

	res.OkWithList(list, count, c)
}
