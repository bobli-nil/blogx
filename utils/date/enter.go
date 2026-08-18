package date

import "time"

func GetNowAfter() time.Time {
	loc, _ := time.LoadLocation("Asia/Shanghai")
	now := time.Now().In(loc)
	endTime := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, loc)
	return endTime
}
