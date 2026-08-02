package banner_api

import (
	"blogx_server/common"
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"errors"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BannerApi struct{}

type BannerCreateRequest struct {
	Cover string `json:"cover"`
	Href  string `json:"href"`
	Show  *bool  `json:"show"`
}

func (BannerApi) BannerCreateView(c *gin.Context) {
	var req BannerCreateRequest
	err := c.ShouldBindJSON(&req)
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	err = global.DB.Create(&models.BannerModel{
		Cover: req.Cover,
		Href:  req.Href,
		Show:  req.Show,
	}).Error
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	res.OkWithMsg("创建Banner成功", c)
}

type BannerListRequest struct {
	common.PageInfo
	Show *bool `form:"show"`
}

func (BannerApi) BannerListView(c *gin.Context) {
	var cr BannerListRequest
	err := c.ShouldBindQuery(&cr)
	fmt.Printf("---->%+v\n", cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	list, count, err := common.ListQuery(models.BannerModel{
		Show: cr.Show,
	}, common.Options{
		PageInfo: cr.PageInfo,
		Debug:    true,
	})
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	res.OkWithList(list, count, c)
}

func (BannerApi) BannerRemoveView(c *gin.Context) {
	var cr models.DeleteRequest
	err := c.ShouldBindJSON(&cr)
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	result := global.DB.Delete(&models.BannerModel{}, cr.IDList)
	if result.Error != nil {
		res.FailWithError(result.Error, c)
		return
	}
	successCount := result.RowsAffected
	failCount := int64(len(cr.IDList)) - result.RowsAffected
	if successCount == 0 {
		res.FailWithMsg("不存在该记录", c)
		return
	}
	res.OkWithMsg(fmt.Sprintf("删除成功，成功删除%d记录，失败删除%d记录", successCount, failCount), c)
}

func (BannerApi) BannerUpdateView(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		res.FailWithError(errors.New("非法的ID"), c)
		return
	}
	var req BannerCreateRequest
	err = c.ShouldBindJSON(&req)
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	db := global.DB.Model(&models.BannerModel{}).Where("id = ?", id).Updates(models.BannerModel{
		Cover: req.Cover,
		Href:  req.Href,
		Show:  req.Show,
	})

	if db.Error != nil {
		res.FailWithError(db.Error, c)
		return
	}
	if db.RowsAffected == 0 {
		res.FailWithMsg("不存在该记录", c)
		return
	}
	res.OkWithMsg("更新成功", c)
}
