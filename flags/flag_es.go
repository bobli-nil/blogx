package flags

import (
	"blogx_server/models"
	"blogx_server/service/es_service"
)

func ESIndex() {
	articleModel := models.ArticleModel{}
	es_service.CreateIndexV2(articleModel.Index(), articleModel.Mapping())
	textModel := models.TextModel{}
	es_service.CreateIndexV2(textModel.Index(), textModel.Mapping())
}
