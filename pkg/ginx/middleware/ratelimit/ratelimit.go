package ratelimit

import (
	_ "embed"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/redis/go-redis/v9"
)

type Builder struct {
	prefix string
	cmd    redis.Cmdable

	slideWindow time.Duration // 窗口大小
	rate        int           // 阈值
}

func NewBuilder(cmd redis.Cmdable, slideWindow time.Duration, rate int) *Builder {
	return &Builder{
		cmd:         cmd,
		slideWindow: slideWindow,
		rate:        rate,
		prefix:      "ip-limiter",
	}
}

func (b *Builder) Prefix(prefix string) *Builder {
	b.prefix = prefix
	return b
}

func (b *Builder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		limited, err := b.limit(ctx)
		if err != nil {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}
		if limited {
			log.Println(err)
			ctx.AbortWithStatus(http.StatusTooManyRequests)
			return
		}
		ctx.Next()
	}
}

//go:embed slide_window.lua
var luaScript string

func (b *Builder) limit(ctx *gin.Context) (bool, error) {
	key := fmt.Sprintf("%s:%s", b.prefix, ctx.ClientIP())
	return b.cmd.Eval(ctx, luaScript, []string{key}, b.slideWindow.Nanoseconds(), b.rate, time.Now().UnixNano()).Bool()
}
