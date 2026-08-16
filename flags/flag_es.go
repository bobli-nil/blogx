package flags

import (
	"blogx_server/models"
	"blogx_server/service/es_service"
)

func ESIndex() {
	articleModel := models.ArticleModel{}
	es_service.CreateIndexV2(articleModel.Index(), articleModel.Mapping())
}
