package main

import (
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
	"blogx_server/service/qiniu_service"
	"fmt"
)

func main() {
	flags.Parse()
	global.Conf = core.ReadConf()

	//url, err := qiNiu.SendFile("uploads/images001/f6d69eefe9bef775c88fdb1deee73582.jpg")
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//fmt.Println("url", url)

	token, err := qiniu_service.GenToken()
	fmt.Println(token, err)
}
