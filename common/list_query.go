package common

import (
	"blogx_server/global"
	"fmt"

	"gorm.io/gorm"
)

type PageInfo struct {
	Page    int    `form:"page"`
	Limit   int    `form:"limit"`
	Keyword string `form:"keyword"`
	Order   string `form:"order"`
}

func (p PageInfo) GetPage() int {
	if p.Page <= 0 {
		return 1
	}
	return p.Page
}

func (p PageInfo) GetLimit() int {
	if p.Limit <= 0 {
		return 10
	}
	return p.Limit
}

func (p PageInfo) GetOffset() int {
	page := p.GetPage()
	limit := p.GetLimit()
	return (page - 1) * limit
}

type Options struct {
	PageInfo PageInfo
	Likes    []string
	PreLoads []string
	Where    *gorm.DB
	Debug    bool
}

func ListQuery[T any](model T, option Options) (list []T, count int, err error) {
	// 基础查询
	query := global.DB.Model(&model).Where(&model)

	// 模糊匹配
	if len(option.Likes) > 0 && option.PageInfo.Keyword != "" {
		likes := global.DB.Where("")
		for _, column := range option.Likes {
			likes.Or(fmt.Sprintf("%s LIKE ?", column), fmt.Sprintf("%%%s%%", option.PageInfo.Keyword))
		}
		query = query.Where(likes)
	}

	// 定制化
	if option.Where != nil {
		query = query.Where(option.Where)
	}

	// Debug
	if option.Debug {
		query = query.Debug()
	}

	// Preload
	if option.PreLoads != nil {
		for _, preload := range option.PreLoads {
			query = query.Preload(preload)
		}
	}

	// 排序
	if option.PageInfo.Order != "" {
		query = query.Order(option.PageInfo.Order)
	}

	// 查总数
	var _c int64
	query.Count(&_c)
	count = int(_c)

	// 分页
	offset := option.PageInfo.GetOffset()
	limit := option.PageInfo.GetLimit()
	err = query.Offset(offset).Limit(limit).Find(&list).Error
	if err != nil {
		return list, count, err
	}
	return
}
