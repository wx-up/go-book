package startup

import (
	"context"
	"log"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB     *gorm.DB
	dbOnce sync.Once
)

func InitTestMysql() *gorm.DB {
	dbOnce.Do(func() {
		db, err := gorm.Open(mysql.New(mysql.Config{
			DSN: "root:root@tcp(localhost:3306)/go_book",
		}), &gorm.Config{})
		if err != nil {
			panic(err)
		}
		sqlDB, err := db.DB()
		if err != nil {
			panic(err)
		}
		for {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			err = sqlDB.PingContext(ctx)
			cancel()
			if err == nil {
				break
			}
			log.Println("等待连接 MySQL", err)
		}

		// model.InitTables(db)

		// 开启 debug 模式
		db.Debug()

		DB = db
	})
	return DB
}
