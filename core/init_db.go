package core

import (
	"blogx_server/global"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

func InitDB() *gorm.DB {
	DB := global.Conf.DB

	if len(DB) == 0 {
		logrus.Fatalf("未配置数据库")
	}

	dc := DB[0] // 读写库

	db, err := gorm.Open(mysql.Open(dc.DSN()), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
		Logger:                                   logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		logrus.Fatalf("数据库连接失败 %s", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logrus.Fatalf("数据库连接池创建失败 %s", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)
	logrus.Infof("数据库连接成功")

	if len(DB) > 1 {
		var readList []gorm.Dialector
		for _, v := range DB[1:] {
			readList = append(readList, mysql.Open(v.DSN()))
		}
		fmt.Printf("%+v\n", readList)
		err = db.Use(dbresolver.Register(dbresolver.Config{
			Sources:           []gorm.Dialector{mysql.Open(dc.DSN())},
			Replicas:          readList,
			Policy:            dbresolver.RandomPolicy{},
			TraceResolverMode: true,
		}))
		if err != nil {
			logrus.Fatalf("读写配置错误 %s", err)
		}
	}

	return db
}
