package data_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"time"

	"github.com/gin-gonic/gin"
)

type ArticleYearDataResponse struct {
	GrowthRate int      `json:"growthRate"` // 今天的增长率
	GrowthNum  int      `json:"growthNum"`  // 今天的增长数
	CountList  []int    `json:"countList"`  // 过去七天的增长数
	DateList   []string `json:"dateList"`   // 过去七天的日期
}

func (DataApi) ArticleYearDataView(c *gin.Context) {
	type Table struct {
		Date  string `gorm:"column:date"`
		Count int    `gorm:"column:count"`
	}
	var dataList []Table
	now := time.Now()
	before12Month := now.AddDate(0, -11, 0)

	global.DB.
		Model(&models.ArticleModel{}).
		Where(
			"created_at >= ? and created_at <= ?",
			before12Month.Format("2006-01-02")+" 00:00:00",
			now.Format("2006-01-02 15:04:05"),
		).
		Select("DATE_FORMAT(created_at, '%Y-%m') as date", "count(id) as count").
		Group("date").
		Scan(&dataList)

	var dateMap = map[string]int{}
	for _, model := range dataList {
		dateMap[model.Date] = model.Count
	}

	response := ArticleYearDataResponse{}
	for i := 0; i < 12; i++ {
		dateStr := before12Month.AddDate(0, i, 0).Format("2006-01")
		count, _ := dateMap[dateStr]
		response.CountList = append(response.CountList, count)
		response.DateList = append(response.DateList, dateStr)
	}

	// 算增长
	response.GrowthNum = response.CountList[11] - response.CountList[10]
	if response.CountList[10] == 0 {
		response.GrowthRate = response.GrowthNum * 100
	} else {
		response.GrowthRate = int(float64(response.GrowthNum) / float64(response.CountList[10]) * 100)
	}

	res.OkWithData(response, c)
}
