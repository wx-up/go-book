package startup

import (
	"context"
	"sync"

	"github.com/redis/go-redis/v9"
)

var (
	redisClient     redis.Cmdable
	redisClientOnce sync.Once
)

func InitTestRedis() redis.Cmdable {
	redisClientOnce.Do(func() {
		redisClient = redis.NewClient(&redis.Options{
			Addr: "localhost:7379",
		})
		if err := redisClient.Ping(context.Background()).Err(); err != nil {
			panic(err)
		}
	})
	return redisClient
}
