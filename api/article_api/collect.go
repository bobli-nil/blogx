package article_api

import (
	"blogx_server/common"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/utils/jwt"
	"fmt"

	"github.com/gin-gonic/gin"
)

type CollectCreateRequest struct {
	ID       uint   `json:"id"`
	Title    string `json:"title" binding:"required,max=32"`
	Abstract string `json:"abstract" binding:"max=256"`
	Cover    string `json:"cover"`
}

func (ArticleApi) CollectCreate(c *gin.Context) {
	cr := middleware.GetBind[CollectCreateRequest](c)
	claims := jwt.GetClaims(c)

	// 创建
	var model models.CollectModel
	if cr.ID == 0 {
		err := global.DB.Take(&model, "user_id = ? and title = ?", claims.UserID, cr.Title).Error
		if err == nil {
			res.FailWithMsg("该收藏夹已存在", c)
			return
		}
		err = global.DB.Create(&models.CollectModel{
			UserID:   claims.UserID,
			Title:    cr.Title,
			Abstract: cr.Abstract,
			Cover:    cr.Cover,
		}).Error
		if err != nil {
			res.FailWithMsg("创建收藏夹失败", c)
			return
		}
		res.OkWithMsg("创建收藏夹成功", c)
		return
	}

	// 更新分类
	var collect models.CollectModel
	err := global.DB.Take(&collect, "user_id = ? and id = ?", claims.UserID, cr.ID).Error
	if err != nil {
		res.FailWithMsg("该收藏夹不存在", c)
		return
	}
	err = global.DB.Model(&collect).Updates(map[string]any{
		"title":    cr.Title,
		"abstract": cr.Abstract,
		"cover":    cr.Cover,
	}).Error
	if err != nil {
		res.FailWithMsg("更新收藏夹错误", c)
		return
	}
	res.OkWithMsg("更新收藏夹成功", c)
}

type CollectListRequest struct {
	common.PageInfo
	UserID uint  `form:"userID"`
	Type   uint8 `form:"type" binding:"required,oneof=1 2 3"` // 1查自己 2查别人 3管理员
}

type CollectListResponse struct {
	models.CollectModel
	ArticleCount int    `json:"articleCount"`
	Nickname     string `json:"nickname,omitempty"`
	Avatar       string `json:"avatar,omitempty"`
}

func (ArticleApi) CollectListView(c *gin.Context) {
	cr := middleware.GetBind[CollectListRequest](c)

	preloads := []string{"ArticleList"}

	switch cr.Type {
	case 1:
		claims, err := jwt.ParseTokenByGin(c)
		if err != nil {
			res.FailWithError(err, c)
			return
		}
		cr.UserID = claims.UserID
	case 2:
		var userConf models.UserConfModel
		err := global.DB.Take(&userConf, "user_id = ?", cr.UserID).Error
		if err != nil {
			res.FailWithMsg("用户不存在", c)
			return
		}
		if !userConf.OpenCollect {
			res.FailWithMsg("用户未开启我的收藏", c)
			return
		}

	case 3:
		claims, err := jwt.ParseTokenByGin(c)
		if err != nil {
			res.FailWithError(err, c)
			return
		}
		if claims.Role != enum.AdminRole {
			res.FailWithMsg("权限错误", c)
			return
		}
		preloads = append(preloads, "UserModel")
	}

	_list, count, _ := common.ListQuery(models.CollectModel{
		UserID: cr.UserID,
	}, common.Options{
		Debug:    true,
		PageInfo: cr.PageInfo,
		Likes:    []string{"title"},
		PreLoads: preloads,
	})

	list := make([]CollectListResponse, 0)
	for _, model := range _list {
		list = append(list, CollectListResponse{
			CollectModel: model,
			ArticleCount: len(model.ArticleList),
			Nickname:     model.UserModel.Nickname,
			Avatar:       model.UserModel.Avatar,
		})
	}

	res.OkWithList(list, count, c)
}

func (ArticleApi) CollectRemoveView(c *gin.Context) {
	cr := middleware.GetBind[models.DeleteRequest](c)
	claims := jwt.GetClaims(c)

	query := global.DB.Where("id in ?", cr.IDList)
	if claims.Role != enum.AdminRole {
		query = query.Where("user_id = ?", claims.UserID)
	}

	var collectList []models.CollectModel
	global.DB.Where(query).Find(&collectList)
	if len(collectList) > 0 {
		if err := global.DB.Delete(&collectList).Error; err != nil {
			res.FailWithMsg(err.Error(), c)
			return
		}
	}

	msg := fmt.Sprintf("删除成功，成功删除%d条数据", len(collectList))
	res.OkWithMsg(msg, c)
}
