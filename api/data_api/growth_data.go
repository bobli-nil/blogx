package data_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/middleware"
	"blogx_server/models"
	"blogx_server/models/enum"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

type GrowthDataRequest struct {
	Type int8 `form:"type" binding:"required,oneof=1 2 3"`
}

type GrowthDataResponse struct {
	GrowthRate int      `json:"growthRate"` // 今天的增长率
	GrowthNum  int      `json:"growthNum"`  // 今天的增长数
	CountList  []int    `json:"countList"`  // 过去七天的增长数
	DateList   []string `json:"dateList"`   // 过去七天的日期
}

func (DataApi) GrowthData(c *gin.Context) {
	cr := middleware.GetBind[GrowthDataRequest](c)

	type Table struct {
		Date  string `gorm:"column:date"`
		Count int    `gorm:"column:count"`
	}
	var dataList []Table
	now := time.Now()
	before7 := now.AddDate(0, 0, -6)

	switch cr.Type {
	case 1:
		global.DB.Debug().
			Model(&models.SiteFlowModel{}).
			Where(
				"created_at >= ? and created_at <= ?",
				before7.Format("2006-01-02")+" 00:00:00",
				now.Format("2006-01-02 15:04:05"),
			).
			Select("DATE_FORMAT(created_at, '%Y-%m-%d') as date", "sum(count) as count").
			Group("date").
			Scan(&dataList)
	case 2:
		global.DB.
			Debug().
			Model(&models.ArticleModel{}).
			Where(
				"created_at >= ? and created_at <= ? and status = ?",
				before7.Format("2006-01-02")+" 00:00:00",
				now.Format("2006-01-02 15:04:05"),
				enum.ArticleStatusPublished,
			).
			Select("DATE_FORMAT(created_at, '%Y-%m-%d') as date", "count(id) as count").
			Group("date").
			Scan(&dataList)
	case 3:
		global.DB.
			Debug().
			Model(&models.UserModel{}).
			Where(
				"created_at >= ? and created_at <= ?",
				before7.Format("2006-01-02")+" 00:00:00",
				now.Format("2006-01-02 15:04:05"),
			).
			Select("DATE_FORMAT(created_at, '%Y-%m-%d') as date", "count(id) as count").
			Group("date").
			Scan(&dataList)
	}

	dateMap := map[string]int{}
	for _, table := range dataList {
		fmt.Println("table row --> ", table.Date, table.Count)
		dateMap[table.Date] = table.Count
	}

	response := GrowthDataResponse{}
	for i := 0; i < 7; i++ {
		dateStr := before7.AddDate(0, 0, i).Format("2006-01-02")
		count, _ := dateMap[dateStr]
		response.CountList = append(response.CountList, count)
		response.DateList = append(response.DateList, dateStr)
	}
	// 算增长
	response.GrowthNum = response.CountList[6] - response.CountList[5]
	if response.CountList[5] == 0 {
		response.GrowthRate = response.GrowthNum * 100
	} else {
		response.GrowthRate = int(float64(response.GrowthNum) / float64(response.CountList[5]) * 100)
	}

	res.OkWithData(response, c)
}
