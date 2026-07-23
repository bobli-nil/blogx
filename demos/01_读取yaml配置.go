package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

type SystemConfig struct {
	IP   string `yaml:"ip"`
	Port int    `yaml:"port"`
}

type Config struct {
	System SystemConfig `yaml:"system"`
}

// os.readFile 读取配置文件
// 使用 gopkg.in/yaml.v2 将yaml文件转化为结构体
func main() {
	byteData, err := os.ReadFile("settings.yaml")
	if err != nil {
		panic(err)
	}

	var config Config
	err = yaml.Unmarshal(byteData, &config)
	if err != nil {
		panic(err)
	}

	fmt.Println(config.System.IP, config.System.Port)
}
