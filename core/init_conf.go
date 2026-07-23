package core

import (
	"blogx_server/flags"
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type Config struct {
	System System
}

type System struct {
	IP   string `yaml:"ip"`
	Port int    `yaml:"port"`
	Env  string `yaml:"env"`
}

func ReadConf() {
	byteData, err := os.ReadFile(flags.FlagOptions.File)
	if err != nil {
		panic(err)
	}

	var config Config
	err = yaml.Unmarshal(byteData, &config)
	if err != nil {
		panic(fmt.Sprintf("yaml配置文件解析错误 %s", err))
	}
	fmt.Printf("读取配置 %v 成功\n", config.System)
}
