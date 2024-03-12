package ioc

import (
	"github.com/redis/go-redis/v9"
	"github.com/wx-up/go-book/config"
)

func CreateRedis() redis.Cmdable {
	return redis.NewClient(&redis.Options{
		Addr:     config.C.Redis.Addr,
		Password: config.C.Redis.Password, // no password set
		DB:       config.C.Redis.DB,       // use default DB
	})
}
