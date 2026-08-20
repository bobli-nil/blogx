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

	r1 := GetRootComment(11)
	fmt.Println(r1.ID)
	r2 := GetRootComment(10)
	fmt.Println(r2.ID)
	r3 := GetRootComment(2)
	fmt.Println(r3.ID)
}

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
