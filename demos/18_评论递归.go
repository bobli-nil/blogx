package main

import (
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
	"blogx_server/models"
	"fmt"
)

func main() {
	// 解析命令行参数
	flags.Parse()
	// 读取配置文件
	global.Conf = core.ReadConf()
	// 日志配置初始化
	core.InitLogrus()
	global.DB = core.InitDB()

	//model := models.CommentModel{
	//	Model: models.Model{ID: 2},
	//}
	//GetCommentTree(&model)

	//model := GetCommentTreeV2(2)

	//for _, v1 := range model.SubCommentList {
	//	fmt.Printf("--%d\n", v1.ID)
	//	for _, v2 := range v1.SubCommentList {
	//		fmt.Printf("----%d\n", v2.ID)
	//		for _, v3 := range v2.SubCommentList {
	//			fmt.Printf("------%d\n", v3.ID)
	//		}
	//	}
	//}

	//list := GetCommentOneDimensional(2)
	//for _, model := range list[1:] {
	//	fmt.Println(model.ID)
	//}

	//res := GetCommentTreeV3(2)
	//bs, _ := json.Marshal(res)
	//fmt.Println(string(bs))

	list := GetParents(30)
	for _, v := range list {
		fmt.Println(v.ID)
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

func GetCommentTree(model *models.CommentModel) {
	global.DB.Preload("SubCommentList").Take(&model)
	for _, v := range model.SubCommentList {
		GetCommentTree(v)
	}
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

type CommentResponse struct {
	ID           uint               `json:"id"`
	Content      string             `json:"content"`
	UserID       uint               `json:"userID"`
	UserNickname string             `json:"userNickname"`
	UserAvatar   string             `json:"userAvatar"`
	ArticleID    uint               `json:"articleID"`
	ParentID     *uint              `json:"parentID"`
	DiggCount    int                `json:"diggCount"`   // 点赞数
	ApplyCount   uint               `json:"applyCount"`  // 回复数
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
