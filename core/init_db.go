package core

import (
	"blogx_server/global"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	DB := global.Conf.DB
	url := DB.DSN()
	fmt.Println("url", url)

	db, err := gorm.Open(mysql.Open(url), &gorm.Config{
		DisableForeignKeyConstraintWhenMigrating: true,
	})
	if err != nil {
		logrus.Fatalf("数据库连接失败 %s \n", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		logrus.Fatalf("数据库连接池创建失败 %s \n", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(50)
	sqlDB.SetConnMaxLifetime(time.Hour)

	logrus.Infof("数据库连接成功 \n")

	return db
}
