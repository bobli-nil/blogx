package search_api

import (
	"blogx_server/common"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"context"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

type AggType struct {
	DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
	SumOtherDocCount        int `json:"sum_other_doc_count"`
	Buckets                 []struct {
		Key      string `json:"key"`
		DocCount int    `json:"doc_count"`
	} `json:"buckets"`
}

type AggTypeTotal struct {
	Value int `json:"value"`
}

type TagAggResponse struct {
	Tag          string `json:"tag"`
	ArticleCount int    `json:"articleCount"`
}

func (SearchApi) TagAggView(c *gin.Context) {
	cr := middleware.GetBind[common.PageInfo](c)

	agg := elastic.NewTermsAggregation().Field("tag_list")
	agg.SubAggregation("page", elastic.NewBucketSortAggregation().
		From(cr.GetOffset()).
		Size(cr.GetLimit()))
	query := elastic.NewBoolQuery()
	query.MustNot(elastic.NewTermQuery("tag_list", ""))
	result, err := global.ESClient.
		Search(models.ArticleModel{}.Index()).
		Query(query).
		Aggregation("tags", agg).
		Aggregation("total", elastic.NewCardinalityAggregation().Field("tag_list")).
		Size(0).
		Do(context.Background())
	if err != nil {
		logrus.Errorf("ES查询失败 %s", err)
		res.FailWithMsg("ES查询失败", c)
		return
	}
	val := result.Aggregations["tags"]
	totalVal := result.Aggregations["total"]
	var t AggType
	var t1 AggTypeTotal
	err = json.Unmarshal(val, &t)
	err = json.Unmarshal(totalVal, &t1)
	if err != nil {
		res.FailWithMsg("解析JSON失败", c)
		logrus.Errorf("解析失败 %s", err)
		return
	}

	list := make([]TagAggResponse, 0)
	for _, bucket := range t.Buckets {
		list = append(list, TagAggResponse{
			Tag:          bucket.Key,
			ArticleCount: bucket.DocCount,
		})
	}
	res.OkWithList(list, t1.Value, c)
}
