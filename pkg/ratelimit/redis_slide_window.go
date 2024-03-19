package ratelimit

import (
	"context"
	_ "embed"
	"time"

	"github.com/redis/go-redis/v9"
)

//go:embed slide_window.lua
var luaScript string

type RedisSlideWindowLimiter struct {
	cmd redis.Cmdable

	// internal 窗口大小
	internal time.Duration
	// rate 阈值
	rate int
	// interval=1s rate= 100 表示一秒钟最多100个请求
}

func NewRedisSlideWindowLimiter(cmd redis.Cmdable, internal time.Duration, rate int) Limiter {
	return &RedisSlideWindowLimiter{
		cmd:      cmd,
		internal: internal,
		rate:     rate,
	}
}

func (r *RedisSlideWindowLimiter) Limit(ctx context.Context, key string) (bool, error) {
	return r.cmd.Eval(ctx, luaScript, []string{key}, r.internal.Nanoseconds(), r.rate, time.Now().Nanosecond()).Bool()
}
