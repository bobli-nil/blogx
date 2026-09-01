package ai_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/service/ai_service"
	"context"
	"encoding/json"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

type ArticleAiRequest struct {
	Content string `form:"content" binding:"required"`
}

func (AiApi) ArticleAiView(c *gin.Context) {
	cr := middleware.GetBind[ArticleAiRequest](c)
	if !global.Conf.Ai.Enable {
		res.SSEFail("站点未启用AI功能", c)
		return
	}

	query := elastic.NewBoolQuery()
	query.Should(
		elastic.NewMatchQuery("title", cr.Content),
		elastic.NewMatchQuery("content", cr.Content),
		elastic.NewMatchQuery("abstract", cr.Content))
	query.Must(elastic.NewTermQuery("status", 3))

	result, err := global.ESClient.
		Search(models.ArticleModel{}.Index()).
		Query(query).
		From(0).
		Size(100).
		Do(context.Background())
	if err != nil {
		source, _ := query.Source()
		byteData, _ := json.Marshal(source)
		logrus.Errorf("ES查询失败 %s \n %s \n", err, string(byteData))
		res.SSEFail("查询失败", c)
		return
	}

	var list []string
	for _, hit := range result.Hits.Hits {
		list = append(list, string(hit.Source))
	}
	content := "[" + strings.Join(list, ",") + "]"
	msgChan, err := ai_service.ChatStream(cr.Content, content)
	if err != nil {
		res.SSEFail("文章分析失败", c)
	}
	for s := range msgChan {
		res.SSEOk(s, c)
	}
}
