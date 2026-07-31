package core

import (
	"blogx_server/conf"
	"blogx_server/flags"
	"blogx_server/global"
	"fmt"
	"os"

	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

func ReadConf() (c *conf.Config) {
	byteData, err := os.ReadFile(flags.FlagOptions.File)
	if err != nil {
		panic(err)
	}

	c = new(conf.Config)

	err = yaml.Unmarshal(byteData, c)
	if err != nil {
		panic(fmt.Sprintf("yaml配置文件解析错误 %s", err))
	}
	fmt.Printf("读取配置System %v, Log %v 成功\n", c.System, c.Log)

	return
}

func SetConf() {
	byteData, err := yaml.Marshal(global.Conf)
	if err != nil {
		logrus.Error("序列化conf失败 %s", err)
		return
	}

	err = os.WriteFile(flags.FlagOptions.File, byteData, 0666)
	if err != nil {
		logrus.Errorf("设置配置文件失败 %s", err)
		return
	}
}
