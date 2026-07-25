package core

import (
	"blogx_server/utils/ip"
	"fmt"
	"strings"

	"github.com/lionsoul2014/ip2region/binding/golang/service"
	"github.com/sirupsen/logrus"
)

var Searcher *service.Ip2Region

func InitIPAddr() {
	v4Config, err := service.NewV4Config(service.VIndexCache, "init/ip2region_v4.xdb", 20)
	if err != nil {
		logrus.Errorf("初始化 ip2region_v4.xdb 失败: %v", err)
		return
	}

	_search, err := service.NewIp2Region(v4Config, nil)
	if err != nil {
		logrus.Errorf("初始化 ip2region_v4.xdb 失败: %v", err)
		return
	}

	Searcher = _search
}

func GetIPAddr(IP string) string {
	region, err := Searcher.Search(IP)
	if err != nil {
		logrus.Errorf("错误的IP地址 %s", err)
		return "错误的IP地址"
	}

	if ip.HasLocalIpAddr(IP) {
		return "内网地址"
	}

	IPInfo := strings.Split(region, "|")
	country := IPInfo[0]
	province := IPInfo[1]
	city := IPInfo[2]

	if province != "0" && city != "0" {
		return fmt.Sprintf("%s · %s", province, city)
	}
	if country != "0" && province != "0" {
		return fmt.Sprint("%s · %s", country, province)
	}
	if country != "0" {
		return country
	}

	return region
}
