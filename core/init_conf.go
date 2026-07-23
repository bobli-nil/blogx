package core

import (
	"blogx_server/conf"
	"blogx_server/flags"
	"fmt"
	"os"

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
