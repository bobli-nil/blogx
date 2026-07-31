package image_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/utils"
	"errors"
	"fmt"
	"strings"

	"github.com/gin-gonic/gin"
)

func (ImageApi) UploadImageView(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	// 判断大小
	if fileHeader.Size > global.Conf.Upload.Size*1024*1024 {
		res.FailWithMsg(fmt.Sprintf("文件大小不能超过%dMB", global.Conf.Upload.Size), c)
		return
	}

	// 后缀判断
	err = ImageSuffixJudge(fileHeader.Filename)
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	filePath := fmt.Sprintf("uploads/%s/%s", global.Conf.Upload.UploadDir, fileHeader.Filename)
	err = c.SaveUploadedFile(fileHeader, filePath)
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	res.Ok("/"+filePath, "上传成功", c)
}

func ImageSuffixJudge(filename string) error {
	_list := strings.Split(filename, ".")
	if len(_list) < 2 {
		return errors.New("文件名不合法")
	}
	suffix := _list[len(_list)-1]
	if utils.InList(suffix, global.Conf.Upload.WhiteList) {
		return nil
	}
	return errors.New("文件后缀不正确")
}
