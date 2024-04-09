package ioc

import (
	"fmt"
	"sync/atomic"
	"time"
	"unsafe"

	"github.com/wx-up/go-book/internal/repository/dao"

	"go.uber.org/zap"

	"go.uber.org/zap/zapcore"

	glogger "gorm.io/gorm/logger"

	"github.com/fsnotify/fsnotify"

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
	err := viper.UnmarshalKey("db.mysql", &c)
	if err != nil {
		panic(fmt.Errorf("初始化数据库配置失败：%w", err))
	}
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN: c.DSN,
	}), &gorm.Config{
		Logger: glogger.New(gormLoggerFunc(zap.L().Info), glogger.Config{
			SlowThreshold:             time.Millisecond * 100, // 慢 SQL 阈值（ insert、select 等语句 ）
			Colorful:                  true,                   // 彩色日志
			IgnoreRecordNotFoundError: false,                  // 是否忽略找不到记录的错误
			ParameterizedQueries:      false,                  // 参数值是否填充到 sql 语句中，true 不填充，线上建议是用 true，填充操作是有性能损耗的
			LogLevel:                  glogger.Info,
		}),
	})
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

type gormLoggerFunc func(msg string, fields ...zapcore.Field)

func (f gormLoggerFunc) Printf(m string, args ...interface{}) {
	f(fmt.Sprintf(m, args...))
}

func CreateDBProvider(db *gorm.DB) dao.DBProvider {
	viper.OnConfigChange(func(in fsnotify.Event) {
		type Config struct {
			DSN string `yaml:"dsn"`
		}
		var c Config
		err := viper.UnmarshalKey("db", &c)
		if err != nil {
			return
		}
		newDB, err := gorm.Open(mysql.New(mysql.Config{
			DSN: c.DSN,
		}), &gorm.Config{})
		if err != nil {
			return
		}
		// 原子操作
		atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&db)), unsafe.Pointer(newDB))
	})
	return func() *gorm.DB {
		return db
	}
}
