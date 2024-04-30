package ioc

import (
	"context"
	"sync"

	"github.com/wx-up/go-book/pkg/redisx"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/spf13/viper"
	"github.com/wx-up/go-book/internal/repository/cache"

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

func CreateArticleRedisCache(client *redis.Client) cache.ArticleCache {
	client.AddHook(redisx.NewPrometheusHook(prometheus.SummaryOpts{
		Namespace: "wx",
		Subsystem: "go_book",
		Name:      "redis_resp_time",
		Help:      "统计缓存服务的性能数据",
		ConstLabels: map[string]string{
			"biz": "article",
		},
		Objectives: map[float64]float64{},
	}))
	return cache.NewRedisArticleCache(client)
}
