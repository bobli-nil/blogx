package search_api

import (
	"blogx_server/common"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/service/text_service"
	"context"
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

type TextSearchRequest struct {
	common.PageInfo
}

type TextSearchResponse struct {
	ArticleID uint   `json:"articleID"`
	Head      string `json:"head"`
	Body      string `json:"body"`
}

func (SearchApi) TextSearchView(c *gin.Context) {
	cr := middleware.GetBind[TextSearchRequest](c)

	query := elastic.NewBoolQuery()
	if cr.Keyword != "" {
		query.Should(
			elastic.NewMatchQuery("head", cr.Keyword),
			elastic.NewMatchQuery("body", cr.Keyword),
		)
	}

	highlight := elastic.NewHighlight()
	highlight.Field("head")
	highlight.Field("body")

	result, err := global.ESClient.
		Search(models.TextModel{}.Index()).
		Query(query).
		From(cr.GetOffset()).
		Size(cr.GetLimit()).
		Highlight(highlight).
		Do(context.Background())
	if err != nil {
		fmt.Println(err)
		res.FailWithMsg("ES查询失败", c)
	}

	source, _ := query.Source()
	byteData, _ := json.Marshal(source)
	fmt.Println("es查询语句", string(byteData))

	count := result.Hits.TotalHits.Value
	fmt.Println("count-->", count)

	var list = make([]TextSearchResponse, 0)
	for _, hit := range result.Hits.Hits {
		var item text_service.TextModel
		err := json.Unmarshal(hit.Source, &item)
		if err != nil {
			logrus.Warnf("解析失败 %s  %s", err, string(hit.Source))
			continue
		}
		if len(hit.Highlight["head"]) > 0 {
			item.Head = hit.Highlight["head"][0]
		}
		if len(hit.Highlight["body"]) > 0 {
			item.Body = hit.Highlight["body"][0]
		}
		list = append(list, TextSearchResponse{
			ArticleID: item.ArticleID,
			Head:      item.Head,
			Body:      item.Body,
		})
	}

	res.OkWithList(list, int(count), c)
}
