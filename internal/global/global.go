package global

import (
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
	return nil
}
