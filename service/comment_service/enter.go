package comment_service

import (
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/service/redis_service/redis_comment"
	"fmt"
)

func GetRootComment(commentID uint) (model *models.CommentModel) {
	var comment models.CommentModel
	err := global.DB.Take(&comment, commentID).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}
	if comment.ParentID != nil {
		return GetRootComment(*comment.ParentID)
	}
	return &comment
}

func GetCommentTree(model *models.CommentModel) {
	global.DB.Preload("SubCommentList").Take(&model)
	for _, v := range model.SubCommentList {
		GetCommentTree(v)
	}
}

func GetCommentTreeV2(id uint) (model *models.CommentModel) {
	model = &models.CommentModel{
		Model: models.Model{ID: id},
	}
	global.DB.Preload("SubCommentList").Take(model, id)

	for i := 0; i < len(model.SubCommentList); i++ {
		model.SubCommentList[i] = GetCommentTreeV2(model.SubCommentList[i].ID)
	}

	return
}

type CommentResponse struct {
	ID           uint               `json:"id"`
	Content      string             `json:"content"`
	UserID       uint               `json:"userID"`
	UserNickname string             `json:"userNickname"`
	UserAvatar   string             `json:"userAvatar"`
	ArticleID    uint               `json:"articleID"`
	ParentID     *uint              `json:"parentID"`
	DiggCount    int                `json:"diggCount"`   // 点赞数
	ApplyCount   int                `json:"applyCount"`  // 回复数
	SubComments  []*CommentResponse `json:"subComments"` // 子评论
}

func GetCommentTreeV3(id uint) (res *CommentResponse) {
	model := models.CommentModel{
		Model: models.Model{ID: id},
	}
	global.DB.Preload("UserModel").Preload("SubCommentList").Take(&model)

	res = &CommentResponse{
		ID:           model.ID,
		Content:      model.Content,
		UserID:       model.UserID,
		UserNickname: model.UserModel.Nickname,
		UserAvatar:   model.UserModel.Avatar,
		ArticleID:    model.ArticleID,
		ParentID:     model.ParentID,
		DiggCount:    model.DiggCount,
		ApplyCount:   0,
		SubComments:  make([]*CommentResponse, 0),
	}
	for _, commentModel := range model.SubCommentList {
		res.SubComments = append(res.SubComments, GetCommentTreeV3(commentModel.ID))
	}

	return
}

func GetCommentTreeV4(id uint) (res *CommentResponse) {
	return getCommentTreeByLine(id, 1)
}

func getCommentTreeByLine(id uint, line int) (res *CommentResponse) {
	model := models.CommentModel{
		Model: models.Model{ID: id},
	}
	global.DB.Preload("UserModel").Preload("SubCommentList").Take(&model)

	res = &CommentResponse{
		ID:           model.ID,
		Content:      model.Content,
		UserID:       model.UserID,
		UserNickname: model.UserModel.Nickname,
		UserAvatar:   model.UserModel.Avatar,
		ArticleID:    model.ArticleID,
		ParentID:     model.ParentID,
		DiggCount:    model.DiggCount + redis_comment.GetCacheDigg(model.ID),
		ApplyCount:   redis_comment.GetCacheApply(model.ID),
		SubComments:  make([]*CommentResponse, 0),
	}
	if line >= global.Conf.Site.Article.CommentLine {
		return
	}
	for _, commentModel := range model.SubCommentList {
		res.SubComments = append(res.SubComments, getCommentTreeByLine(commentModel.ID, line+1))
	}

	return
}

// GetParents 获取该评论所有的祖先评论
func GetParents(commentID uint) (list []*models.CommentModel) {
	var comment models.CommentModel
	if err := global.DB.Take(&comment, commentID).Error; err != nil {
		return
	}
	list = append(list, &comment)
	if comment.ParentID != nil {
		list = append(list, GetParents(*comment.ParentID)...)
	}
	return
}

func GetCommentOneDimensional(id uint) (list []models.CommentModel) {
	model := models.CommentModel{
		Model: models.Model{ID: id},
	}
	global.DB.Preload("SubCommentList").Take(&model)
	list = append(list, model)
	for _, commentModel := range model.SubCommentList {
		subList := GetCommentOneDimensional(commentModel.ID)
		list = append(list, subList...)
	}
	return
}
