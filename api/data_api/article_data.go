package data_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum"
	"time"

	"github.com/gin-gonic/gin"
)

type ArticleDataRequest struct {
	GrowthRate int      `json:"growthRate"` // 今天的增长率
	GrowthNum  int      `json:"growthNum"`  // 今天的增长数
	CountList  []int    `json:"countList"`  // 过去七天的增长数
	DateList   []string `json:"dateList"`   // 过去七天的日期
}

func (DataApi) ArticleDataView(c *gin.Context) {
	now := time.Now()
	before7 := now.AddDate(0, 0, -7)

	var articleList []models.ArticleModel
	global.DB.Find(
		&articleList,
		"created_at >= ? and created_at <= ? and Status = ?",
		before7.Format("2006-01-02")+" 00:00:00",
		now.Format("2006-01-02 15:04:05"),
		enum.ArticleStatusPublished,
	)

	var dateMap = map[string]int{}
	for _, model := range articleList {
		key := model.CreatedAt.Format("2006-01-02")
		count, ok := dateMap[key]
		if !ok {
			dateMap[key] = 1
		} else {
			dateMap[key] = count + 1
		}
	}

	response := ArticleDataRequest{}
	for i := 0; i < 7; i++ {
		dateStr := before7.AddDate(0, 0, i).Format("2006-01-02")
		count, _ := dateMap[dateStr]
		response.CountList = append(response.CountList, count)
		response.DateList = append(response.DateList, dateStr)
	}

	// 算增长
	response.GrowthNum = response.CountList[6] - response.CountList[5]
	response.GrowthRate = int(float64(response.GrowthNum) / float64(response.CountList[5]) * 100)

	res.OkWithData(response, c)
}
