package global

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Init() error {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: fmt.Sprintf(
			"%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=True&multiStatements=true&loc=Local",
			"root",
			"root",
			"localhost",
			"13316",
			"go_book",
		),
	}), &gorm.Config{})
	if err != nil {
		return err
	}
	DB = db
	return nil
}
