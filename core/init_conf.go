package core

import (
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
}

var confPath = "settings.yaml"

func ReadConf() {
	byteData, err := os.ReadFile(confPath)
	if err != nil {
		panic(err)
	}

	var config Config
	err = yaml.Unmarshal(byteData, &config)
	if err != nil {
		panic(fmt.Sprintf("yaml配置文件解析错误 %s", err))
	}
	fmt.Println(config.System.IP, config.System.Port)
}
