package main

import (
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
)

func main() {
	// 解析命令行参数
	flags.Parse()
	// 读取配置文件
	global.Conf = core.ReadConf()
	// 日志配置初始化
	core.InitLogrus()
	// 数据库连接初始化
	global.DB = core.InitDB()
	// IP地址映射文件初始化
	core.InitIPAddr()
	// 数据表迁移
	flags.Run()
}
