package article_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/utils/jwt"

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
