package ioc

import (
	"fmt"
	"time"

	"github.com/spf13/viper"

	"github.com/wx-up/go-book/internal/repository/dao/model"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func CreateMysql() *gorm.DB {
	type Config struct {
		DSN string `yaml:"dsn"`
	}
	var c Config
	err := viper.UnmarshalKey("db", &c)
	if err != nil {
		panic(fmt.Errorf("初始化数据库配置失败：%w", err))
	}
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: c.DSN,
	}), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	_ = model.InitTables(db)
	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}
	sqlDB.SetMaxIdleConns(200)
	sqlDB.SetConnMaxIdleTime(time.Minute * 10)
	sqlDB.SetMaxOpenConns(1200)
	return db
}
