package ioc

import (
	"time"

	"github.com/wx-up/go-book/config"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func CreateMysql() *gorm.DB {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: config.C.Mysql.DSN,
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(200)
	sqlDB.SetConnMaxIdleTime(time.Minute * 10)
	sqlDB.SetMaxOpenConns(1200)
	return db
}
