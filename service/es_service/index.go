package es_service

import (
	"blogx_server/global"
	"context"

	"github.com/sirupsen/logrus"
)

func CreateIndexV2(index, mapping string) {
	if ExistIndex(index) {
		DeleteIndex(index)
	}
	CreateIndex(index, mapping)
}

// CreateIndex 创建索引
func CreateIndex(index, mapping string) {
	_, err := global.ESClient.CreateIndex(index).BodyString(mapping).Do(context.Background())
	if err != nil {
		logrus.Errorf("%s 索引创建失败 %s", index, err)
		return
	}
	logrus.Infof("%s 索引创建成功", index)
}

// ExistIndex 判断索引是否存在
func ExistIndex(index string) bool {
	exists, _ := global.ESClient.IndexExists(index).Do(context.Background())
	return exists
}

func DeleteIndex(index string) {
	_, err := global.ESClient.DeleteIndex(index).Do(context.Background())
	if err != nil {
		logrus.Errorf("%s 删除索引失败 %s", index, err)
		return
	}
	logrus.Infof("%s 索引删除成功", index)
}
