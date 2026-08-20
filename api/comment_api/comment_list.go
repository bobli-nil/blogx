package comment_api

import (
	"blogx_server/common"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/utils/jwt"
	"time"

	"github.com/gin-gonic/gin"
)

type CommentListRequest struct {
	common.PageInfo
	UserID uint  `form:"userID"`
	Type   uint8 `form:"type" binding:"required,oneof=1 2 3"` // 1查我发布的文章的评论 2查我发布的评论 3管理员看所有
}

type CommentListResponse struct {
	ID           uint       `json:"id"`
	CreatedAt    *time.Time `json:"createdAt"`
	Content      string     `json:"content"`
	UserID       uint       `json:"userID"`
	UserNickname string     `json:"userNickname"`
	UserAvatar   string     `json:"userAvatar"`
	ArticleID    uint       `json:"articleID"`
	ArticleTitle string     `json:"articleTitle"`
	ArticleCover string     `json:"articleCover"`
	DiggCount    int        `json:"diggCount"`
}

func (CommentApi CommentApi) CommentListView(c *gin.Context) {
	cr := middleware.GetBind[CommentListRequest](c)
	query := global.DB.Where("")
	claims := jwt.GetClaims(c)

	switch cr.Type {
	case 1:
		var articleIDList []uint
		global.DB.Model(&models.ArticleModel{}).
			Find(&[]models.ArticleModel{}, "user_id = ? and status = ?", claims.UserID, enum.ArticleStatusPublished).
			Select("id").
			Scan(&articleIDList)
		query.Where("article_id in ?", articleIDList)
		cr.UserID = 0
	case 2:
		cr.UserID = claims.UserID
	case 3:
		cr.UserID = 0
	}

	_list, count, _ := common.ListQuery(models.CommentModel{
		UserID: cr.UserID,
	}, common.Options{
		Debug:    true,
		PageInfo: cr.PageInfo,
		Likes:    []string{"content"},
		PreLoads: []string{"ArticleModel", "UserModel"},
		Where:    query,
	})
	var list = make([]CommentListResponse, 0)
	for _, model := range _list {
		list = append(list, CommentListResponse{
			ID:           model.ID,
			CreatedAt:    model.CreatedAt,
			Content:      model.Content,
			UserID:       model.UserID,
			UserNickname: model.UserModel.Nickname,
			UserAvatar:   model.UserModel.Avatar,
			ArticleID:    model.ArticleID,
			ArticleTitle: model.ArticleModel.Title,
			DiggCount:    model.DiggCount,
		})
	}

	res.OkWithList(list, count, c)
}
