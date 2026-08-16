package enum

type ArticleStatus uint8

const (
	ArticleStatusDraft     ArticleStatus = 1 // 草稿
	ArticleStatusExamined  ArticleStatus = 2 // 审核中
	ArticleStatusPublished ArticleStatus = 3 // 已发布
)
