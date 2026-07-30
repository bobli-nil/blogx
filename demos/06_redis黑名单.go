package main

import (
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
	"blogx_server/models/enum"
	"blogx_server/service/redis_service/redis_jwt"
	"blogx_server/utils/jwt"
	"fmt"
)

func main() {
	// 解析命令行参数
	flags.Parse()
	// 读取配置文件
	global.Conf = core.ReadConf()
	// redis初始化
	global.Redis = core.InitRedis()

	core.InitLogrus()

	token, err := jwt.GenerateToken(123, "zhangsan", enum.AdminRole)
	if err != nil {
		fmt.Println("token生成失败")
		return
	}
	redis_jwt.TokenBlack(token, redis_jwt.UserBlackType)
	if blk, ok := redis_jwt.HasTokenBlack(token); ok {
		fmt.Println(blk)
	}

}
