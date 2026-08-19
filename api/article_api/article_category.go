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

type CategoryCreateRequest struct {
	ID    uint64 `json:"id"`
	Title string `json:"title"`
}

func (ArticleApi) ArticleCategoryCreate(c *gin.Context) {
	cr := middleware.GetBind[CategoryCreateRequest](c)
	claims := jwt.GetClaims(c)

	// 创建
	if cr.ID == 0 {
		err := global.DB.Take(&models.CategoryModel{}, "user_id = ? and title = ?", claims.UserID, cr.Title).Error
		if err == nil {
			res.FailWithMsg("该分类已存在", c)
			return
		}
		err = global.DB.Create(&models.CategoryModel{
			UserID: claims.UserID,
			Title:  cr.Title,
		}).Error
		if err != nil {
			res.FailWithMsg("创建分类失败", c)
			return
		}
		res.OkWithMsg("创建分类成功", c)
		return
	}

	// 更新分类
	var category models.CategoryModel
	err := global.DB.Take(&category, "user_id = ? and id = ?", claims.UserID, cr.ID).Error
	if err != nil {
		res.FailWithMsg("该分类不存在", c)
		return
	}
	err = global.DB.Model(&category).Update("title", cr.Title).Error
	if err != nil {
		res.FailWithMsg("更新分类错误", c)
		return
	}
	res.OkWithMsg("更新分类成功", c)
}

type CategoryListRequest struct {
	common.PageInfo
	UserID uint  `form:"userID"`
	Type   uint8 `form:"type" binding:"required,oneof=1 2 3"` // 1查自己 2查别人 3管理员
}

type CategoryListResponse struct {
	models.CategoryModel
	ArticleCount int    `json:"articleCount"`
	Nickname     string `json:"nickname,omitempty"`
	Avatar       string `json:"avatar,omitempty"`
}

func (ArticleApi) CategoryListView(c *gin.Context) {
	cr := middleware.GetBind[CategoryListRequest](c)

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

	_list, count, _ := common.ListQuery(models.CategoryModel{
		UserID: cr.UserID,
	}, common.Options{
		Debug:    true,
		PageInfo: cr.PageInfo,
		Likes:    []string{"title"},
		PreLoads: preloads,
	})

	list := make([]CategoryListResponse, 0)
	for _, model := range _list {
		list = append(list, CategoryListResponse{
			CategoryModel: model,
			ArticleCount:  len(model.ArticleList),
			Nickname:      model.UserModel.Nickname,
			Avatar:        model.UserModel.Avatar,
		})
	}

	res.OkWithList(list, count, c)
}

func (ArticleApi) CategoryRemoveView(c *gin.Context) {
	cr := middleware.GetBind[models.DeleteRequest](c)
	claims := jwt.GetClaims(c)

	query := global.DB.Where("id in ?", cr.IDList)
	if claims.Role != enum.AdminRole {
		query = query.Where("user_id = ?", claims.UserID)
	}

	var categoryList []models.CategoryModel
	global.DB.Where(query).Find(&categoryList)
	if len(categoryList) > 0 {
		if err := global.DB.Delete(&categoryList).Error; err != nil {
			res.FailWithMsg(err.Error(), c)
			return
		}
	}

	msg := fmt.Sprintf("删除成功，成功删除%d条数据", len(categoryList))
	res.OkWithMsg(msg, c)
}
