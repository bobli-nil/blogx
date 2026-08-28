package main

import (
	"blogx_server/common"
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
	"blogx_server/models"
	"fmt"
)

func main() {
	flags.Parse()
	global.Conf = core.ReadConf()
	core.InitLogrus()
	global.DB = core.InitDB()

	mps := common.ScanMapV2(models.UserModel{}, common.ScanOption{
		Where: global.DB.Where("id in ?", []uint{6}),
	})
	fmt.Printf("%+v\n", mps)
}
