package image_api

import (
	"blogx_server/common"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/models/enum"
	"blogx_server/service/log_service"
	"fmt"

	"github.com/gin-gonic/gin"
)

type ImageApi struct{}

type ImageListResponse struct {
	models.ImageModel
	WebPath string `json:"webPath"`
}

func (ImageApi) ImageListView(c *gin.Context) {
	var cr common.PageInfo
	err := c.ShouldBindQuery(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	list, count, err := common.ListQuery(models.ImageModel{}, common.Options{
		PageInfo: cr,
		Likes:    []string{"file_name"},
		Debug:    true,
	})
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	var ResponseList []ImageListResponse
	for _, v := range list {
		ResponseList = append(ResponseList, ImageListResponse{
			ImageModel: v,
			WebPath:    v.WebPath(),
		})
	}
	res.OkWithList(ResponseList, count, c)
}

func (ImageApi) ImageRemoveView(c *gin.Context) {
	log := log_service.GetLogFromGinContext(c)
	log.SetTitle("图片删除")
	log.SetLevel(enum.LogInfoLevel)
	log.SetRequest()

	// 绑定参数
	var cr models.DeleteRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	// 删除
	var images []models.ImageModel
	var successCount, errorCount int64
	global.DB.Where("id in ?", cr.IDList).Find(&images)
	if len(images) > 0 {
		successCount = global.DB.Delete(&images).RowsAffected
	}
	errorCount = int64(len(images)) - successCount
	msg := fmt.Sprintf("操作成功，成功删除%d条，失败删除%d条", successCount, errorCount)
	res.OkWithMsg(msg, c)
}
