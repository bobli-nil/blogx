package core

import (
	"blogx_server/global"

	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

func EsConnect() *elastic.Client {
	esconfig := global.Conf.ES
	client, err := elastic.NewClient(
		elastic.SetURL(esconfig.Url),
		elastic.SetSniff(false),
		elastic.SetBasicAuth(esconfig.Username, esconfig.Password))
	if err != nil {
		logrus.Fatalf("连接ES失败: %v", err)
		return nil
	}
	logrus.Info("ES连接成功")
	return client
}
