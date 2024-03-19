package ioc

import (
	"context"
	"sync"

	"github.com/redis/go-redis/v9"
	"github.com/wx-up/go-book/config"
)

var (
	redisClient     redis.Cmdable
	redisClientOnce sync.Once
)

func CreateRedis() redis.Cmdable {
	redisClientOnce.Do(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr:     config.C.Redis.Addr,
			Password: config.C.Redis.Password, // no password set
			DB:       config.C.Redis.DB,       // use default DB
		})
		if err := redisClient.Ping(context.Background()).Err(); err != nil {
			panic(err)
		}
	})
	return redisClient
}
