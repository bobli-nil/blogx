package image_api

import (
	"blogx_server/common/res"
	"blogx_server/middleware"
	"blogx_server/service/qiniu_service"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type TransferDepositRequest struct {
	Url string `json:"url" binding:"required"`
}

func (ImageApi) TransferDepositView(c *gin.Context) {
	cr := middleware.GetBind[TransferDepositRequest](c)
	fmt.Println("cr", cr)

	response, err := http.Get(cr.Url)
	if err != nil {
		logrus.Errorf("Transfer deposit view fail, %v", err)
		res.FailWithMsg("图片请求错误", c)
		return
	}
	byteData, _ := io.ReadAll(response.Body)
	qiNiuUrl, err := qiniu_service.SendFileWithoutSuffix(byteData)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	res.OkWithData(qiNiuUrl, c)
}
