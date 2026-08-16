package main

import (
	"blogx_server/core"
	"blogx_server/flags"
	"blogx_server/global"
	"blogx_server/models"
	"context"
	"fmt"

	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

func create() {
	var article = models.ArticleModel{
		Model: models.Model{
			ID: 1,
		},
		Title:   "标题",
		Content: "内容",
		UserID:  1,
		Status:  1,
	}

	indexResponse, err := global.ESClient.Index().Index(article.Index()).BodyJson(article).Do(context.Background())
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%#v\n", indexResponse)
}

func list() {
	limit := 2
	page := 1
	from := (page - 1) * limit

	query := elastic.NewBoolQuery()
	res, err := global.ESClient.Search(models.ArticleModel{}.Index()).Query(query).From(from).Size(limit).Do(context.Background())
	if err != nil {
		logrus.Errorf("ES查询失败 %s", err)
		return
	}
	count := res.Hits.TotalHits.Value
	fmt.Printf("结果总数 %d \n", count)

	for _, hit := range res.Hits.Hits {
		fmt.Println(string(hit.Source))
	}
}

func DocDelete() {
	deleteResponse, err := global.ESClient.Delete().Index(models.ArticleModel{}.Index()).Id("8e0t_58BDsyOIJ87nUSj").Refresh("true").Do(context.Background())
	if err != nil {
		logrus.Errorf("删除失败 %s", err)
		return
	}
	fmt.Println(deleteResponse)
}

func Update() {
	updateResponse, err := global.ESClient.Update().
		Index(models.ArticleModel{}.Index()).
		Id("8u1m_58BDsyOIJ87zkT7").
		Refresh("true").
		Doc(map[string]any{
			"content": "新的内容",
		}).
		Do(context.Background())
	if err != nil {
		logrus.Errorf("ES 更新失败 %s", err)
		return
	}
	fmt.Printf("%#v\n", updateResponse)
}

func main() {
	flags.Parse()
	global.Conf = core.ReadConf()
	core.InitLogrus()
	global.ESClient = core.EsConnect()

	//create()
	//list()
	DocDelete()
	//Update()
}
