package image_api

import (
	"blogx_server/common/res"
	"blogx_server/global"
	"blogx_server/models"
	"blogx_server/service/qiniu_service"
	"blogx_server/utils/file"
	"blogx_server/utils/hash"
	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
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
	suffix, err := file.ImageSuffixJudge(fileHeader.Filename)
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	// 文件哈希
	file, err := fileHeader.Open()
	if err != nil {
		res.FailWithError(err, c)
		return
	}
	byteData, _ := io.ReadAll(file)
	md5String := hash.Md5(byteData)
	fmt.Println("md5String", md5String)

	// 判断哈希是否已经在库中
	var model models.ImageModel
	err = global.DB.Model(&model).Where("hash = ?", md5String).First(&model).Error
	if err == nil {
		logrus.Infof("图片重复 %s %s", fileHeader.Filename, md5String)
		res.Ok(model.Path, "上传成功", c)
		return
	}

	// 入库
	filePath := fmt.Sprintf("%s/%s/%s.%s", global.Conf.QiNiu.Uri, global.Conf.QiNiu.Prefix, md5String, suffix)
	model = models.ImageModel{
		FileName: fileHeader.Filename,
		Path:     filePath,
		Size:     fileHeader.Size,
		Hash:     md5String,
	}
	err = global.DB.Create(&model).Error
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	// 存储图片
	//err = c.SaveUploadedFile(fileHeader, filePath)
	//if err != nil {
	//	res.FailWithError(err, c)
	//	return
	//}

	// 上传到七牛云
	url, err := qiniu_service.SendFileByteData(byteData, fileHeader.Filename)
	if err != nil {
		res.FailWithError(err, c)
		return
	}

	res.Ok(url, "上传成功", c)
}
