package qiniu_service

import (
	"blogx_server/global"
	"context"
	"time"

	"github.com/qiniu/go-sdk/v7/storagev2/credentials"
	"github.com/qiniu/go-sdk/v7/storagev2/uptoken"
)

func GenToken() (token string, err error) {
	qiNiu := global.Conf.QiNiu
	mac := credentials.NewCredentials(qiNiu.AccessKey, qiNiu.SecretKey)
	putPolicy, err := uptoken.NewPutPolicy(qiNiu.Bucket, time.Now().Add(time.Duration(qiNiu.Expiry)*time.Second))
	if err != nil {
		return
	}
	token, err = uptoken.NewSigner(putPolicy, mac).GetUpToken(context.Background())
	if err != nil {
		return
	}
	return
}
