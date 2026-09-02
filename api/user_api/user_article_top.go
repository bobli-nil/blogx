package user_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/utils/jwt"

	"github.com/gin-gonic/gin"
)

type UserArticleTopRequest struct {
	ArticleID uint `json:"articleID"`
	Type      int8 `json:"type" binding:"required,oneof=1 2"` // 1普通用户 2管理员
}

func (UserApi) UserArticleTopView(c *gin.Context) {
	cr := middleware.GetBind[UserArticleTopRequest](c)

	var model models.ArticleModel
	if err := global.DB.Take(&model, cr.ArticleID).Error; err != nil {
		res.FailWithMsg("文章不存在", c)
		return
	}

	claims := jwt.GetClaims(c)

	switch cr.Type {
	case 1:
		// 用户置顶文章
		// 验证文章是不是自己且已发布
		if model.UserID != claims.UserID {
			res.FailWithMsg("用户只能置顶自己的文章", c)
			return
		}
		if model.Status != enum.ArticleStatusPublished {
			res.FailWithMsg("用户只能置顶已发布的文章", c)
			return
		}
		// 判断之前自己有没有置顶过，用户只能置顶一篇文章
		var userTopArticleList []models.UserTopArticleModel
		global.DB.Find(&userTopArticleList, "user_id = ?", claims.UserID)

		// 如果查不到，就是没置顶过，置顶
		if len(userTopArticleList) == 0 {
			err := global.DB.Create(&models.UserTopArticleModel{
				UserID:    claims.UserID,
				ArticleID: cr.ArticleID,
			}).Error
			if err != nil {
				res.FailWithMsg("置顶失败", c)
				return
			}
			res.OkWithMsg("置顶成功", c)
			return
		}

		// 查到一个，如果是这篇文章就取消置顶；如果不是，提示超过最大置顶数量
		if len(userTopArticleList) == 1 {
			uta := userTopArticleList[0]
			if uta.ArticleID != cr.ArticleID {
				res.FailWithMsg("用户只能置顶一篇文章", c)
				return
			}
			if err := global.DB.Delete(&uta).Error; err != nil {
				res.FailWithMsg("取消置顶失败", c)
				return
			}
			res.OkWithMsg("取消置顶成功", c)
		}
	case 2:
		// 管理员置顶文章
		if claims.Role != enum.AdminRole {
			res.FailWithMsg("权限错误", c)
			return
		}
		if model.Status != enum.ArticleStatusPublished {
			res.FailWithMsg("管理员只能置顶已经发布的文章", c)
			return
		}
		var userTopArticle models.UserTopArticleModel
		if err := global.DB.Take(&userTopArticle, "user_id = ? and article_id = ?", claims.UserID, cr.ArticleID).Error; err != nil {
			if err := global.DB.Create(&models.UserTopArticleModel{
				UserID:    claims.UserID,
				ArticleID: cr.ArticleID,
			}).Error; err != nil {
				res.FailWithMsg("置顶失败", c)
				return
			}
			res.OkWithMsg("置顶成功", c)
			return
		}
		if err := global.DB.Delete(&userTopArticle).Error; err != nil {
			res.FailWithMsg("取消置顶失败", c)
			return
		}
		res.OkWithMsg("取消置顶成功", c)
	}
}
