package ioc

import (
	"context"
	"sync"

	"github.com/spf13/viper"

	"github.com/redis/go-redis/v9"
	"github.com/wx-up/go-book/config"
)

var (
	redisClient     redis.Cmdable
	redisClientOnce sync.Once
)

func CreateRedis() redis.Cmdable {
	redisClientOnce.Do(func() {
		type C struct {
			Addr     string `yaml:"addr"`
			Password string `yaml:"password"`
			DB       int    `yaml:"db"`
		}
		var c C
		if err := viper.UnmarshalKey("redis", &c); err != nil {
			panic(err)
		}
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
