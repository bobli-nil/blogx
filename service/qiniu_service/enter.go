package qiniu_service

import (
	"blogx_server/global"
	"blogx_server/utils/file"
	"blogx_server/utils/hash"
	"bytes"
	"context"
	"fmt"

	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/http_client"
	"github.com/qiniu/go-sdk/v7/storagev2/uploader"
)

// SendFile 根据路径上传
func SendFile(path string) (string, error) {
	qiNiu := global.Conf.QiNiu
	hash, err := hash.FileMD5(path)
	if err != nil {
		return "", err
	}
	suffix, err := file.ImageSuffixJudge(path)
	if err != nil {
		return "", err
	}
	fileName := fmt.Sprintf("%s.%s", hash, suffix)
	key := fmt.Sprintf("%s/%s", qiNiu.Prefix, fileName)

	mac := credentials.NewCredentials(qiNiu.AccessKey, qiNiu.SecretKey)
	uploadManager := uploader.NewUploadManager(&uploader.UploadManagerOptions{
		Options: http_client.Options{
			Credentials: mac,
		},
	})
	err = uploadManager.UploadFile(context.Background(), path, &uploader.ObjectOptions{
		BucketName: qiNiu.Bucket,
		ObjectName: &key,
		FileName:   fileName,
	}, nil)

	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", qiNiu.Uri, key), nil
}

// SendFileByteData 根据[]byte和fileName上传
func SendFileByteData(byteData []byte, fileName string) (string, error) {
	qiNiu := global.Conf.QiNiu
	hashString := hash.Md5(byteData)
	suffix, err := file.ImageSuffixJudge(fileName)
	if err != nil {
		return "", err
	}
	_fileName := fmt.Sprintf("%s.%s", hashString, suffix)
	key := fmt.Sprintf("%s/%s", qiNiu.Prefix, _fileName)

	mac := credentials.NewCredentials(qiNiu.AccessKey, qiNiu.SecretKey)
	uploadManager := uploader.NewUploadManager(&uploader.UploadManagerOptions{
		Options: http_client.Options{
			Credentials: mac,
		},
	})
	err = uploadManager.UploadReader(context.Background(), bytes.NewReader(byteData), &uploader.ObjectOptions{
		BucketName: qiNiu.Bucket,
		ObjectName: &key,
		FileName:   _fileName,
	}, nil)

	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", qiNiu.Uri, key), nil
}

// SendFileWithoutSuffix 直接使用hash作为文件名，不加后缀
func SendFileWithoutSuffix(byteData []byte) (string, error) {
	qiNiu := global.Conf.QiNiu
	hashString := hash.Md5(byteData)
	_fileName := fmt.Sprintf("%s", hashString)
	key := fmt.Sprintf("%s/%s", qiNiu.Prefix, _fileName)

	mac := credentials.NewCredentials(qiNiu.AccessKey, qiNiu.SecretKey)
	uploadManager := uploader.NewUploadManager(&uploader.UploadManagerOptions{
		Options: http_client.Options{
			Credentials: mac,
		},
	})
	err := uploadManager.UploadReader(context.Background(), bytes.NewReader(byteData), &uploader.ObjectOptions{
		BucketName: qiNiu.Bucket,
		ObjectName: &key,
		FileName:   _fileName,
	}, nil)

	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s/%s", qiNiu.Uri, key), nil
}
