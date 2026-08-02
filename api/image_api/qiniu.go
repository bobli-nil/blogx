package image_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/service/qiniu_service"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type QiNiuGenTokenRequest struct {
	Suffix string
}
type QiNiuGenTokenResponse struct {
	Token  string `json:"token"`
	Region string `json:"region"`
	Key    string `json:"key"`
	Url    string `json:"url"`
	Size   int    `json:"size"`
}

func (ImageApi) QiNiuGenToken(c *gin.Context) {
	var cr QiNiuGenTokenRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	token, err := qiniu_service.GenToken()
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	qiNiu := global.Conf.QiNiu
	id, err := uuid.NewUUID()
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	key := fmt.Sprintf("%s/%s.%s", qiNiu.Prefix, id.String(), cr.Suffix)
	url := fmt.Sprintf("%s/%s", qiNiu.Uri, key)
	res.OkWithData(QiNiuGenTokenResponse{
		Token:  token,
		Region: qiNiu.Region,
		Key:    key,
		Url:    url,
		Size:   qiNiu.Size,
	}, c)
	return
}
