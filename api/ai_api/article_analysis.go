package ai_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/service/ai_service"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type ArticleAnalysisRequest struct {
	Content string `json:"content" binding:"required"`
}

type ArticleAnalysisResponse struct {
	Title    string   `json:"title"`
	Abstract string   `json:"abstract"`
	Category string   `json:"category"`
	Tag      []string `json:"tag"`
}

func (AiApi) ArticleAnalysisView(c *gin.Context) {
	cr := middleware.GetBind[ArticleAnalysisRequest](c)

	if !global.Conf.Ai.Enable {
		res.FailWithMsg("站点未启用AI功能", c)
		return
	}

	msg, err := ai_service.Chat(cr.Content)
	if err != nil {
		logrus.Errorf("AI分析失败 %s %s", err, cr.Content)
		res.FailWithMsg("AI分析失败", c)
		return
	}
	var data ArticleAnalysisResponse
	err = json.Unmarshal([]byte(msg), &data)
	if err != nil {
		logrus.Errorf("数据解析失败 %s %s", err, msg)
		res.FailWithMsg("数据解析失败", c)
		return
	}

	res.OkWithData(data, c)
}
