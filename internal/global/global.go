package global

import (
	"time"

	"github.com/wx-up/go-book/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() error {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: config.C.Mysql.DSN,
	}), &gorm.Config{})
	if err != nil {
		return err
	}
	DB = db
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}
	sqlDB.SetMaxIdleConns(200)
	sqlDB.SetConnMaxIdleTime(time.Minute * 10)
	sqlDB.SetMaxOpenConns(1200)
	return nil
}
