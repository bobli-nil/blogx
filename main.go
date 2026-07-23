package main

import (
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"

	"github.com/sirupsen/logrus"
)

func main() {
	flags.Parse()
	global.Conf = core.ReadConf()
	core.InitLogrus()

	logrus.Debugln("这是一个测试")
	logrus.Infoln("这是一个Info")
	logrus.Warnln("这是一个警告")
	logrus.Errorln("这是一个错误")
}
