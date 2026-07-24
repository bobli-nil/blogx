package main

import (
	"fmt"
	"strings"

	"github.com/lionsoul2014/ip2region/binding/golang/service"
)

func main() {
	v4Config, err := service.NewV4Config(service.VIndexCache, "init/ip2region_v4.xdb", 20)
	if err != nil {
		fmt.Printf("创建v4Config失败 %s \n", err)
		return
	}

	IP2Region, err := service.NewIp2Region(v4Config, nil)
	if err != nil {
		fmt.Errorf("创建IP2Region失败 %s \n", err)
		return
	}
	defer IP2Region.Close()

	search, err := IP2Region.Search("113.92.157.29")
	if err != nil {
		fmt.Errorf("解析IP失败 %s \n", err)
		return
	}
	split := strings.Split(search, "|")
	fmt.Println(split)

}
