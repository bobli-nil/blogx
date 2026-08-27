package main

import (
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
	"blogx_server/service/chat_service"
)

func main() {
	flags.Parse()
	// 读取配置文件
	global.Conf = core.ReadConf()
	// 日志配置初始化
	core.InitLogrus()
	// redis初始化
	global.Redis = core.InitRedis()
	// 数据库连接初始化
	global.DB = core.InitDB()

	//r1 := models.CalcUserPatchRelationship2(1, []uint{2})
	//fmt.Println(r1)
	//r2 := models.CalcUserPatchRelationship2(2, []uint{1})
	//fmt.Println(r2)

	chat_service.ToTextChat(1, 2, "在干嘛")
	chat_service.ToTextChat(2, 1, "在吃饭")
}
